#! /usr/bin/env bash

set -e -o pipefail

CHARTS=$(find . -iname "Chart.yaml" | awk -F "/" '{print $(NF-1)}')

echo "[Charts to push]"
echo $CHARTS

echo "[Using repo ${PLUGIN_HELM_REPO}]"
helm repo add chartmuseum ${PLUGIN_HELM_REPO}

for CHART in ${CHARTS}; do
    echo "Pushing ${CHART}"
    # helm dependency update
    helm push --force ${CHART}/ chartmuseum
done

exit 0
