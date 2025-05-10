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

# az containerapp update \
#   --name capp-acc-laptime-tracker-api \
#   --resource-group rg-acc-laptime-tracker \
#   --set-env-vars ACC_COSMOS_CONNECTION_STRING=sessions

# az containerapp env create -n menv-acc-laptime-tracker -g rg-acc-laptime-tracker \
#             --logs-workspace-id e3834547-305e-40b9-8585-5fe92b4bba1d \
#             --location southafricanorth

rm Dockerfile



