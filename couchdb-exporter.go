package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/urfave/cli"
	"k8s.io/klog"

	"github.com/gesellix/couchdb-prometheus-exporter/lib"
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

type loggingConfigType struct {
	toStderr        bool   // The -logtostderr flag.
	alsoToStderr    bool   // The -alsologtostderr flag.
	verbosity       int    // V logging level, the value of the -v flag/
	stderrThreshold int    // The -stderrthreshold flag.
	logDir          string // The -log_dir flag.
}

var exporterConfig exporterConfigType

// custom exposed (but hidden) logging config flags
var loggingConfig loggingConfigType

var appFlags []cli.Flag

// TODO graceful migration to new parameter names
// 1) Warn, that these parameters are deprecated and will be removed/renamed
// 2) Fail at startup, when deprecated parameters are used. Maybe allow override by explicit "i-know-what-i-am-doing"-parameter
// 3) Remove (ignore) deprecated parameters
func init() {
	// TODO replace mechanism with urfave/cli
	//flag.String(flag.DefaultConfigFlagname, "", "path to config file")

	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:  "telemetry.address",
			Usage: "Address on which to expose metrics",
			//EnvVar:      "TELEMETRY_ADDRESS",
			Hidden:      false,
			Value:       "localhost:9984",
			Destination: &exporterConfig.listenAddress,
		},
		cli.StringFlag{
			Name:  "telemetry.endpoint",
			Usage: "Path under which to expose metrics",
			//EnvVar:      "TELEMETRY_ENDPOINT",
			Hidden:      false,
			Value:       "/metrics",
			Destination: &exporterConfig.metricsEndpoint,
		},
		cli.StringFlag{
			Name:  "couchdb.uri",
			Usage: "URI to the CouchDB instance",
			//EnvVar:      "COUCHDB_URI",
			Hidden:      false,
			Value:       "http://localhost:5984",
			Destination: &exporterConfig.couchdbURI,
		},
		cli.StringFlag{
			Name:  "couchdb.username",
			Usage: "Basic auth username for the CouchDB instance",
			//EnvVar:      "COUCHDB_USERNAME",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbUsername,
		},
		cli.StringFlag{
			Name:  "couchdb.password",
			Usage: "Basic auth password for the CouchDB instance",
			//EnvVar:      "COUCHDB_PASSWORD",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbPassword,
		},
		// TODO doesn't print the default when showing the command help
		cli.BoolTFlag{
			Name:  "couchdb.insecure",
			Usage: "Ignore server certificate if using https",
			//EnvVar:      "COUCHDB_INSECURE",
			Hidden:      false,
			Destination: &exporterConfig.couchdbInsecure,
		},
		// TODO use cli.StringSliceFlag?
		cli.StringFlag{
			Name:  "databases",
			Usage: fmt.Sprintf("Comma separated list of database names, or '%s'", lib.AllDbs),
			//EnvVar:      "DATABASES",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.databases,
		},
		// TODO doesn't print the default when showing the command help
		cli.BoolTFlag{
			Name:  "databases.views",
			Usage: "Collect view details of every observed database",
			//EnvVar:      "DATABASES_VIEWS",
			Hidden:      false,
			Destination: &exporterConfig.databaseViews,
		},
		// TODO doesn't print the default when showing the command help
		cli.BoolFlag{
			Name:  "scheduler.jobs",
			Usage: "Collect active replication jobs (CouchDB 2.x+ only)",
			//EnvVar:      "SCHEDULER_JOBS",
			Hidden:      false,
			Destination: &exporterConfig.schedulerJobs,
		},

		cli.BoolTFlag{
			Name:        "logtostderr",
			Usage:       "log to standard error instead of files",
			Hidden:      true,
			Destination: &loggingConfig.toStderr,
		},
		cli.BoolFlag{
			Name:        "alsologtostderr",
			Usage:       "log to standard error as well as files",
			Hidden:      true,
			Destination: &loggingConfig.alsoToStderr,
		},
		// TODO `v` clashed with urfave/cli's `--version` shortcut `-v`.
		cli.IntFlag{
			Name:        "verbosity",
			Usage:       "log level for V logs",
			Value:       0,
			Hidden:      true,
			Destination: &loggingConfig.verbosity,
		},
		cli.IntFlag{
			Name:        "stderrthreshold",
			Usage:       "logs at or above this threshold go to stderr",
			Value:       2,
			Hidden:      true,
			Destination: &loggingConfig.stderrThreshold,
		},
		cli.StringFlag{
			Name:        "log_dir",
			Usage:       "If non-empty, write log files in this directory",
			Hidden:      true,
			Destination: &loggingConfig.logDir,
		},
	}
}

func main() {
	var appAction = func(c *cli.Context) error {
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
			_, err := fmt.Fprint(w, "OK")
			if err != nil {
				klog.Error(err)
			}
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("Please GET %s", *&exporterConfig.metricsEndpoint), http.StatusNotFound)
		})

		klog.Infof("Starting exporter at '%s' to read from CouchDB at '%s'", *&exporterConfig.listenAddress, *&exporterConfig.couchdbURI)
		err := http.ListenAndServe(*&exporterConfig.listenAddress, nil)
		if err != nil {
			klog.Fatal(err)
		}
		return err
	}

	app := cli.NewApp()
	app.Name = "CouchDB Prometheus Exporter"
	//app.Usage = ""
	app.Description = "CouchDB stats exporter for Prometheus"
	//app.Version = ""
	app.Flags = appFlags
	app.Before = initKlogFlags
	app.Action = appAction

	defer klog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		klog.Fatal(err)
	}
}

func initKlogFlags(_ *cli.Context) error {
	klogFlags := flag.NewFlagSet("klog", flag.ContinueOnError)
	klog.InitFlags(klogFlags)

	flags := map[string]string{
		"logtostderr":     strconv.FormatBool(loggingConfig.toStderr),
		"alsologtostderr": strconv.FormatBool(loggingConfig.alsoToStderr),
		"stderrthreshold": strconv.Itoa(loggingConfig.stderrThreshold),
		"v":               strconv.Itoa(loggingConfig.verbosity),
		"log_dir":         loggingConfig.logDir,
	}
	for k, v := range flags {
		if err := klogFlags.Set(k, v); err != nil {
			return err
		}
	}

	klog.Infof("adopted logging config: %+v\n", loggingConfig)
	return nil
}
