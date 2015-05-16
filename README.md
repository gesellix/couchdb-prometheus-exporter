# CouchDB Exporter
[CouchDB](http://couchdb.apache.org/) stats exporter for [Prometheus](http://prometheus.io/)

The CouchDB Exporter requests the CouchDB stats from the `/_stats` endpoint and 
exposes them for Prometheus consumption.

## Work in progress

*this is a work in progress: the current implementation only works as proof of concept
and only returns a gauge with the "couchdb_up" state of a CouchDB instance*

## Run it as container

```
docker run -p 9984:9984 gesellix/couchdb-exporter -logtostderr -couchdb.uri=http://192.168.59.103:5984
```
