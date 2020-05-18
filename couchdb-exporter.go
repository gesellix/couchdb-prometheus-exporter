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
	"github.com/urfave/cli/altsrc"
	"k8s.io/klog"

	"github.com/gesellix/couchdb-prometheus-exporter/fileutil"
	"github.com/gesellix/couchdb-prometheus-exporter/lib"
)

var (
	version = "dev"
	commit  = "none"
	date    = "unknown"
)

type exporterConfigType struct {
	listenAddress              string
	metricsEndpoint            string
	couchdbURI                 string
	couchdbUsername            string
	couchdbPassword            string
	couchdbInsecure            bool
	databases                  string
	databaseViews              bool
	databaseConcurrentRequests uint
	schedulerJobs              bool
}

type loggingConfigType struct {
	toStderr        bool   // The -logtostderr flag.
	alsoToStderr    bool   // The -alsologtostderr flag.
	verbosity       int    // V logging level, the value of the -v flag/
	stderrThreshold int    // The -stderrthreshold flag.
	logDir          string // The -log_dir flag.
}

var exporterConfig exporterConfigType

var configFileFlagname = "config"

// custom exposed (but hidden) logging config flags
var loggingConfig loggingConfigType

var appFlags []cli.Flag

// TODO graceful migration to new parameter names
// 1) Warn, for deprecated parameters to be removed/renamed
// 2) Fail at startup, when deprecated parameters are used. Maybe allow override by explicit "i-know-what-i-am-doing"-parameter
// 3) Remove (ignore) deprecated parameters
func init() {
	appFlags = []cli.Flag{
		cli.StringFlag{
			Name:   configFileFlagname,
			Usage:  "Path to config file",
			EnvVar: "CONFIG",
			Hidden: false,
		},
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "telemetry.address",
			Usage:       "Address on which to expose metrics",
			EnvVar:      "TELEMETRY.ADDRESS,TELEMETRY_ADDRESS",
			Hidden:      false,
			Value:       "localhost:9984",
			Destination: &exporterConfig.listenAddress,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "telemetry.endpoint",
			Usage:       "Path under which to expose metrics",
			EnvVar:      "TELEMETRY.ENDPOINT,TELEMETRY_ENDPOINT",
			Hidden:      false,
			Value:       "/metrics",
			Destination: &exporterConfig.metricsEndpoint,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "couchdb.uri",
			Usage:       "URI to the CouchDB instance",
			EnvVar:      "COUCHDB.URI,COUCHDB_URI",
			Hidden:      false,
			Value:       "http://localhost:5984",
			Destination: &exporterConfig.couchdbURI,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "couchdb.username",
			Usage:       "Basic auth username for the CouchDB instance",
			EnvVar:      "COUCHDB.USERNAME,COUCHD_USERNAME",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbUsername,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "couchdb.password",
			Usage:       "Basic auth password for the CouchDB instance",
			EnvVar:      "COUCHDB.PASSWORD,COUCHDB_PASSWORD",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbPassword,
		}),
		// TODO doesn't print the default when showing the command help
		altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:        "couchdb.insecure",
			Usage:       "Ignore server certificate if using https",
			EnvVar:      "COUCHDB.INSECURE,COUCHDB_INSECURE",
			Hidden:      false,
			Destination: &exporterConfig.couchdbInsecure,
		}),
		// TODO use cli.StringSliceFlag?
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "databases",
			Usage:       fmt.Sprintf("Comma separated list of database names, or '%s'", lib.AllDbs),
			EnvVar:      "DATABASES",
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.databases,
		}),
		// TODO doesn't print the default when showing the command help
		altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:        "databases.views",
			Usage:       "Collect view details of every observed database",
			EnvVar:      "DATABASES.VIEWS,DATABASES_VIEWS",
			Hidden:      false,
			Destination: &exporterConfig.databaseViews,
		}),
		altsrc.NewUintFlag(cli.UintFlag{
			Name:        "database.concurrent.requests",
			Usage:       "maximum concurrent calls to CouchDB, or 0 for unlimited",
			Value:       0,
			Hidden:      false,
			Destination: &exporterConfig.databaseConcurrentRequests,
		}),
		// TODO doesn't print the default when showing the command help
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "scheduler.jobs",
			Usage:       "Collect active replication jobs (CouchDB 2.x+ only)",
			EnvVar:      "SCHEDULER.JOBS,SCHEDULER_JOBS",
			Hidden:      false,
			Destination: &exporterConfig.schedulerJobs,
		}),

		altsrc.NewBoolTFlag(cli.BoolTFlag{
			Name:        "logtostderr",
			Usage:       "log to standard error instead of files",
			Hidden:      true,
			Destination: &loggingConfig.toStderr,
		}),
		altsrc.NewBoolFlag(cli.BoolFlag{
			Name:        "alsologtostderr",
			Usage:       "log to standard error as well as files",
			Hidden:      true,
			Destination: &loggingConfig.alsoToStderr,
		}),
		// TODO `v` clashed with urfave/cli's "version" shortcut `-v`.
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "verbosity",
			Usage:       "log level for V logs",
			Value:       0,
			Hidden:      true,
			Destination: &loggingConfig.verbosity,
		}),
		altsrc.NewIntFlag(cli.IntFlag{
			Name:        "stderrthreshold",
			Usage:       "logs at or above this threshold go to stderr",
			Value:       2,
			Hidden:      true,
			Destination: &loggingConfig.stderrThreshold,
		}),
		altsrc.NewStringFlag(cli.StringFlag{
			Name:        "log_dir",
			Usage:       "If non-empty, write log files in this directory",
			Hidden:      true,
			Destination: &loggingConfig.logDir,
		}),
	}
}

