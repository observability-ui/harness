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

| ID           | CVE/Aliases    | Module                   | Summary                                                                                                                 |
| ------------ | -------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2026-4599 | CVE-2026-27137 | stdlib                   | Incorrect enforcement of email constraints in crypto/x509                                                               |
| GO-2026-4600 | CVE-2026-27138 | stdlib                   | Panic in name constraint checking for malformed certificates in crypto/x509                                             |
| GO-2026-4601 | CVE-2026-25679 | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139 | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142 | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282 | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289 | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4866 | CVE-2026-33810 | stdlib                   | Case-sensitive excludedSubtrees name constraints cause Auth Bypass in crypto/x509                                       |
| GO-2026-4869 | CVE-2026-32288 | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283 | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814 | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281 | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280 | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836 | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825 | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499 | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826 | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811 | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823 | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820 | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |

## monitoring-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `524eb72 Merge pull request #1018 from jgbernalp/update-vulnerable-dependencies-4.15-26-06-2026`
- `02e9fbd fix: update vulnerable dependencies`
- `a3ff6a4 Merge pull request #1011 from openshift-cherrypick-robot/cherry-pick-1009-to-release-coo-ocp-4.15`
- `f13fb8e fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package               | Severity | Title                                                                                                                                                |
| ------------------------------------------------- | -------------- | --------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it            | high     | LinkifyIt#match scan loop has quadratic algorithmic complexity                                                                                       |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router          | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation                            |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                       |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                    | moderate | ws: Uninitialized memory disclosure                                                                                                                  |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                         |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server    | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                              |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                    | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                        |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                       |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                                                        |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data             | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                             |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server    | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                              |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                           |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                   | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                  |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                            |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend  | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                                                       |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                    | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor         | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                               |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                                                   |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                    |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote           | critical | shell-quote quote() does not escape newlines in object .op values                                                                                    |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2025-3488 | CVE-2025-22868                      | golang.org/x/oauth2      | Unexpected memory consumption during token parsing in golang.org/x/oauth2                                               |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## monitoring-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `b4c81a8 Merge pull request #1016 from jgbernalp/update-vulnerable-dependencies-4.19-26-06-2026`
- `f119969 fix: update vulnerable dependencies`
- `f1f7895 Merge pull request #1009 from openshift-cherrypick-robot/cherry-pick-1005-to-release-coo-ocp-4.19`
- `c1eaf8e fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package               | Severity | Title                                                                                                                                                        |
| ------------------------------------------------- | -------------- | --------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it            | high     | LinkifyIt#match scan loop has quadratic algorithmic complexity                                                                                               |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router          | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation                                    |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                               |
| https://github.com/advisories/GHSA-39q2-94rc-95cp | N/A            | dompurify             | moderate | DOMPurify's ADD_TAGS function form bypasses FORBID_TAGS due to short-circuit evaluation                                                                      |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core           | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                                |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                    | moderate | ws: Uninitialized memory disclosure                                                                                                                          |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                                 |
| https://github.com/advisories/GHSA-76mc-f452-cxcm | N/A            | dompurify             | moderate | DOMPurify: Hook mutation of `data.allowedTags` / `data.allowedAttributes` permanently pollutes `DEFAULT_ALLOWED_TAGS` / `DEFAULT_ALLOWED_ATTR`               |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                    | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                                |
| https://github.com/advisories/GHSA-cj63-jhhr-wcxv | N/A            | dompurify             | moderate | DOMPurify USE_PROFILES prototype pollution allows event handlers                                                                                             |
| https://github.com/advisories/GHSA-cjmm-f4jc-qw8r | N/A            | dompurify             | moderate | DOMPurify ADD_ATTR predicate skips URI validation                                                                                                            |
| https://github.com/advisories/GHSA-cmwh-pvxp-8882 | N/A            | dompurify             | moderate | DOMPurify: Permanent `ALLOWED_ATTR` pollution via `setConfig()` bypassing the hook clone-guard (incomplete fix of the 3.4.7 hook-pollution patch)            |
| https://github.com/advisories/GHSA-crv5-9vww-q3g8 | CVE-2026-41239 | dompurify             | moderate | DOMPurify has a SAFE_FOR_TEMPLATES bypass in RETURN_DOM mode                                                                                                 |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                               |
| https://github.com/advisories/GHSA-gvmj-g25r-r7wr | N/A            | dompurify             | low      | DOMPurify: SAFE_FOR_TEMPLATES bypass - template expressions survive sanitization inside <template> content when using DOM output modes                       |
| https://github.com/advisories/GHSA-h7mw-gpvr-xq4m | CVE-2026-41240 | dompurify             | moderate | DOMPurify: FORBID_TAGS bypassed by function-based ADD_TAGS predicate (asymmetry with FORBID_ATTR fix)                                                        |
| https://github.com/advisories/GHSA-h8r8-wccr-v5f2 | N/A            | dompurify             | moderate | DOMPurify is vulnerable to mutation-XSS via Re-Contextualization                                                                                             |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                                                                |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data             | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                                     |
| https://github.com/advisories/GHSA-hpcv-96wg-7vj8 | CVE-2026-49458 | dompurify             | moderate | DOMPurify: Cross-realm IN_PLACE sanitization leaves executable markup intact via realm-bound `instanceof` checks                                             |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server    | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                                      |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                                   |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                   | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                          |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                                    |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend  | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                                                               |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                    | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set         |
| https://github.com/advisories/GHSA-r47g-fvhr-h676 | CVE-2026-49459 | dompurify             | moderate | DOMPurify: IN_PLACE mode preserves attributes of a clobbered root element, allowing XSS via attacker-controlled root DOM                                     |
| https://github.com/advisories/GHSA-rp9w-3fw7-7cwq | CVE-2026-49978 | dompurify             | moderate | DOMPurify IN_PLACE Sanitization Bypass via Attached Shadow Root Inside <template>.content                                                                    |
| https://github.com/advisories/GHSA-v2wj-7wpq-c8vv | CVE-2026-0540  | dompurify             | moderate | DOMPurify contains a Cross-site Scripting vulnerability                                                                                                      |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor         | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                                       |
| https://github.com/advisories/GHSA-v9jr-rg53-9pgp | CVE-2026-41238 | dompurify             | moderate | DOMPurify: Prototype Pollution to XSS Bypass via CUSTOM_ELEMENT_HANDLING Fallback                                                                            |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                                                           |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                            |
| https://github.com/advisories/GHSA-vxr8-fq34-vvx9 | N/A            | dompurify             | low      | DOMPurify: Trusted Types policy survives `clearConfig()` and can poison later `RETURN_TRUSTED_TYPE` output                                                   |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote           | critical | shell-quote quote() does not escape newlines in object .op values                                                                                            |
| https://github.com/advisories/GHSA-x4vx-rjvf-j5p4 | N/A            | dompurify             | low      | DOMPurify: `IN_PLACE` mode trusts attacker-controlled `nodeName` on live non-form nodes, allowing script retention and XSS via attacker-supplied DOM objects |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2025-3488 | CVE-2025-22868                      | golang.org/x/oauth2      | Unexpected memory consumption during token parsing in golang.org/x/oauth2                                               |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## monitoring-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `c6e230f Merge pull request #1017 from jgbernalp/update-vulnerable-dependencies-4.22-26-06-2026`
- `bf0f8ab fix: update vulnerable dependencies`
- `f2d03bd Merge pull request #1005 from openshift-cherrypick-robot/cherry-pick-998-to-release-coo-ocp-4.22`
- `8fd3eb2 fix(perses): fall back to metadata.name when dashboard display name is missing`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                                                        |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------------------------------------------ |
| https://github.com/advisories/GHSA-22p9-wv53-3rq4 | CVE-2026-48801 | linkify-it                               | high     | LinkifyIt#match scan loop has quadratic algorithmic complexity                                                                                               |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                               |
| https://github.com/advisories/GHSA-39q2-94rc-95cp | N/A            | dompurify                                | moderate | DOMPurify's ADD_TAGS function form bypasses FORBID_TAGS due to short-circuit evaluation                                                                      |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                                |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                                                          |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                                 |
| https://github.com/advisories/GHSA-76mc-f452-cxcm | N/A            | dompurify                                | moderate | DOMPurify: Hook mutation of `data.allowedTags` / `data.allowedAttributes` permanently pollutes `DEFAULT_ALLOWED_TAGS` / `DEFAULT_ALLOWED_ATTR`               |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                                      |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                                |
| https://github.com/advisories/GHSA-cj63-jhhr-wcxv | N/A            | dompurify                                | moderate | DOMPurify USE_PROFILES prototype pollution allows event handlers                                                                                             |
| https://github.com/advisories/GHSA-cjmm-f4jc-qw8r | N/A            | dompurify                                | moderate | DOMPurify ADD_ATTR predicate skips URI validation                                                                                                            |
| https://github.com/advisories/GHSA-cmwh-pvxp-8882 | N/A            | dompurify                                | moderate | DOMPurify: Permanent `ALLOWED_ATTR` pollution via `setConfig()` bypassing the hook clone-guard (incomplete fix of the 3.4.7 hook-pollution patch)            |
| https://github.com/advisories/GHSA-crv5-9vww-q3g8 | CVE-2026-41239 | dompurify                                | moderate | DOMPurify has a SAFE_FOR_TEMPLATES bypass in RETURN_DOM mode                                                                                                 |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                                                             |
| https://github.com/advisories/GHSA-g7r4-m6w7-qqqr | N/A            | esbuild                                  | low      | esbuild allows arbitrary file read when running the development server on Windows                                                                            |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                               |
| https://github.com/advisories/GHSA-gvmj-g25r-r7wr | N/A            | dompurify                                | low      | DOMPurify: SAFE_FOR_TEMPLATES bypass - template expressions survive sanitization inside <template> content when using DOM output modes                       |
| https://github.com/advisories/GHSA-h7mw-gpvr-xq4m | CVE-2026-41240 | dompurify                                | moderate | DOMPurify: FORBID_TAGS bypassed by function-based ADD_TAGS predicate (asymmetry with FORBID_ATTR fix)                                                        |
| https://github.com/advisories/GHSA-h8r8-wccr-v5f2 | N/A            | dompurify                                | moderate | DOMPurify is vulnerable to mutation-XSS via Re-Contextualization                                                                                             |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                                                                |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                                     |
| https://github.com/advisories/GHSA-hpcv-96wg-7vj8 | CVE-2026-49458 | dompurify                                | moderate | DOMPurify: Cross-realm IN_PLACE sanitization leaves executable markup intact via realm-bound `instanceof` checks                                             |
| https://github.com/advisories/GHSA-jxxr-4gwj-5jf2 | CVE-2026-45149 | brace-expansion                          | moderate | brace-expansion: Large numeric range defeats documented `max` DoS protection                                                                                 |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                                      |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                                   |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                          |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                                    |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                                                       |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                                       | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set         |
| https://github.com/advisories/GHSA-r47g-fvhr-h676 | CVE-2026-49459 | dompurify                                | moderate | DOMPurify: IN_PLACE mode preserves attributes of a clobbered root element, allowing XSS via attacker-controlled root DOM                                     |
| https://github.com/advisories/GHSA-rp9w-3fw7-7cwq | CVE-2026-49978 | dompurify                                | moderate | DOMPurify IN_PLACE Sanitization Bypass via Attached Shadow Root Inside <template>.content                                                                    |
| https://github.com/advisories/GHSA-v2wj-7wpq-c8vv | CVE-2026-0540  | dompurify                                | moderate | DOMPurify contains a Cross-site Scripting vulnerability                                                                                                      |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                                                               |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                                       |
| https://github.com/advisories/GHSA-v9jr-rg53-9pgp | CVE-2026-41238 | dompurify                                | moderate | DOMPurify: Prototype Pollution to XSS Bypass via CUSTOM_ELEMENT_HANDLING Fallback                                                                            |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                                                           |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                            |
| https://github.com/advisories/GHSA-vxr8-fq34-vvx9 | N/A            | dompurify                                | low      | DOMPurify: Trusted Types policy survives `clearConfig()` and can poison later `RETURN_TRUSTED_TYPE` output                                                   |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                                                            |
| https://github.com/advisories/GHSA-x4vx-rjvf-j5p4 | N/A            | dompurify                                | low      | DOMPurify: `IN_PLACE` mode trusts attacker-controlled `nodeName` on live non-form nodes, allowing script retention and XSS via attacker-supplied DOM objects |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2025-3488 | CVE-2025-22868                      | golang.org/x/oauth2      | Unexpected memory consumption during token parsing in golang.org/x/oauth2                                               |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## logging-view-plugin (release-coo-ocp-4.12)

