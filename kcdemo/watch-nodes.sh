#!/bin/bash
watch --color "kubectl get no -o json | jq -C -r '.items[] | {name:.metadata.name, labels:.metadata.labels, annotations:.metadata.annotations}' | grep -v -e '\"beta.kubernetes.io' -e '.alpha.kubernetes.io' -e 'volumes.kubernetes.io' -e '\"kubernetes.io'"
