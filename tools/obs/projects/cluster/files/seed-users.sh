#!/bin/bash
set -e -o pipefail

HTPASSWD_FILE="$1"
OAUTH_FILE="$2"

echo "Creating htpasswd secret in openshift-config..."
if oc get secret htpass-secret -n openshift-config &>/dev/null; then
  echo "Secret already exists, updating..."
  oc create secret generic htpass-secret \
    --from-file=htpasswd="$HTPASSWD_FILE" \
    -n openshift-config \
    --dry-run=client -o yaml | oc apply -f -
else
  oc create secret generic htpass-secret \
    --from-file=htpasswd="$HTPASSWD_FILE" \
    -n openshift-config
fi
echo "htpasswd secret created/updated"

echo "Configuring OAuth identity provider..."
oc apply -f "$OAUTH_FILE"
echo "OAuth configured"

echo "Waiting for OAuth pods to restart..."
oc rollout status deployment/oauth-openshift -n openshift-authentication --timeout=120s || \
  echo "OAuth rollout may still be in progress"

echo "Users seeded. Password = username for all users."
echo "Example: oc login -u admin1 -p admin1"
