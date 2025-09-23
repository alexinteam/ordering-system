#!/bin/bash

kubectl delete namespace api-gateway-system --ignore-not-found=true
kubectl delete namespace notification-system --ignore-not-found=true
kubectl delete namespace billing-system --ignore-not-found=true
kubectl delete namespace order-system --ignore-not-found=true
kubectl delete namespace infrastructure --ignore-not-found=true