### Commits analyzed

- `5d3f32f Merge pull request #388 from jgbernalp/update-vulnerable-dependencies-4.12-26-06-2026`
- `72c3401 fix: update vulnerable dependencies`
- `8bc9d84 Merge pull request #384 from jgbernalp/fix-load-more-and-time-precision-coo-4.12`
- `1085301 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`
- `6df12dd refactor: remove unused currentTime`
- `a45ad60 bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                                                |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-2328-f5f3-gj25 | CVE-2026-33896 | node-forge                               | high     | Forge has a basicConstraints bypass in its certificate chain verification (RFC 5280 violation)                                                       |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted                                  | high     | flatted vulnerable to unbounded recursion DoS in parse() revive phase                                                                                |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router                             | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation                            |
| https://github.com/advisories/GHSA-2mjp-6q6p-2qxm | CVE-2026-1525  | undici                                   | moderate | Undici has an HTTP Request/Response Smuggling issue                                                                                                  |
| https://github.com/advisories/GHSA-2qvq-rjwj-gvw9 | CVE-2026-33916 | handlebars                               | moderate | Handlebars.js has Prototype Pollution Leading to XSS through Partial Template Injection                                                              |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars                               | critical | Handlebars.js has JavaScript Injection via AST Type Confusion                                                                                        |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                       |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867  | path-to-regexp                           | high     | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters                                                      |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block                                                            |
| https://github.com/advisories/GHSA-3v7f-55p6-f55p | CVE-2026-33672 | picomatch                                | moderate | Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching                                                                |
| https://github.com/advisories/GHSA-442j-39wm-28r2 | N/A            | handlebars                               | low      | Handlebars.js has a Property Access Validation Bypass in container.lookup                                                                            |
| https://github.com/advisories/GHSA-4992-7rv2-5pvq | CVE-2026-1527  | undici                                   | moderate | Undici has CRLF Injection in undici via `upgrade` option                                                                                             |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                        |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                                                  |
| https://github.com/advisories/GHSA-5c6j-r48x-rmvq | N/A            | serialize-javascript                     | high     | Serialize JavaScript is Vulnerable to RCE via RegExp.flags and Date.prototype.toISOString()                                                          |
| https://github.com/advisories/GHSA-5m6q-g25r-mvwx | CVE-2026-33891 | node-forge                               | high     | Forge has Denial of Service via Infinite Loop in BigInteger.modInverse() with Zero Input                                                             |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                         |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                              |
| https://github.com/advisories/GHSA-7rx3-28cr-v5wh | N/A            | handlebars                               | moderate | Handlebars.js has a Prototype Method Access Control Gap via Missing **lookupSetter** Blocklist Entry                                                 |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                        |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars                               | high     | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation                                                           |
| https://github.com/advisories/GHSA-c27g-q93r-2cwf | CVE-2024-52011 | launch-editor                            | high     | launch-editor vulnerable to command injection via the crafted request on Windows                                                                     |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch                                | high     | Picomatch has a ReDoS vulnerability via extglob quantifiers                                                                                          |
| https://github.com/advisories/GHSA-f23m-r3pf-42rh | CVE-2026-2950  | lodash                                   | moderate | lodash vulnerable to Prototype Pollution via array path bypass in `_.unset` and `_.omit`                                                             |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528  | undici                                   | high     | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client                                                                    |
| https://github.com/advisories/GHSA-f886-m6hf-6m8v | CVE-2026-33750 | brace-expansion                          | moderate | brace-expansion: Zero-step sequence causes process hang and memory exhaustion                                                                        |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                                                     |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                       |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                                                         |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                             |
| https://github.com/advisories/GHSA-mh29-5h37-fv8m | CVE-2025-64718 | js-yaml                                  | moderate | js-yaml has prototype pollution in merge (<<)                                                                                                        |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                              |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                           |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                  |
| https://github.com/advisories/GHSA-ppp5-5v6c-4jwp | CVE-2026-33894 | node-forge                               | high     | Forge has signature forgery in RSA-PKCS due to ASN.1 extra field                                                                                     |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                                               |
| https://github.com/advisories/GHSA-q67f-28xg-22rw | CVE-2026-33895 | node-forge                               | high     | Forge has signature forgery in Ed25519 due to missing S > L check                                                                                    |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend                     | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                                                       |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                                       | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set |
| https://github.com/advisories/GHSA-qj8w-gfj5-8c6v | CVE-2026-34043 | serialize-javascript                     | moderate | Serialize JavaScript has CPU Exhaustion Denial of Service via crafted array-like objects                                                             |
| https://github.com/advisories/GHSA-qx2v-qp2m-jg93 | CVE-2026-41305 | postcss                                  | moderate | PostCSS has XSS via Unescaped </style> in its CSS Stringify Output                                                                                   |
| https://github.com/advisories/GHSA-r4q5-vmmm-2653 | N/A            | follow-redirects                         | moderate | follow-redirects leaks Custom Authentication Headers to Cross-Domain Redirect Targets                                                                |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800  | lodash                                   | high     | lodash vulnerable to Code Injection via `_.template` imports key names                                                                               |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted                                  | high     | Prototype Pollution via parse() in NodeJS flatted                                                                                                    |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                                                       |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                               |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229  | undici                                   | high     | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation                                                  |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526  | undici                                   | high     | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression                                                                |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                    |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                                                    |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable                                | high     | Immutable is vulnerable to Prototype Pollution                                                                                                       |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial                                              |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars                               | high     | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options                                                            |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2024-3333 | CVE-2024-45338, GHSA-w32m-9786-jp63 | golang.org/x/net         | Non-linear parsing of case-insensitive content in golang.org/x/net/html                                                 |
| GO-2025-3488 | CVE-2025-22868                      | golang.org/x/oauth2      | Unexpected memory consumption during token parsing in golang.org/x/oauth2                                               |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## logging-view-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `65042b0 Merge pull request #387 from jgbernalp/update-vulnerable-dependencies-4.15-26-06-2026`
- `e726f63 fix: update vulnerable dependencies`
- `352a5e6 Merge pull request #383 from jgbernalp/fix-load-more-and-time-precision-coo-4.15`
- `c1d6289 bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`
- `dec4cdb refactor: remove unused currentTime`
- `5243694 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                     |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted                                  | high     | flatted vulnerable to unbounded recursion DoS in parse() revive phase                                                     |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router                             | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation |
| https://github.com/advisories/GHSA-2mjp-6q6p-2qxm | CVE-2026-1525  | undici                                   | moderate | Undici has an HTTP Request/Response Smuggling issue                                                                       |
| https://github.com/advisories/GHSA-2qvq-rjwj-gvw9 | CVE-2026-33916 | handlebars                               | moderate | Handlebars.js has Prototype Pollution Leading to XSS through Partial Template Injection                                   |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars                               | critical | Handlebars.js has JavaScript Injection via AST Type Confusion                                                             |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                            |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867  | path-to-regexp                           | high     | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters                           |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block                                 |
| https://github.com/advisories/GHSA-3v7f-55p6-f55p | CVE-2026-33672 | picomatch                                | moderate | Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching                                     |
| https://github.com/advisories/GHSA-442j-39wm-28r2 | N/A            | handlebars                               | low      | Handlebars.js has a Property Access Validation Bypass in container.lookup                                                 |
| https://github.com/advisories/GHSA-4992-7rv2-5pvq | CVE-2026-1527  | undici                                   | moderate | Undici has CRLF Injection in undici via `upgrade` option                                                                  |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                             |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                       |
| https://github.com/advisories/GHSA-5c6j-r48x-rmvq | N/A            | serialize-javascript                     | high     | Serialize JavaScript is Vulnerable to RCE via RegExp.flags and Date.prototype.toISOString()                               |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass              |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                   |
| https://github.com/advisories/GHSA-7rx3-28cr-v5wh | N/A            | handlebars                               | moderate | Handlebars.js has a Prototype Method Access Control Gap via Missing **lookupSetter** Blocklist Entry                      |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                             |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars                               | high     | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation                                |
| https://github.com/advisories/GHSA-c27g-q93r-2cwf | CVE-2024-52011 | launch-editor                            | high     | launch-editor vulnerable to command injection via the crafted request on Windows                                          |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch                                | high     | Picomatch has a ReDoS vulnerability via extglob quantifiers                                                               |
| https://github.com/advisories/GHSA-f23m-r3pf-42rh | CVE-2026-2950  | lodash                                   | moderate | lodash vulnerable to Prototype Pollution via array path bypass in `_.unset` and `_.omit`                                  |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528  | undici                                   | high     | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client                                         |
| https://github.com/advisories/GHSA-f886-m6hf-6m8v | CVE-2026-33750 | brace-expansion                          | moderate | brace-expansion: Zero-step sequence causes process hang and memory exhaustion                                             |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                          |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                            |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                  |
| https://github.com/advisories/GHSA-jxxr-4gwj-5jf2 | CVE-2026-45149 | brace-expansion                          | moderate | brace-expansion: Large numeric range defeats documented `max` DoS protection                                              |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                   |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                       |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                    |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend                     | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                            |
| https://github.com/advisories/GHSA-qj8w-gfj5-8c6v | CVE-2026-34043 | serialize-javascript                     | moderate | Serialize JavaScript has CPU Exhaustion Denial of Service via crafted array-like objects                                  |
| https://github.com/advisories/GHSA-qx2v-qp2m-jg93 | CVE-2026-41305 | postcss                                  | moderate | PostCSS has XSS via Unescaped </style> in its CSS Stringify Output                                                        |
| https://github.com/advisories/GHSA-r4q5-vmmm-2653 | N/A            | follow-redirects                         | moderate | follow-redirects leaks Custom Authentication Headers to Cross-Domain Redirect Targets                                     |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800  | lodash                                   | high     | lodash vulnerable to Code Injection via `_.template` imports key names                                                    |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted                                  | high     | Prototype Pollution via parse() in NodeJS flatted                                                                         |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                            |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                    |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229  | undici                                   | high     | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation                       |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526  | undici                                   | high     | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression                                     |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                         |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                         |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable                                | high     | Immutable is vulnerable to Prototype Pollution                                                                            |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial                   |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars                               | high     | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options                                 |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2025-3488 | CVE-2025-22868                      | golang.org/x/oauth2      | Unexpected memory consumption during token parsing in golang.org/x/oauth2                                               |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## logging-view-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `7edfaf8 Merge pull request #386 from jgbernalp/update-vulnerable-dependencies-4.22-26-06-2026`
- `377462b fix: update vulnerable dependencies`
- `2da80bd Merge pull request #382 from openshift-cherrypick-robot/cherry-pick-381-to-release-coo-ocp-4.22`
- `d24db30 refactor: add nano seconds unit to variable names, use nanoseconds in all loki client functions`
- `233a1ba refactor: remove unused currentTime`
- `71d13ff bugfix: add nanosecond precision to load more calculation to avoid loosing logs on high volume results`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                                                |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                       |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                        |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                                                  |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                         |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                              |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                        |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                                                     |
| https://github.com/advisories/GHSA-g7r4-m6w7-qqqr | N/A            | esbuild                                  | low      | esbuild allows arbitrary file read when running the development server on Windows                                                                    |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                       |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                                                        |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                             |
| https://github.com/advisories/GHSA-jxxr-4gwj-5jf2 | CVE-2026-45149 | brace-expansion                          | moderate | brace-expansion: Large numeric range defeats documented `max` DoS protection                                                                         |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                              |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                           |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                  |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                            |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                                               |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                                       | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                                                       |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                               |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                                                   |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                    |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                                                    |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                    | Module | Summary                                                                                                  |
| ------------ | ------------------------------ | ------ | -------------------------------------------------------------------------------------------------------- |
| GO-2026-4337 | CVE-2025-68121                 | stdlib | Unexpected session resumption in crypto/tls                                                              |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730 | stdlib | Handshake messages may be processed at the incorrect encryption level in crypto/tls                      |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726 | stdlib | Memory exhaustion in query parameter parsing in net/url                                                  |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728 | stdlib | Excessive CPU consumption when building archive index in archive/zip                                     |
| GO-2026-4601 | CVE-2026-25679                 | stdlib | Incorrect parsing of IPv6 host literals in net/url                                                       |
| GO-2026-4602 | CVE-2026-27139                 | stdlib | FileInfo can escape from a Root in os                                                                    |
| GO-2026-4603 | CVE-2026-27142                 | stdlib | URLs in meta content attribute actions are not escaped in html/template                                  |
| GO-2026-4864 | CVE-2026-32282                 | stdlib | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                        |
| GO-2026-4865 | CVE-2026-32289                 | stdlib | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                |
| GO-2026-4869 | CVE-2026-32288                 | stdlib | Unbounded allocation for old GNU sparse in archive/tar                                                   |
| GO-2026-4870 | CVE-2026-32283                 | stdlib | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls |
| GO-2026-4946 | CVE-2026-32281                 | stdlib | Inefficient policy validation in crypto/x509                                                             |
| GO-2026-4947 | CVE-2026-32280                 | stdlib | Unexpected work during chain building in crypto/x509                                                     |
| GO-2026-4971 | CVE-2026-39836                 | stdlib | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                    |
| GO-2026-4976 | CVE-2026-39825                 | stdlib | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil           |
| GO-2026-4977 | CVE-2026-42499                 | stdlib | Quadratic string concatenation in consumePhrase in net/mail                                              |
| GO-2026-4980 | CVE-2026-39826                 | stdlib | Escaper bypass leads to XSS in html/template                                                             |
| GO-2026-4981 | CVE-2026-33811                 | stdlib | Crash when handling long CNAME response in net                                                           |
| GO-2026-4982 | CVE-2026-39823                 | stdlib | Bypass of meta content URL escaping causes XSS in html/template                                          |
| GO-2026-4986 | CVE-2026-39820                 | stdlib | Quadratic string concatentation in consumeComment in net/mail                                            |

