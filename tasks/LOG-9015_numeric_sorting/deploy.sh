#!/usr/bin/env bash
set -euo pipefail

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
NAMESPACE="log-9015-numeric-sorting"

echo "Cleaning up previous deployment (if any)..."
oc delete -f "${SCRIPT_DIR}/manifests/" --ignore-not-found

echo "Applying manifests..."
oc apply -f "${SCRIPT_DIR}/manifests/"

echo "Waiting for deployment rollout..."
oc rollout status deployment/log-generator -n "${NAMESPACE}" --timeout=120s

echo "Tailing logs (Ctrl+C to stop)..."
oc logs -n "${NAMESPACE}" deployment/log-generator -f
