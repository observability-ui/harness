# Fixed CVEs Report

Generated: 2026-06-29

## perses-operator (release-coo-1.5)

### Commits analyzed

- `f7d3b85 [FORK] downgrade to go 1.26.3 to match konflux`
- `7df0910 Merge pull request #8 from rhobs/optimize-watcher-mermory-usage`
- `9cbc2dc Merge pull request #10 from perses/fix-vulnerable-dependencies`
- `ffe645d [FORK] fix vulnerable dependencies`
- `31fee9b [ENHANCEMENT] use only metadata to watch for resources to optimize memory consumption`

### NPM Vulnerabilities Fixed

N/A — no frontend in this project.

### Go Vulnerabilities Fixed

No high/critical Go vulnerabilities were fixed.

## monitoring-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `524eb72 Merge pull request #1018 from jgbernalp/update-vulnerable-dependencies-4.15-26-06-2026`
- `02e9fbd fix: update vulnerable dependencies`
- `a3ff6a4 Merge pull request #1011 from openshift-cherrypick-robot/cherry-pick-1009-to-release-coo-ocp-4.15`
- `f13fb8e fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it | high | LinkifyIt#match scan loop has quadratic algorithmic complexity |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2025-3488 | CVE-2025-22868 | golang.org/x/oauth2 | N/A | Unexpected memory consumption during token parsing in golang.org/x/oauth2 |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## monitoring-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `b4c81a8 Merge pull request #1016 from jgbernalp/update-vulnerable-dependencies-4.19-26-06-2026`
- `f119969 fix: update vulnerable dependencies`
- `f1f7895 Merge pull request #1009 from openshift-cherrypick-robot/cherry-pick-1005-to-release-coo-ocp-4.19`
- `c1eaf8e fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it | high | LinkifyIt#match scan loop has quadratic algorithmic complexity |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2025-3488 | CVE-2025-22868 | golang.org/x/oauth2 | N/A | Unexpected memory consumption during token parsing in golang.org/x/oauth2 |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## monitoring-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `c6e230f Merge pull request #1017 from jgbernalp/update-vulnerable-dependencies-4.22-26-06-2026`
- `bf0f8ab fix: update vulnerable dependencies`
- `f2d03bd Merge pull request #1005 from openshift-cherrypick-robot/cherry-pick-998-to-release-coo-ocp-4.22`
- `8fd3eb2 fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it | high | LinkifyIt#match scan loop has quadratic algorithmic complexity |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2025-3488 | CVE-2025-22868 | golang.org/x/oauth2 | N/A | Unexpected memory consumption during token parsing in golang.org/x/oauth2 |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## logging-view-plugin (release-coo-ocp-4.12)

### Commits analyzed

- `5d3f32f Merge pull request #388 from jgbernalp/update-vulnerable-dependencies-4.12-26-06-2026`
- `72c3401 fix: update vulnerable dependencies`
- `8bc9d84 Merge pull request #384 from jgbernalp/fix-load-more-and-time-precision-coo-4.12`
- `1085301 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`
- `6df12dd refactor: remove unused currentTime`
- `a45ad60 bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-2328-f5f3-gj25 | CVE-2026-33896 | node-forge | high | Forge has a basicConstraints bypass in its certificate chain verification (RFC 5280 violation) |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted | high | flatted vulnerable to unbounded recursion DoS in parse() revive phase |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars | critical | Handlebars.js has JavaScript Injection via AST Type Confusion |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867 | path-to-regexp | high | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block |
| https://github.com/advisories/GHSA-5c6j-r48x-rmvq | N/A | serialize-javascript | high | Serialize JavaScript is Vulnerable to RCE via RegExp.flags and Date.prototype.toISOString() |
| https://github.com/advisories/GHSA-5m6q-g25r-mvwx | CVE-2026-33891 | node-forge | high | Forge has Denial of Service via Infinite Loop in BigInteger.modInverse() with Zero Input |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars | high | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation |
| https://github.com/advisories/GHSA-c27g-q93r-2cwf | CVE-2024-52011 | launch-editor | high | launch-editor vulnerable to command injection via the crafted request on Windows |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch | high | Picomatch has a ReDoS vulnerability via extglob quantifiers |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528 | undici | high | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-ppp5-5v6c-4jwp | CVE-2026-33894 | node-forge | high | Forge has signature forgery in RSA-PKCS due to ASN.1 extra field   |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-q67f-28xg-22rw | CVE-2026-33895 | node-forge | high | Forge has signature forgery in Ed25519 due to missing S > L check |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800 | lodash | high | lodash vulnerable to Code Injection via `_.template` imports key names |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted | high | Prototype Pollution via parse() in NodeJS flatted |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229 | undici | high | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526 | undici | high | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable | high | Immutable is vulnerable to Prototype Pollution |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars | high | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2024-3333 | CVE-2024-45338, GHSA-w32m-9786-jp63 | golang.org/x/net | N/A | Non-linear parsing of case-insensitive content in golang.org/x/net/html |
| GO-2025-3488 | CVE-2025-22868 | golang.org/x/oauth2 | N/A | Unexpected memory consumption during token parsing in golang.org/x/oauth2 |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## logging-view-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `65042b0 Merge pull request #387 from jgbernalp/update-vulnerable-dependencies-4.15-26-06-2026`
- `e726f63 fix: update vulnerable dependencies`
- `352a5e6 Merge pull request #383 from jgbernalp/fix-load-more-and-time-precision-coo-4.15`
- `c1d6289 bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`
- `dec4cdb refactor: remove unused currentTime`
- `5243694 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted | high | flatted vulnerable to unbounded recursion DoS in parse() revive phase |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars | critical | Handlebars.js has JavaScript Injection via AST Type Confusion |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867 | path-to-regexp | high | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block |
| https://github.com/advisories/GHSA-5c6j-r48x-rmvq | N/A | serialize-javascript | high | Serialize JavaScript is Vulnerable to RCE via RegExp.flags and Date.prototype.toISOString() |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars | high | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation |
| https://github.com/advisories/GHSA-c27g-q93r-2cwf | CVE-2024-52011 | launch-editor | high | launch-editor vulnerable to command injection via the crafted request on Windows |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch | high | Picomatch has a ReDoS vulnerability via extglob quantifiers |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528 | undici | high | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800 | lodash | high | lodash vulnerable to Code Injection via `_.template` imports key names |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted | high | Prototype Pollution via parse() in NodeJS flatted |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229 | undici | high | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526 | undici | high | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable | high | Immutable is vulnerable to Prototype Pollution |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars | high | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2025-3488 | CVE-2025-22868 | golang.org/x/oauth2 | N/A | Unexpected memory consumption during token parsing in golang.org/x/oauth2 |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## logging-view-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `7edfaf8 Merge pull request #386 from jgbernalp/update-vulnerable-dependencies-4.22-26-06-2026`
- `377462b fix: update vulnerable dependencies`
- `2da80bd Merge pull request #382 from openshift-cherrypick-robot/cherry-pick-381-to-release-coo-ocp-4.22`
- `d24db30 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`
- `233a1ba refactor: remove unused currentTime`
- `71d13ff bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

No high/critical Go vulnerabilities were fixed.

## distributed-tracing-plugin (release-coo-ocp-4.12)

### Commits analyzed

- `82b87f6 Merge pull request #289 from openshift/update-vulnerable-dependencies-4.12-26-06-2026`
- `d7ad2f5 fix: update vulnerable dependencies`
- `96ba973 Merge pull request #278 from andreasgerstmayr/backport-more-traces-avail`
- `e7eb82a TRACING-5589: Show a note if there are more traces available`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-7p7h-4mm5-852v | CVE-2021-33623 | trim-newlines | high | Uncontrolled Resource Consumption in trim-newlines |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-grv7-fg5c-xmjg | CVE-2024-4068 | braces | high | Uncontrolled resource consumption in braces |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w5p7-h5w8-2hfq | CVE-2020-7753 | trim | high | Regular Expression Denial of Service in trim |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | N/A | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net |
| GO-2025-3595 | CVE-2025-22872 | golang.org/x/net | N/A | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net |
| GO-2026-4440 | CVE-2025-47911 | golang.org/x/net | N/A | Quadratic parsing complexity in golang.org/x/net/html |
| GO-2026-4441 | CVE-2025-58190 | golang.org/x/net | N/A | Infinite parsing loop in golang.org/x/net |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## distributed-tracing-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `a695d12 Merge pull request #288 from openshift/update-vulnerable-dependencies-4.15-26-06-2026`
- `d76d152 fix: update vulnerable dependencies`
- `ee610cc Merge pull request #281 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.15`
- `c781a14 Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## distributed-tracing-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `ab35720 Merge pull request #287 from openshift/update-vulnerable-dependencies-4.19-26-06-2026`
- `803ea83 fix: update vulnerable dependencies`
- `5efd64c Merge pull request #280 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.19`
- `d68e1b9 Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## distributed-tracing-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `d36dd66 Merge pull request #291 from openshift-cherrypick-robot/cherry-pick-285-to-release-coo-ocp-4.22`
- `5dfcc0c fix: update vulnerable dependencies`
- `a1fec23 Merge pull request #279 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.22`
- `6ff2f4c Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734 | undici | high | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697 | undici | high | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |

### Go Vulnerabilities Fixed

| ID | CVE/Aliases | Module | Severity | Summary |
| -- | ----------- | ------ | -------- | ------- |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | N/A | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-5024 | CVE-2026-39824 | golang.org/x/sys | N/A | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows |
| GO-2026-5025 | CVE-2026-42506 | golang.org/x/net | N/A | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html |
| GO-2026-5026 | CVE-2026-39821 | golang.org/x/net | N/A | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna |
| GO-2026-5027 | CVE-2026-42502 | golang.org/x/net | N/A | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html |
| GO-2026-5028 | CVE-2026-25680 | golang.org/x/net | N/A | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html |
| GO-2026-5029 | CVE-2026-25681 | golang.org/x/net | N/A | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html |
| GO-2026-5030 | CVE-2026-27136 | golang.org/x/net | N/A | Invoking duplicate attributes can cause XSS in golang.org/x/net/html |

## troubleshooting-panel-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `4a04933 Merge pull request #254 from openshift/update-vulnerable-dependencies-4.19-26-06-2026`
- `44bd915 fix: update vulnerable dependencies`
- `83f2504 Merge pull request #233 from openshift-cherrypick-robot/cherry-pick-232-to-release-coo-ocp-4.19`
- `b357576 COO-1819:fix: Add TLS min version and cipher suite configuration support`