## distributed-tracing-plugin (release-coo-ocp-4.12)

### Commits analyzed

- `82b87f6 Merge pull request #289 from openshift/update-vulnerable-dependencies-4.12-26-06-2026`
- `d7ad2f5 fix: update vulnerable dependencies`
- `96ba973 Merge pull request #278 from andreasgerstmayr/backport-more-traces-avail`
- `e7eb82a TRACING-5589: Show a note if there are more traces available`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                        |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------ |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                               |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                |
| https://github.com/advisories/GHSA-566m-qj78-rww5 | CVE-2021-23382 | postcss                                  | moderate | Regular Expression Denial of Service in postcss                                                              |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                          |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass |
| https://github.com/advisories/GHSA-67mh-4wv8-2f99 | N/A            | esbuild                                  | moderate | esbuild enables any website to send any requests to the development server and read the response             |
| https://github.com/advisories/GHSA-7fh5-64p2-3v2j | CVE-2023-44270 | postcss                                  | moderate | PostCSS line return parsing error                                                                            |
| https://github.com/advisories/GHSA-7p7h-4mm5-852v | CVE-2021-33623 | trim-newlines                            | high     | Uncontrolled Resource Consumption in trim-newlines                                                           |
| https://github.com/advisories/GHSA-952p-6rrq-rcjv | CVE-2024-4067  | micromatch                               | moderate | Regular Expression Denial of Service (ReDoS) in micromatch                                                   |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input             |
| https://github.com/advisories/GHSA-g3ch-rx76-35fx | CVE-2024-6783  | vue-template-compiler                    | moderate | vue-template-compiler vulnerable to client-side Cross-Site Scripting (XSS)                                   |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching               |
| https://github.com/advisories/GHSA-grv7-fg5c-xmjg | CVE-2024-4068  | braces                                   | high     | Uncontrolled resource consumption in braces                                                                  |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                 |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                     |
| https://github.com/advisories/GHSA-mh29-5h37-fv8m | CVE-2025-64718 | js-yaml                                  | moderate | js-yaml has prototype pollution in merge (<<)                                                                |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                   |
| https://github.com/advisories/GHSA-pfrx-2q88-qq97 | CVE-2022-33987 | got                                      | moderate | Got allows a redirect to a UNIX socket                                                                       |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                          |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                    |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                       |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend                     | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                               |
| https://github.com/advisories/GHSA-qx2v-qp2m-jg93 | CVE-2026-41305 | postcss                                  | moderate | PostCSS has XSS via Unescaped </style> in its CSS Stringify Output                                           |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                               |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                       |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent           |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                            |
| https://github.com/advisories/GHSA-w5p7-h5w8-2hfq | CVE-2020-7753  | trim                                     | high     | Regular Expression Denial of Service in trim                                                                 |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                            |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                         | Module                   | Summary                                                                                                                 |
| ------------ | ----------------------------------- | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2025-3503 | CVE-2025-22870, GHSA-qxp5-gwg8-xv66 | golang.org/x/net, stdlib | HTTP Proxy bypass using IPv6 Zone IDs in golang.org/x/net                                                               |
| GO-2025-3595 | CVE-2025-22872                      | golang.org/x/net         | Incorrect Neutralization of Input During Web Page Generation in x/net in golang.org/x/net                               |
| GO-2026-4337 | CVE-2025-68121                      | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730      | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726      | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728      | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4440 | CVE-2025-47911                      | golang.org/x/net         | Quadratic parsing complexity in golang.org/x/net/html                                                                   |
| GO-2026-4441 | CVE-2025-58190                      | golang.org/x/net         | Infinite parsing loop in golang.org/x/net                                                                               |
| GO-2026-4601 | CVE-2026-25679                      | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                      | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                      | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                      | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                      | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                      | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                      | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                      | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                      | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                      | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                      | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                      | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                      | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                      | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                      | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                      | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                      | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                      | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                      | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                      | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                      | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                      | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                      | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                      | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## distributed-tracing-plugin (release-coo-ocp-4.15)

