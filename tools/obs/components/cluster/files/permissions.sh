#!/bin/bash
set -e -o pipefail

NS="${PERMISSIONS_NAMESPACE:-openshift-monitoring}"

assign() {
  local role="$1" user="$2" scope="$3"
  if [ "$scope" = "cluster" ]; then
    oc adm policy add-cluster-role-to-user "$role" "$user" && \
      echo "  $user → $role (cluster)" || \
      echo "  WARN: failed to assign $role to $user"
  else
    oc -n "$NS" policy add-role-to-user "$role" "$user" && \
      echo "  $user → $role ($NS)" || \
      echo "  WARN: failed to assign $role to $user in $NS"
  fi
}

echo "=== Cluster Admins ==="
for u in admin1 admin2 admin3 admin4 admin5; do
  assign cluster-admin "$u" cluster
done

echo ""
echo "=== Monitoring Viewers ==="
for u in mon-viewer1 mon-viewer2; do
  assign view "$u" ns
  assign cluster-monitoring-view "$u" cluster
done

echo ""
echo "=== Monitoring Editors ==="
for u in mon-editor1 mon-editor2; do
  assign view "$u" ns
  assign cluster-monitoring-view "$u" cluster
  assign monitoring-rules-edit "$u" ns
done

echo ""
echo "=== Logging Viewers ==="
for u in log-viewer1 log-viewer2; do
  assign view "$u" ns
  assign cluster-logging-application-view "$u" cluster
done

echo ""
echo "=== Logging Editors ==="
for u in log-editor1 log-editor2; do
  assign view "$u" ns
  assign cluster-logging-application-view "$u" cluster
  assign monitoring-rules-edit "$u" ns
done

echo ""
echo "=== Perses Viewers ==="
for u in perses-viewer1 perses-viewer2; do
  assign view "$u" ns
  assign persesdashboard-viewer-role "$u" ns
  assign persesdatasource-viewer-role "$u" ns
done

echo ""
echo "=== Perses Editors ==="
for u in perses-editor1 perses-editor2; do
  assign view "$u" ns
  assign persesdashboard-editor-role "$u" ns
  assign persesdatasource-editor-role "$u" ns
done

echo ""
echo "All permissions assigned."
