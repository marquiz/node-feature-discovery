#!/bin/bash
kubectl_cmd="kubectl get no -o json"
jq_cmd="jq -C -r '.items[] | {name:.metadata.name, labels:.metadata.labels, annotations:.metadata.annotations, allocatable: {\"gpu.intel.com/i915\": .status.allocatable.\"gpu.intel.com/i915\"}}'"
grep_cmd="grep -v -e '\"beta.kubernetes.io' -e '.alpha.kubernetes.io' -e 'volumes.kubernetes.io' -e '\"kubernetes.io'"
watch --color "$kubectl_cmd | $jq_cmd | $grep_cmd"
