# ACM (Advanced Cluster Management) Testing for runbook_url Sanitization

This directory contains scripts to test the runbook_url sanitization fix in ACM Observability alerts.

## Files

| File | Purpose |
|------|---------|
| `acm1-runbook-test.sh` | Full ACM setup script with test alerts (use for new clusters) |
| `update-acm-rules.sh` | Update only the alert rules (use for existing ACM clusters) ⭐ **RECOMMENDED** |
| `cleanup-acm-test-rules.sh` | Restore original rules (remove test alerts) |
| `ACM-TESTING-README.md` | This file |

## Prerequisites

- OpenShift cluster with ACM (Advanced Cluster Management) installed
- ACM Observability configured and running
- `oc` CLI logged in with cluster-admin permissions
- The `thanos-ruler-custom-rules` ConfigMap already exists in `open-cluster-management-observability` namespace

## Quick Start (Existing ACM Cluster)

Since you mentioned you already have ACM running with custom rules, use this approach:

```bash
# 1. Navigate to the task directory
cd /Users/emurasak/workspace/harness-official/tasks/OU-1488-sanitize-runbook_url/

# 2. Make the script executable
chmod +x update-acm-rules.sh

# 3. Run the script to add test alerts
./update-acm-rules.sh

# 4. Wait 1-2 minutes for Thanos Ruler to reload
# The test alerts will start firing automatically (they use vector(1) which always evaluates to true)

# 5. Test in ACM Console (see "Testing in ACM Console" section below)

# 6. When done testing, clean up
chmod +x cleanup-acm-test-rules.sh
./cleanup-acm-test-rules.sh
```

## What Gets Added

The `update-acm-rules.sh` script adds a new rule group called `runbook-url-sanitization-tests` with 8 test alerts:

| Alert Name | runbook_url | Expected Behavior |
|------------|-------------|-------------------|
| `TestRunbookValidHTTPS` | `https://runbooks.example.com/cluster-health` | ✅ Should be clickable link |
| `TestRunbookJavaScript` | `javascript:alert('XSS Attack from ACM!')` | ❌ Should be plain text |
| `TestRunbookDataURI` | `data:text/html,<script>alert('XSS from ACM')</script>` | ❌ Should be plain text |
| `TestRunbookCaseSensitive` | `JaVaScRiPt:alert('Case test')` | ❌ Should be plain text |
| `TestRunbookVBScript` | `vbscript:msgbox('ACM XSS')` | ❌ Should be plain text |
| `TestRunbookFileProtocol` | `file:///etc/passwd` | ❌ Should be plain text |
| `TestRunbookRelativeURL` | `/etc/passwd` | ❌ Should be plain text |
| `TestRunbookProtocolRelative` | `//evil.example.com/xss` | ❌ Should be plain text |

**Important**: Each alert also includes the URL in its `description` field to test the `LinkifyExternal` component.

## Testing in ACM Console

### Step 1: Access ACM Observability Alerts

1. Open the ACM Console in your browser
   - Usually at: `https://multicloud-console.apps.<cluster-domain>`
2. Navigate to **Overview** in the left sidebar
3. Click on the **Alerts** tab or card
4. You should see alerts from all your managed clusters

### Step 2: Filter for Test Alerts

1. In the alerts view, use the filter/search box
2. Type: `TestRunbook` or filter by label `test=runbook-sanitization`
3. You should see 8 test alerts (all firing)

### Step 3: Test Each Alert

For each test alert:

1. **Click on the alert name** to view details
2. **Locate the Runbook URL field**
   - It may appear in the "Annotations" section
   - Or as a dedicated "Runbook" field/button
3. **Verify the behavior**:
   - ✅ `TestRunbookValidHTTPS`: URL should be a **clickable link** with external link icon
   - ❌ All other test alerts: URLs should appear as **plain text** (NOT clickable)
4. **Check the Description field**:
   - The description contains the URL as well
   - Valid URLs should be converted to clickable links
   - Malicious URLs should remain as plain text
5. **Security check**:
   - ❌ No JavaScript alerts should pop up
   - ❌ No script execution should occur
   - ❌ No browser errors in console (F12 > Console)

### Step 4: Test in Different Locations

The ACM UI may display alerts in multiple places:

**Location 1: Alerts Overview Page**
- Shows a list/table of all alerts
- May have tooltips or preview panels

**Location 2: Alert Detail Page**
- Full details when you click on an alert
- Shows all annotations including runbook_url and description

**Location 3: Notification Bell** (if available)
- Top-right corner of the ACM Console
- Dropdown showing recent alerts

**Location 4: Grafana** (if integrated)
- ACM Observability may integrate with Grafana
- Alerts may appear in Grafana dashboards
- Verify sanitization works there too

### Step 5: Browser DevTools Inspection

For thorough testing:

