#!/bin/bash
while true
do
  kubectl delete pod -n crossplane-system $(kubectl get pod -n crossplane-system | grep aws | awk '{print $1;}')
  sleep 1
done