### Commits analyzed

- `a695d12 Merge pull request #288 from openshift/update-vulnerable-dependencies-4.15-26-06-2026`
- `d76d152 fix: update vulnerable dependencies`
- `ee610cc Merge pull request #281 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.15`
- `c781a14 Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                     |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router                             | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                            |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                             |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                       |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass              |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                   |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                             |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                          |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                            |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                              |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                             |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                  |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                   |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                       |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                 |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                    |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend                     | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                            |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                            |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                    |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                        |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                         |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                         |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                    | Module                   | Summary                                                                                                                 |
| ------------ | ------------------------------ | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2026-4337 | CVE-2025-68121                 | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730 | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726 | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728 | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4601 | CVE-2026-25679                 | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                 | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                 | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                 | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                 | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                 | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                 | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                 | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                 | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                 | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                 | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                 | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                 | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                 | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                 | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                 | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                 | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                 | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                 | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                 | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                 | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                 | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                 | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                 | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## distributed-tracing-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `ab35720 Merge pull request #287 from openshift/update-vulnerable-dependencies-4.19-26-06-2026`
- `803ea83 fix: update vulnerable dependencies`
- `5efd64c Merge pull request #280 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.19`
- `d68e1b9 Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                     |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router                             | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                            |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                             |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                       |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass              |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                   |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                             |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                          |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                            |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                              |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                             |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                  |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                   |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                       |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                 |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                    |
| https://github.com/advisories/GHSA-q89c-q3h5-w34g | CVE-2026-41691 | i18next-http-backend                     | moderate | i18next-http-backend has Path Traversal & URL Injection via Unsanitised lng/ns                                            |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                            |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                    |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                        |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                         |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                         |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                    | Module                   | Summary                                                                                                                 |
| ------------ | ------------------------------ | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2026-4337 | CVE-2025-68121                 | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730 | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726 | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728 | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4601 | CVE-2026-25679                 | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                 | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                 | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                 | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                 | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                 | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                 | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                 | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                 | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                 | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                 | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                 | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                 | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                 | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                 | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                 | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                 | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                 | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                 | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                 | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                 | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                 | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                 | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                 | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## distributed-tracing-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `d36dd66 Merge pull request #291 from openshift-cherrypick-robot/cherry-pick-285-to-release-coo-ocp-4.22`
- `5dfcc0c fix: update vulnerable dependencies`
- `a1fec23 Merge pull request #279 from openshift-cherrypick-robot/cherry-pick-275-to-release-coo-ocp-4.22`
- `6ff2f4c Fix: Open documentation link in a new browser tab`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                                                |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                       |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                        |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                                                  |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                         |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                              |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                        |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                                                     |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                       |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                                                         |
| https://github.com/advisories/GHSA-hm92-r4w5-c3mj | CVE-2026-6734  | undici                                   | high     | undici vulnerable to cross-origin request routing via SOCKS5 proxy pool reuse                                                                        |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                             |
| https://github.com/advisories/GHSA-jxxr-4gwj-5jf2 | CVE-2026-45149 | brace-expansion                          | moderate | brace-expansion: Large numeric range defeats documented `max` DoS protection                                                                         |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                              |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                           |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                  |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                            |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                                               |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                                       | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                                                       |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                               |
| https://github.com/advisories/GHSA-vmh5-mc38-953g | CVE-2026-9697  | undici                                   | high     | undici vulnerable to TLS certificate validation bypass via dropped requestTls in SOCKS5 ProxyAgent                                                   |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                    |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                                                    |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                    | Module                   | Summary                                                                                                                 |
| ------------ | ------------------------------ | ------------------------ | ----------------------------------------------------------------------------------------------------------------------- |
| GO-2026-4337 | CVE-2025-68121                 | stdlib                   | Unexpected session resumption in crypto/tls                                                                             |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730 | stdlib                   | Handshake messages may be processed at the incorrect encryption level in crypto/tls                                     |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726 | stdlib                   | Memory exhaustion in query parameter parsing in net/url                                                                 |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728 | stdlib                   | Excessive CPU consumption when building archive index in archive/zip                                                    |
| GO-2026-4601 | CVE-2026-25679                 | stdlib                   | Incorrect parsing of IPv6 host literals in net/url                                                                      |
| GO-2026-4602 | CVE-2026-27139                 | stdlib                   | FileInfo can escape from a Root in os                                                                                   |
| GO-2026-4603 | CVE-2026-27142                 | stdlib                   | URLs in meta content attribute actions are not escaped in html/template                                                 |
| GO-2026-4864 | CVE-2026-32282                 | stdlib                   | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                                       |
| GO-2026-4865 | CVE-2026-32289                 | stdlib                   | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                               |
| GO-2026-4869 | CVE-2026-32288                 | stdlib                   | Unbounded allocation for old GNU sparse in archive/tar                                                                  |
| GO-2026-4870 | CVE-2026-32283                 | stdlib                   | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls                |
| GO-2026-4918 | CVE-2026-33814                 | golang.org/x/net, stdlib | Infinite loop in HTTP/2 transport when given bad SETTINGS_MAX_FRAME_SIZE in net/http/internal/http2 in golang.org/x/net |
| GO-2026-4946 | CVE-2026-32281                 | stdlib                   | Inefficient policy validation in crypto/x509                                                                            |
| GO-2026-4947 | CVE-2026-32280                 | stdlib                   | Unexpected work during chain building in crypto/x509                                                                    |
| GO-2026-4971 | CVE-2026-39836                 | stdlib                   | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                                   |
| GO-2026-4976 | CVE-2026-39825                 | stdlib                   | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil                          |
| GO-2026-4977 | CVE-2026-42499                 | stdlib                   | Quadratic string concatenation in consumePhrase in net/mail                                                             |
| GO-2026-4980 | CVE-2026-39826                 | stdlib                   | Escaper bypass leads to XSS in html/template                                                                            |
| GO-2026-4981 | CVE-2026-33811                 | stdlib                   | Crash when handling long CNAME response in net                                                                          |
| GO-2026-4982 | CVE-2026-39823                 | stdlib                   | Bypass of meta content URL escaping causes XSS in html/template                                                         |
| GO-2026-4986 | CVE-2026-39820                 | stdlib                   | Quadratic string concatentation in consumeComment in net/mail                                                           |
| GO-2026-5024 | CVE-2026-39824                 | golang.org/x/sys         | Invoking integer overflow in NewNTUnicodeString in golang.org/x/sys/windows                                             |
| GO-2026-5025 | CVE-2026-42506                 | golang.org/x/net         | Invoking incorrect handling of namespaced elements in foreign content in golang.org/x/net/html                          |
| GO-2026-5026 | CVE-2026-39821                 | golang.org/x/net         | Invoking failure to reject ASCII-only Punycode-encoded labels in golang.org/x/net/idna                                  |
| GO-2026-5027 | CVE-2026-42502                 | golang.org/x/net         | Invoking incorrect handling of HTML elements in foreign content in golang.org/x/net/html                                |
| GO-2026-5028 | CVE-2026-25680                 | golang.org/x/net         | Invoking denial of service when parsing arbitrary HTML in golang.org/x/net/html                                         |
| GO-2026-5029 | CVE-2026-25681                 | golang.org/x/net         | Invoking incorrect handling of character references in DOCTYPE nodes in golang.org/x/net/html                           |
| GO-2026-5030 | CVE-2026-27136                 | golang.org/x/net         | Invoking duplicate attributes can cause XSS in golang.org/x/net/html                                                    |

