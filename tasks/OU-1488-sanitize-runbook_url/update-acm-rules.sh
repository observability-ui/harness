#!/bin/bash

# Update only the custom rules ConfigMap with runbook_url sanitization tests
# Use this script if ACM is already running and you just want to update the alert rules

echo "Updating thanos-ruler-custom-rules ConfigMap with runbook_url test alerts..."

oc apply -f - <<EOF
apiVersion: v1
kind: ConfigMap
metadata:
  name: thanos-ruler-custom-rules
  namespace: open-cluster-management-observability
  labels:
    cluster.open-cluster-management.io/backup: ""
data:
  custom_rules.yaml: |
    groups:
      - name: alertrule-testing
        rules:
        - alert: Watchdog
          annotations:
            summary: An alert that should always be firing to certify that Alertmanager is working properly.
            description: This is an alert meant to ensure that the entire alerting pipeline is functional.
          expr: vector(1)
          labels:
            instance: "local"
            cluster: "local"
            clusterID: "111111111"
            severity: info
        - alert: Watchdog-spoke
          annotations:
            summary: An alert that should always be firing to certify that Alertmanager is working properly.
            description: This is an alert meant to ensure that the entire alerting pipeline is functional.
          expr: vector(1)
          labels:
            instance: "spoke"
            cluster: "spoke"
            clusterID: "22222222"
            severity: warn
      - name: cluster-health
        rules:
        - alert: ClusterCPUHealth-jb
          annotations:
            summary: Notify when CPU utilization on a cluster is greater than the defined utilization limit
            description: "The cluster has a high CPU usage: core for."
          expr: |
            max(cluster:cpu_usage_cores:sum) by (clusterID, cluster, prometheus) > 0
          labels:
            cluster: "{{ \$labels.cluster }}"
            prometheus: "{{ \$labels.prometheus }}"
            severity: critical
      - name: runbook-url-sanitization-tests
        rules:
        - alert: TestRunbookValidHTTPS
          annotations:
            runbook_url: "https://runbooks.example.com/cluster-health"
            summary: Test alert with valid HTTPS runbook URL
            description: "This alert tests that valid HTTPS URLs render as clickable links. Runbook: https://runbooks.example.com/cluster-health"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookJavaScript
          annotations:
            runbook_url: "javascript:alert('XSS Attack from ACM!')"
            summary: Test alert with malicious javascript protocol (should be blocked)
            description: "This alert tests that javascript: URLs are sanitized and render as plain text. Malicious URL: javascript:alert('XSS Attack from ACM!')"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookDataURI
          annotations:
            runbook_url: "data:text/html,<script>alert('XSS from ACM')</script>"
            summary: Test alert with malicious data URI (should be blocked)
            description: "This alert tests that data: URLs are sanitized and render as plain text. Malicious URL: data:text/html,<script>alert('XSS from ACM')</script>"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookCaseSensitive
          annotations:
            runbook_url: "JaVaScRiPt:alert('Case test')"
            summary: Test alert with case-insensitive javascript protocol
            description: "This tests that protocol validation is case-insensitive. Malicious URL: JaVaScRiPt:alert('Case test')"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookVBScript
          annotations:
            runbook_url: "vbscript:msgbox('ACM XSS')"
            summary: Test alert with vbscript protocol (should be blocked)
            description: "This tests that vbscript: URLs are blocked. Malicious URL: vbscript:msgbox('ACM XSS')"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookFileProtocol
          annotations:
            runbook_url: "file:///etc/passwd"
            summary: Test alert with file protocol (should be blocked)
            description: "This tests that file: URLs are blocked. Malicious URL: file:///etc/passwd"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookRelativeURL
          annotations:
            runbook_url: "/etc/passwd"
            summary: Test alert with relative URL (should be blocked)
            description: "This tests that relative URLs are blocked. Test URL: /etc/passwd"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
        - alert: TestRunbookProtocolRelative
          annotations:
            runbook_url: "//evil.example.com/xss"
            summary: Test alert with protocol-relative URL (should be blocked)
            description: "This tests that protocol-relative URLs are blocked. Test URL: //evil.example.com/xss"
          expr: vector(1)
          labels:
            cluster: "test"
            severity: warning
            test: runbook-sanitization
EOF

echo ""
echo "✅ ConfigMap updated successfully!"
echo ""
echo "The custom rules have been updated with runbook_url sanitization test alerts."
echo ""
echo "Wait 1-2 minutes for Thanos Ruler to reload the configuration, then:"
echo "1. Go to ACM Console > Overview > Alerts"
echo "2. Look for alerts starting with 'TestRunbook*'"
echo "3. Click on each alert to view details"
echo "4. Verify:"
echo "   - TestRunbookValidHTTPS: runbook_url is a clickable link"
echo "   - TestRunbookJavaScript: runbook_url is plain text (not clickable)"
echo "   - TestRunbookDataURI: runbook_url is plain text (not clickable)"
echo "   - Other test alerts: URLs are plain text (not clickable)"
echo "5. Check the Description field for each alert - URLs should be sanitized there too"
echo ""
echo "To verify the ConfigMap was updated:"
echo "  oc get configmap thanos-ruler-custom-rules -n open-cluster-management-observability -o yaml"
echo ""
echo "To clean up test alerts later:"
echo "  # Edit the ConfigMap and remove the 'runbook-url-sanitization-tests' group"
echo ""
