#!/usr/bin/env bash

cd examples/compose

cat << EOB > .couchdb-env
COUCHDB.USERNAME=root
COUCHDB.PASSWORD=a-secret
EOB

# We want to use overlay networks to make them `attachable`.
# Overlay networks are only available in swarm mode, though.
# See https://docs.docker.com/compose/compose-file/compose-file-v2/#network-configuration-reference
docker swarm init
docker-compose -f docker-compose-cluster.yml -p db up -d

curl -X GET "http://localhost:9984/metrics" -v

docker-compose -f docker-compose-cluster.yml -p db down
