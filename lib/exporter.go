package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const (
	namespace = "couchdb"
)

type Exporter struct {
	client *CouchdbClient
	mutex  sync.RWMutex

	up prometheus.Gauge

	authCacheHits   *prometheus.GaugeVec
	authCacheMisses *prometheus.GaugeVec
	databaseReads   *prometheus.GaugeVec
	databaseWrites  *prometheus.GaugeVec
	openDatabases   *prometheus.GaugeVec
	openOsFiles     *prometheus.GaugeVec
	requestTime     *prometheus.GaugeVec

	httpdStatusCodes    *prometheus.GaugeVec
	httpdRequestMethods *prometheus.GaugeVec

	clientsRequestingChanges *prometheus.GaugeVec
	temporaryViewReads       *prometheus.GaugeVec
	requests                 *prometheus.GaugeVec
	bulkRequests             *prometheus.GaugeVec
	viewReads                *prometheus.GaugeVec
}

func NewExporter(uri string, basicAuth BasicAuth) *Exporter {

	return &Exporter{
		client: NewCouchdbClient(uri, basicAuth),

		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "up",
				Help:      "Was the last query of CouchDB stats successful.",
			}),

		authCacheHits: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "auth_cache_hits",
				Help:      "number of authentication cache hits",
			},
			[]string{"name"}),
		authCacheMisses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "auth_cache_misses",
				Help:      "number of authentication cache misses",
			},
			[]string{"name"}),
		databaseReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "database_reads",
				Help:      "number of times a document was read from a database",
			},
			[]string{"name"}),
		databaseWrites: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "database_writes",
				Help:      "number of times a database was changed",
			},
			[]string{"name"}),
		openDatabases: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "open_databases",
				Help:      "number of open databases",
			},
			[]string{"name"}),
		openOsFiles: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "open_os_files",
				Help:      "number of file descriptors CouchDB has open",
			},
			[]string{"name"}),
		requestTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "request_time",
				Help:      "length of a request inside CouchDB without MochiWeb",
			},
			[]string{"name"}),

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
			[]string{"method", "node"}),

		bulkRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "bulk_requests",
				Help:      "number of bulk requests",
			},
			[]string{"name"}),
		clientsRequestingChanges: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "clients_requesting_changes",
				Help:      "number of clients for continuous _changes",
			},
			[]string{"name"}),
		requests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "requests",
				Help:      "number of HTTP requests",
			},
			[]string{"name"}),
		temporaryViewReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "temporary_view_reads",
				Help:      "number of temporary view reads",
			},
			[]string{"name"}),
		viewReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "view_reads",
				Help:      "number of view reads",
			},
			[]string{"name"}),
	}
}
