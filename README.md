# CouchDB Prometheus Exporter

[![Build Status](https://travis-ci.org/gesellix/couchdb-prometheus-exporter.svg?branch=master)](https://travis-ci.org/gesellix/couchdb-prometheus-exporter)

[CouchDB](http://couchdb.apache.org/) metrics exporter for [Prometheus](http://prometheus.io/)

The CouchDB metrics exporter requests the CouchDB stats from the `/_stats` and `/_active_tasks` endpoints and 
exposes them for Prometheus consumption. You can optionally monitor detailed database stats like
disk and data size to monitor the storage overhead. The exporter can be configured via program parameters,
environment variables, and config file.

## Build the binary

You can find pre-build releases for different platforms at [our GitHub Releases](https://github.com/gesellix/couchdb-prometheus-exporter/releases).

If you prefer to build your own binary or in case you'd like to build from the current `master`,
you'll have to get and install [a recent version of Golang](https://golang.org/dl/) for your platform, first.
Then, you have to perform the following commands in the cloned repository:

````shell script
export GO111MODULE=on  # in a Windows shell, please replace `export` with `set`
go get github.com/gesellix/couchdb-prometheus-exporter
````

Those commands will install the binary in your local `GOBIN` directory, usually something like
`$HOME/go/bin`. Please ensure that the directory is in your system's `PATH`. Then the following
should work:

````shell script
couchdb-prometheus-exporter --help
````

## Run the binary

You can get an overview over possible configuration options with their defaults in the help screen:

    couchdb-prometheus-exporter --help

Configuration is possible via:

- environment variables (e.g. `COUCHDB_USERNAME=admin`)
- command line parameters (e.g. `--couchdb.username admin`)
- configuration file (e.g. `--config=config.ini`)

The configuration file format is the "properties" file format, e.g. like this:

````properties
couchdb.username=admin
couchdb.password=a-secret
````

## Using TLS and/or Basic authentication

TLS and/or Basic authentication is supported via `--web.config` parameter:

    couchdb-prometheus-exporter --config=config.ini --web.config=web-config.yaml

A complete `web-config.yml` might look like this:

````yaml
---
tls_server_config :
  cert_file : "path/to/https/server.crt"
  key_file : "path/to/https/server.key"
basic_auth_users:
  alice: $2y$12$1DpfPeqF9HzHJt.EWswy1exHluGfbhnn3yXhR7Xes6m3WJqFg0Wby
  bob: $2y$18$4VeFDzXIoPHKnKTU3O3GH.N.vZu06CVqczYZ8WvfzrddFU6tGqjR.
  carol: $2y$10$qRTBuFoULoYNA7AQ/F3ck.trZBPyjV64.oA4ZsSBCIWvXuvQlQTuu
  dave: $2y$10$2UXri9cIDdgeKjBo4Rlpx.U3ZLDV8X1IxKmsfOvhcM5oXQt/mLmXq
...
````

For further information about TLS and/or Basic auth,
please visit: [exporter-toolkit/https](https://pkg.go.dev/github.com/prometheus/exporter-toolkit@v0.4.0/https)
or [github.com/prometheus/exporter-toolkit](https://github.com/prometheus/exporter-toolkit).

## Run it as container

    docker run --rm -p 9984:9984 gesellix/couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --logtostderr

Please note that host names like `localhost` won't leave the container, so you have to use non-loopback
dns names or ip addresses when configuring the CouchDB URI.

## Logging

The couchdb-exporter uses the [glog](https://godoc.org/github.com/golang/glog) library for logging.
With the default parameters everything will be logged to `/tmp/`.
Use `--logtostderr` to enable logging to stderr and `--help` to see all options.

## CouchDB 2+ clusters

For CouchDB 2.x, you should configure the exporter to fetch the stats from one node, to get
a complete cluster overview. In contrast to CouchDB 1.x you'll need to configure the admin
credentials, e.g. like this:

    couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --couchdb.username=root --couchdb.password=a-secret

## Database disk usage stats

If you need database disk usage stats, add a comma separated list of database names like this:

    couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --databases=db-1,db-2 --couchdb.username=root --couchdb.password=a-secret

Or, if you want to get stats for every database, please use `_all_dbs` as database name:

    couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --databases=_all_dbs --couchdb.username=root --couchdb.password=a-secret

## Monitoring CouchDB with Prometheus, Grafana and Docker

For a step-by-step guide, see [Monitoring CouchDB with Prometheus, Grafana and Docker](https://medium.com/@redgeoff/monitoring-couchdb-with-prometheus-grafana-and-docker-4693bc8408f0)

A complete example stack including multiple CouchDB instances, exporters, Prometheus, Grafana, etc. is available at `examples/grafana/`, and can be deployed locally:

````bash
cd examples/grafana
docker swarm init
docker stack deploy --compose-file docker-traefik-stack.yml example
````

## Examples

The `examples` directory in this repository contains ready-to-run examples for

- [Docker Compose](examples/compose/README.md)
- [Docker Swarm/Stack](examples/stack/README.md)

## Credits

Thanks go to the Prometheus team, which is very active and responsive!

I also have to admit that the couchdb-prometheus-exporter code is heavily inspired by 
the other [available exporters](http://prometheus.io/docs/instrumenting/exporters/), 
and that some ideas have just been copied from them.

Last but not least, this project wouldn't be possible without users submitting issues,
feature requests and adding [code contributions](https://github.com/gesellix/couchdb-prometheus-exporter/graphs/contributors).
Thanks a lot!

## Metrics Overview
The file [README_metrics.md](https://github.com/gesellix/couchdb-prometheus-exporter/blob/master/README_metrics.md) gives you an overview on the currently exposed metrics.
