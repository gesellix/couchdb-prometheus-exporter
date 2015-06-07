package main

import (
	"flag"
	"fmt"
	"net/http"

	"github.com/gesellix/couchdb-exporter/lib"
	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	listenAddress   = flag.String("telemetry.address", "localhost:9984", "Address on which to expose metrics.")
	metricsEndpoint = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	couchdbURI      = flag.String("couchdb.uri", "http://localhost:5984", "URI to the CouchDB instance")
	couchdbUsername = flag.String("couchdb.username", "", "Basic auth username for the CouchDB instance")
	couchdbPassword = flag.String("couchdb.password", "", "Basic auth password for the CouchDB instance")
)

func main() {
	flag.Parse()

	exporter := lib.NewExporter(*couchdbURI, lib.BasicAuth{Username: *couchdbUsername, Password: *couchdbPassword})
	prometheus.MustRegister(exporter)

	http.Handle(*metricsEndpoint, prometheus.Handler())
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "OK")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *metricsEndpoint, http.StatusMovedPermanently)
	})

	glog.Infof("Starting exporter at %s to read from CouchDB at %s", *listenAddress, *couchdbURI)
	err := http.ListenAndServe(*listenAddress, nil)
	if err != nil {
		glog.Fatal(err)
	}
}
