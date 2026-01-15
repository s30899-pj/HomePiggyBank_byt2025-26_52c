#!/bin/sh
set -e

cd /app

echo "Starting in production mode..."
make build-for-docker
./bin/HomePiggyBank