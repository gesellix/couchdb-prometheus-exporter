package main

import (
	"fmt"
	goflag "flag"
	"net/http"
	"strconv"
	"github.com/gesellix/couchdb-exporter/lib"
	"github.com/golang/glog"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"strings"
)

type exporterConfigType struct {
	listenAddress   string
	metricsEndpoint string
	couchdbURI      string
	couchdbUsername string
	couchdbPassword string
	databases       string
}

var exporterConfig exporterConfigType

func init() {
	flag.StringVar(&exporterConfig.listenAddress, "telemetry.address", "localhost:9984", "Address on which to expose metrics.")
	flag.StringVar(&exporterConfig.metricsEndpoint, "telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&exporterConfig.couchdbURI, "couchdb.uri", "http://localhost:5984", "URI to the CouchDB instance")
	flag.StringVar(&exporterConfig.couchdbUsername, "couchdb.username", "", "Basic auth username for the CouchDB instance")
	flag.StringVar(&exporterConfig.couchdbPassword, "couchdb.password", "", "Basic auth password for the CouchDB instance")
	flag.StringVar(&exporterConfig.databases, "databases", "", "Comma separated list of database names")

	flag.BoolVar(&logging.toStderr, "logtostderr", false, "log to standard error instead of files")
	flag.BoolVar(&logging.alsoToStderr, "alsologtostderr", false, "log to standard error as well as files")
	flag.Var(&logging.verbosity, "v", "log level for V logs")
	flag.Var(&logging.stderrThreshold, "stderrthreshold", "logs at or above this threshold go to stderr")

	// Default stderrThreshold is ERROR.
	logging.stderrThreshold = errorLog
}

func main() {
	flag.Parse()
	// Convinces goflags that we have called Parse() to avoid noisy logs.
	// Necessary due to https://github.com/golang/glog/commit/65d674618f712aa808a7d0104131b9206fc3d5ad
	// and us using another flags package.
	goflag.CommandLine.Parse([]string{})
	goflag.Lookup("logtostderr").Value.Set(strconv.FormatBool(*&logging.toStderr))
	goflag.Lookup("alsologtostderr").Value.Set(strconv.FormatBool(*&logging.alsoToStderr))
	goflag.Lookup("v").Value.Set(logging.verbosity.String())
	goflag.Lookup("stderrthreshold").Value.Set(logging.stderrThreshold.String())

	databases := strings.Split(*&exporterConfig.databases, ",")

	exporter := lib.NewExporter(
		*&exporterConfig.couchdbURI,
		lib.BasicAuth{Username: *&exporterConfig.couchdbUsername, Password: *&exporterConfig.couchdbPassword},
		databases)
	prometheus.MustRegister(exporter)

	http.Handle(*&exporterConfig.metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "OK")
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, *&exporterConfig.metricsEndpoint, http.StatusMovedPermanently)
	})

	glog.Infof("Starting exporter at %s to read from CouchDB at %s", *&exporterConfig.listenAddress, *&exporterConfig.couchdbURI)
	err := http.ListenAndServe(*&exporterConfig.listenAddress, nil)
	if err != nil {
		glog.Fatal(err)
	}
}
