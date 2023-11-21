# Setup an AWS provider to a cluster without crossplane

kubectl apply -f crds
kubectl apply -f clusterrolebinding.yaml
kubectl apply -f deployment.yaml
kubectl apply -f clusterrole.yaml

NAMESPACE=crossplane-config

kubectl create ns ${NAMESPACE}
kubectl --namespace ${NAMESPACE} \
    create secret generic aws-creds \
    --from-file creds=${HOME}/.aws/credentials

cat <<EOF > /tmp/aws-cred.yaml
apiVersion: aws.gke.cloud.google.com/v1beta1
kind: ProviderConfig
metadata:
  name: default-aws
spec:
  credentials:
    source: Secret
    secretRef:
      namespace: ${NAMESPACE}
      name: aws-creds
      key: creds
EOF

kubectl apply -f /tmp/aws-cred.yaml

cat << EOF > /tmp/vpc.yaml
apiVersion: ec2.aws.gke.cloud.google.com/v1beta1
kind: VPC
metadata:
  name: zicong-cp-vpc
spec:
  providerConfigRef:
    name: default-aws
  forProvider:
    cidrBlock: 10.0.0.0/16
    region: us-east-1
    tags:
      Name: zicong-crossplane-test
EOF

kubectl apply -f /tmp/vpc.yaml


