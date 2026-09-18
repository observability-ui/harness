for file in test-rule-*.yaml; do
  oc apply -f "$file"
done