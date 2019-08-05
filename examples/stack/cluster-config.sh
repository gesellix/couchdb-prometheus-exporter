#!/usr/bin/env sh

set -o errexit
set -o nounset
#set -o xtrace

COUCHDB_ADDRESS="${1:-http://couchdb:5984}"
CONFIG_PROPERTIES="${2:-config.properties}"

# concept taken from https://stackoverflow.com/a/38096496
prop() {
    grep "${1}" "${CONFIG_PROPERTIES}" | cut -d'=' -f2
}

curl -v -X POST "${COUCHDB_ADDRESS}/_cluster_setup" \
     -H "Content-Type: application/json" \
     -d "{\"action\":\"enable_cluster\", \"username\":\"$(prop couchdb.username)\", \"password\":\"$(prop couchdb.password)\", \"bind_address\":\"0.0.0.0\", \"port\":5984, \"node_count\": 1}"
curl -v -X POST "${COUCHDB_ADDRESS}/_session" \
     -H "Content-Type: application/x-www-form-urlencoded; charset=UTF-8" \
     -d "name=$(prop couchdb.username)&password=$(prop couchdb.password)" \
     --cookie-jar couchdb-cookies.txt
curl -v -X POST "${COUCHDB_ADDRESS}/_cluster_setup" \
     -H "Content-Type: application/json" \
     -d '{"action":"finish_cluster"}' \
     --cookie couchdb-cookies.txt
