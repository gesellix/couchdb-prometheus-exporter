#!/usr/bin/env bash

set -o nounset
set -o errexit
#set -o xtrace

docker-compose -f examples/compose/docker-compose-integrationtest.yml -p db down
