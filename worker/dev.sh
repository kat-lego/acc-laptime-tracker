#!/bin/bash

set -e

# Colors
YELLOW="\033[1;33m"
GREEN="\033[1;32m"
BLUE="\033[1;34m"
CYAN="\033[1;36m"
RED="\033[1;31m"
RESET="\033[0m"

# Emoji prefixes
INFO="🛈"
OK="✅"
RUN="🚀"
COPY="📁"
WARN="⚠️"
BUILD="🔧"

# Step 0: Validate command
CMD="$1"
if [ -n "$CMD" ] && [ "$CMD" != "run" ] && [ "$CMD" != "build" ]; then
  echo -e "${RED}${WARN} Invalid command: $CMD${RESET}"
  echo -e "${YELLOW}Usage: $0 [run|build]${RESET}"
  echo -e "  ${CYAN}run${RESET}   - Run 'go run .' in worker/ using PowerShell"
  echo -e "  ${CYAN}build${RESET} - Build worker to Windows home directory and copy setup.ps1"
  exit 1
fi

# Step 1: Get raw Windows TEMP and HOME paths
echo -e "${BLUE}${INFO} Gathering Windows environment paths...${RESET}"
RAW_WINDOWS_TEMP=$(powershell.exe -Command "[System.Environment]::GetEnvironmentVariable('TEMP', 'User')")
RAW_WIN_HOME=$(powershell.exe -Command "[System.Environment]::GetEnvironmentVariable('USERPROFILE')")

WINDOWS_TEMP_WIN=$(echo "$RAW_WINDOWS_TEMP" | tr -d '\r' | sed 's/\\$//')
WINDOWS_TEMP_WSL=$(echo "$WINDOWS_TEMP_WIN" | sed 's/\\/\//g' | sed -E 's/^([A-Z]):/\/mnt\/\L\1/')

WIN_HOME_WIN=$(echo "$RAW_WIN_HOME" | tr -d '\r' | sed 's/\\$//')
WIN_HOME_WSL=$(echo "$WIN_HOME_WIN" | sed 's/\\/\//g' | sed -E 's/^([A-Z]):/\/mnt\/\L\1/')

# Step 2: Locate go.mod
echo -e "${BLUE}${INFO} Locating go.mod...${RESET}"
CURRENT_DIR="$(pwd)"
while [ "$CURRENT_DIR" != "/" ]; do
  if [ -f "$CURRENT_DIR/go.mod" ]; then
    SOURCE_MODULE="$CURRENT_DIR"
    echo -e "${OK} Found go.mod in ${GREEN}$SOURCE_MODULE${RESET}"
    break
  fi
  CURRENT_DIR="$(dirname "$CURRENT_DIR")"
done

if [ -z "$SOURCE_MODULE" ]; then
  echo -e "${RED}${WARN} Error: Could not find go.mod in parent directories.${RESET}"
  exit 1
fi

# Step 3: Copy module to temp
DEST_DIR_WSL="$WINDOWS_TEMP_WSL/acc-laptime-tracker"
DEST_DIR_WIN="$WINDOWS_TEMP_WIN\\acc-laptime-tracker"
WORKER_SUBDIR="worker"

echo -e "${COPY} Copying module to ${GREEN}$DEST_DIR_WSL${RESET}"
mkdir -p "$DEST_DIR_WSL"
cp -r "$SOURCE_MODULE"/* "$DEST_DIR_WSL"

# Step 4: Handle commands
case "$CMD" in
  run)
    echo -e "${RUN} Running 'go run .' in PowerShell inside ${CYAN}worker/${RESET}"
    powershell.exe -Command "cd '$DEST_DIR_WIN\\$WORKER_SUBDIR'; go run ."
    ;;
  build)
    BUILD_OUT_DIR="$WIN_HOME_WIN\\acc-laptime-tracker"
    BUILD_OUT_DIR_WSL="$WIN_HOME_WSL/acc-laptime-tracker"
    mkdir -p "$BUILD_OUT_DIR_WSL"

    echo -e "${BUILD} Building worker to ${GREEN}$BUILD_OUT_DIR\\acc-laptime-tracker.exe${RESET}"
    powershell.exe -Command "cd '$DEST_DIR_WIN\\$WORKER_SUBDIR'; go build -o '$BUILD_OUT_DIR\\acc-laptime-tracker.exe'"

    ;;
  *)
    echo -e "${OK} Copy complete. No command executed.${RESET}"
    ;;
esac

