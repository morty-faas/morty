#!/bin/bash

set -Eeuo pipefail

GLOBAL_DEST="pkg/client"

GH_ORG="morty-faas"
GH_HOST="github.com"
GH_REPO="morty"

generate_client() {
    local client_name=$1
    local spec_file=$2
    local dest="${GLOBAL_DEST}/${client_name}"

    echo "Starting '$client_name' client generation for OpenAPI spec '$spec_file' into '$dest'"

    mkdir -p "${dest}"
    openapi-generator generate -i "${spec_file}" \
        -g go \
        -o "${dest}" \
        --git-user-id "${GH_ORG}" \
        --git-repo-id "${GH_REPO}" \
        --git-host "${GH_HOST}/${dest}" \
        -c ./tools/openapi-generator/config.yml

    # cp ./tools/openapi-generator/README.md "${DEST}"

    rm "${dest}/git_push.sh" || true
    rm "${dest}/.travis.yml" || true
    rm -rf "${dest}/test" || true
    rm "${dest}/go.mod" || true
    rm "${dest}/go.sum" || true
}

rm -rf "${GLOBAL_DEST}" || true

for spec in $(ls api/openapi-spec); do
    file=$(basename -- "$spec")
    filename="${file%.*}"
    generate_client $filename api/openapi-spec/$spec
done