1. Open browser DevTools (F12)
2. Go to **Elements/Inspector** tab
3. Find and click on a test alert
4. **For TestRunbookValidHTTPS**:
   - Inspect the runbook URL element
   - Should see: `<a href="https://..." target="_blank" rel="noopener noreferrer">...</a>`
5. **For TestRunbookJavaScript**:
   - Inspect the runbook URL element
   - Should see: `<span>javascript:alert(...)</span>` (or similar non-link element)
   - **NO `<a>` tag, NO `href` attribute**
6. Check the browser **Console** tab:
   - Should be clean (no errors)
   - No JavaScript execution

## Verification Commands

```bash
# View the current ConfigMap
oc get configmap thanos-ruler-custom-rules -n open-cluster-management-observability -o yaml

# Check if test alerts are firing in Thanos Ruler
oc get pods -n open-cluster-management-observability | grep thanos-ruler

# View Thanos Ruler logs (to see if it picked up the new rules)
oc logs -n open-cluster-management-observability <thanos-ruler-pod-name> --tail=50

# Port-forward to Thanos Ruler to check alerts directly (optional)
oc port-forward -n open-cluster-management-observability svc/observability-thanos-ruler 9090:9090
# Then open http://localhost:9090/alerts in your browser
```

## Troubleshooting

### Issue: Test alerts don't appear in ACM Console

**Possible causes:**
1. Thanos Ruler hasn't reloaded the configuration yet
   - **Solution**: Wait 2-3 minutes
2. Thanos Ruler had errors loading the rules
   - **Solution**: Check logs: `oc logs -n open-cluster-management-observability <thanos-ruler-pod>`
3. ConfigMap wasn't updated
   - **Solution**: Verify with `oc get configmap thanos-ruler-custom-rules -n open-cluster-management-observability -o yaml`

### Issue: All URLs appear as plain text (including valid HTTPS)

**This indicates a bug** - the fix may have broken legitimate URLs.
- Report this as a regression
- Check if the monitoring-plugin version is correct
- Check browser console for JavaScript errors

### Issue: JavaScript/Data URI URLs are clickable

**This is a critical security issue** - the sanitization is not working.
- Verify the monitoring-plugin/ACM console version
- Check if this is the patched version
- Report immediately

### Issue: ConfigMap update fails

**Error**: `the server could not find the requested resource`
- The ConfigMap doesn't exist yet
- You may need to run the full `acm1-runbook-test.sh` script first

**Error**: `forbidden: User cannot update configmaps`
- You don't have permissions
- Login as cluster-admin: `oc login --token=... --server=...`

## Cleanup

After testing is complete:

```bash
# Remove test alerts and restore original rules
./cleanup-acm-test-rules.sh

# Verify cleanup
oc get configmap thanos-ruler-custom-rules -n open-cluster-management-observability -o yaml | grep -i testrunbook
# Should return nothing
```

## Differences Between OpenShift Monitoring and ACM Observability

| Aspect | OpenShift Monitoring | ACM Observability |
|--------|---------------------|-------------------|
| **Resource Type** | `PrometheusRule` | `ConfigMap` with embedded YAML |
| **Namespace** | `openshift-monitoring` | `open-cluster-management-observability` |
| **Alert Manager** | Prometheus/Alertmanager | Thanos Ruler |
| **UI** | OpenShift Console > Observe | ACM Console > Overview > Alerts |
| **Scope** | Single cluster | Multi-cluster (hub + managed clusters) |
| **Update Method** | `oc apply -f rule.yaml` | `oc apply -f configmap.yaml` |
| **Reload Time** | ~30 seconds | 1-2 minutes |

## Success Criteria for ACM Testing

- ✅ All 8 test alerts appear in ACM Console
- ✅ `TestRunbookValidHTTPS` shows clickable link in runbook_url field
- ✅ `TestRunbookValidHTTPS` shows clickable link in description field
- ✅ All malicious URL tests (`TestRunbookJavaScript`, `TestRunbookDataURI`, etc.) show plain text in runbook_url field
- ✅ All malicious URLs in description field are plain text (not clickable)
- ✅ No JavaScript execution occurs
- ✅ No XSS attacks can be triggered
- ✅ Browser console shows no errors
- ✅ Sanitization works in all UI locations (alerts list, detail page, notifications)
- ✅ Both OpenShift Monitoring and ACM Observability are protected

## Notes

- The test alerts use `expr: vector(1)` which always evaluates to true, so they will fire immediately
- All test alerts have the label `test: runbook-sanitization` for easy filtering
- The original Watchdog and ClusterCPUHealth-jb alerts are preserved
- Thanos Ruler reloads the configuration automatically when the ConfigMap changes
- You can re-run `update-acm-rules.sh` multiple times safely

## Related Files

- Main test case documentation: `test-case.md`
- OpenShift monitoring test rules: `test-rule-*.yaml`
- Create all test files script: `create-all-test-rules.sh`
