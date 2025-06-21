#!/bin/bash

DOCKERHUB_USER="99percentmoses"
IMAGE_NAME="acc-laptime-tracker-api"
FULL_IMAGE_NAME="${DOCKERHUB_USER}/${IMAGE_NAME}:latest"

GCP_PROJECT="acc-laptime-tracker-460418"
CLOUD_RUN_SERVICE="acc-laptime-tracker-api"
REGION="africa-south1"
PLATFORM="managed"

ACC_CORS_ORIGINS="https://acc-laptime-tracker-460418.web.app"
ACC_CORS_ORIGINS="$ACC_CORS_ORIGINS|https://acc.katlegomodupi.com"

ENV_VARS="ACC_FIREBASE_PROJECT_ID=acc-laptime-tracker-460418"
ENV_VARS="$ENV_VARS,ACC_FIREBASE_DATABASE=acclaptimetracker"
ENV_VARS="$ENV_VARS,ACC_FIREBASE_COLLECTION=session"
ENV_VARS="$ENV_VARS,ACC_CORS_ORIGINS=$ACC_CORS_ORIGINS"

echo $ENV_VARS

cp Dockerfile ../
cd ..

echo "📦 Building Docker image: $FULL_IMAGE_NAME"
docker build -t "$FULL_IMAGE_NAME" .

echo "🚀 Pushing Docker image to Docker Hub..."
docker push "$FULL_IMAGE_NAME"

echo "🎯 Setting GCP project: $GCP_PROJECT"
gcloud config set project "$GCP_PROJECT"

echo "🚀 Deploying $CLOUD_RUN_SERVICE to Cloud Run..."
gcloud run deploy "$CLOUD_RUN_SERVICE" \
  --image "docker.io/$FULL_IMAGE_NAME" \
  --platform "$PLATFORM" \
  --region "$REGION" \
  --allow-unauthenticated \
  --set-env-vars "$ENV_VARS" \
  --project "$GCP_PROJECT" \
  --max-instances 1 \

rm Dockerfile

