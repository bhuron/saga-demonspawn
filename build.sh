#!/bin/bash
set -e

echo "Cleaning old binary..."
rm -f saga

echo "Building saga..."
go build -o saga ./cmd/saga

echo "Verifying binary exists..."
ls -lh saga

echo ""
echo "✅ Build complete! Run with: ./saga"
