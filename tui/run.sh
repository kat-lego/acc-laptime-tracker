#!/bin/bash

# set some environment variables
export SSH_HOST_KEY_PEM=$(cat ~/.ssh/id_acc_laptime_tracker_tui )

go run main.go
