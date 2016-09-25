# CouchDB Exporter

[![Build Status](https://travis-ci.org/gesellix/couchdb-exporter.svg?branch=master)](https://travis-ci.org/gesellix/couchdb-exporter)

[CouchDB](http://couchdb.apache.org/) stats exporter for [Prometheus](http://prometheus.io/)

The CouchDB Exporter requests the CouchDB stats from the `/_stats` endpoint and 
exposes them for Prometheus consumption.

## Run it as container

```
docker run -p 9984:9984 gesellix/couchdb-exporter -couchdb.uri=http://couchdb:5984
```

The couchdb-exporter uses the [glog](https://godoc.org/github.com/golang/glog) library for logging.
With the default parameters nothing will be logged.
Use `-logtostderr` to enable logging to stderr and `--help` to see all options.

For CouchDB 2.x, you should configure the exporter to fetch the stats on port 15984, where you'll
get a complete cluster overview. In contrast to CouchDB 1.x you'll need to configure the admin
credentials, e.g. like this:

```
docker run -p 9984:9984 gesellix/couchdb-exporter -couchdb.uri=http://couchdb:5984 -couchdb.username=root -couchdb.password=a-secret
```

## Metrics Overview
The following list gives you an overview on the currently exposed metrics.
Please note that beyond the complete CouchDB stats a metric `couchdb_up` has been
added as minimal connection health check.

The `node_name` label defaults to `master`, but makes more sense in a clustered CouchDB 2.x environment.
Then, you'll get separated stats for each node. 

```
# HELP couchdb_auth_cache_hits number of authentication cache hits
# TYPE couchdb_auth_cache_hits gauge
couchdb_auth_cache_hits 0
# HELP couchdb_auth_cache_misses number of authentication cache misses
# TYPE couchdb_auth_cache_misses gauge
couchdb_auth_cache_misses 0
# HELP couchdb_database_reads number of times a document was read from a database
# TYPE couchdb_database_reads gauge
couchdb_database_reads 3
# HELP couchdb_database_writes number of times a database was changed
# TYPE couchdb_database_writes gauge
couchdb_database_writes 2
# HELP couchdb_httpd_bulk_requests number of bulk requests
# TYPE couchdb_httpd_bulk_requests gauge
couchdb_httpd_bulk_requests 0
# HELP couchdb_httpd_clients_requesting_changes number of clients for continuous _changes
# TYPE couchdb_httpd_clients_requesting_changes gauge
couchdb_httpd_clients_requesting_changes 0
# HELP couchdb_httpd_request_methods number of HTTP requests by method
# TYPE couchdb_httpd_request_methods gauge
couchdb_httpd_request_methods{method="COPY"} 0
couchdb_httpd_request_methods{method="DELETE"} 0
couchdb_httpd_request_methods{method="GET"} 136
couchdb_httpd_request_methods{method="HEAD"} 0
couchdb_httpd_request_methods{method="POST"} 0
couchdb_httpd_request_methods{method="PUT"} 3
# HELP couchdb_httpd_requests number of HTTP requests
# TYPE couchdb_httpd_requests gauge
couchdb_httpd_requests 139
# HELP couchdb_httpd_status_codes number of HTTP responses by status code
# TYPE couchdb_httpd_status_codes gauge
couchdb_httpd_status_codes{code="200"} 68
couchdb_httpd_status_codes{code="201"} 3
couchdb_httpd_status_codes{code="202"} 0
couchdb_httpd_status_codes{code="301"} 0
couchdb_httpd_status_codes{code="304"} 0
couchdb_httpd_status_codes{code="400"} 0
couchdb_httpd_status_codes{code="401"} 0
couchdb_httpd_status_codes{code="403"} 0
couchdb_httpd_status_codes{code="404"} 0
couchdb_httpd_status_codes{code="405"} 0
couchdb_httpd_status_codes{code="409"} 0
couchdb_httpd_status_codes{code="412"} 0
couchdb_httpd_status_codes{code="500"} 0
# HELP couchdb_httpd_temporary_view_reads number of temporary view reads
# TYPE couchdb_httpd_temporary_view_reads gauge
couchdb_httpd_temporary_view_reads 0
# HELP couchdb_httpd_view_reads number of view reads
# TYPE couchdb_httpd_view_reads gauge
couchdb_httpd_view_reads 0
# HELP couchdb_open_databases number of open databases
# TYPE couchdb_open_databases gauge
couchdb_open_databases 1
# HELP couchdb_open_os_files number of file descriptors CouchDB has open
# TYPE couchdb_open_os_files gauge
couchdb_open_os_files 2
# HELP couchdb_request_time length of a request inside CouchDB without MochiWeb
# TYPE couchdb_request_time gauge
couchdb_request_time 121.732
# HELP couchdb_up Was the last query of CouchDB stats successful.
# TYPE couchdb_up gauge
couchdb_up 1
```

## Thanks

Thanks go to the Prometheus team, which is very active and responsive!

I also have to admit that the couchdb-exporter code is heavily inspired by 
the other [available exporters](http://prometheus.io/docs/instrumenting/exporters/), 
and that some ideas have just been copied from them.