### NPM Vulnerabilities Fixed

| Advisory | CVE | Package | Severity | Title |
| -------- | --- | ------- | -------- | ----- |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted | high | flatted vulnerable to unbounded recursion DoS in parse() revive phase |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars | critical | Handlebars.js has JavaScript Injection via AST Type Confusion |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867 | path-to-regexp | high | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws | high | ws: Memory exhaustion DoS from tiny fragments and data chunks |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars | high | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch | high | Picomatch has a ReDoS vulnerability via extglob quantifiers |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528 | undici | high | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data | high | form-data: CRLF injection in form-data via unescaped multipart field names and filenames |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp | high | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321 | fast-uri | high | fast-uri vulnerable to path traversal via percent-encoded dot segments |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800 | lodash | high | lodash vulnerable to Code Injection via `_.template` imports key names |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted | high | Prototype Pollution via parse() in NodeJS flatted |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322 | fast-uri | high | fast-uri vulnerable to host confusion via percent-encoded authority delimiters |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229 | undici | high | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526 | undici | high | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici | high | undici WebSocket client vulnerable to denial of service via fragment count bypass |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277 | shell-quote | critical | shell-quote quote() does not escape newlines in object .op values |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable | high | Immutable is vulnerable to Prototype Pollution |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars | high | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars | high | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options |

