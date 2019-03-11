package main

import (
	goflag "flag"
	"fmt"
	"github.com/gesellix/couchdb-prometheus-exporter/glogadapt"
	"github.com/gesellix/couchdb-prometheus-exporter/lib"
	"github.com/golang/glog"
	"github.com/namsral/flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"net/http"
	"strconv"
	"strings"
)

type exporterConfigType struct {
	listenAddress   string
	metricsEndpoint string
	couchdbURI      string
	couchdbUsername string
	couchdbPassword string
	couchdbInsecure bool
	databases       string
	databaseViews   bool
	schedulerJobs   bool
}

var exporterConfig exporterConfigType

func init() {
	flag.String(flag.DefaultConfigFlagname, "", "path to config file")
	flag.StringVar(&exporterConfig.listenAddress, "telemetry.address", "localhost:9984", "Address on which to expose metrics.")
	flag.StringVar(&exporterConfig.metricsEndpoint, "telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	flag.StringVar(&exporterConfig.couchdbURI, "couchdb.uri", "http://localhost:5984", "URI to the CouchDB instance")
	flag.StringVar(&exporterConfig.couchdbUsername, "couchdb.username", "", "Basic auth username for the CouchDB instance")
	flag.StringVar(&exporterConfig.couchdbPassword, "couchdb.password", "", "Basic auth password for the CouchDB instance")
	flag.BoolVar(&exporterConfig.couchdbInsecure, "couchdb.insecure", true, "Ignore server certificate if using https")
	flag.StringVar(&exporterConfig.databases, "databases", "", fmt.Sprintf("Comma separated list of database names, or '%s'", lib.AllDbs))
	flag.BoolVar(&exporterConfig.databaseViews, "databases.views", true, "Collect view details of every observed database")
	flag.BoolVar(&exporterConfig.schedulerJobs, "scheduler.jobs", false, "Collect active replication jobs (CouchDB 2.x+ only)")

	flag.BoolVar(&glogadapt.Logging.ToStderr, "logtostderr", false, "log to standard error instead of files")
	flag.BoolVar(&glogadapt.Logging.AlsoToStderr, "alsologtostderr", false, "log to standard error as well as files")
	flag.Var(&glogadapt.Logging.Verbosity, "v", "log level for V logs")
	flag.Var(&glogadapt.Logging.StderrThreshold, "stderrthreshold", "logs at or above this threshold go to stderr")
	flag.StringVar(&glogadapt.Logging.LogDir, "log_dir", "", "If non-empty, write log files in this directory")
}

func main() {
	err := initFlags()
	if err != nil {
		glog.Fatal(err)
	}

	var databases []string
	if *&exporterConfig.databases != "" {
		databases = strings.Split(*&exporterConfig.databases, ",")
	}

	exporter := lib.NewExporter(
		*&exporterConfig.couchdbURI,
		lib.BasicAuth{
			Username: *&exporterConfig.couchdbUsername,
			Password: *&exporterConfig.couchdbPassword},
		lib.CollectorConfig{
			Databases:            databases,
			CollectViews:         *&exporterConfig.databaseViews,
			CollectSchedulerJobs: *&exporterConfig.schedulerJobs,
		},
		*&exporterConfig.couchdbInsecure)
	prometheus.MustRegister(exporter)

	http.Handle(*&exporterConfig.metricsEndpoint, promhttp.Handler())
	http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
		_, err = fmt.Fprint(w, "OK")
		if err != nil {
			glog.Error(err)
		}
	})
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, fmt.Sprintf("Please GET %s", *&exporterConfig.metricsEndpoint), http.StatusNotFound)
	})

	glog.Infof("Starting exporter at '%s' to read from CouchDB at '%s'", *&exporterConfig.listenAddress, *&exporterConfig.couchdbURI)
	err = http.ListenAndServe(*&exporterConfig.listenAddress, nil)
	if err != nil {
		glog.Fatal(err)
	}
}

func initFlags() error {
	flag.Parse()
	// Convinces goflags that we have called Parse() to avoid noisy logs.
	// Necessary due to https://github.com/golang/glog/commit/65d674618f712aa808a7d0104131b9206fc3d5ad
	// and us using another flags package.
	err := goflag.CommandLine.Parse([]string{})
	if err != nil {
		return err
	}
	err = goflag.Lookup("logtostderr").Value.Set(strconv.FormatBool(*&glogadapt.Logging.ToStderr))
	if err != nil {
		return err
	}
	err = goflag.Lookup("alsologtostderr").Value.Set(strconv.FormatBool(*&glogadapt.Logging.AlsoToStderr))
	if err != nil {
		return err
	}
	err = goflag.Lookup("v").Value.Set(glogadapt.Logging.Verbosity.String())
	if err != nil {
		return err
	}
	err = goflag.Lookup("stderrthreshold").Value.Set(glogadapt.Logging.StderrThreshold.String())
	if err != nil {
		return err
	}
	err = goflag.Lookup("log_dir").Value.Set(glogadapt.Logging.LogDir)
	if err != nil {
		return err
	}
	return nil
}
