#!/bin/bash

VERSION="1.0.0"
OS=$(uname -s)
ARCH=$(uname -m)

if [ "$OS" == "Linux" ]; then
  if [ "$ARCH" == "x86_64" ]; then
    URL="https://github.com/yourusername/stargate/releases/download/v${VERSION}/stargate-linux-amd64"
  else
    echo "Unsupported architecture: $ARCH"
    exit 1
  fi
elif [ "$OS" == "Darwin" ]; then
  if [ "$ARCH" == "x86_64" ]; then
    URL="https://github.com/yourusername/stargate/releases/download/v${VERSION}/stargate-darwin-amd64"
  else
    echo "Unsupported architecture: $ARCH"
    exit 1
  fi
else
  echo "Unsupported OS: $OS"
  exit 1
fi

echo "Downloading Stargate CLI..."
curl -LO $URL

echo "Making Stargate CLI executable..."
chmod +x stargate-*

echo "Moving Stargate CLI to /usr/local/bin..."
sudo mv stargate-* /usr/local/bin/stargate

echo "Stargate CLI installed successfully!"