## troubleshooting-panel-plugin (release-coo-ocp-4.19)

### Commits analyzed

- `4a04933 Merge pull request #254 from openshift/update-vulnerable-dependencies-4.19-26-06-2026`
- `44bd915 fix: update vulnerable dependencies`
- `83f2504 Merge pull request #233 from openshift-cherrypick-robot/cherry-pick-232-to-release-coo-ocp-4.19`
- `b357576 COO-1819:fix: Add TLS min version and cipher suite configuration support`

### NPM Vulnerabilities Fixed

| Advisory                                          | CVE            | Package                                  | Severity | Title                                                                                                                                                |
| ------------------------------------------------- | -------------- | ---------------------------------------- | -------- | ---------------------------------------------------------------------------------------------------------------------------------------------------- |
| https://github.com/advisories/GHSA-25h7-pfq9-p65f | CVE-2026-32141 | flatted                                  | high     | flatted vulnerable to unbounded recursion DoS in parse() revive phase                                                                                |
| https://github.com/advisories/GHSA-2j2x-hqr9-3h42 | CVE-2026-40181 | react-router                             | moderate | React Router's same-origin redirect with path starting // causes open redirect via protocol-relative URL reinterpretation                            |
| https://github.com/advisories/GHSA-2mjp-6q6p-2qxm | CVE-2026-1525  | undici                                   | moderate | Undici has an HTTP Request/Response Smuggling issue                                                                                                  |
| https://github.com/advisories/GHSA-2qvq-rjwj-gvw9 | CVE-2026-33916 | handlebars                               | moderate | Handlebars.js has Prototype Pollution Leading to XSS through Partial Template Injection                                                              |
| https://github.com/advisories/GHSA-2w6w-674q-4c4q | CVE-2026-33937 | handlebars                               | critical | Handlebars.js has JavaScript Injection via AST Type Confusion                                                                                        |
| https://github.com/advisories/GHSA-35p6-xmwp-9g52 | CVE-2026-6733  | undici                                   | low      | undici vulnerable to HTTP response queue poisoning via keep-alive socket reuse                                                                       |
| https://github.com/advisories/GHSA-37ch-88jc-xwx2 | CVE-2026-4867  | path-to-regexp                           | high     | path-to-regexp vulnerable to Regular Expression Denial of Service via multiple route parameters                                                      |
| https://github.com/advisories/GHSA-3mfm-83xf-c92r | CVE-2026-33938 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion by tampering @partial-block                                                            |
| https://github.com/advisories/GHSA-3v7f-55p6-f55p | CVE-2026-33672 | picomatch                                | moderate | Picomatch: Method Injection in POSIX Character Classes causes incorrect Glob Matching                                                                |
| https://github.com/advisories/GHSA-442j-39wm-28r2 | N/A            | handlebars                               | low      | Handlebars.js has a Property Access Validation Bypass in container.lookup                                                                            |
| https://github.com/advisories/GHSA-4992-7rv2-5pvq | CVE-2026-1527  | undici                                   | moderate | Undici has CRLF Injection in undici via `upgrade` option                                                                                             |
| https://github.com/advisories/GHSA-4x5r-pxfx-6jf8 | CVE-2026-49356 | @babel/core                              | low      | @babel/core: Arbitrary File Read via sourceMappingURL Comment                                                                                        |
| https://github.com/advisories/GHSA-58qx-3vcg-4xpx | CVE-2026-45736 | ws                                       | moderate | ws: Uninitialized memory disclosure                                                                                                                  |
| https://github.com/advisories/GHSA-64mm-vxmg-q3vj | CVE-2026-55602 | http-proxy-middleware                    | moderate | http-proxy-middleware `router` host+path substring matching allows Host-header-driven backend routing bypass                                         |
| https://github.com/advisories/GHSA-79cf-xcqc-c78w | CVE-2026-6402  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to cross-origin source code exposure on non-HTTPS origins                                                              |
| https://github.com/advisories/GHSA-7rx3-28cr-v5wh | N/A            | handlebars                               | moderate | Handlebars.js has a Prototype Method Access Control Gap via Missing **lookupSetter** Blocklist Entry                                                 |
| https://github.com/advisories/GHSA-96hv-2xvq-fx4p | CVE-2026-48779 | ws                                       | high     | ws: Memory exhaustion DoS from tiny fragments and data chunks                                                                                        |
| https://github.com/advisories/GHSA-9cx6-37pm-9jff | CVE-2026-33939 | handlebars                               | high     | Handlebars.js has Denial of Service via Malformed Decorator Syntax in Template Compilation                                                           |
| https://github.com/advisories/GHSA-c2c7-rcm5-vvqj | CVE-2026-33671 | picomatch                                | high     | Picomatch has a ReDoS vulnerability via extglob quantifiers                                                                                          |
| https://github.com/advisories/GHSA-f23m-r3pf-42rh | CVE-2026-2950  | lodash                                   | moderate | lodash vulnerable to Prototype Pollution via array path bypass in `_.unset` and `_.omit`                                                             |
| https://github.com/advisories/GHSA-f269-vfmq-vjvj | CVE-2026-1528  | undici                                   | high     | Undici: Malicious WebSocket 64-bit length overflows parser and crashes the client                                                                    |
| https://github.com/advisories/GHSA-f886-m6hf-6m8v | CVE-2026-33750 | brace-expansion                          | moderate | brace-expansion: Zero-step sequence causes process hang and memory exhaustion                                                                        |
| https://github.com/advisories/GHSA-fv7c-fp4j-7gwp | CVE-2026-44728 | @babel/plugin-transform-modules-systemjs | high     | @babel/plugin-transform-modules-systemjs generates arbitrary code when compiling malicious input                                                     |
| https://github.com/advisories/GHSA-g8m3-5g58-fq7m | CVE-2026-11525 | undici                                   | low      | undici vulnerable to Set-Cookie SameSite attribute downgrade via permissive substring matching                                                       |
| https://github.com/advisories/GHSA-h67p-54hq-rp68 | CVE-2026-53550 | js-yaml                                  | moderate | JS-YAML: Quadratic-complexity DoS in merge key handling via repeated aliases                                                                         |
| https://github.com/advisories/GHSA-hmw2-7cc7-3qxx | CVE-2026-12143 | form-data                                | high     | form-data: CRLF injection in form-data via unescaped multipart field names and filenames                                                             |
| https://github.com/advisories/GHSA-mh29-5h37-fv8m | CVE-2025-64718 | js-yaml                                  | moderate | js-yaml has prototype pollution in merge (<<)                                                                                                        |
| https://github.com/advisories/GHSA-mx8g-39q3-5c79 | CVE-2026-9595  | webpack-dev-server                       | moderate | webpack-dev-server vulnerable to HMR WebSocket interception via permissive user proxies                                                              |
| https://github.com/advisories/GHSA-p88m-4jfj-68fv | CVE-2026-9679  | undici                                   | moderate | undici vulnerable to HTTP header injection via Set-Cookie percent-decoding                                                                           |
| https://github.com/advisories/GHSA-ph9p-34f9-6g65 | CVE-2026-44705 | tmp                                      | high     | tmp has Path Traversal via unsanitized prefix/postfix that enables directory escape                                                                  |
| https://github.com/advisories/GHSA-phc3-fgpg-7m6h | CVE-2026-2581  | undici                                   | moderate | Undici has Unbounded Memory Consumption in its DeduplicationHandler via Response Buffering that leads to DoS                                         |
| https://github.com/advisories/GHSA-pr7r-676h-xcf6 | CVE-2026-9678  | undici                                   | moderate | undici vulnerable to cross-user information disclosure via shared cache whitespace bypass                                                            |
| https://github.com/advisories/GHSA-q3j6-qgpj-74h6 | CVE-2026-6321  | fast-uri                                 | high     | fast-uri vulnerable to path traversal via percent-encoded dot segments                                                                               |
| https://github.com/advisories/GHSA-q8mj-m7cp-5q26 | CVE-2026-8723  | qs                                       | moderate | qs has a remotely triggerable DoS: qs.stringify crashes with TypeError on null/undefined entries in comma-format arrays when encodeValuesOnly is set |
| https://github.com/advisories/GHSA-qx2v-qp2m-jg93 | CVE-2026-41305 | postcss                                  | moderate | PostCSS has XSS via Unescaped </style> in its CSS Stringify Output                                                                                   |
| https://github.com/advisories/GHSA-r4q5-vmmm-2653 | N/A            | follow-redirects                         | moderate | follow-redirects leaks Custom Authentication Headers to Cross-Domain Redirect Targets                                                                |
| https://github.com/advisories/GHSA-r5fr-rjxr-66jc | CVE-2026-4800  | lodash                                   | high     | lodash vulnerable to Code Injection via `_.template` imports key names                                                                               |
| https://github.com/advisories/GHSA-rf6f-7fwh-wjgh | CVE-2026-33228 | flatted                                  | high     | Prototype Pollution via parse() in NodeJS flatted                                                                                                    |
| https://github.com/advisories/GHSA-v39h-62p7-jpjc | CVE-2026-6322  | fast-uri                                 | high     | fast-uri vulnerable to host confusion via percent-encoded authority delimiters                                                                       |
| https://github.com/advisories/GHSA-v6wh-96g9-6wx3 | CVE-2026-53632 | launch-editor                            | moderate | launch-editor: NTLMv2 hash disclosure via UNC path handling on Windows                                                                               |
| https://github.com/advisories/GHSA-v9p9-hfj2-hcw8 | CVE-2026-2229  | undici                                   | high     | Undici has Unhandled Exception in WebSocket Client Due to Invalid server_max_window_bits Validation                                                  |
| https://github.com/advisories/GHSA-vpq2-c234-7xj6 | CVE-2026-3449  | @tootallnate/once                        | low      | @tootallnate/once vulnerable to Incorrect Control Flow Scoping                                                                                       |
| https://github.com/advisories/GHSA-vrm6-8vpv-qv8q | CVE-2026-1526  | undici                                   | high     | Undici has Unbounded Memory Consumption in WebSocket permessage-deflate Decompression                                                                |
| https://github.com/advisories/GHSA-vxpw-j846-p89q | CVE-2026-12151 | undici                                   | high     | undici WebSocket client vulnerable to denial of service via fragment count bypass                                                                    |
| https://github.com/advisories/GHSA-w7jw-789q-3m8p | CVE-2026-9277  | shell-quote                              | critical | shell-quote quote() does not escape newlines in object .op values                                                                                    |
| https://github.com/advisories/GHSA-wf6x-7x77-mvgw | CVE-2026-29063 | immutable                                | high     | Immutable is vulnerable to Prototype Pollution                                                                                                       |
| https://github.com/advisories/GHSA-xhpv-hc6g-r9c6 | CVE-2026-33940 | handlebars                               | high     | Handlebars.js has JavaScript Injection via AST Type Confusion when passing an object as dynamic partial                                              |
| https://github.com/advisories/GHSA-xjpj-3mr7-gcpf | CVE-2026-33941 | handlebars                               | high     | Handlebars.js has JavaScript Injection in CLI Precompiler via Unescaped Names and Options                                                            |

