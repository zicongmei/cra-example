#!/bin/bash

OP=apply

cat <<EOF | kubectl ${OP} -f -

apiVersion: pkg.crossplane.io/v1
kind: Provider
metadata:
  name: provider-gcp
spec:
  package: xpkg.upbound.io/upbound/provider-gcp:v0.31.0
EOF

