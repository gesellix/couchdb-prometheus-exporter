# CouchDB Exporter
[CouchDB](http://couchdb.apache.org/) stats exporter for [Prometheus](http://prometheus.io/)

The CouchDB Exporter requests the CouchDB stats from the `/_stats` endpoint and 
exposes them for Prometheus consumption.

## Run it as container

```
docker run -p 9984:9984 gesellix/couchdb-exporter -logtostderr -couchdb.uri=http://192.168.59.103:5984
```