### Go Vulnerabilities Fixed

| ID           | CVE/Aliases                    | Module | Summary                                                                                                  |
| ------------ | ------------------------------ | ------ | -------------------------------------------------------------------------------------------------------- |
| GO-2026-4337 | CVE-2025-68121                 | stdlib | Unexpected session resumption in crypto/tls                                                              |
| GO-2026-4340 | CVE-2025-61730, CVE-2025-61730 | stdlib | Handshake messages may be processed at the incorrect encryption level in crypto/tls                      |
| GO-2026-4341 | CVE-2025-61726, CVE-2025-61726 | stdlib | Memory exhaustion in query parameter parsing in net/url                                                  |
| GO-2026-4342 | CVE-2025-61728, CVE-2025-61728 | stdlib | Excessive CPU consumption when building archive index in archive/zip                                     |
| GO-2026-4601 | CVE-2026-25679                 | stdlib | Incorrect parsing of IPv6 host literals in net/url                                                       |
| GO-2026-4602 | CVE-2026-27139                 | stdlib | FileInfo can escape from a Root in os                                                                    |
| GO-2026-4603 | CVE-2026-27142                 | stdlib | URLs in meta content attribute actions are not escaped in html/template                                  |
| GO-2026-4864 | CVE-2026-32282                 | stdlib | TOCTOU permits root escape on Linux via Root.Chmod in os in internal/syscall/unix                        |
| GO-2026-4865 | CVE-2026-32289                 | stdlib | JsBraceDepth Context Tracking Bugs (XSS) in html/template                                                |
| GO-2026-4869 | CVE-2026-32288                 | stdlib | Unbounded allocation for old GNU sparse in archive/tar                                                   |
| GO-2026-4870 | CVE-2026-32283                 | stdlib | Unauthenticated TLS 1.3 KeyUpdate record can cause persistent connection retention and DoS in crypto/tls |
| GO-2026-4946 | CVE-2026-32281                 | stdlib | Inefficient policy validation in crypto/x509                                                             |
| GO-2026-4947 | CVE-2026-32280                 | stdlib | Unexpected work during chain building in crypto/x509                                                     |
| GO-2026-4971 | CVE-2026-39836                 | stdlib | Panic in Dial and LookupPort when handling NUL byte on Windows in net                                    |
| GO-2026-4976 | CVE-2026-39825                 | stdlib | ReverseProxy forwards queries with more than urlmaxqueryparams parameters in net/http/httputil           |
| GO-2026-4977 | CVE-2026-42499                 | stdlib | Quadratic string concatenation in consumePhrase in net/mail                                              |
| GO-2026-4980 | CVE-2026-39826                 | stdlib | Escaper bypass leads to XSS in html/template                                                             |
| GO-2026-4981 | CVE-2026-33811                 | stdlib | Crash when handling long CNAME response in net                                                           |
| GO-2026-4982 | CVE-2026-39823                 | stdlib | Bypass of meta content URL escaping causes XSS in html/template                                          |
| GO-2026-4986 | CVE-2026-39820                 | stdlib | Quadratic string concatentation in consumeComment in net/mail                                            |

