# Deploy gke-on-aws with CRA

Use can deploy:
1. All clusters in a same VPC
2. Each cluster has its own VPC

## All clusters in same VPC

### 1. Deploy VPC and other prerequisites
```shell
pushd prerequisites
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

Then replace the info in [cluster-1.yaml](cluster%2Fcluster-1.yaml) and [cluster-2.yaml](cluster%2Fcluster-2.yaml).

Finally, apply the clusters

```shell
pushd cluster
kubectl apply -f xrd.yaml
kubectl apply -f composition.yaml
kubectl apply -f cluster-1.yaml
kubectl apply -f cluster-2.yaml
popd
```

## Each cluster has its own VPC
TODO