func main() {
	var appAction = func(c *cli.Context) error {
		var databases []string
		if exporterConfig.databases != "" {
			databases = strings.Split(exporterConfig.databases, ",")
		}

		exporter := lib.NewExporter(
			exporterConfig.couchdbURI,
			lib.BasicAuth{
				Username: exporterConfig.couchdbUsername,
				Password: exporterConfig.couchdbPassword},
			lib.CollectorConfig{
				Databases:            databases,
				CollectViews:         exporterConfig.databaseViews,
				CollectSchedulerJobs: exporterConfig.schedulerJobs,
				ConcurrentRequests:   exporterConfig.databaseConcurrentRequests,
			},
			exporterConfig.couchdbInsecure)
		prometheus.MustRegister(exporter)

		http.Handle(exporterConfig.metricsEndpoint, promhttp.Handler())
		http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, "OK")
			if err != nil {
				klog.Error(err)
			}
		})
		http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
			http.Error(w, fmt.Sprintf("Please GET %s", exporterConfig.metricsEndpoint), http.StatusNotFound)
		})

		klog.Infof("Starting exporter at '%s' to read from CouchDB at '%s'", exporterConfig.listenAddress, exporterConfig.couchdbURI)
		err := http.ListenAndServe(exporterConfig.listenAddress, nil)
		if err != nil {
			klog.Fatal(err)
		}
		return err
	}

	app := cli.NewApp()
	app.Name = "CouchDB Prometheus Exporter"
	//app.Usage = ""
	app.Description = "CouchDB stats exporter for Prometheus"
	app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	app.Flags = appFlags
	app.Before = beforeApp(appFlags)
	app.Action = appAction

	defer klog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		klog.Fatal(err)
	}
}

func beforeApp(appFlags []cli.Flag) cli.BeforeFunc {
	return func(context *cli.Context) error {
		// TODO decide on a preferred config file format, maybe support different ones.
		inputSource := fileutil.NewPropertiesSourceFromFlagFunc(configFileFlagname)
		//inputSource := altsrc.NewYamlSourceFromFlagFunc(configFileFlagname)
		//inputSource := altsrc.NewTomlSourceFromFlagFunc(configFileFlagname)
		if err := altsrc.InitInputSourceWithContext(appFlags, inputSource)(context); err != nil {
			return err
		}
		return initKlogFlags(context)
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
