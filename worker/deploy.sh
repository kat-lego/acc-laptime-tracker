#!/bin/bash

set -e

YELLOW="\033[1;33m"
GREEN="\033[1;32m"
BLUE="\033[1;34m"
CYAN="\033[1;36m"
RED="\033[1;31m"
RESET="\033[0m"

INFO="🛈"
OK="✅"
RUN="🚀"
COPY="📁"
WARN="⚠️"
BUILD="🔧"

echo -e "${BLUE}${INFO} Gathering Windows environment paths...${RESET}"
RAW_WINDOWS_TEMP=$(powershell.exe -Command "[System.Environment]::GetEnvironmentVariable('TEMP', 'User')")
RAW_WIN_HOME=$(powershell.exe -Command "[System.Environment]::GetEnvironmentVariable('USERPROFILE')")

WINDOWS_TEMP_WIN=$(echo "$RAW_WINDOWS_TEMP" | tr -d '\r' | sed 's/\\$//')
WINDOWS_TEMP_WSL=$(echo "$WINDOWS_TEMP_WIN" | sed 's/\\/\//g' | sed -E 's/^([A-Z]):/\/mnt\/\L\1/')

WIN_HOME_WIN=$(echo "$RAW_WIN_HOME" | tr -d '\r' | sed 's/\\$//')
WIN_HOME_WSL=$(echo "$WIN_HOME_WIN" | sed 's/\\/\//g' | sed -E 's/^([A-Z]):/\/mnt\/\L\1/')

TMP_DIR_WSL="$WINDOWS_TEMP_WSL/acc-laptime-tracker"
TMP_DIR_WIN="$WINDOWS_TEMP_WIN\\acc-laptime-tracker"

echo -e "${COPY} Copying module to ${GREEN}$TMP_DIR_WSL${RESET}"
mkdir -p "$TMP_DIR_WSL"
rsync -av --progress ../* "$TMP_DIR_WSL" --exclude web

BUILD_OUT_DIR="$WIN_HOME_WIN\\acc-laptime-tracker"
BUILD_OUT_DIR_WSL="$WIN_HOME_WSL/acc-laptime-tracker"
mkdir -p "$BUILD_OUT_DIR_WSL"

echo -e "${BUILD} Building worker to ${GREEN}$BUILD_OUT_DIR\\acc-laptime-tracker.exe${RESET}"
powershell.exe -Command "cd '$TMP_DIR_WIN\\worker'; go build -o '$BUILD_OUT_DIR\\acc-laptime-tracker.exe'"

