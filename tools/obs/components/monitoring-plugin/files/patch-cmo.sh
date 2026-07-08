set -e
INDEX=$(oc get deploy cluster-monitoring-operator -n openshift-monitoring -o json | \
  jq -r '.spec.template.spec.containers[0].args | to_entries[] | select(.value | contains("monitoring-plugin")) | .key')
oc patch deploy cluster-monitoring-operator -n openshift-monitoring --type=json \
  -p "[{\"op\":\"replace\",\"path\":\"/spec/template/spec/containers/0/args/$INDEX\",\"value\":\"--images=monitoring-plugin=$MP_IMAGE\"}]"
