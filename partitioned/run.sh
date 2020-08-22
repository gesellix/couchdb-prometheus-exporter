#!/usr/bin/env bash

# helps with debugging for https://github.com/gesellix/couchdb-prometheus-exporter/issues/73

set -euf -o pipefail

COUCHDB_USER="admin"
COUCHDB_PASSWORD="admin"

function finish {
  docker stop couchdb
}
trap finish EXIT

docker run --rm -d \
  --name couchdb \
  -p 5984:5984 \
  -v "$(pwd)"/data:/opt/couchdb/data \
  -e COUCHDB_USER=$COUCHDB_USER \
  -e COUCHDB_PASSWORD=$COUCHDB_PASSWORD \
  couchdb:3.1.0

sleep 5

curl --user $COUCHDB_USER:$COUCHDB_PASSWORD -X PUT http://localhost:5984/_users http://localhost:5984/_replicator

curl --user $COUCHDB_USER:$COUCHDB_PASSWORD -X PUT http://localhost:5984/foo?partitioned=false
curl --user $COUCHDB_USER:$COUCHDB_PASSWORD -X POST http://localhost:5984/foo -H 'Content-Type:application/json' -d @view-non-partitioned.json

curl --user $COUCHDB_USER:$COUCHDB_PASSWORD -X PUT http://localhost:5984/bar?partitioned=true
curl --user $COUCHDB_USER:$COUCHDB_PASSWORD -X POST http://localhost:5984/bar -H 'Content-Type:application/json' -d @view-partitioned.json

docker logs -f couchdb
