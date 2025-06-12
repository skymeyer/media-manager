#!/usr/bin/env bash

set -o errexit
set -o nounset
set -o pipefail

export DOCKER_DEFAULT_PLATFORM=linux/amd64

docker build -t media-manager:dev .
