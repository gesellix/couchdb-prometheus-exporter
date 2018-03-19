package lib

import (
	"github.com/prometheus/client_golang/prometheus"
	"sync"
)

const (
	namespace = "couchdb"
)

type Exporter struct {
	client    *CouchdbClient
	databases []string
	mutex     sync.RWMutex

	up     prometheus.Gauge
	nodeUp *prometheus.GaugeVec

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

	diskSize         *prometheus.GaugeVec
	dataSize         *prometheus.GaugeVec
	diskSizeOverhead *prometheus.GaugeVec

	activeTasks                      *prometheus.GaugeVec
	activeTasksDatabaseCompaction    *prometheus.GaugeVec
	activeTasksViewCompaction        *prometheus.GaugeVec
	activeTasksIndexer               *prometheus.GaugeVec
	activeTasksReplication           *prometheus.GaugeVec
	activeTasksReplicationLastUpdate *prometheus.GaugeVec
}

func NewExporter(uri string, basicAuth BasicAuth, databases []string, insecure bool) *Exporter {

	return &Exporter{
		client:    NewCouchdbClient(uri, basicAuth, databases, insecure),
		databases: databases,

		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "up",
				Help:      "Was the last query of CouchDB stats successful.",
			}),
		nodeUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "node_up",
				Help:      "Is the node available.",
			},
			[]string{"node_name"}),

		authCacheHits: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "auth_cache_hits",
				Help:      "number of authentication cache hits",
			},
			[]string{"node_name"}),
		authCacheMisses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "auth_cache_misses",
				Help:      "number of authentication cache misses",
			},
			[]string{"node_name"}),
		databaseReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "database_reads",
				Help:      "number of times a document was read from a database",
			},
			[]string{"node_name"}),
		databaseWrites: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "database_writes",
				Help:      "number of times a database was changed",
			},
			[]string{"node_name"}),
		openDatabases: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "open_databases",
				Help:      "number of open databases",
			},
			[]string{"node_name"}),
		openOsFiles: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "open_os_files",
				Help:      "number of file descriptors CouchDB has open",
			},
			[]string{"node_name"}),
		requestTime: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "request_time",
				Help:      "length of a request inside CouchDB without MochiWeb",
			},
			[]string{"node_name"}),

		httpdStatusCodes: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "status_codes",
				Help:      "number of HTTP responses by status code",
			},
			[]string{"code", "node_name"}),
		httpdRequestMethods: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "request_methods",
				Help:      "number of HTTP requests by method",
			},
			[]string{"method", "node_name"}),

		bulkRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "bulk_requests",
				Help:      "number of bulk requests",
			},
			[]string{"node_name"}),
		clientsRequestingChanges: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "clients_requesting_changes",
				Help:      "number of clients for continuous _changes",
			},
			[]string{"node_name"}),
		requests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "requests",
				Help:      "number of HTTP requests",
			},
			[]string{"node_name"}),
		temporaryViewReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "temporary_view_reads",
				Help:      "number of temporary view reads",
			},
			[]string{"node_name"}),
		viewReads: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "view_reads",
				Help:      "number of view reads",
			},
			[]string{"node_name"}),

		diskSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "disk_size",
				Help:      "disk size",
			},
			[]string{"db_name"}),
		dataSize: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "data_size",
				Help:      "data size",
			},
			[]string{"db_name"}),
		diskSizeOverhead: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "overhead",
				Help:      "disk size overhead",
			},
			[]string{"db_name"}),

		activeTasks: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks",
				Help:      "active tasks",
			},
			[]string{"node_name"}),

		activeTasksDatabaseCompaction: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks_database_compaction",
				Help:      "active tasks",
			},
			[]string{"node_name"}),
		activeTasksViewCompaction: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks_view_compaction",
				Help:      "active tasks",
			},
			[]string{"node_name"}),
		activeTasksIndexer: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks_indexer",
				Help:      "active tasks",
			},
			[]string{"node_name"}),
		activeTasksReplication: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks_replication",
				Help:      "active tasks",
			},
			[]string{"node_name"}),

		activeTasksReplicationLastUpdate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "active_tasks_replication_updated_on",
				Help:      "active tasks",
			},
			[]string{"node_name", "doc_id", "continuous", "source", "target"}),
	}
}
