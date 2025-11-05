#!/bin/bash

version=$1
helm_chart_version=${1#v}

if [ -z "$version" ]; then
  echo "Usage: $0 <version>"
  exit 1
fi

for variant in "" -minimal -full; do
    img=node-feature-discovery:${version}${variant}
    printf "%-40s %s\n" ${img}: $(oras discover --format json gcr.io/k8s-staging-nfd/${img} | jq .digest)
done

img=charts/node-feature-discovery:${helm_chart_version}
printf "%-40s %s\n " ${img}: $(oras discover --format json gcr.io/k8s-staging-nfd/${img} | jq .digest)
