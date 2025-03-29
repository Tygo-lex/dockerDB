#!/bin/bash

# Colors for prettier output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Absolute path to the dockerdb binary
# This gets the directory where the script is located
SCRIPT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )" &> /dev/null && pwd )"
BINARY_PATH="$SCRIPT_DIR/bin/dockerdb"

# Check if the binary exists
if [ ! -f "$BINARY_PATH" ]; then
    echo -e "${RED}Error: dockerdb binary not found at $BINARY_PATH${NC}"
    echo -e "${YELLOW}Building dockerdb...${NC}"
    
    # Build the binary
    cd "$SCRIPT_DIR"
    go build -o bin/dockerdb
    
    if [ $? -ne 0 ]; then
        echo -e "${RED}Failed to build dockerdb. Please build it manually with 'go build -o bin/dockerdb'${NC}"
        exit 1
    fi
    
    echo -e "${GREEN}Successfully built dockerdb binary${NC}"
fi

# Target location in PATH
TARGET_PATH="/usr/local/bin/dockerdb"

# Check if the symlink already exists
if [ -L "$TARGET_PATH" ]; then
    echo -e "${YELLOW}dockerdb is already installed in PATH. Reinstalling...${NC}"
    sudo rm "$TARGET_PATH"
elif [ -f "$TARGET_PATH" ]; then
    echo -e "${YELLOW}A file already exists at $TARGET_PATH. Replacing it...${NC}"
    sudo rm "$TARGET_PATH"
fi

# Create the symlink
echo -e "Creating symlink from $BINARY_PATH to $TARGET_PATH"
sudo ln -s "$BINARY_PATH" "$TARGET_PATH"

if [ $? -eq 0 ]; then
    echo -e "${GREEN}dockerdb successfully installed to $TARGET_PATH${NC}"
    echo -e "${GREEN}You can now run 'dockerdb' from anywhere.${NC}"
else
    echo -e "${RED}Failed to create symlink. Try running the following command manually:${NC}"
    echo -e "sudo ln -s $BINARY_PATH $TARGET_PATH"
    exit 1
fi

# Check if dockerdb is in PATH
which dockerdb > /dev/null 2>&1
if [ $? -eq 0 ]; then
    echo -e "${GREEN}Verified: dockerdb is in your PATH${NC}"
    dockerdb --version || echo -e "${YELLOW}Note: dockerdb command found but returned non-zero exit code when checking version${NC}"
else
    echo -e "${YELLOW}Warning: dockerdb was installed but isn't in your PATH.${NC}"
    echo -e "You may need to restart your terminal or add /usr/local/bin to your PATH."
fi

exit 0