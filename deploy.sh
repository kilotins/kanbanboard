#!/usr/bin/env bash
set -euo pipefail

echo "Loading kanbanboard image..."
docker load -i kanbanboard-image.tar

echo "Starting services..."
docker compose -f docker-compose.deploy.yml up -d

echo "Done. App available at http://localhost:8080"
