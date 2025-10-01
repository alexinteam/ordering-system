#!/bin/bash

set -e

kubectl apply -f k8s/namespaces.yaml

kubectl apply -f k8s/rabbitmq.yaml
kubectl apply -f k8s/rabbitmq-ingress.yaml
kubectl apply -f k8s/rabbitmq-static-ingress.yaml

kubectl wait --namespace infrastructure --for=condition=ready pod --selector=app=rabbitmq --timeout=60s

kubectl apply -f k8s/billing-service.yaml
kubectl apply -f k8s/warehouse-service.yaml
kubectl apply -f k8s/delivery-service.yaml
kubectl apply -f k8s/notification-service.yaml
kubectl apply -f k8s/order-service.yaml
kubectl apply -f k8s/api-gateway.yaml

kubectl wait --namespace billing-system --for=condition=ready pod --selector=app=billing-service --timeout=60s
kubectl wait --namespace warehouse-system --for=condition=ready pod --selector=app=warehouse-service --timeout=60s
kubectl wait --namespace delivery-system --for=condition=ready pod --selector=app=delivery-service --timeout=60s
kubectl wait --namespace notification-system --for=condition=ready pod --selector=app=notification-service --timeout=60s
kubectl wait --namespace order-system --for=condition=ready pod --selector=app=order-service --timeout=60s
kubectl wait --namespace api-gateway-system --for=condition=ready pod --selector=app=api-gateway --timeout=60s
