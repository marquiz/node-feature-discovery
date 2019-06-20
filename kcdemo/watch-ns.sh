#!/bin/bash
kubectl -n $1 get ds,svc -o custom-columns='KIND:kind,NAME:metadata.name'

echo
kubectl -n $1 get po -o custom-columns='KIND:kind,NAME:metadata.name,STATUS:status.phase,NODE:spec.nodeName'
