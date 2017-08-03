# CouchDB Prometheus Exporter

[![Build Status](https://travis-ci.org/gesellix/couchdb-prometheus-exporter.svg?branch=master)](https://travis-ci.org/gesellix/couchdb-prometheus-exporter)

[CouchDB](http://couchdb.apache.org/) metrics exporter for [Prometheus](http://prometheus.io/)

The CouchDB metrics exporter requests the CouchDB stats from the `/_stats` and `/_active_tasks` endpoints and 
exposes them for Prometheus consumption. You can optionally monitor detailed database stats like
disk and data size to monitor the storage overhead. The exporter can be configured via program parameters,
environment variables, and config file.

## Run it as container

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter -couchdb.uri=http://couchdb:5984

The couchdb-exporter uses the [glog](https://godoc.org/github.com/golang/glog) library for logging.
With the default parameters nothing will be logged.
Use `-logtostderr` to enable logging to stderr and `--help` to see all options.

For CouchDB 2.x, you should configure the exporter to fetch the stats from one node, to get
a complete cluster overview. In contrast to CouchDB 1.x you'll need to configure the admin
credentials, e.g. like this:

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter -couchdb.uri=http://couchdb:5984 -couchdb.username=root -couchdb.password=a-secret

If you need database disk usage stats, simply add a comma separated list of database names like this:

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter -couchdb.uri=http://couchdb:5984 -databases=db-1,db-2 -couchdb.username=root -couchdb.password=a-secret
 

## Metrics Overview
The following list gives you an overview on the currently exposed metrics.
Please note that beyond the complete CouchDB stats a metric `couchdb_up` has been
added as minimal connection health check.

The `node_name` label defaults to `master` for CouchDB 1.x and `nonode@nohost` for a single CouchDB 2.x setup.
The node name makes more sense in a clustered CouchDB 2.x environment with several nodes to provide separated stats for each node. 

    # HELP couchdb_database_data_size data size
    # TYPE couchdb_database_data_size gauge
    couchdb_database_data_size{db_name="db-1",node_name="nonode@nohost"} 1.764941e+06
    couchdb_database_data_size{db_name="db-2",node_name="nonode@nohost"} 4.93011369e+08
    # HELP couchdb_database_disk_size disk size
    # TYPE couchdb_database_disk_size gauge
    couchdb_database_disk_size{db_name="db-1",node_name="nonode@nohost"} 2.098663e+06
    couchdb_database_disk_size{db_name="db-2",node_name="nonode@nohost"} 6.60108847e+08
    # HELP couchdb_database_overhead disk size overhead
    # TYPE couchdb_database_overhead gauge
    couchdb_database_overhead{db_name="db-1",node_name="nonode@nohost"} 333722
    couchdb_database_overhead{db_name="db-2",node_name="nonode@nohost"} 1.67097478e+08
    # HELP couchdb_httpd_auth_cache_hits number of authentication cache hits
    # TYPE couchdb_httpd_auth_cache_hits gauge
    couchdb_httpd_auth_cache_hits{node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_auth_cache_misses number of authentication cache misses
    # TYPE couchdb_httpd_auth_cache_misses gauge
    couchdb_httpd_auth_cache_misses{node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_bulk_requests number of bulk requests
    # TYPE couchdb_httpd_bulk_requests gauge
    couchdb_httpd_bulk_requests{node_name="nonode@nohost"} 1871
    # HELP couchdb_httpd_clients_requesting_changes number of clients for continuous _changes
    # TYPE couchdb_httpd_clients_requesting_changes gauge
    couchdb_httpd_clients_requesting_changes{node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_database_reads number of times a document was read from a database
    # TYPE couchdb_httpd_database_reads gauge
    couchdb_httpd_database_reads{node_name="nonode@nohost"} 8
    # HELP couchdb_httpd_database_writes number of times a database was changed
    # TYPE couchdb_httpd_database_writes gauge
    couchdb_httpd_database_writes{node_name="nonode@nohost"} 14955
    # HELP couchdb_httpd_open_databases number of open databases
    # TYPE couchdb_httpd_open_databases gauge
    couchdb_httpd_open_databases{node_name="nonode@nohost"} 16
    # HELP couchdb_httpd_open_os_files number of file descriptors CouchDB has open
    # TYPE couchdb_httpd_open_os_files gauge
    couchdb_httpd_open_os_files{node_name="nonode@nohost"} 16
    # HELP couchdb_httpd_request_methods number of HTTP requests by method
    # TYPE couchdb_httpd_request_methods gauge
    couchdb_httpd_request_methods{method="COPY",node_name="nonode@nohost"} 0
    couchdb_httpd_request_methods{method="DELETE",node_name="nonode@nohost"} 0
    couchdb_httpd_request_methods{method="GET",node_name="nonode@nohost"} 1465
    couchdb_httpd_request_methods{method="HEAD",node_name="nonode@nohost"} 0
    couchdb_httpd_request_methods{method="POST",node_name="nonode@nohost"} 2995
    couchdb_httpd_request_methods{method="PUT",node_name="nonode@nohost"} 152
    # HELP couchdb_httpd_request_time length of a request inside CouchDB without MochiWeb
    # TYPE couchdb_httpd_request_time gauge
    couchdb_httpd_request_time{node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_requests number of HTTP requests
    # TYPE couchdb_httpd_requests gauge
    couchdb_httpd_requests{node_name="nonode@nohost"} 4611
    # HELP couchdb_httpd_status_codes number of HTTP responses by status code
    # TYPE couchdb_httpd_status_codes gauge
    couchdb_httpd_status_codes{code="200",node_name="nonode@nohost"} 2422
    couchdb_httpd_status_codes{code="201",node_name="nonode@nohost"} 2171
    couchdb_httpd_status_codes{code="202",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="301",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="304",node_name="nonode@nohost"} 4
    couchdb_httpd_status_codes{code="400",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="401",node_name="nonode@nohost"} 2
    couchdb_httpd_status_codes{code="403",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="404",node_name="nonode@nohost"} 9
    couchdb_httpd_status_codes{code="405",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="409",node_name="nonode@nohost"} 0
    couchdb_httpd_status_codes{code="412",node_name="nonode@nohost"} 1
    couchdb_httpd_status_codes{code="500",node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_temporary_view_reads number of temporary view reads
    # TYPE couchdb_httpd_temporary_view_reads gauge
    couchdb_httpd_temporary_view_reads{node_name="nonode@nohost"} 0
    # HELP couchdb_httpd_up Was the last query of CouchDB stats successful.
    # TYPE couchdb_httpd_up gauge
    couchdb_httpd_up 1
    # HELP couchdb_httpd_view_reads number of view reads
    # TYPE couchdb_httpd_view_reads gauge
    couchdb_httpd_view_reads{node_name="nonode@nohost"} 0
    # HELP couchdb_server_active_tasks active tasks
    # TYPE couchdb_server_active_tasks gauge
    couchdb_server_active_tasks{database_compaction="0",indexer="0",node_name="nonode@nohost",replication="1",view_compaction="0"} 1
    ... (more go internal stats)

## Thanks

Thanks go to the Prometheus team, which is very active and responsive!

I also have to admit that the couchdb-prometheus-exporter code is heavily inspired by 
the other [available exporters](http://prometheus.io/docs/instrumenting/exporters/), 
and that some ideas have just been copied from them.

## Monitoring CouchDB with Prometheus, Grafana and Docker

For a step-by-step guide, see [Monitoring CouchDB with Prometheus, Grafana and Docker](https://medium.com/@redgeoff/monitoring-couchdb-with-prometheus-grafana-and-docker-4693bc8408f0)
