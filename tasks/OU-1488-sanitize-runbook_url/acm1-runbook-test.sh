#!/bin/bash

# Modified version of acm1.sh to test runbook_url sanitization
# This script includes test alerts with both valid and malicious runbook URLs

oc apply -f - <<EOF
apiVersion: operator.open-cluster-management.io/v1
kind: MultiClusterHub
metadata:
  name: multiclusterhub
  namespace: open-cluster-management
spec: {}
EOF


oc apply -f - <<EOF
apiVersion: observability.open-cluster-management.io/v1beta2
kind: MultiClusterObservability
metadata:
  name: observability
spec:
  observabilityAddonSpec: {}
  storageConfig:
    metricObjectStorage:
      name: thanos-object-storage
      key: thanos.yaml
EOF

oc apply -f - <<EOF
apiVersion: v1
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
            cluster: "{{ $labels.cluster }}"
            prometheus: "{{ $labels.prometheus }}"
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
kind: ConfigMap
metadata:
  annotations:
    kubectl.kubernetes.io/last-applied-configuration: |
      {"apiVersion":"v1","data":{"custom_rules.yaml":"groups:\n  - name: alertrule-testing\n    rules:\n    - alert: Watchdog\n      annotations:\n        summary: An alert that should always be firing to certify that Alertmanager is working properly.\n        description: This is an alert meant to ensure that the entire alerting pipeline is functional.\n      expr: vector(1)\n      labels:\n        instance: \"local\"\n        cluster: \"local\"\n        clusterID: \"111111111\"\n        severity: info\n    - alert: Watchdog-spoke\n      annotations:\n        summary: An alert that should always be firing to certify that Alertmanager is working properly.\n        description: This is an alert meant to ensure that the entire alerting pipeline is functional.\n      expr: vector(1)\n      labels:\n        instance: \"spoke\"\n        cluster: \"spoke\"\n        clusterID: \"22222222\"\n        severity: warn\n  - name: cluster-health\n    rules:\n    - alert: ClusterCPUHealth-jb\n      annotations:\n        summary: Notify when CPU utilization on a cluster is greater than the defined utilization limit\n        description: \"The cluster has a high CPU usage: {{  }} core for {{ .cluster }} {{ .clusterID }}.\"\n      expr: |\n        max(cluster:cpu_usage_cores:sum) by (clusterID, cluster, prometheus) > 0\n      labels:\n        cluster: \"{{ .cluster }}\"\n        prometheus: \"{{ .prometheus }}\"\n        severity: critical\n"},"kind":"ConfigMap","metadata":{"annotations":{},"name":"thanos-ruler-custom-rules","namespace":"open-cluster-management-observability"}}
  creationTimestamp: "2025-06-01T06:16:10Z"
  labels:
    cluster.open-cluster-management.io/backup: ""
  name: thanos-ruler-custom-rules
  namespace: open-cluster-management-observability
  resourceVersion: "192432"
  uid: 969bf381-0963-45fb-b1cb-11270c982ea2
EOF
