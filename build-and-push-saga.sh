#!/bin/bash

set -e

echo "Сборка и пуш Docker образов для Saga Order System..."

# Сборка и пуш Billing Service
echo "Сборка Billing Service..."
cd billing-service
docker build -t alexinteam/billing-service:latest .
docker push alexinteam/billing-service:latest
cd ..

# Сборка и пуш Warehouse Service
echo "Сборка Warehouse Service..."
cd warehouse-service
docker build -t alexinteam/warehouse-service:latest .
docker push alexinteam/warehouse-service:latest
cd ..

# Сборка и пуш Delivery Service
echo "Сборка Delivery Service..."
cd delivery-service
docker build -t alexinteam/delivery-service:latest .
docker push alexinteam/delivery-service:latest
cd ..

# Сборка и пуш Order Service (обновленный)
echo "Сборка Order Service..."
cd order-service
docker build -t alexinteam/order-service:latest .
docker push alexinteam/order-service:latest
cd ..

# Сборка и пуш Notification Service (обновленный)
echo "Сборка Notification Service..."
cd notification-service
docker build -t alexinteam/notification-service:latest .
docker push alexinteam/notification-service:latest
cd ..

# Сборка и пуш API Gateway (обновленный)
echo "Сборка API Gateway..."
cd api-gateway
docker build -t alexinteam/api-gateway:latest .
docker push alexinteam/api-gateway:latest
cd ..

echo "Все образы успешно собраны и запушены!"
echo ""
echo "Доступные образы:"
echo "- alexinteam/billing-service:latest"
echo "- alexinteam/warehouse-service:latest"
echo "- alexinteam/delivery-service:latest"
echo "- alexinteam/order-service:latest"
echo "- alexinteam/notification-service:latest"
echo "- alexinteam/api-gateway:latest"
