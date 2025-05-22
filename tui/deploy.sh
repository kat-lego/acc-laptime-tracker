#!/bin/bash

DOCKERHUB_USER="99percentmoses"
IMAGE_NAME="acc-laptime-tracker-tui"
FULL_IMAGE_NAME="${DOCKERHUB_USER}/${IMAGE_NAME}:latest"

GCP_PROJECT="acc-laptime-tracker-460418"
CLOUD_RUN_SERVICE="acc-laptime-tracker-tui"
REGION="africa-south1"
PLATFORM="managed"

ENV_VARS="SSH_HOST_KEY_PEM=$(cat ~/.ssh/id_acc_laptime_tracker_tui)"

cp Dockerfile ../
cd ..

echo "📦 Building Docker image: $FULL_IMAGE_NAME"
docker build -t "$FULL_IMAGE_NAME" .

echo "🚀 Pushing Docker image to Docker Hub..."
docker push "$FULL_IMAGE_NAME"

echo "🚀 Deploying $CLOUD_RUN_SERVICE to Cloud Run..."
gcloud run deploy "$CLOUD_RUN_SERVICE" \
  --image "docker.io/$FULL_IMAGE_NAME" \
  --platform "$PLATFORM" \
  --region "$REGION" \
  --allow-unauthenticated \
  --set-env-vars "$ENV_VARS" \
  --project "$GCP_PROJECT"

rm Dockerfile

