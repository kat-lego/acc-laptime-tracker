#!/bin/bash

# set some environment variables
export ACC_COSMOS_CONNECTION_STRING=$(pass cosmos/acclaptracker)
export ACC_COSMOS_DATABASE=sessions
export ACC_COSMOS_CONTAINER=sessions
export ACC_CORS_ORIGINS="http://localhost:3000"
export PORT=8080

go run main.go
