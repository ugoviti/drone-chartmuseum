#! /usr/bin/env bash

set -e -o pipefail

echo "[Charts]"
# This will not work correctly if a build was restarted or a build has failed, was skipped, etc."
# See https://discourse.drone.io/t/how-to-get-the-commit-range-of-a-push-pr-event/1716/7
# CHARTS=$(git diff $DRONE_PREV_COMMIT_SHA $DRONE_COMMIT_SHA --name-only | awk -F '/' '{ print $1 }' | uniq)
CHARTS=$(find . -iname "Chart.yaml" | awk -F "/" '{print $(NF-1)}')
for CHART in $CHARTS; do echo $CHART; done

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
