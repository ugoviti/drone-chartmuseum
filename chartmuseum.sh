#! /usr/bin/env bash

set -e -o pipefail

echo "[Finding changes]"
CHARTS=$(git diff $DRONE_PREV_COMMIT_SHA $DRONE_COMMIT_SHA --name-only | awk -F '/' '{ print $1 }' | uniq)
echo "$CHARTS"

echo "[Using repo ${PLUGIN_HELM_REPO}]"
helm repo add chartmuseum ${PLUGIN_HELM_REPO}

echo "[Pushing charts]"
for CHART in ${CHARTS}; do

    if [ ! -d $CHART ]; then 
        echo "skip $CHART, not a directory"
        continue
    fi

    if [ ! -f $CHART/Chart.yaml ]; then
        echo "skip $CHART, no Chart.yaml"
        continue
    fi

    echo "Pushing ${CHART}"
    helm push ${CHART}/ chartmuseum
done
