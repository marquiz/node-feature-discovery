#!/bin/sh -ex

go-to-protobuf \
   --output-base=. \
   --go-header-file ../../../../hack/boilerplate.go.txt \
   --proto-import ../../../../vendor/ \
   --packages +sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha2=v1alpha2 \
   --keep-gogoproto=false \
   --apimachinery-packages "-k8s.io/apimachinery/pkg/util/intstr"

mv sigs.k8s.io/node-feature-discovery/pkg/apis/nfd/v1alpha2/* .
rm -rf sigs.k8s.io
