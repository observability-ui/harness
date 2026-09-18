#!/bin/bash

# Restore the original custom rules (without runbook_url test alerts)
# This removes the 'runbook-url-sanitization-tests' group from the ConfigMap

echo "Restoring original custom rules ConfigMap (removing test alerts)..."

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
EOF

echo ""
echo "✅ ConfigMap restored to original state!"
echo ""
echo "Wait 1-2 minutes for Thanos Ruler to reload the configuration."
echo "The TestRunbook* alerts should disappear from the ACM Console."
echo ""
echo "To verify:"
echo "  oc get configmap thanos-ruler-custom-rules -n open-cluster-management-observability -o yaml"
echo ""
