# Test Case: Sanitize runbook_url in Alert Detail Page

## Background

**Jira**: [OU-1488](https://redhat.atlassian.net/browse/OU-1488)  
**PR**: [monitoring-plugin#1255](https://github.com/openshift/monitoring-plugin/pull/1255)

### Issue Summary
The alert and rule detail pages previously accepted `annotations.runbook_url` from Prometheus/Thanos API responses and passed them verbatim to the `ExternalLink` component without validation. This created a potential XSS vulnerability where malicious scripts could be injected via unsafe URL protocols.

### Fix Summary
Added URL validation to the `ExternalLink` component that:
- Only accepts `http:` and `https:` protocols
- Rejects unsafe protocols (`javascript:`, `data:`, `vbscript:`, `mailto:`, etc.)
- Rejects relative URLs and protocol-relative URLs
- Renders invalid/unsafe URLs as plain text instead of clickable links

---

## Quick Reference: What to Look For

When testing, here's what you should see in the Alert Details page:

**⚠️ IMPORTANT**: Test alerts include URLs in TWO places:
1. **Runbook URL field** - Tests the `ExternalLink` component directly
2. **Description field** - Tests the `LinkifyExternal` component (auto-linkifies URLs in text)

Both fields must be tested for each alert!

### ✅ VALID URL (Should Render as Link)
**Input**: `https://runbooks.example.com/alert`

**Expected Appearance**:
```
Runbook URL: [https://runbooks.example.com/alert 🔗]
              ↑ Blue, underlined, clickable with external link icon
```

**DOM Structure**:
```html
<a href="https://runbooks.example.com/alert" 
   target="_blank" 
   rel="noopener noreferrer">
  https://runbooks.example.com/alert
</a>
```

---

### ❌ MALICIOUS URL (Should Render as Plain Text)
**Input**: `javascript:alert('XSS')`

**Expected Appearance**:
```
Runbook URL: javascript:alert('XSS')
             ↑ Black, plain text, NOT clickable, NO icon
```

**DOM Structure**:
```html
<span>javascript:alert('XSS')</span>
<!-- or similar non-link element -->
<!-- NO <a> tag, NO href attribute -->
```

---

## Testing Workflow Overview

```
┌─────────────────────────────────────────────────────────────────┐
│                        PREREQUISITES                             │
│  • OpenShift cluster access                                      │
│  • oc CLI installed and logged in                                │
│  • Permissions to create PrometheusRules                         │
│  • Browser with developer tools                                  │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                   CREATE TEST PROMETHEUSRULES                    │
│  1. Valid HTTPS URL → test-runbook-valid-https.yaml              │
│  2. JavaScript protocol → test-runbook-javascript.yaml           │
│  3. Data URI → test-runbook-data-uri.yaml                        │
│                                                                   │
│  Apply with: oc apply -f <filename>.yaml                         │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    WAIT FOR ALERTS TO FIRE                       │
│  • Takes 30-60 seconds                                           │
│  • Verify: oc exec prometheus-k8s-0 -- check alerts              │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                 OPEN OPENSHIFT CONSOLE UI                        │
│  Navigate to: Observe > Alerting                                 │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                    TEST EACH ALERT                               │
│                                                                   │
│  Test 1: TestAlertValidHTTPS                                     │
│  ✅ Verify: URL is clickable link with icon                      │
│  ✅ Verify: Opens in new tab                                     │
│                                                                   │
│  Test 2: TestAlertJavaScript                                     │
│  ✅ Verify: URL is plain text (not clickable)                    │
│  ✅ Verify: No JavaScript execution                              │
│  ✅ Verify: No alert dialog appears                              │
│                                                                   │
│  Test 3: TestAlertDataURI                                        │
│  ✅ Verify: URL is plain text (not clickable)                    │
│  ✅ Verify: No script execution                                  │
│                                                                   │
└────────────────────────┬────────────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────────────┐
│                       CLEANUP                                    │
│  Delete test PrometheusRules:                                    │
│  oc delete prometheusrule -l test=runbook-sanitization           │
└─────────────────────────────────────────────────────────────────┘
```

---

## Testing Checklist

Use this checklist during manual testing:

- [ ] Prerequisites completed (cluster access, permissions, tools)
- [ ] Created PrometheusRule with valid HTTPS URL
- [ ] Created PrometheusRule with `javascript:` protocol
- [ ] Created PrometheusRule with `data:` URI
- [ ] All alerts are firing (visible in Alerting page)
- [ ] Valid HTTPS URL renders as clickable link with icon in **runbook_url** field
- [ ] Valid HTTPS URL renders as clickable link in **description** field
- [ ] Valid HTTPS link opens in new tab
- [ ] JavaScript protocol renders as plain text (not clickable) in **runbook_url** field
- [ ] JavaScript protocol renders as plain text (not clickable) in **description** field
- [ ] No JavaScript alert dialogs appear
- [ ] Data URI renders as plain text (not clickable) in **runbook_url** field
- [ ] Data URI renders as plain text (not clickable) in **description** field
- [ ] No script execution occurs
- [ ] Browser console shows no errors
- [ ] Verified URLs in notification bell icon (🔔) are properly sanitized
- [ ] Verified URLs in alerts list/tooltips are properly sanitized
- [ ] Cleanup completed (test resources deleted)

---

## Test Scenarios

### Scenario 1: Valid HTTPS URL
**Objective**: Verify that valid HTTPS URLs render as clickable external links

**Input**: `runbook_url: "https://runbooks.example.com/alert"`

**Expected Result**:
- URL renders as a clickable link with external link icon
- Link has `target="_blank"` attribute
- Link has `rel="noopener noreferrer"` attribute
- Clicking the link opens in a new tab

---

### Scenario 2: Valid HTTP URL
**Objective**: Verify that valid HTTP URLs render as clickable external links

**Input**: `runbook_url: "http://runbooks.example.com/alert?severity=high"`

**Expected Result**:
- URL renders as a clickable link with external link icon
- Link has `target="_blank"` attribute
- Link has `rel="noopener noreferrer"` attribute
- Clicking the link opens in a new tab

---

### Scenario 3: JavaScript Protocol Attack
**Objective**: Verify that `javascript:` protocol URLs are blocked

**Input**: `runbook_url: "javascript:alert(1)"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created
- Text displays the literal string "javascript:alert(1)"
- No JavaScript execution occurs

---

### Scenario 4: Case-Insensitive JavaScript Protocol
**Objective**: Verify that protocol validation is case-insensitive

**Input**: `runbook_url: "JaVaScRiPt:alert(1)"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created
- No JavaScript execution occurs

---

### Scenario 5: Data URI Attack
**Objective**: Verify that `data:` protocol URLs are blocked

**Input**: `runbook_url: "data:text/html,<script>alert(1)</script>"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created
- No script execution occurs

---

### Scenario 6: VBScript Protocol Attack
**Objective**: Verify that `vbscript:` protocol URLs are blocked

**Input**: `runbook_url: "vbscript:msgbox(1)"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created

---

### Scenario 7: Mailto Protocol
**Objective**: Verify that `mailto:` URLs are blocked

**Input**: `runbook_url: "mailto:security@example.com"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created

---

### Scenario 8: Relative URL
**Objective**: Verify that relative URLs are blocked

**Input**: `runbook_url: "/runbooks/alert"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created

---

### Scenario 9: Protocol-Relative URL
**Objective**: Verify that protocol-relative URLs are blocked

**Input**: `runbook_url: "//runbooks.example.com/alert"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created

---

### Scenario 10: Whitespace Prefix Attack
**Objective**: Verify that URLs with leading whitespace are blocked

**Input**: `runbook_url: "\tjavascript:alert(1)"`

**Expected Result**:
- URL renders as plain text (not clickable)
- No `<a>` element is created
- No JavaScript execution occurs

---

### Scenario 11: Invalid URL String
**Objective**: Verify that non-URL strings are handled gracefully

**Input**: `runbook_url: "not a URL"`

**Expected Result**:
- Text renders as plain text (not clickable)
- No `<a>` element is created

---

### Scenario 12: Missing runbook_url
**Objective**: Verify behavior when runbook_url is not present

**Input**: Alert with no `annotations.runbook_url`

**Expected Result**:
- No runbook link is displayed
- No errors are thrown

---

## Manual Testing Steps

### Prerequisites

#### 1. OpenShift Cluster Access
- Access to an OpenShift 4.x cluster with the monitoring stack deployed
- Cluster should have the updated monitoring-plugin version with the fix deployed
- Verify monitoring is enabled:
  ```bash
  oc get pod -n openshift-monitoring
  ```
  Expected: You should see prometheus, alertmanager, and monitoring-plugin pods running

#### 2. Required Permissions
You need permissions to create PrometheusRule resources. Verify with:
```bash
oc auth can-i create prometheusrules -n openshift-monitoring
```
Expected output: `yes`

If you don't have permissions, request one of the following roles:
- `cluster-admin` (full access)
- `monitoring-rules-edit` (monitoring-specific)
- Custom role with `prometheusrules` resource access

#### 3. Command Line Tools
- `oc` CLI installed and configured
- Login to your cluster:
  ```bash
  oc login <cluster-url> -u <username> -p <password>
  # or using token
  oc login --token=<token> --server=<cluster-url>
  ```

#### 4. Browser Setup
- Modern web browser (Chrome, Firefox, Edge)
- Browser developer tools enabled (F12 or right-click > Inspect)
- Access to OpenShift Console (usually at `https://console-openshift-console.apps.<cluster-domain>`)

---

### Manual Test Execution

#### Part A: Setup and Create Test PrometheusRules

##### Step 1: Create a Test Namespace (Optional)
You can use `openshift-monitoring` namespace or create your own:

```bash
# Option 1: Use existing namespace
NAMESPACE="openshift-monitoring"

# Option 2: Create your own namespace (requires user workload monitoring enabled)
NAMESPACE="test-monitoring"
oc create namespace ${NAMESPACE}
oc label namespace ${NAMESPACE} openshift.io/cluster-monitoring="true"
```

##### Step 2: Create PrometheusRule with Valid HTTPS URL

Create a file named `test-rule-valid-https.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-valid-https
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertValidHTTPS
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "https://runbooks.example.com/test-alert"
        summary: "Test alert with valid HTTPS runbook URL"
        description: "This alert tests that valid HTTPS URLs render as clickable links. The URL is: https://runbooks.example.com/test-alert"
```

Apply the rule:
```bash
oc apply -f test-rule-valid-https.yaml
```

Verify the rule was created:
```bash
oc get prometheusrule test-runbook-valid-https -n openshift-monitoring
```

##### Step 3: Create PrometheusRule with JavaScript Protocol (Attack Scenario)

Create a file named `test-rule-javascript.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-javascript
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertJavaScript
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "javascript:alert('XSS Attack!')"
        summary: "Test alert with javascript protocol (should be blocked)"
        description: "This alert tests that javascript: URLs are sanitized and render as plain text. Malicious URL: javascript:alert('XSS Attack!')"
```

Apply the rule:
```bash
oc apply -f test-rule-javascript.yaml
```

##### Step 4: Create PrometheusRule with Data URI (Attack Scenario)

Create a file named `test-rule-data-uri.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-data-uri
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertDataURI
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "data:text/html,<script>alert('XSS')</script>"
        summary: "Test alert with data URI (should be blocked)"
        description: "This alert tests that data: URLs are sanitized and render as plain text. Malicious URL: data:text/html,<script>alert('XSS')</script>"
```

Apply the rule:
```bash
oc apply -f test-rule-data-uri.yaml
```

##### Step 5: Wait for Alerts to Fire

The alerts use `expr: vector(1)` which always evaluates to true, so they should fire immediately.

Wait 30-60 seconds, then verify alerts are firing:
```bash
oc exec -n openshift-monitoring prometheus-k8s-0 -- promtool query instant http://localhost:9090 'ALERTS{test="runbook-sanitization"}'
```

Or check in the Prometheus UI:
```bash
# Port-forward to Prometheus
oc port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090
# Open browser to http://localhost:9090/alerts
```

---

#### Part B: Test in OpenShift Console UI

##### Step 6: Access the OpenShift Console

1. Open your browser and navigate to the OpenShift Console
2. Login with your credentials
3. You should see the main dashboard

##### Step 7: Navigate to Alerts Page

1. In the left navigation menu (vertical sidebar on the left), click **Observe**
   - Look for the bar chart icon 📊
   - The menu will expand showing sub-options
2. Click **Alerting**
   - Look for the bell icon 🔔
3. You should see a page titled "Alerting" with:
   - A filter/search box at the top
   - Tabs: "Alerts" (default), "Silences", "Alerting Rules"
   - A table showing firing alerts with columns: Name, Labels, Severity, State, Alerting Rule

**Note**: The page URL should be something like:
```
https://console-openshift-console.apps.<cluster>/monitoring/alerts
```

##### Step 8: Test Valid HTTPS URL

1. In the alerts list, find **TestAlertValidHTTPS** (use the filter box if needed)
   - Type "TestAlertValidHTTPS" in the filter/search box at the top of the table
   - The table should filter to show only matching alerts
2. Click on the alert name **TestAlertValidHTTPS** (it's a blue clickable link in the Name column)
3. You'll be taken to the Alert Details page
   - Page title: "TestAlertValidHTTPS"
   - URL format: `https://console.../monitoring/alerts/TestAlertValidHTTPS?...`
4. In the alert details page, you'll see several sections:
   - **Alert Details** (top section with severity, state, etc.)
   - **Labels** (key-value pairs)
   - **Annotations** section
5. Locate the **Runbook URL** field in the Annotations section
   - It's typically displayed as: `Runbook URL: <value>`
   - The value will be either a link or plain text depending on validation

6. **IMPORTANT**: Also check the **Description** field
   - The test alerts include the URL in the description text as well
   - This tests the `LinkifyExternal` component which automatically converts URLs in text to links
   - The description appears in multiple places: alert details page, alert list, notification bell icon area
   - Valid URLs in the description should be converted to clickable links
   - Malicious URLs in the description should remain as plain text
4. **Verify**:
   - ✅ The URL `https://runbooks.example.com/test-alert` appears as a **clickable link**
   - ✅ The link has an external link icon (🔗) next to it
   - ✅ The link text is blue/underlined (standard link styling)
   - ✅ When you hover over it, your mouse cursor changes to a pointer (hand icon)

5. **Inspect the DOM** (this verifies the security attributes):
   - Right-click directly on the link text
   - Select **Inspect** or **Inspect Element** from the context menu
   - The browser Developer Tools will open (usually at the bottom or right side)
   - The **Elements** or **Inspector** tab will highlight the HTML element
   
6. In the developer tools, look at the highlighted HTML code and verify it shows:
   ```html
   <a href="https://runbooks.example.com/test-alert" 
      target="_blank" 
      rel="noopener noreferrer"
      class="pf-c-button pf-m-link pf-m-inline">
     https://runbooks.example.com/test-alert
   </a>
   ```
   
   **Key things to check**:
   - ✅ The element is an `<a>` tag (anchor/link)
   - ✅ `href` attribute contains the HTTPS URL
   - ✅ `target="_blank"` is present (opens in new tab)
   - ✅ `rel="noopener noreferrer"` is present (security attributes)
7. Click the link and verify:
   - ✅ It opens in a **new tab**
   - ✅ No security warnings appear
   - ✅ The browser attempts to navigate to the URL (will fail since it's a fake domain, but that's expected)

8. **Test the Description field**:
   - Scroll to the **Description** field in the Annotations section
   - You should see: "This alert tests that valid HTTPS URLs render as clickable links. The URL is: https://runbooks.example.com/test-alert"
   - The URL at the end should ALSO be a clickable link (converted by LinkifyExternal)
   - Click on it to verify it opens in a new tab
   - This confirms URL linkification works in text fields too

##### Step 9: Test JavaScript Protocol (Should Be Blocked)

1. Navigate back to **Observe > Alerting**
   - Click the back button or click **Alerting** in the left menu
2. Find **TestAlertJavaScript** in the alerts list
   - Type "TestAlertJavaScript" in the filter box
3. Click on the alert name to open the alert details page
4. In the alert details, locate the **Runbook URL** section

5. **Visual Verification**:
   - ✅ The text `javascript:alert('XSS Attack!')` appears as **plain text** (not a link)
   - ✅ The text is **NOT blue or underlined**
   - ✅ There is **NO external link icon** (🔗) next to it
   - ✅ The text is **NOT clickable** (cursor remains as arrow, not pointer)
   - ✅ The text appears in the same style as the "Summary" or "Description" fields (regular text)

6. **Inspect the DOM**:
   - Right-click directly on the text `javascript:alert('XSS Attack!')`
   - Select **Inspect** or **Inspect Element**
   - The developer tools will open and highlight the element

7. In the developer tools, verify the HTML structure:
   ```html
   <!-- Should see something like this: -->
   <span>javascript:alert('XSS Attack!')</span>
   
   <!-- Or possibly wrapped in a different non-link element -->
   <div>javascript:alert('XSS Attack!')</div>
   ```
   
   **Key things to verify**:
   - ✅ There is **NO `<a>` tag** wrapping the text
   - ✅ There is **NO `href` attribute** anywhere
   - ✅ The element is a plain text container like `<span>`, `<div>`, or `<p>`
   - ✅ The text content exactly matches: `javascript:alert('XSS Attack!')`

8. **Test the Description field**:
   - Find the **Description** field in the Annotations section
   - You should see: "This alert tests that javascript: URLs are sanitized and render as plain text. Malicious URL: javascript:alert('XSS Attack!')"
   - **Verify the URL in the description is ALSO plain text**:
     - ✅ The text `javascript:alert('XSS Attack!')` is **NOT clickable**
     - ✅ It appears as regular text (not blue/underlined)
     - ✅ No link is created
   - This confirms LinkifyExternal component also blocks malicious URLs

9. **Security Check** (CRITICAL):
   - ✅ **No JavaScript alert dialog appears** when the page loads
   - ✅ **No alert appears when you click** on the text (in either runbook_url or description)
   - ✅ Open the browser console (F12 > Console tab) and verify:
     - ✅ **No JavaScript errors** appear
     - ✅ **No warning messages** about blocked scripts
     - ✅ Console is clean
   - ✅ Try clicking on the text multiple times in both fields - nothing should happen

##### Step 10: Test Data URI (Should Be Blocked)

1. Navigate back to **Observe > Alerting**
2. Find **TestAlertDataURI** in the alerts list
3. Click on the alert name to open the alert details page
4. In the alert details, locate the **Runbook URL** section

5. **Visual Verification**:
   - ✅ The text `data:text/html,<script>alert('XSS')</script>` appears as **plain text**
   - ✅ The text is **NOT clickable**
   - ✅ **No external link icon** is present
   - ✅ Text appears in regular font (not blue/underlined)
   - ✅ Cursor remains as arrow when hovering over it

6. **Inspect the DOM**:
   - Right-click on the data URI text
   - Select **Inspect Element**
   - Verify in the HTML:
     ```html
     <span>data:text/html,&lt;script&gt;alert('XSS')&lt;/script&gt;</span>
     ```
   - **Key checks**:
     - ✅ **NO `<a>` element**
     - ✅ **NO `href` attribute**
     - ✅ Script tags should be HTML-encoded as `&lt;` and `&gt;`

7. **Test the Description field**:
   - Find the **Description** field in the Annotations section
   - You should see: "This alert tests that data: URLs are sanitized and render as plain text. Malicious URL: data:text/html,&lt;script&gt;alert('XSS')&lt;/script&gt;"
   - **Verify the URL in the description is ALSO plain text**:
     - ✅ The data URI is **NOT clickable**
     - ✅ It appears as regular text
     - ✅ No link is created
     - ✅ Script tags are HTML-encoded as `&lt;script&gt;`
   - This confirms both ExternalLink and LinkifyExternal components sanitize data URIs

8. **Security Check** (CRITICAL):
   - ✅ **No script execution occurs** when the page loads
   - ✅ **No alert popup appears**
   - ✅ **No JavaScript runs** when clicking the text (in either runbook_url or description)
   - ✅ Check browser console (F12 > Console) - should be clean
   - ✅ The actual script code `<script>alert('XSS')</script>` is never executed in any field

---

##### Step 11a: Test Description Field in Other UI Locations

The description field (which now contains the URLs) appears in several locations in the OpenShift Console. Verify sanitization works everywhere:

**Location 1: Notification Bell Icon**
1. Click the **Bell icon** (🔔) in the top-right corner of the OpenShift Console
2. This shows a dropdown with recent alerts
3. Each alert shows its description text
4. **Verify**:
   - ✅ For TestAlertValidHTTPS: The URL in the description should be a **clickable link**
   - ✅ For TestAlertJavaScript: The URL in the description should be **plain text** (not clickable)
   - ✅ For TestAlertDataURI: The URL in the description should be **plain text** (not clickable)

**Location 2: Alerts List Page**
1. Go to **Observe > Alerting**
2. Find your test alerts in the list
3. Hover over an alert to see a tooltip/preview that may include the description
4. **Verify**: URLs are properly sanitized in the tooltip/preview

**Location 3: Alert Details Summary**
1. On the alert details page (where you've been testing)
2. The description appears in the main content area
3. Already verified in previous steps

**Location 4: Prometheus UI** (Optional, advanced)
1. Port-forward to Prometheus: `oc port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090`
2. Open browser to `http://localhost:9090/alerts`
3. Find your test alerts
4. Click to expand alert details
5. Verify URLs are sanitized in the Prometheus UI as well (though this may not use the same React components)

**Why This Matters**:
- The `LinkifyExternal` component automatically converts URLs in text fields to links
- It's used throughout the console wherever text might contain URLs
- All these locations must sanitize malicious URLs consistently
- A failure in any location could be an XSS vulnerability

---

#### Part C: Additional Edge Case Testing

##### Step 11: Test Other Malicious Protocols (Optional)

Create additional test rules to verify other attack vectors. Each test follows the same pattern: create the YAML file, apply it, wait for the alert to fire, then verify in the UI.

---

**Test Case A: Case-Insensitive JavaScript Protocol**

Create a file named `test-rule-javascript-case.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-javascript-case
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertJavaScriptCase
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "JaVaScRiPt:alert(1)"
        summary: "Test alert with case-insensitive javascript protocol"
        description: "This tests that protocol validation is case-insensitive. Malicious URL: JaVaScRiPt:alert(1)"
```

Apply and verify:
```bash
oc apply -f test-rule-javascript-case.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `JaVaScRiPt:alert(1)`
- ✅ NOT clickable
- ✅ No JavaScript execution

---

**Test Case B: VBScript Protocol**

Create a file named `test-rule-vbscript.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-vbscript
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertVBScript
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "vbscript:msgbox(1)"
        summary: "Test alert with vbscript protocol"
        description: "This tests that vbscript: URLs are blocked. Malicious URL: vbscript:msgbox(1)"
```

Apply and verify:
```bash
oc apply -f test-rule-vbscript.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `vbscript:msgbox(1)`
- ✅ NOT clickable
- ✅ No script execution

---

**Test Case C: Mailto Protocol**

Create a file named `test-rule-mailto.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-mailto
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertMailto
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "mailto:security@example.com"
        summary: "Test alert with mailto protocol"
        description: "This tests that mailto: URLs are blocked. Test URL: mailto:security@example.com"
```

Apply and verify:
```bash
oc apply -f test-rule-mailto.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `mailto:security@example.com`
- ✅ NOT clickable
- ✅ No email client opens

---

**Test Case D: Relative URL**

Create a file named `test-rule-relative-url.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-relative-url
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertRelativeURL
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "/runbooks/alert"
        summary: "Test alert with relative URL"
        description: "This tests that relative URLs are blocked. Test URL: /runbooks/alert"
```

Apply and verify:
```bash
oc apply -f test-rule-relative-url.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `/runbooks/alert`
- ✅ NOT clickable
- ✅ Does not navigate anywhere

---

**Test Case E: Protocol-Relative URL**

Create a file named `test-rule-protocol-relative.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-protocol-relative
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertProtocolRelative
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "//evil.example.com/malicious"
        summary: "Test alert with protocol-relative URL"
        description: "This tests that protocol-relative URLs are blocked. Test URL: //evil.example.com/malicious"
```

Apply and verify:
```bash
oc apply -f test-rule-protocol-relative.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `//evil.example.com/malicious`
- ✅ NOT clickable
- ✅ Does not navigate anywhere

---

**Test Case F: File Protocol**

Create a file named `test-rule-file-protocol.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-file-protocol
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertFileProtocol
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "file:///etc/passwd"
        summary: "Test alert with file protocol"
        description: "This tests that file: URLs are blocked. Malicious URL: file:///etc/passwd"
```

Apply and verify:
```bash
oc apply -f test-rule-file-protocol.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `file:///etc/passwd`
- ✅ NOT clickable
- ✅ Does not access local filesystem

---

**Test Case G: Whitespace Prefix Attack**

Create a file named `test-rule-whitespace.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-whitespace
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertWhitespace
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "  javascript:alert(1)"
        summary: "Test alert with whitespace prefix"
        description: "This tests that leading whitespace doesn't bypass validation. Malicious URL:   javascript:alert(1)"
```

Apply and verify:
```bash
oc apply -f test-rule-whitespace.yaml
```

**Expected Result in UI**:
- ✅ URL renders as plain text: `  javascript:alert(1)`
- ✅ NOT clickable
- ✅ No JavaScript execution

---

**Test Case H: Invalid/Non-URL Text**

Create a file named `test-rule-invalid-text.yaml`:

```yaml
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-invalid-text
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertInvalidText
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "not a valid URL at all"
        summary: "Test alert with invalid text"
        description: "This tests that non-URL strings are handled gracefully. Test text: not a valid URL at all"
```

Apply and verify:
```bash
oc apply -f test-rule-invalid-text.yaml
```

**Expected Result in UI**:
- ✅ Text renders as plain text: `not a valid URL at all`
- ✅ NOT clickable
- ✅ No errors displayed

---

##### Testing All Edge Cases in the UI

After creating all the above rules, navigate through each alert in the console:

1. Go to **Observe > Alerting**
2. For each test alert (TestAlertJavaScriptCase, TestAlertVBScript, etc.):
   - Click on the alert name
   - Verify the runbook URL appears as **plain text** (not a link)
   - Verify **no external link icon** appears
   - Right-click and inspect to verify **no `<a>` tag** exists
   - Verify **no security issues** occur (no script execution, no unwanted behavior)

**Quick Verification Command**:
```bash
# Apply all edge case tests at once
for file in test-rule-javascript-case.yaml \
            test-rule-vbscript.yaml \
            test-rule-mailto.yaml \
            test-rule-relative-url.yaml \
            test-rule-protocol-relative.yaml \
            test-rule-file-protocol.yaml \
            test-rule-whitespace.yaml \
            test-rule-invalid-text.yaml; do
  oc apply -f "$file"
done

# Verify all alerts are created
oc get prometheusrule -n openshift-monitoring -l test=runbook-sanitization
```

---

---

#### Part D: Summary of All Test Files

For your reference, here's a complete list of all test files you can create:

| File Name | Alert Name | Test Scenario | Should Render As |
|-----------|------------|---------------|------------------|
| `test-rule-valid-https.yaml` | TestAlertValidHTTPS | Valid HTTPS URL | ✅ Clickable link |
| `test-rule-javascript.yaml` | TestAlertJavaScript | `javascript:` protocol | ❌ Plain text |
| `test-rule-data-uri.yaml` | TestAlertDataURI | `data:` URI | ❌ Plain text |
| `test-rule-javascript-case.yaml` | TestAlertJavaScriptCase | Case-insensitive JS | ❌ Plain text |
| `test-rule-vbscript.yaml` | TestAlertVBScript | `vbscript:` protocol | ❌ Plain text |
| `test-rule-mailto.yaml` | TestAlertMailto | `mailto:` protocol | ❌ Plain text |
| `test-rule-relative-url.yaml` | TestAlertRelativeURL | Relative URL | ❌ Plain text |
| `test-rule-protocol-relative.yaml` | TestAlertProtocolRelative | Protocol-relative URL | ❌ Plain text |
| `test-rule-file-protocol.yaml` | TestAlertFileProtocol | `file:` protocol | ❌ Plain text |
| `test-rule-whitespace.yaml` | TestAlertWhitespace | Whitespace prefix | ❌ Plain text |
| `test-rule-invalid-text.yaml` | TestAlertInvalidText | Invalid text | ❌ Plain text |

**Minimum Required Tests** (for basic validation):
- `test-rule-valid-https.yaml` ✅
- `test-rule-javascript.yaml` ❌
- `test-rule-data-uri.yaml` ❌

**Full Test Suite** (for comprehensive validation):
- All 11 tests listed above

---

#### Part E: Cleanup

##### Step 12: Remove Test Resources

After testing is complete, clean up the test PrometheusRules:

**Option 1: Delete all test rules at once (Recommended)**
```bash
# Delete all test rules by label
oc delete prometheusrule -n openshift-monitoring -l test=runbook-sanitization

# Confirm deletion
oc get prometheusrule -n openshift-monitoring -l test=runbook-sanitization
```
Expected output: `No resources found in openshift-monitoring namespace.`

**Option 2: Delete individual rules**
```bash
# Core test rules
oc delete prometheusrule test-runbook-valid-https -n openshift-monitoring
oc delete prometheusrule test-runbook-javascript -n openshift-monitoring
oc delete prometheusrule test-runbook-data-uri -n openshift-monitoring

# Edge case test rules (if you created them)
oc delete prometheusrule test-runbook-javascript-case -n openshift-monitoring
oc delete prometheusrule test-runbook-vbscript -n openshift-monitoring
oc delete prometheusrule test-runbook-mailto -n openshift-monitoring
oc delete prometheusrule test-runbook-relative-url -n openshift-monitoring
oc delete prometheusrule test-runbook-protocol-relative -n openshift-monitoring
oc delete prometheusrule test-runbook-file-protocol -n openshift-monitoring
oc delete prometheusrule test-runbook-whitespace -n openshift-monitoring
oc delete prometheusrule test-runbook-invalid-text -n openshift-monitoring
```

**Option 3: Delete using a loop**
```bash
# List of all test rule names
for rule in test-runbook-valid-https \
            test-runbook-javascript \
            test-runbook-data-uri \
            test-runbook-javascript-case \
            test-runbook-vbscript \
            test-runbook-mailto \
            test-runbook-relative-url \
            test-runbook-protocol-relative \
            test-runbook-file-protocol \
            test-runbook-whitespace \
            test-runbook-invalid-text; do
  oc delete prometheusrule "$rule" -n openshift-monitoring 2>/dev/null || true
done
```

**Clean up test namespace** (if you created one):
```bash
oc delete namespace test-monitoring
```

**Verify complete cleanup**:
```bash
# Should return no results
oc get prometheusrule -n openshift-monitoring | grep test-runbook

# Check that no test alerts are firing
oc exec -n openshift-monitoring prometheus-k8s-0 -- \
  promtool query instant http://localhost:9090 'ALERTS{test="runbook-sanitization"}'
```

Expected: No results or empty output

**Clean up local YAML files** (optional):
```bash
# If you want to remove the test files from your local directory
rm -f test-rule-*.yaml
```

---

### Troubleshooting Common Issues

#### Issue 1: "oc auth can-i" returns "no"
**Problem**: You don't have permissions to create PrometheusRules

**Solution**:
```bash
# Check your current permissions
oc auth can-i --list

# Request permissions from cluster admin, or ask them to create a role binding:
oc adm policy add-role-to-user monitoring-rules-edit <your-username> -n openshift-monitoring
```

#### Issue 2: Alerts Not Firing
**Problem**: After creating PrometheusRule, alerts don't appear in the UI

**Troubleshooting Steps**:
1. Check if the rule was created successfully:
   ```bash
   oc get prometheusrule test-runbook-valid-https -n openshift-monitoring -o yaml
   ```

2. Check Prometheus configuration:
   ```bash
   oc exec -n openshift-monitoring prometheus-k8s-0 -- promtool check config /etc/prometheus/prometheus.yml
   ```

3. Check Prometheus logs for errors:
   ```bash
   oc logs -n openshift-monitoring prometheus-k8s-0 | grep -i error
   ```

4. Verify the rule is loaded in Prometheus:
   ```bash
   oc port-forward -n openshift-monitoring prometheus-k8s-0 9090:9090
   # Navigate to http://localhost:9090/rules and search for your alert
   ```

5. Wait longer - alerts may take 1-2 minutes to appear

#### Issue 3: Cannot Access OpenShift Console
**Problem**: Browser cannot reach the console URL

**Solution**:
1. Get the correct console URL:
   ```bash
   oc get route console -n openshift-console -o jsonpath='{.spec.host}'
   ```

2. Verify you can reach it:
   ```bash
   curl -k https://$(oc get route console -n openshift-console -o jsonpath='{.spec.host}')
   ```

3. Check if you need VPN or network access to the cluster

#### Issue 4: Alerts Appear But No Runbook URL Section
**Problem**: Alert details page doesn't show runbook URL

**Possible Causes**:
1. The monitoring-plugin might not be updated to the version with the fix
2. The alert annotation key might be misspelled (should be `runbook_url` not `runbookUrl`)

**Solution**:
```bash
# Check monitoring-plugin version
oc get deployment monitoring-plugin -n openshift-monitoring -o jsonpath='{.spec.template.spec.containers[0].image}'

# Verify the annotation in your PrometheusRule
oc get prometheusrule test-runbook-valid-https -n openshift-monitoring -o jsonpath='{.spec.groups[0].rules[0].annotations}'
```

#### Issue 5: All URLs Appear as Plain Text (Including Valid Ones)
**Problem**: Even `https://` URLs are not rendered as links

**This indicates a problem** - the fix may have broken legitimate URLs

**Steps**:
1. Verify the monitoring-plugin version includes the fix
2. Check browser console for JavaScript errors (F12 > Console)
3. Inspect the DOM to see what HTML is being generated
4. Report this as a regression bug

#### Issue 6: PrometheusRule Creation Fails with Validation Error
**Problem**: `oc apply` returns validation errors

**Common Issues**:
- YAML indentation errors (use spaces, not tabs)
- Missing required fields
- Incorrect API version

**Solution**:
```bash
# Validate your YAML syntax
yamllint test-rule-valid-https.yaml

# Or use a YAML validator online
# Check the exact error message from oc apply
```

#### Issue 7: User Workload Monitoring Not Enabled
**Problem**: Want to test in your own namespace but user workload monitoring is disabled

**Solution**:
```bash
# Enable user workload monitoring (requires cluster-admin)
oc apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: cluster-monitoring-config
  namespace: openshift-monitoring
data:
  config.yaml: |
    enableUserWorkload: true
EOF

# Wait for user workload monitoring pods to start
oc get pods -n openshift-user-workload-monitoring
```

---

## Quick Start: Automated Test File Generation

For convenience, you can use this script to automatically create all test YAML files:

### Script: `create-all-test-rules.sh`

```bash
#!/bin/bash
# create-all-test-rules.sh
# Creates all PrometheusRule test files for runbook_url sanitization testing

set -e

NAMESPACE="openshift-monitoring"

echo "Creating test PrometheusRule YAML files..."
echo ""

# Create all test files in the current directory
# Test 1: Valid HTTPS URL
cat > test-rule-valid-https.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-valid-https
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertValidHTTPS
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "https://runbooks.example.com/test-alert"
        summary: "Test alert with valid HTTPS runbook URL"
        description: "This alert tests that valid HTTPS URLs render as clickable links. The URL is: https://runbooks.example.com/test-alert"
EOF

# Test 2: JavaScript Protocol
cat > test-rule-javascript.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-javascript
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertJavaScript
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "javascript:alert('XSS Attack!')"
        summary: "Test alert with javascript protocol (should be blocked)"
        description: "This alert tests that javascript: URLs are sanitized and render as plain text. Malicious URL: javascript:alert('XSS Attack!')"
EOF

# Test 3: Data URI
cat > test-rule-data-uri.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-data-uri
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertDataURI
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "data:text/html,<script>alert('XSS')</script>"
        summary: "Test alert with data URI (should be blocked)"
        description: "This alert tests that data: URLs are sanitized and render as plain text. Malicious URL: data:text/html,<script>alert('XSS')</script>"
EOF

echo "✅ Created 3 core test files"
echo ""

# Ask if user wants to create edge case tests too
read -p "Create additional edge case test files? (y/n) " -n 1 -r
echo
if [[ $REPLY =~ ^[Yy]$ ]]; then
  
  # Additional edge case tests
  cat > test-rule-javascript-case.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-javascript-case
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertJavaScriptCase
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "JaVaScRiPt:alert(1)"
        summary: "Test alert with case-insensitive javascript protocol"
        description: "This tests that protocol validation is case-insensitive. Malicious URL: JaVaScRiPt:alert(1)"
EOF

  cat > test-rule-vbscript.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-vbscript
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertVBScript
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "vbscript:msgbox(1)"
        summary: "Test alert with vbscript protocol"
        description: "This tests that vbscript: URLs are blocked. Malicious URL: vbscript:msgbox(1)"
EOF

  cat > test-rule-mailto.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-mailto
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertMailto
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "mailto:security@example.com"
        summary: "Test alert with mailto protocol"
        description: "This tests that mailto: URLs are blocked. Test URL: mailto:security@example.com"
EOF

  cat > test-rule-relative-url.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-relative-url
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertRelativeURL
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "/runbooks/alert"
        summary: "Test alert with relative URL"
        description: "This tests that relative URLs are blocked. Test URL: /runbooks/alert"
EOF

  cat > test-rule-protocol-relative.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-protocol-relative
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertProtocolRelative
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "//evil.example.com/malicious"
        summary: "Test alert with protocol-relative URL"
        description: "This tests that protocol-relative URLs are blocked. Test URL: //evil.example.com/malicious"
EOF

  cat > test-rule-file-protocol.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-file-protocol-1
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertFileProtocol
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
        label: file:///etc/passwd
      annotations:
        runbook_url: "file:///etc/passwd"
        summary: "Test alert with file protocol file:///etc/passwd"
        description: "This tests that file: URLs are blocked. Malicious URL: file:///etc/passwd"
EOF

  cat > test-rule-whitespace.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-whitespace
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertWhitespace
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "  javascript:alert(1)"
        summary: "Test alert with whitespace prefix"
        description: "This tests that leading whitespace doesn't bypass validation. Malicious URL:   javascript:alert(1)"
EOF

  cat > test-rule-invalid-text.yaml <<'EOF'
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: test-runbook-invalid-text
  namespace: openshift-monitoring
spec:
  groups:
  - name: test-runbook-sanitization
    interval: 30s
    rules:
    - alert: TestAlertInvalidText
      expr: vector(1)
      for: 0m
      labels:
        severity: warning
        test: runbook-sanitization
      annotations:
        runbook_url: "not a valid URL at all"
        summary: "Test alert with invalid text"
        description: "This tests that non-URL strings are handled gracefully. Test text: not a valid URL at all"
EOF

  echo "✅ Created 8 additional edge case test files"
fi

echo ""
echo "========================================="
echo "Test files created successfully!"
echo "========================================="
echo ""
echo "Files created in current directory:"
ls -1 test-rule-*.yaml
echo ""
echo "Next steps:"
echo "1. Apply all tests: for f in test-rule-*.yaml; do oc apply -f \$f; done"
echo "2. Wait 60 seconds for alerts to fire"
echo "3. Navigate to OpenShift Console > Observe > Alerting"
echo "4. Verify each alert according to the test plan"
echo "5. Clean up: oc delete prometheusrule -n ${NAMESPACE} -l test=runbook-sanitization"
echo ""
```

### Usage

**Make the script executable and run it**:
```bash
chmod +x create-all-test-rules.sh
./create-all-test-rules.sh
```

**Apply all created test files**:
```bash
# Apply all test files at once
for file in test-rule-*.yaml; do
  oc apply -f "$file"
done
```

**Or apply them individually**:
```bash
# Core tests only (minimum required)
oc apply -f test-rule-valid-https.yaml
oc apply -f test-rule-javascript.yaml
oc apply -f test-rule-data-uri.yaml

# All tests (comprehensive)
for file in test-rule-*.yaml; do
  oc apply -f "$file"
done
```

---

## Automated Testing

### Unit Tests
The fix includes comprehensive unit tests in `Link.spec.tsx`:

```bash
# Run unit tests
cd web
npm test -- Link.spec.tsx
```

**Expected Output**: All tests pass, including:
- Valid HTTP/HTTPS URLs render as links
- Invalid/unsafe URLs render as text
- LinkifyExternal converts URLs in text to safe links

### Automated Shell Script Test

```bash
#!/bin/bash
# test-runbook-url-sanitization.sh

set -e

NAMESPACE="openshift-monitoring"
TEST_RULE_NAME="test-runbook-sanitization"

echo "========================================="
echo "Testing runbook_url Sanitization"
echo "========================================="

# Test Case 1: Valid HTTPS URL
echo ""
echo "Test 1: Valid HTTPS URL"
cat <<EOF | oc apply -f -
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: ${TEST_RULE_NAME}-valid-https
  namespace: ${NAMESPACE}
spec:
  groups:
  - name: test
    rules:
    - alert: TestAlertValidHTTPS
      expr: vector(1)
      annotations:
        runbook_url: "https://runbooks.example.com/alert"
        summary: "Test alert with valid HTTPS runbook"
EOF

# Test Case 2: JavaScript Protocol (Should be sanitized)
echo ""
echo "Test 2: JavaScript Protocol (Should be blocked)"
cat <<EOF | oc apply -f -
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: ${TEST_RULE_NAME}-javascript
  namespace: ${NAMESPACE}
spec:
  groups:
  - name: test
    rules:
    - alert: TestAlertJavaScript
      expr: vector(1)
      annotations:
        runbook_url: "javascript:alert(1)"
        summary: "Test alert with javascript protocol"
EOF

# Test Case 3: Data URI (Should be sanitized)
echo ""
echo "Test 3: Data URI (Should be blocked)"
cat <<EOF | oc apply -f -
apiVersion: monitoring.coreos.com/v1
kind: PrometheusRule
metadata:
  name: ${TEST_RULE_NAME}-data-uri
  namespace: ${NAMESPACE}
spec:
  groups:
  - name: test
    rules:
    - alert: TestAlertDataURI
      expr: vector(1)
      annotations:
        runbook_url: "data:text/html,<script>alert(1)</script>"
        summary: "Test alert with data URI"
EOF

echo ""
echo "========================================="
echo "Test PrometheusRules created successfully"
echo "========================================="
echo ""
echo "Manual Verification Steps:"
echo "1. Wait for alerts to fire (they will fire immediately due to vector(1))"
echo "2. Navigate to OpenShift Console > Observe > Alerting"
echo "3. Find and click on each test alert"
echo "4. Verify:"
echo "   - TestAlertValidHTTPS: runbook_url is a clickable link"
echo "   - TestAlertJavaScript: runbook_url is plain text (not clickable)"
echo "   - TestAlertDataURI: runbook_url is plain text (not clickable)"
echo ""
echo "Cleanup:"
echo "  oc delete prometheusrule -n ${NAMESPACE} -l test=runbook-sanitization"
echo ""
```

### Running the Automated Test

```bash
# Make script executable
chmod +x test-runbook-url-sanitization.sh

# Run the test (requires oc login and cluster-admin permissions)
./test-runbook-url-sanitization.sh
```

---

## Browser-Based E2E Test

For a more complete end-to-end test, use Cypress or Playwright:

```javascript
// Example Cypress test
describe('Runbook URL Sanitization', () => {
  it('renders valid HTTPS URLs as links', () => {
    cy.visit('/monitoring/alerts/TestAlertValidHTTPS');
    cy.get('a[href="https://runbooks.example.com/alert"]')
      .should('exist')
      .should('have.attr', 'target', '_blank')
      .should('have.attr', 'rel', 'noopener noreferrer');
  });

  it('renders javascript: URLs as plain text', () => {
    cy.visit('/monitoring/alerts/TestAlertJavaScript');
    cy.contains('javascript:alert(1)')
      .should('exist')
      .should('not.have.attr', 'href');
  });

  it('renders data: URLs as plain text', () => {
    cy.visit('/monitoring/alerts/TestAlertDataURI');
    cy.contains('data:text/html')
      .should('exist')
      .should('not.have.attr', 'href');
  });
});
```

---

## Success Criteria

✅ All unit tests pass  
✅ Valid HTTP/HTTPS URLs render as clickable external links in **runbook_url** field  
✅ Valid HTTP/HTTPS URLs render as clickable links in **description** field (LinkifyExternal)  
✅ Valid links work correctly in all UI locations: alert details, notification bell, alerts list  
✅ JavaScript protocol URLs are blocked and rendered as text in **all fields**  
✅ Data URI protocols are blocked and rendered as text in **all fields**  
✅ Other unsafe protocols (vbscript, mailto, file) are blocked in **all fields**  
✅ Relative and protocol-relative URLs are blocked in **all fields**  
✅ No XSS vulnerabilities can be triggered via runbook_url or description annotations  
✅ User experience for legitimate runbook URLs remains unchanged  
✅ Both `ExternalLink` and `LinkifyExternal` components properly sanitize URLs  

---

## Regression Testing

Ensure existing functionality is not broken:

1. **Alert Rules with valid runbook URLs**: Verify existing alerts with legitimate runbook URLs still work correctly
2. **Alert Rules without runbook URLs**: Verify alerts without runbook URLs display correctly without errors
3. **Multiple alerts**: Verify the fix works consistently across multiple alerts on the same page
4. **LinkifyExternal component**: Verify other uses of the LinkifyExternal component still work correctly

---

## Security Verification

**XSS Attack Vectors to Test**:
- [ ] `javascript:` protocol
- [ ] `javascript:` with various case combinations
- [ ] `data:` protocol with embedded HTML/scripts
- [ ] `vbscript:` protocol
- [ ] `file:` protocol
- [ ] URLs with embedded newlines or special characters
- [ ] Unicode-escaped javascript protocol
- [ ] URLs with null bytes

If any of these render as clickable links, the fix is incomplete.

---

## Appendix A: Browser Developer Tools Guide

For testers new to browser developer tools, here's a quick guide:

### Opening Developer Tools

**Method 1: Keyboard Shortcut**
- Windows/Linux: Press `F12` or `Ctrl + Shift + I`
- Mac: Press `Cmd + Option + I`

**Method 2: Right-Click Menu**
- Right-click anywhere on the page
- Select "Inspect" or "Inspect Element"

**Method 3: Browser Menu**
- Chrome: Menu (⋮) > More Tools > Developer Tools
- Firefox: Menu (≡) > More Tools > Web Developer Tools
- Edge: Menu (…) > More Tools > Developer Tools

### Understanding the Developer Tools Layout

When DevTools open, you'll see several tabs:

1. **Elements (Chrome/Edge)** or **Inspector (Firefox)**
   - Shows the HTML structure (DOM tree)
   - This is where you'll inspect link elements
   - Selected elements are highlighted in the browser window

2. **Console**
   - Shows JavaScript errors, warnings, and logs
   - Used to verify no scripts are executing

3. **Network**, **Sources**, etc.
   - Not needed for this test

### How to Inspect an Element

**Step-by-step**:
1. Right-click directly on the text/link you want to inspect
2. Select "Inspect" or "Inspect Element"
3. DevTools opens with the element already highlighted
4. The highlighted line shows the HTML for that element

**Example of what you'll see**:

```html
<!-- This is what a valid link looks like in DevTools: -->
<a href="https://example.com" target="_blank" rel="noopener noreferrer">
  https://example.com
</a>

<!-- This is what blocked text looks like: -->
<span>javascript:alert(1)</span>
```

### Reading HTML Attributes

When inspecting an element, you'll see attributes inside the opening tag:

```html
<a href="..." target="..." rel="...">
   ↑    ↑       ↑         ↑
   |    |       |         └─ These are attributes
   |    |       └─────────── attribute name="value"
   |    └─────────────────── href attribute (the URL)
   └──────────────────────── element type (anchor/link)
```

**What to look for**:
- `<a>` = anchor tag (link element)
- `href` = the destination URL
- `target="_blank"` = opens in new tab
- `rel="noopener noreferrer"` = security attributes

### Using the Console Tab

1. Click the **Console** tab in DevTools
2. Look for messages in red (errors) or yellow (warnings)
3. For this test, the console should be **empty** or show only unrelated messages
4. If you see errors like "Refused to execute...", that might indicate a security issue

### Closing Developer Tools

- Press `F12` again, or
- Click the X button in the DevTools panel, or
- Press `Esc` (if DevTools are docked at bottom)

---

## Appendix B: Common OpenShift Console Locations

To help you navigate the OpenShift Console:

### Main Navigation Structure
```
OpenShift Console
├── Home
├── Operators
├── Workloads
├── Networking
├── Storage
├── Builds
├── Pipelines
├── Observe ← START HERE
│   ├── Dashboards
│   ├── Metrics
│   ├── Alerting ← GO HERE
│   ├── Logs
│   └── Distributed Tracing
├── Compute
└── Administration
```

### Alerting Page Layout
```
┌──────────────────────────────────────────────────────────────┐
│ Alerting                                           [Filter ▼]│
├──────────────────────────────────────────────────────────────┤
│ [Alerts] [Silences] [Alerting Rules]                         │
├──────────────────────────────────────────────────────────────┤
│ Filter: [_________________________] [State ▼] [Severity ▼]  │
├──────────────────────────────────────────────────────────────┤
│ Name               │ Labels    │ Severity │ State │ Rule     │
├────────────────────┼───────────┼──────────┼───────┼──────────┤
│ TestAlertValidHTTP │ test=r... │ warning  │ Firing│ Test...  │
│ TestAlertJavaScrip │ test=r... │ warning  │ Firing│ Test...  │
│ TestAlertDataURI   │ test=r... │ warning  │ Firing│ Test...  │
└──────────────────────────────────────────────────────────────┘
```

### Alert Details Page Layout
```
┌──────────────────────────────────────────────────────────────┐
│ < Back to Alerting                                            │
├──────────────────────────────────────────────────────────────┤
│ TestAlertValidHTTPS                           [⋮ Actions]    │
├──────────────────────────────────────────────────────────────┤
│ Alert Details                                                 │
│   Severity: warning                                          │
│   State: Firing                                              │
│   Alerting rule: TestAlertValidHTTPS                         │
├──────────────────────────────────────────────────────────────┤
│ Labels                                                        │
│   severity: warning                                          │
│   test: runbook-sanitization                                 │
├──────────────────────────────────────────────────────────────┤
│ Annotations                                                   │
│   Summary: Test alert with valid HTTPS runbook URL          │
│   Description: This alert tests that valid HTTPS URLs...    │
│   Runbook URL: [https://runbooks.example.com/test-alert 🔗]  │
│                 ↑↑↑ THIS IS WHAT YOU'RE TESTING ↑↑↑          │
└──────────────────────────────────────────────────────────────┘
```

---

## Appendix C: Expected Test Results Summary

Quick reference table for expected results:

**Note**: These results apply to URLs in BOTH the `runbook_url` field AND the `description` field.

| Test Case | runbook_url Value | Should be Link? | Should Execute? |
|-----------|-------------------|-----------------|-----------------|
| Valid HTTPS | `https://runbooks.example.com/test` | ✅ YES | N/A |
| Valid HTTP | `http://runbooks.example.com/test` | ✅ YES | N/A |
| JavaScript | `javascript:alert(1)` | ❌ NO (plain text) | ❌ NO |
| JavaScript (case) | `JaVaScRiPt:alert(1)` | ❌ NO (plain text) | ❌ NO |
| Data URI | `data:text/html,<script>...` | ❌ NO (plain text) | ❌ NO |
| VBScript | `vbscript:msgbox(1)` | ❌ NO (plain text) | ❌ NO |
| Mailto | `mailto:test@example.com` | ❌ NO (plain text) | N/A |
| Relative URL | `/runbooks/alert` | ❌ NO (plain text) | N/A |
| Protocol-relative | `//evil.example.com/alert` | ❌ NO (plain text) | N/A |
| Invalid text | `not a URL` | ❌ NO (plain text) | N/A |

**Legend**:
- ✅ YES = Expected behavior
- ❌ NO = Should be blocked
- N/A = Not applicable
