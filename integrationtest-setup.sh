#!/usr/bin/env bash

set -o nounset
set -o errexit
#set -o xtrace

# We want to use overlay networks to make them `attachable`.
# Overlay networks are only available in swarm mode, though.
# See https://docs.docker.com/compose/compose-file/compose-file-v2/#network-configuration-reference
if [ "$(docker info --format '{{ .Swarm.LocalNodeState }}')" == "inactive" ];then
  echo "Initialize swarm manager."
  docker swarm init
fi
docker-compose -f examples/compose/docker-compose-integrationtest.yml -p db up -d

curl -X GET "http://localhost:9984/metrics" -v
