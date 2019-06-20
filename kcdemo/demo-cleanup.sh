#!/bin/bash -x

ns="nfd"
if [ -n "$1" ]; then
    ns=$1
fi

# Remove objects
kubectl -n $ns delete ds/nfd-master \
               ds/nfd-worker \
               svc/nfd-master \
               clusterrolebinding/nfd-master \
               clusterrole/nfd-master \
               sa/nfd-master

# Remove annotations and labels
kubectl get no -o yaml | sed -e '/^\s*nfd.node.kubernetes.io/d' -e '/^\s*feature.node.kubernetes.io/d' | kubectl replace -f -

kubectl delete ns demo
kubectl delete ns nfd
