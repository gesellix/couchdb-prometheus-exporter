#!/usr/bin/env bash

cd examples/compose

cat << EOB > .couchdb-env
COUCHDB.USERNAME=root
COUCHDB.PASSWORD=a-secret
EOB
docker-compose -d -f docker-compose-cluster.yml -p db up

curl -X GET "http://localhost:9984/metrics" -v

docker-compose -f docker-compose-cluster.yml -p db down
