#!/bin/bash

OP=apply
#OP=delete

XRD_NUM=4
COMPOSITE_PER_XRD=3
CLAIM_PER_COMPOSITE=2


xrd_i=0
composite_i=0
claim_i=0

xrd_file=/tmp/xrds.yaml
composition_file=/tmp/compositions.yaml
claim_file=/tmp/claims.yaml
rm -f ${xrd_file}
rm -f ${composition_file}
rm -f ${claim_file}

for (( xrd_i=0; xrd_i<$XRD_NUM; xrd_i++ ))
do
export xrd_i composite_i claim_i
cat xrd.yaml.template | envsubst >> ${xrd_file}
for (( composite_i=0; composite_i<$COMPOSITE_PER_XRD; composite_i++ ))
do
export xrd_i composite_i claim_i
cat composition.yaml.template | envsubst >> ${composition_file}
for (( claim_i=0; claim_i<$CLAIM_PER_COMPOSITE; claim_i++ ))
do
export xrd_i composite_i claim_i
cat claim.yaml.template | envsubst >> ${claim_file}
done
done
done

kubectl ${OP} -f ${xrd_file}
kubectl ${OP} -f ${composition_file}
kubectl ${OP} -f ${claim_file}
