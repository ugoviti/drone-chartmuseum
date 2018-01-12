#! /usr/bin/env bash

CURL=$(which curl)
DIR=$(find . -iname "Chart.yaml" | awk -F "/" '{print $(NF-1)}')

helm init --client-only

for TARGET in ${DIR}
do  
    pushd ${TARGET}
    ls requirements.yaml > /dev/null 2>&1
    
    if [ $? -eq 0 ]; then
        helm repo add ${PLUGIN_HELM_REPO_NAME} ${PLUGIN_HELM_REPO_ADD}
        helm dependency update
    fi

    echo ""
    echo "Create package of ${TARGET} !"
    MSG=$(helm package .)
    PACKAGE=$(echo ${MSG} | awk -F "/" '{print $NF}')
    
    echo "Upload charts ${PACKAGE}..."
    RESULT="$("${CURL}" -s --data-binary "@${PACKAGE}" "${PLUGIN_UPLOAD_REPO_URL}/${PLUGIN_PATH}" | jq -r '.[]')"

    if [ "${RESULT}" == "true" ] || [ "${RESULT}" == "file already exists" ]; then
        echo "Upload ${PACKAGE} complete !"
    else 
        echo ${RESULT}
        exit 1
    fi
    
    popd
done

exit 0
