#!/bin/sh

# Run in the root folder in this repo, invoke like:
# ./scripts/build.sh

REGISTRY=zzn2
COMMIT=$(git rev-parse HEAD | cut -c 1-8)

# Build docker image
IMAGE_NAME=$REGISTRY/app-store-demo:$COMMIT
docker build . -t $IMAGE_NAME
docker push $IMAGE_NAME

# Replace the values in template using sed for a while.
# TODO: Learn how to deploy with Helm.
cp k8s/deployment.yaml .
sed -i -e "s/{{.Values.registry}}/$REGISTRY/g" deployment.yaml
sed -i -e "s/{{.Values.dockerTag}}/$COMMIT/g" deployment.yaml

# Deploy to kubernetes cluster
kubectl apply -f deployment.yaml
rm deployment.yaml