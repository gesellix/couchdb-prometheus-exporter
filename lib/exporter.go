package lib

import (
	"sync"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "couchdb"
)

type Exporter struct {
	client    *CouchdbClient
	mutex     sync.RWMutex

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

func NewExporter(uri string, basicAuth BasicAuth) *Exporter {

	return &Exporter{
		client: NewCouchdbClient(uri, basicAuth),

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
