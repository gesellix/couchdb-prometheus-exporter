package main

import (
	"fmt"
	"log"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/exporter-toolkit/web"
	"github.com/urfave/cli/v2"
	"github.com/urfave/cli/v2/altsrc"

	"github.com/gesellix/couchdb-prometheus-exporter/v30/fileutil"
	"github.com/gesellix/couchdb-prometheus-exporter/v30/lib"
)

var (
	version = "dev"
	commit  = "none"
	date    = time.Now().Format(time.RFC3339)
)

type webConfigType struct {
	listenAddress   string
	metricsEndpoint string
}

type exporterConfigType struct {
	couchdbURI                 string
	couchdbUsername            string
	couchdbPassword            string
	couchdbInsecure            bool
	scrapeInterval             time.Duration
	scrapeLocalNode            bool
	databases                  string
	databaseViews              bool
	databaseConcurrentRequests uint
	schedulerJobs              bool
}

var exporterConfig exporterConfigType
var webConfig webConfigType

var configFileFlagname = "config"
var webConfigFile = ""

var appFlags []cli.Flag

// TODO graceful migration to new parameter names
// 1) Warn, for deprecated parameters to be removed/renamed
// 2) Fail at startup, when deprecated parameters are used. Maybe allow override by explicit "i-know-what-i-am-doing"-parameter
// 3) Remove (ignore) deprecated parameters
func init() {
	appFlags = []cli.Flag{
		&cli.StringFlag{
			Name:    configFileFlagname,
			Usage:   "Path to config ini file that configures the CouchDB connection",
			EnvVars: []string{"CONFIG"},
			Hidden:  false,
		},
		&cli.StringFlag{
			Name:        "web.config",
			Usage:       "Path to config yaml file that can enable TLS or authentication",
			EnvVars:     []string{"WEB_CONFIG"},
			Hidden:      false,
			Value:       "",
			Destination: &webConfigFile,
		},
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "telemetry.address",
			Usage:       "Address on which to expose metrics",
			EnvVars:     []string{"TELEMETRY.ADDRESS", "TELEMETRY_ADDRESS"},
			Hidden:      false,
			Value:       "localhost:9984",
			Destination: &webConfig.listenAddress,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "telemetry.endpoint",
			Usage:       "Path under which to expose metrics",
			EnvVars:     []string{"TELEMETRY.ENDPOINT", "TELEMETRY_ENDPOINT"},
			Hidden:      false,
			Value:       "/metrics",
			Destination: &webConfig.metricsEndpoint,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "couchdb.uri",
			Usage:       "URI to the CouchDB instance",
			EnvVars:     []string{"COUCHDB.URI", "COUCHDB_URI"},
			Hidden:      false,
			Value:       "http://localhost:5984",
			Destination: &exporterConfig.couchdbURI,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "couchdb.username",
			Usage:       "Basic auth username for the CouchDB instance",
			EnvVars:     []string{"COUCHDB.USERNAME", "COUCHDB_USERNAME"},
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbUsername,
		}),
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "couchdb.password",
			Usage:       "Basic auth password for the CouchDB instance",
			EnvVars:     []string{"COUCHDB.PASSWORD", "COUCHDB_PASSWORD"},
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.couchdbPassword,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "couchdb.insecure",
			Usage:       "Ignore server certificate if using https",
			EnvVars:     []string{"COUCHDB.INSECURE", "COUCHDB_INSECURE"},
			Hidden:      false,
			Value:       true,
			Destination: &exporterConfig.couchdbInsecure,
		}),
		altsrc.NewDurationFlag(&cli.DurationFlag{
			Name:        "scrape.interval",
			Usage:       fmt.Sprintf("Duration between metrics collection from the CouchDB cluster. '0s' collects only on Prometheus scrapes"),
			EnvVars:     []string{"SCRAPE_INTERVAL"},
			Hidden:      false,
			Value:       0 * time.Second,
			Destination: &exporterConfig.scrapeInterval,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "scrape.localnode",
			Usage:       fmt.Sprintf("Whether collect metrics from the whole cluster or local instance"),
			EnvVars:     []string{"SCRAPE_MODE"},
			Hidden:      false,
			Value:       false,
			Destination: &exporterConfig.scrapeLocalNode,
		}),
		// TODO use cli.StringSliceFlag?
		altsrc.NewStringFlag(&cli.StringFlag{
			Name:        "databases",
			Usage:       fmt.Sprintf("Comma separated list of database names, or '%s'", lib.AllDbs),
			EnvVars:     []string{"DATABASES"},
			Hidden:      false,
			Value:       "",
			Destination: &exporterConfig.databases,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "databases.views",
			Usage:       "Collect view details of every observed database",
			EnvVars:     []string{"DATABASES.VIEWS", "DATABASES_VIEWS"},
			Hidden:      false,
			Value:       true,
			Destination: &exporterConfig.databaseViews,
		}),
		altsrc.NewUintFlag(&cli.UintFlag{
			Name:        "database.concurrent.requests",
			Usage:       "maximum concurrent calls to CouchDB, or 0 for unlimited",
			Value:       0,
			Hidden:      false,
			Destination: &exporterConfig.databaseConcurrentRequests,
		}),
		altsrc.NewBoolFlag(&cli.BoolFlag{
			Name:        "scheduler.jobs",
			Usage:       "Collect active replication jobs (CouchDB 2.x+ only)",
			EnvVars:     []string{"SCHEDULER.JOBS", "SCHEDULER_JOBS"},
			Hidden:      false,
			Destination: &exporterConfig.schedulerJobs,
		}),
	}
}

