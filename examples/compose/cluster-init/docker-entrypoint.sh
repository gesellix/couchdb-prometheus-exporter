#!/usr/bin/env sh

#cat << EOB > .couchdb-env
#COUCHDB.USERNAME=admin
#COUCHDB.PASSWORD=password
#COUCHDB_USER=admin
#COUCHDB_PASSWORD=password
#EOB
#docker-compose -f docker-compose-cluster.yml -p db up

/await-nodes.sh
/init-cluster.sh
