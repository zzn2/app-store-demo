#!/bin/sh

# Run in the root folder in this repo, invoke like:
# ./deploy.sh

REGISTRY=zzn2
COMMIT=$(git rev-parse HEAD | cut -c 1-8)

# Build docker image
IMAGE_NAME=$REGISTRY/app-store-demo:$COMMIT
docker build . -t $IMAGE_NAME
docker push $IMAGE_NAME

# Deploy with Helm
helm upgrade --install app-store ./chart/app-store/ --set image.registry=$REGISTRY,image.tag=$COMMIT
