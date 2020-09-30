#!/bin/bash

# This script is intended to stress test prometheus.
# It will create specified number of namespaces and deploy specified number
# of eagle replicas in those namespaces.
# Example usage:
#
#   ./cardinality.sh ../manifests/deploy.yaml 50 20
# 
# This will create 50 namespaces and deploy 20 replicas of eagle application
# in each namespace. Effectively creating:
#   - 50 namespaces
#   - 50 deployments and services
#   - 50 service monitors
#   - 1000 eagle pods

BASE="${1:-'../manifests/deploy.yaml'}"
TOTAL_NAMESPACES="${2:-1}"
REPLICAS_PER_NAMESPACE="${3:-1}"

for i in $(seq 1 $TOTAL_NAMESPACES); do
    NAMESPACE="test-${i}"
    kubectl create namespace ${NAMESPACE}
    kubectl apply -n ${NAMESPACE} -f "${BASE}"
done
