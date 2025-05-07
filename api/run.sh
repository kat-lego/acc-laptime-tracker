#!/bin/bash

# set some environment variables
export ACCLTRCR_POSTGRES_CONNECTION_STRING=$(pass postgres/acc)
export PORT="8080"

go run main.go