func ofBool(i bool) *bool {
	return &i
}
func ofString(i string) *string {
	return &i
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	var appAction = func(c *cli.Context) error {
		var databases []string
		if exporterConfig.databases != "" {
			databases = strings.Split(exporterConfig.databases, ",")
		}

		exporter := lib.NewExporter(
			exporterConfig.couchdbURI,
			exporterConfig.scrapeLocalNode,
			lib.BasicAuth{
				Username: exporterConfig.couchdbUsername,
				Password: exporterConfig.couchdbPassword},
			lib.CollectorConfig{
				ScrapeInterval:       exporterConfig.scrapeInterval,
				Databases:            databases,
				CollectViews:         exporterConfig.databaseViews,
				CollectSchedulerJobs: exporterConfig.schedulerJobs,
				ConcurrentRequests:   exporterConfig.databaseConcurrentRequests,
			},
			exporterConfig.couchdbInsecure)
		prometheus.MustRegister(exporter)

		http.Handle(webConfig.metricsEndpoint, promhttp.Handler())
		http.HandleFunc("/status", func(w http.ResponseWriter, r *http.Request) {
			_, err := fmt.Fprint(w, "OK")
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
			}
		})
		redirectToMetricsHandler := func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Location", webConfig.metricsEndpoint)
			w.WriteHeader(http.StatusFound)
			_, err := w.Write([]byte(`<html>
			<head><title>CouchDB Prometheus Exporter</title></head>
			<body>
			<h1>CouchDB Prometheus Exporter</h1>
			<p><a href="` + webConfig.metricsEndpoint + `">Metrics</a></p>
			</body>
			</html>`))
			if err != nil {
				slog.Error(fmt.Sprintf("%v", err))
			}
		}
		landingPageHandler, err := web.NewLandingPage(web.LandingConfig{
			Name:        "CouchDB Prometheus Exporter",
			Description: "CouchDB metrics exporter for Prometheus",
			Version:     fmt.Sprintf("%s (%s, %s)", version, commit, date),
			Links: []web.LandingLinks{
				{
					Address:     "./metrics",
					Text:        "Metrics",
					Description: "Metrics, as scraped by Prometheus",
				},
				{
					Address:     "https://github.com/gesellix/couchdb-prometheus-exporter",
					Text:        "Project Homepage",
					Description: "Source repository and primary project home page",
				},
				{
					Address:     "https://couchdb.apache.org/",
					Text:        "CouchDB Homepage",
					Description: "CouchDB home page",
				},
				{
					Address:     "https://prometheus.io/",
					Text:        "Prometheus Homepage",
					Description: "Prometheus home page",
				},
			},
		})
		if err != nil {
			log.Printf("error creating landing page %v\n", err)
			http.HandleFunc("/", redirectToMetricsHandler)
		} else {
			http.HandleFunc("/", landingPageHandler.ServeHTTP)
		}

		slog.Info(fmt.Sprintf("Starting exporter version %s at '%s' to read from CouchDB at '%s'", version, webConfig.listenAddress, exporterConfig.couchdbURI))
		server := &http.Server{Addr: webConfig.listenAddress}
		flags := web.FlagConfig{
			WebListenAddresses: &([]string{webConfig.listenAddress}),
			WebSystemdSocket:   ofBool(false),
			WebConfigFile:      ofString(webConfigFile),
		}
		if err := web.ListenAndServe(server, &flags, logger); err != nil {
			slog.Error("Failed to start the server", "err", err)
			os.Exit(1)
		}
		return nil
	}

	app := cli.NewApp()
	app.Name = "CouchDB Prometheus Exporter"
	//app.Usage = ""
	app.Description = "CouchDB stats exporter for Prometheus"
	app.Version = fmt.Sprintf("%s (%s, %s)", version, commit, date)
	app.Flags = appFlags
	app.Before = beforeApp(appFlags)
	app.Action = appAction

	//defer klog.Flush()

	err := app.Run(os.Args)
	if err != nil {
		slog.Error(fmt.Sprintf("%v", err))
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
		return nil
	}
}
