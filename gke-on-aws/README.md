# Deploy gke-on-aws with CRA

## All clusters in same VPC

### 1. Deploy VPC and other prerequisites
```shell
pushd [prerequisites](prerequisites)
kubectl apply -f xrd.yaml
kubectl apply -f composition.yaml
kubectl apply -f claim.yaml
popd
```

### 2. Deploy the clusters

Get the parameters from [prerequisites](prerequisites):


```shell
kubectl get  compositeprerequisites.demo.anthos.com/vpc-1 -o json | jq .status
```

Then apply the clusters

```shell
pushd [cluster](cluster)
kubectl apply -f xrd.yaml
kubectl apply -f composition.yaml
kubectl apply -f cluster-1.yaml
kubectl apply -f cluster-2.yaml
popd
```