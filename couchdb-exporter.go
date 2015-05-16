package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "couchdb"
)

var (
	listeningAddress = flag.String("telemetry.address", "localhost:9984", "Address on which to expose metrics.")
	metricsEndpoint = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	couchdbURI = flag.String("couchdb.uri", "http://localhost:5984", "URI to the CouchDB instance")
	insecure = flag.Bool("insecure", true, "Ignore server certificate if using https")
)

var (
	exposedHttpStatusCodes = []string{"200", "201", "202", "301", "304", "400", "401", "403", "404", "405", "409", "412", "500"}
)

// Exporter collects couchdb stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	URI    string
	client *http.Client
	mutex  sync.RWMutex

	up prometheus.Gauge

	authCacheHits   prometheus.Gauge
	authCacheMisses prometheus.Gauge
	databaseReads   prometheus.Gauge
	databaseWrites  prometheus.Gauge
	openDatabases   prometheus.Gauge
	openOsFiles     prometheus.Gauge
	requestTime     prometheus.Gauge

	httpdStatusCodes    *prometheus.GaugeVec
	httpdRequestMethods *prometheus.GaugeVec

	clientsRequestingChanges prometheus.Gauge
	temporaryViewReads       prometheus.Gauge
	requests                 prometheus.Gauge
	bulkRequests             prometheus.Gauge
	viewReads                prometheus.Gauge
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri string) *Exporter {
	couchDbClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		},
	}
	return &Exporter{
		URI:    fmt.Sprintf("%s/_stats", uri),
		client: couchDbClient,

		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "up",
				Help:      "Was the last query of CouchDB stats successful.",
			}),

		authCacheHits: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "auth_cache_hits",
				Help:      "number of authentication cache hits",
			}),
		authCacheMisses: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "auth_cache_misses",
				Help:      "number of authentication cache misses",
			}),
		databaseReads: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "database_reads",
				Help:      "number of times a document was read from a database",
			}),
		databaseWrites: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "database_writes",
				Help:      "number of times a database was changed",
			}),
		openDatabases: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "open_databases",
				Help:      "number of open databases",
			}),
		openOsFiles: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "open_os_files",
				Help:      "number of file descriptors CouchDB has open",
			}),
		requestTime: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      "request_time",
				Help:      "length of a request inside CouchDB without MochiWeb",
			}),

		httpdStatusCodes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "status_codes",
				Help:      "number of HTTP responses by status code",
			},
			[]string{"code"}),
		httpdRequestMethods: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "request_methods",
				Help:      "number of HTTP requests by method",
			},
			[]string{"method"}),


		bulkRequests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "bulk_requests",
				Help:      "number of bulk requests",
			}),
		clientsRequestingChanges: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "clients_requesting_changes",
				Help:      "number of clients for continuous _changes",
			}),
		requests: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "requests",
				Help:      "number of HTTP requests",
			}),
		temporaryViewReads: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "temporary_view_reads",
				Help:      "number of temporary view reads",
			}),
		viewReads: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "view_reads",
				Help:      "number of view reads",
			}),
	}
}

// Describe describes all the metrics ever exported by the couchdb exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
	ch <- e.up.Desc()

	ch <- e.authCacheHits.Desc()
	ch <- e.authCacheMisses.Desc()
	ch <- e.databaseReads.Desc()
	ch <- e.databaseWrites.Desc()
	ch <- e.openDatabases.Desc()
	ch <- e.openOsFiles.Desc()
	ch <- e.requestTime.Desc()

	e.httpdStatusCodes.Describe(ch)
	e.httpdRequestMethods.Describe(ch)

	ch <- e.bulkRequests.Desc()
	ch <- e.clientsRequestingChanges.Desc()
	ch <- e.requests.Desc()
	ch <- e.temporaryViewReads.Desc()
	ch <- e.viewReads.Desc()
}

func (e *Exporter) collect(ch chan <- prometheus.Metric) error {
	sendStatus := func() {
		ch <- e.up
	}
	defer sendStatus()

	e.up.Set(0)
	resp, err := e.client.Get(e.URI)
	if err != nil {
		return fmt.Errorf("Error reading couchdb stats: %v", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	e.up.Set(1)

	var stats StatsResponse
	err = json.Unmarshal(data, &stats)
	//	glog.Info(fmt.Sprintf("stats: %v\n", stats))

	e.authCacheHits.Set(stats.Couchdb.AuthCacheHits.Current)
	ch <- e.authCacheHits
	e.authCacheMisses.Set(stats.Couchdb.AuthCacheMisses.Current)
	ch <- e.authCacheMisses
	e.databaseReads.Set(stats.Couchdb.DatabaseReads.Current)
	ch <- e.databaseReads
	e.databaseWrites.Set(stats.Couchdb.DatabaseWrites.Current)
	ch <- e.databaseWrites
	e.openDatabases.Set(stats.Couchdb.OpenDatabases.Current)
	ch <- e.openDatabases
	e.openOsFiles.Set(stats.Couchdb.OpenOsFiles.Current)
	ch <- e.openOsFiles
	e.requestTime.Set(stats.Couchdb.RequestTime.Current)
	ch <- e.requestTime

	for _, code := range exposedHttpStatusCodes {
		if _, ok := stats.HttpdStatusCodes[code]; ok {
			e.httpdStatusCodes.WithLabelValues(code).Set(stats.HttpdStatusCodes[code].Current)
		}
	}
	e.httpdStatusCodes.Collect(ch)

	e.httpdRequestMethods.WithLabelValues("COPY").Set(stats.HttpdRequestMethods.COPY.Current)
	e.httpdRequestMethods.WithLabelValues("DELETE").Set(stats.HttpdRequestMethods.DELETE.Current)
	e.httpdRequestMethods.WithLabelValues("GET").Set(stats.HttpdRequestMethods.GET.Current)
	e.httpdRequestMethods.WithLabelValues("HEAD").Set(stats.HttpdRequestMethods.HEAD.Current)
	e.httpdRequestMethods.WithLabelValues("POST").Set(stats.HttpdRequestMethods.POST.Current)
	e.httpdRequestMethods.WithLabelValues("PUT").Set(stats.HttpdRequestMethods.PUT.Current)
	e.httpdRequestMethods.Collect(ch)

	e.bulkRequests.Set(stats.Httpd.BulkRequests.Current)
	ch <- e.bulkRequests
	e.clientsRequestingChanges.Set(stats.Httpd.ClientsRequestingChanges.Current)
	ch <- e.clientsRequestingChanges
	e.requests.Set(stats.Httpd.Requests.Current)
	ch <- e.requests
	e.temporaryViewReads.Set(stats.Httpd.TemporaryViewReads.Current)
	ch <- e.temporaryViewReads
	e.viewReads.Set(stats.Httpd.ViewReads.Current)
	ch <- e.viewReads

	return nil
}

// Collect fetches the stats from configured couchdb location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		glog.Error(fmt.Sprintf("Error collecting stats: %s", err))
	}
	return
}

func main() {
	flag.Parse()

	exporter := NewExporter(*couchdbURI)
	prometheus.MustRegister(exporter)

	log.Printf("Starting exporter at %s to read from CouchDB at %s", *listeningAddress, *couchdbURI)
	http.Handle(*metricsEndpoint, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
