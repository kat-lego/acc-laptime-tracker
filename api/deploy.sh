#!/bin/bash

DOCKERHUB_USER="99percentmoses"
IMAGE_NAME="acc-laptime-tracker-api"
FULL_IMAGE_NAME="${DOCKERHUB_USER}/${IMAGE_NAME}:latest"

cp Dockerfile ../
cd ..

echo "Building Docker image..."
docker build -t "$FULL_IMAGE_NAME" .

echo "Deploying to Azure Container Apps using Docker Hub image..."

echo "Pushing image to Docker Hub..."
docker push "$FULL_IMAGE_NAME"

echo "Deploying image to Azure Container Apps..."
az containerapp up \
  --name capp-acc-laptime-tracker-api \
  --resource-group rg-acc-laptime-tracker \
  --location southafricanorth \
  --environment menv-acc-laptime-tracker \
  --image "$FULL_IMAGE_NAME"

rm Dockerfile
