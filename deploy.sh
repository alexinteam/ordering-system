#!/bin/bash

set -e

kubectl apply -f k8s/namespaces.yaml

kubectl apply -f k8s/rabbitmq.yaml
kubectl apply -f k8s/rabbitmq-ingress.yaml
kubectl apply -f k8s/rabbitmq-static-ingress.yaml

echo "Развертывание сервисов..."
kubectl apply -f k8s/order-service.yaml
kubectl apply -f k8s/billing-service.yaml
kubectl apply -f k8s/notification-service.yaml
kubectl apply -f k8s/api-gateway.yaml

echo "Ожидание готовности сервисов..."
kubectl wait --namespace order-system --for=condition=ready pod --selector=app=order-service --timeout=30s
kubectl wait --namespace billing-system --for=condition=ready pod --selector=app=billing-service --timeout=30s
kubectl wait --namespace notification-system --for=condition=ready pod --selector=app=notification-service --timeout=30s
kubectl wait --namespace api-gateway-system --for=condition=ready pod --selector=app=api-gateway --timeout=30s