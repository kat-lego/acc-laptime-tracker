#!/bin/bash

export GOOGLE_APPLICATION_CREDENTIALS="/mnt/c/Users/katlego/acc-laptime-tracker/acc-laptime-tracker-460418-e1d56ffa02e5.json"
export ACC_FIREBASE_PROJECT_ID=acc-laptime-tracker-460418
export ACC_FIREBASE_DATABASE=acclaptimetracker
export ACC_FIREBASE_COLLECTION=session
export ACC_CORS_ORIGINS="http://localhost:8081"
export GIN_MODE=release
export PORT=8080

go run main.go