## troubleshooting-panel-plugin (release-coo-ocp-4.22)

### Commits analyzed

- `94d11ba Merge pull request #247 from openshift-cherrypick-robot/cherry-pick-246-to-release-coo-ocp-4.22`
- `f48f34a fix: COO-1850: Minor UI fixes`
- `ae1a576 Merge pull request #244 from openshift-cherrypick-robot/cherry-pick-243-to-release-coo-ocp-4.22`
- `5365208 fix: COO-1841: prevent Error object from rendering as React child in AgentMenu`

### NPM Vulnerabilities Fixed

No NPM vulnerabilities were fixed.

### Go Vulnerabilities Fixed

No Go vulnerabilities were fixed.

## Summary

| Project                      | Branch               | NPM Fixed | Go Fixed | Total | Status |
| ---------------------------- | -------------------- | --------- | -------- | ----- | ------ |
| perses-operator              | release-coo-1.5      | 0         | 20       | 20    | ok     |
| monitoring-plugin            | release-coo-ocp-4.15 | 20        | 33       | 53    | ok     |
| monitoring-plugin            | release-coo-ocp-4.19 | 36        | 33       | 69    | ok     |
| monitoring-plugin            | release-coo-ocp-4.22 | 40        | 33       | 73    | ok     |
| logging-view-plugin          | release-coo-ocp-4.12 | 53        | 34       | 87    | ok     |
| logging-view-plugin          | release-coo-ocp-4.15 | 47        | 33       | 80    | ok     |
| logging-view-plugin          | release-coo-ocp-4.22 | 23        | 20       | 43    | ok     |
| distributed-tracing-plugin   | release-coo-ocp-4.12 | 31        | 32       | 63    | ok     |
| distributed-tracing-plugin   | release-coo-ocp-4.15 | 23        | 28       | 51    | ok     |
| distributed-tracing-plugin   | release-coo-ocp-4.19 | 23        | 28       | 51    | ok     |
| distributed-tracing-plugin   | release-coo-ocp-4.22 | 23        | 28       | 51    | ok     |
| troubleshooting-panel-plugin | release-coo-ocp-4.19 | 48        | 20       | 68    | ok     |
| troubleshooting-panel-plugin | release-coo-ocp-4.22 | 0         | 0        | 0     | ok     |
