#!/usr/bin/env node

// Reads npm audit --json from stdin, cross-references with package.json and package-lock.json
// to produce a structured markdown summary table with false positive detection and dev/prod classification.

const fs = require('fs');
const path = require('path');

function readInput() {
  const filePath = process.argv[2];
  if (filePath) {
    return Promise.resolve(fs.readFileSync(filePath, 'utf8'));
  }
  return new Promise((resolve, reject) => {
    let data = '';
    process.stdin.setEncoding('utf8');
    process.stdin.on('data', chunk => data += chunk);
    process.stdin.on('end', () => resolve(data));
    process.stdin.on('error', reject);
  });
}

function satisfiesRange(version, range) {
  const clean = v => v.replace(/^[~^>=<! ]+/, '').split('-')[0];
  const parts = v => clean(v).split('.').map(Number);

  const ver = parts(version);

  if (range.startsWith('<=')) {
    const upper = parts(range.slice(2));
    return compareParts(ver, upper) <= 0;
  }
  if (range.startsWith('<')) {
    const upper = parts(range.slice(1));
    return compareParts(ver, upper) < 0;
  }
  if (range.includes(' - ')) {
    const [lo, hi] = range.split(' - ').map(s => parts(s.trim()));
    return compareParts(ver, lo) >= 0 && compareParts(ver, hi) <= 0;
  }

  return true;
}

function compareParts(a, b) {
  for (let i = 0; i < 3; i++) {
    const av = a[i] || 0;
    const bv = b[i] || 0;
    if (av !== bv) return av - bv;
  }
  return 0;
}

function findInstalledVersions(lockData, pkgName) {
  const versions = new Set();

  if (lockData.packages) {
    for (const [depPath, info] of Object.entries(lockData.packages)) {
      const segments = depPath.split('node_modules/');
      const name = segments[segments.length - 1];
      if (name === pkgName && info.version) {
        versions.add(info.version);
      }
    }
  }

  return [...versions];
}

function traceParentChain(lockData, pkgName) {
  const parents = new Set();

  if (lockData.packages) {
    for (const [depPath, info] of Object.entries(lockData.packages)) {
      const allDeps = { ...info.dependencies, ...info.devDependencies, ...info.optionalDependencies };
      if (allDeps && allDeps[pkgName]) {
        const segments = depPath.split('node_modules/');
        const parentName = segments[segments.length - 1] || '(root)';
        parents.add(parentName);
      }
    }
  }

  return [...parents].slice(0, 3).join(' > ') || '(unknown)';
}

const SEVERITY_ORDER = { critical: 0, high: 1, moderate: 2, low: 3, info: 4 };

