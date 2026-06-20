#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

mkdir -p download/descriptions
rm -f download/descriptions/*

docker run -it --entrypoint /bin/sh \
    --rm --name mm-downloader \
    -v "$(pwd)/download:/home/runner/download" \
    media-manager:dev

# Example run:
# mm download --from-file=download.txt --yt-format="bestaudio[ext=m4a]"
# mm download --from-file=download.txt --yt-args="--list-formats"