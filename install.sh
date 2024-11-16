#!/bin/bash

# Fetch the latest release tag from GitHub API
LATEST_TAG=$(curl -s https://api.github.com/repos/frsfahd/go-proxy/releases/latest | grep '"tag_name":' | sed -E 's/.*"([^"]+)".*/\1/')

# Construct the URL for the latest release binary
BINARY_URL="https://github.com/frsfahd/go-proxy/releases/download/$LATEST_TAG/go-proxy"

# Download the binary from GitHub
sudo curl -L -o /usr/local/bin/go-proxy $BINARY_URL

# Make the binary executable
sudo chmod +x /usr/local/bin/go-proxy

echo "Installation complete!"
