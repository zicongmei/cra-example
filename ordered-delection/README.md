# deploy

```shell
kubectl apply -f xrd.yaml
kubectl apply -f composition.yaml
kubectl apply -f claim-background.yaml
kubectl apply -f claim-foreground.yaml
```

# add some dependency to block the deletction

```shell
SUBNET_B_ID=$(kubectl get cvpc vpc-b -ojson| jq -r .status.subnetId)

aws ec2 create-network-interface --subnet-id ${SUBNET_B_ID} --description "my network interface"
```


```shell
SUBNET_F_ID=$(kubectl get cvpc vpc-f -ojson| jq -r .status.subnetId)

aws ec2 create-network-interface --subnet-id ${SUBNET_F_ID} --description "my network interface"
```

# remove dependency

```shell
IF_B_ID=$(aws ec2 describe-network-interfaces --filters Name=subnet-id,Values=${SUBNET_B_ID} |jq -r .NetworkInterfaces[0].NetworkInterfaceId)
aws ec2 delete-network-interface --network-interface-id ${IF_B_ID}
```

```shell
IF_F_ID=$(aws ec2 describe-network-interfaces --filters Name=subnet-id,Values=${SUBNET_F_ID} |jq -r .NetworkInterfaces[0].NetworkInterfaceId)
aws ec2 delete-network-interface --network-interface-id ${IF_F_ID}
```