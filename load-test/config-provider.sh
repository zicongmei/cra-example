#!/bin/bash

OP=apply


#gcloud iam service-accounts keys create \
#      --iam-account=node-sa@zicong-gke-multi-cloud-dev-2.iam.gserviceaccount.com \
#      ~/service_account_keys/node-sa-key.json


kubectl create ns crossplane-config
kubectl --namespace crossplane-config \
    create secret generic gcp-creds \
    --from-file creds=~/service_account_keys/node-sa-key.json

PROJECT_ID=`gcloud config get-value project`

cat <<EOF | kubectl ${OP} -f -

apiVersion: gcp.upbound.io/v1beta1
kind: ProviderConfig
metadata:
  name: default
spec:
  projectID: $PROJECT_ID
  credentials:
    source: Secret
    secretRef:
      namespace: crossplane-config
      name: gcp-creds
      key: creds
EOF