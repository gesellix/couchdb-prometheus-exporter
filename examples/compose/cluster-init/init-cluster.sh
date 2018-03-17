#!/usr/bin/env sh

# http://docs.couchdb.org/en/stable/cluster/setup.html
# http://docs.couchdb.org/en/latest/cluster/setup.html

curl -X PUT http://172.16.238.11:5984/_users
curl -X PUT http://172.16.238.11:5984/_replicator
curl -X PUT http://172.16.238.12:5984/_users
curl -X PUT http://172.16.238.12:5984/_replicator
curl -X PUT http://172.16.238.13:5984/_users
curl -X PUT http://172.16.238.13:5984/_replicator

#USERNAME=$(awk 'BEGIN {print ENVIRON["COUCHDB.USERNAME"]}')
#PASSWORD=$(awk 'BEGIN {print ENVIRON["COUCHDB.PASSWORD"]}')

curl -X PUT http://172.16.238.11:5984/_node/couchdb@172.16.238.11/_config/admins/root -d '"a-secret"'
curl -X PUT http://172.16.238.12:5984/_node/couchdb@172.16.238.12/_config/admins/root -d '"a-secret"'
curl -X PUT http://172.16.238.13:5984/_node/couchdb@172.16.238.13/_config/admins/root -d '"a-secret"'

#curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "content-type:application/json" -d '{"action":"enable_cluster","username":"root","password":"a-secret","bind_address":"0.0.0.0","node_count":3}'

curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "content-type:application/json" -d '{"remote_node":"172.16.238.12","port":"5984","action":"enable_cluster","username":"root","password":"a-secret","bind_address":"0.0.0.0","node_count":3}'
curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "content-type:application/json" -d '{"remote_node":"172.16.238.13","port":"5984","action":"enable_cluster","username":"root","password":"a-secret","bind_address":"0.0.0.0","node_count":3}'

curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "content-type:application/json" -d '{"action":"add_node","host":"172.16.238.12","port":"5984","username":"root","password":"a-secret"}'
curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "content-type:application/json" -d '{"action":"add_node","host":"172.16.238.13","port":"5984","username":"root","password":"a-secret"}'

curl --user root:a-secret -X POST http://172.16.238.11:5984/_cluster_setup -H "Content-Type: application/json" -d '{"action": "finish_cluster"}'

curl --user root:a-secret -X GET http://172.16.238.11:5984/_membership
