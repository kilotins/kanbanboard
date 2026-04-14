#!/usr/bin/env bash
set -euo pipefail

VERSION="1.2.2"
APP_IMAGE="kanbanboard:${VERSION}"
IMAGE_TAR="kanbanboard-image.tar"
PACKAGE_NAME="kanbanboard-${VERSION}.tar.gz"

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
ROOT_DIR="$(dirname "${SCRIPT_DIR}")"

cd "${ROOT_DIR}"

echo "Building app image ${APP_IMAGE}..."
docker build \
  --build-arg VERSION="${VERSION}" \
  -t "${APP_IMAGE}" \
  -f backend/Dockerfile \
  .

echo "Saving image to ${IMAGE_TAR}..."
docker save "${APP_IMAGE}" -o "${IMAGE_TAR}"

echo "Creating package ${PACKAGE_NAME}..."
tar czf "${PACKAGE_NAME}" \
  "${IMAGE_TAR}" \
  docker-compose.deploy.yml \
  deploy.sh

rm "${IMAGE_TAR}"

echo "Done. Deploy package: ${PACKAGE_NAME}"
echo "Transfer it to the target machine and run: tar xzf ${PACKAGE_NAME} && bash deploy.sh"
