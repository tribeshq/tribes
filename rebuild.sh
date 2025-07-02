#!/bin/bash

# Default values
VERSION=${1:-v0.1.0}
SNAPSHOT_NAME=${2:-default}

echo "🔄 Rebuilding Tribes Node..."
echo "📦 Version: $VERSION"
echo "📁 Snapshot: $SNAPSHOT_NAME"

# Remove container se estiver rodando
docker stop tribes-node 2>/dev/null || true
docker rm tribes-node 2>/dev/null || true

# Remove imagem antiga
docker rmi tribes-node 2>/dev/null || true

# Build nova imagem com argumentos
echo "📦 Building new image..."
docker build \
  --build-arg SNAPSHOT_URL=https://github.com/tribeshq/tribes/releases/download/${VERSION}/tribes-snapshot.tar.gz \
  --build-arg SNAPSHOT_NAME=${SNAPSHOT_NAME} \
  -f build/Dockerfile.node \
  -t tribes-node:${VERSION} \
  -t tribes-node:latest \
  .

if [ $? -eq 0 ]; then
    echo "✅ Build successful!"
    echo "🚀 Run with: docker run -it --rm tribes-node:${VERSION}"
    echo "📋 Available tags: tribes-node:${VERSION}, tribes-node:latest"
else
    echo "❌ Build failed!"
    exit 1
fi 