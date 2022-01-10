#!/bin/sh

# Run in the root folder in this repo, invoke like:
# ./scripts/build.sh

IMAGE_NAME=zzn2/app-store-demo

docker build . -t $IMAGE_NAME
docker push $IMAGE_NAME

kubectl apply -f k8s/deployment.yaml