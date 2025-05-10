#!/bin/bash

# set some environment variables
export ACC_COSMOS_CONNECTION_STRING=$(pass cosmos/acclaptracker)
export ACC_COSMOS_DATABASE=sessions
export ACC_COSMOS_CONTAINER=sessions
export PORT=8080

go run main.go