### Go Vulnerabilities Fixed

No high/critical Go vulnerabilities were fixed.

## troubleshooting-panel-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `94d11ba Merge pull request #247 from openshift-cherrypick-robot/cherry-pick-246-to-release-coo-ocp-4.22`
- `f48f34a fix: COO-1850: Minor UI fixes`
- `ae1a576 Merge pull request #244 from openshift-cherrypick-robot/cherry-pick-243-to-release-coo-ocp-4.22`
- `5365208 fix: COO-1841: prevent Error object from rendering as React child in AgentMenu`

### NPM Vulnerabilities Fixed

No high/critical NPM vulnerabilities were fixed.

### Go Vulnerabilities Fixed

No high/critical Go vulnerabilities were fixed.

## Summary

| Project | Branch | NPM Fixed | Go Fixed | Total | Status |
| ------- | ------ | --------- | -------- | ----- | ------ |
| perses-operator | release-coo-1.5 | 0 | 0 | 0 | ok |
| monitoring-plugin | release-coo-ocp-4.15 | 8 | 13 | 21 | ok |
| monitoring-plugin | release-coo-ocp-4.19 | 8 | 13 | 21 | ok |
| monitoring-plugin | release-coo-ocp-4.22 | 11 | 13 | 24 | ok |
| logging-view-plugin | release-coo-ocp-4.12 | 28 | 14 | 42 | ok |
| logging-view-plugin | release-coo-ocp-4.15 | 24 | 13 | 37 | ok |
| logging-view-plugin | release-coo-ocp-4.22 | 10 | 0 | 10 | ok |
| distributed-tracing-plugin | release-coo-ocp-4.12 | 13 | 12 | 25 | ok |
| distributed-tracing-plugin | release-coo-ocp-4.15 | 10 | 8 | 18 | ok |
| distributed-tracing-plugin | release-coo-ocp-4.19 | 10 | 8 | 18 | ok |
| distributed-tracing-plugin | release-coo-ocp-4.22 | 10 | 8 | 18 | ok |
| troubleshooting-panel-plugin | release-coo-ocp-4.19 | 22 | 0 | 22 | ok |
| troubleshooting-panel-plugin | release-coo-ocp-4.22 | 0 | 0 | 0 | ok |

