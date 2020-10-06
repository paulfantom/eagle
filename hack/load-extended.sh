#!/bin/bash

BASE="${1:-'manifests/extended.yaml'}"
TOTAL_NAMESPACES="${2:-1}"
REPLICAS_PER_NAMESPACE="${3:-1}"
TMP=$(mktemp -d)

trap "rm -rf $TMP" EXIT

for i in $(seq 1 $REPLICAS_PER_NAMESPACE); do
	cat "${BASE}" | sed "s/{{replica}}/${i}/g" > "${TMP}/${i}.yaml"
done	

for i in $(seq 1 "$TOTAL_NAMESPACES"); do
	NAMESPACE="test-${i}"
	kubectl create namespace "${NAMESPACE}"
	kubectl apply -n "${NAMESPACE}" -f "${TMP}/"
done
