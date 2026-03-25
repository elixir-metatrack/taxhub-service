#!/usr/bin/env bash

set -euo pipefail

echo "Starting minikube..."
minikube start

echo "Building image..."
eval $(minikube docker-env)
docker build -t taxhub:local ..

echo "Applying manifests..."
kubectl apply -f taxhub.yaml

echo "Waiting for pods to be ready..."
kubectl rollout status deployment/taxhub-api

echo "Forwarding http://localhost:8080 -> taxhub-api:8080 (Ctrl+C to stop)"
kubectl port-forward service/taxhub-api 8080:8080
