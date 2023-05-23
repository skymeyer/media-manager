#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p download/descriptions

docker run --rm --name mm-downloader \
    -v "$(pwd)/download:/home/runner/download" \
    media-manager:dev \
    download --from-file=download.txt

mv download/*.description download/descriptions || true