async function main() {
  let auditJson;
  try {
    const raw = await readInput();
    auditJson = JSON.parse(raw);
  } catch (e) {
    console.error('Error: Could not parse npm audit --json from stdin');
    process.exit(1);
  }

  let pkgJson = {};
  try {
    pkgJson = JSON.parse(fs.readFileSync(path.join(process.cwd(), 'package.json'), 'utf8'));
  } catch (e) {
    console.error('Warning: Could not read package.json');
  }

  let lockData = {};
  try {
    lockData = JSON.parse(fs.readFileSync(path.join(process.cwd(), 'package-lock.json'), 'utf8'));
  } catch (e) {
    console.error('Warning: Could not read package-lock.json');
  }

  const devDeps = new Set(Object.keys(pkgJson.devDependencies || {}));
  const prodDeps = new Set(Object.keys(pkgJson.dependencies || {}));

  const vulns = auditJson.vulnerabilities || {};

  const rows = [];

  function isDevOnly(pkgName, visited = new Set()) {
    if (visited.has(pkgName)) return true;
    visited.add(pkgName);

    if (prodDeps.has(pkgName)) return false;
    if (devDeps.has(pkgName)) return true;

    const parents = [];
    if (lockData.packages) {
      for (const [depPath, info] of Object.entries(lockData.packages)) {
        const allDeps = { ...info.dependencies, ...info.optionalDependencies };
        const allDevDeps = info.devDependencies || {};
        if (allDeps[pkgName]) {
          const segments = depPath.split('node_modules/');
          const parentName = segments[segments.length - 1] || '(root)';
          if (parentName === '(root)') {
            if (prodDeps.has(pkgName)) return false;
            if (devDeps.has(pkgName)) return true;
          } else {
            parents.push(parentName);
          }
        }
        if (allDevDeps[pkgName]) {
          const segments = depPath.split('node_modules/');
          const parentName = segments[segments.length - 1] || '(root)';
          if (parentName === '(root)') return true;
          parents.push(parentName);
        }
      }
    }

    if (parents.length === 0) return !prodDeps.has(pkgName);
    return parents.every(p => isDevOnly(p, new Set(visited)));
  }

  for (const [name, vuln] of Object.entries(vulns)) {
    if (!vuln.range) continue;

    const installedVersions = findInstalledVersions(lockData, name);
    const severity = vuln.severity || 'unknown';
    const range = vuln.range;
    const fixAvailable = vuln.fixAvailable;

    let isDirect = prodDeps.has(name) || devDeps.has(name);
    let isDev = isDevOnly(name);

    const declaredVersion = (pkgJson.dependencies || {})[name] || (pkgJson.devDependencies || {})[name] || '';

    const advisories = (vuln.via || [])
      .filter(v => typeof v === 'object' && v.url)
      .map(v => v.url.split('/').pop())
      .slice(0, 3)
      .join(', ') || '';

    for (const installed of installedVersions) {
      const isFalsePositive = !satisfiesRange(installed, range);

      let fixInfo = 'No fix';
      if (fixAvailable === true) {
        fixInfo = 'npm audit fix';
      } else if (typeof fixAvailable === 'object' && fixAvailable.version) {
        fixInfo = `${fixAvailable.name}@${fixAvailable.version}${fixAvailable.isSemVerMajor ? ' (BREAKING)' : ''}`;
      }

      const parentChain = isDirect ? '(direct)' : traceParentChain(lockData, name);

      rows.push({
        name,
        installed,
        declaredVersion,
        range,
        severity,
        type: isDev ? 'dev' : 'prod',
        parentChain,
        isFalsePositive,
        fixInfo,
        advisories,
        severityOrder: SEVERITY_ORDER[severity] || 5,
      });
    }
  }

  rows.sort((a, b) => {
    if (a.type !== b.type) return a.type === 'prod' ? -1 : 1;
    if (a.severityOrder !== b.severityOrder) return a.severityOrder - b.severityOrder;
    return a.name.localeCompare(b.name);
  });

  const summary = { critical: 0, high: 0, moderate: 0, low: 0 };
  const falsePositives = [];
  const realVulns = [];

  for (const row of rows) {
    if (row.isFalsePositive) {
      falsePositives.push(row);
    } else {
      realVulns.push(row);
      summary[row.severity] = (summary[row.severity] || 0) + 1;
    }
  }

  console.log('## npm Audit Summary\n');
  console.log(`| Severity | Count |`);
  console.log(`|----------|-------|`);
  for (const sev of ['critical', 'high', 'moderate', 'low']) {
    if (summary[sev]) console.log(`| ${sev} | ${summary[sev]} |`);
  }
  console.log(`| **Total** | **${realVulns.length}** |`);

  console.log('\n### Vulnerabilities\n');
  console.log('| Package | Declared | Installed | Vuln Range | Severity | Prod/Dev | Parent Chain | Fix | Advisories |');
  console.log('|---------|----------|-----------|------------|----------|----------|-------------|-----|------------|');
  for (const row of realVulns) {
    console.log(`| ${row.name} | ${row.declaredVersion || '-'} | ${row.installed} | ${row.range} | ${row.severity} | ${row.type} | ${row.parentChain} | ${row.fixInfo} | ${row.advisories || '-'} |`);
  }

  if (falsePositives.length > 0) {
    console.log('\n### False Positives (installed version outside vulnerable range)\n');
    console.log('| Package | Installed | Vuln Range | Reason |');
    console.log('|---------|-----------|------------|--------|');
    for (const row of falsePositives) {
      console.log(`| ${row.name} | ${row.installed} | ${row.range} | ${row.installed} is outside ${row.range} |`);
    }
  }
}

main().catch(e => {
  console.error('Error:', e.message);
  process.exit(1);
});
