
Get the values from prerequisites:

```shell
kubectl get  compositeprerequisites.demo.anthos.com/vpc-1 -o json | jq .status
```

Then apply the yamls

1. `kubectl apply -f xrd.yaml`
2. `kubectl apply -f composition.yaml`
3. `kubectl apply -f claim.yaml`

