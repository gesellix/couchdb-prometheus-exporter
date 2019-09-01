# CouchDB Prometheus Exporter

[![Build Status](https://travis-ci.org/gesellix/couchdb-prometheus-exporter.svg?branch=master)](https://travis-ci.org/gesellix/couchdb-prometheus-exporter)

[CouchDB](http://couchdb.apache.org/) metrics exporter for [Prometheus](http://prometheus.io/)

The CouchDB metrics exporter requests the CouchDB stats from the `/_stats` and `/_active_tasks` endpoints and 
exposes them for Prometheus consumption. You can optionally monitor detailed database stats like
disk and data size to monitor the storage overhead. The exporter can be configured via program parameters,
environment variables, and config file.

## Run it as container

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --logtostderr

The couchdb-exporter uses the [glog](https://godoc.org/github.com/golang/glog) library for logging.
With the default parameters everything will be logged to `/tmp/`.
Use `--logtostderr` to enable logging to stderr and `--help` to see all options.

For CouchDB 2.x, you should configure the exporter to fetch the stats from one node, to get
a complete cluster overview. In contrast to CouchDB 1.x you'll need to configure the admin
credentials, e.g. like this:

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --couchdb.username=root --couchdb.password=a-secret

## Database disk usage stats

If you need database disk usage stats, add a comma separated list of database names like this:

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --databases=db-1,db-2 --couchdb.username=root --couchdb.password=a-secret

Or, if you want to get stats for every database, please use `_all_dbs` as database name:

    docker run -p 9984:9984 gesellix/couchdb-prometheus-exporter --couchdb.uri=http://couchdb:5984 --databases=_all_dbs --couchdb.username=root --couchdb.password=a-secret


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
