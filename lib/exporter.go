package lib

import (
	"sync"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "couchdb"
)

type Exporter struct {
	client          *CouchdbClient
	collectorConfig CollectorConfig
	mutex           sync.RWMutex

	requestCount prometheus.Gauge

	up             prometheus.Gauge
	databasesTotal prometheus.Gauge
	nodeUp         *prometheus.GaugeVec
	nodeInfo       *prometheus.GaugeVec

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

	dbInfo           *prometheus.GaugeVec
	diskSize         *prometheus.GaugeVec
	dataSize         *prometheus.GaugeVec
	docCount         *prometheus.GaugeVec
	docDelCount      *prometheus.GaugeVec
	compactRunning   *prometheus.GaugeVec
	diskSizeOverhead *prometheus.GaugeVec

	activeTasks                      *prometheus.GaugeVec
	activeTasksDatabaseCompaction    *prometheus.GaugeVec
	activeTasksViewCompaction        *prometheus.GaugeVec
	activeTasksIndexer               *prometheus.GaugeVec
	activeTasksReplication           *prometheus.GaugeVec
	activeTasksReplicationLastUpdate *prometheus.GaugeVec

	couchLog *prometheus.GaugeVec

	fabricWorker      *prometheus.GaugeVec
	fabricOpenShard   *prometheus.GaugeVec
	fabricReadRepairs *prometheus.GaugeVec
	fabricDocUpdate   *prometheus.GaugeVec

	couchReplicatorChangesReadFailures  *prometheus.GaugeVec
	couchReplicatorChangesReaderDeaths  *prometheus.GaugeVec
	couchReplicatorChangesManagerDeaths *prometheus.GaugeVec
	couchReplicatorChangesQueueDeaths   *prometheus.GaugeVec
	couchReplicatorCheckpoints          *prometheus.GaugeVec
	couchReplicatorFailedStarts         *prometheus.GaugeVec
	couchReplicatorRequests             *prometheus.GaugeVec
	couchReplicatorResponses            *prometheus.GaugeVec
	couchReplicatorStreamResponses      *prometheus.GaugeVec
	couchReplicatorWorkerDeaths         *prometheus.GaugeVec
	couchReplicatorWorkersStarted       *prometheus.GaugeVec
	couchReplicatorClusterIsStable      *prometheus.GaugeVec
	couchReplicatorDbScans              *prometheus.GaugeVec
	couchReplicatorDocs                 *prometheus.GaugeVec
	couchReplicatorJobs                 *prometheus.GaugeVec
	couchReplicatorConnection           *prometheus.GaugeVec

	nodeMemoryOther         *prometheus.GaugeVec
	nodeMemoryAtom          *prometheus.GaugeVec
	nodeMemoryAtomUsed      *prometheus.GaugeVec
	nodeMemoryProcesses     *prometheus.GaugeVec
	nodeMemoryProcessesUsed *prometheus.GaugeVec
	nodeMemoryBinary        *prometheus.GaugeVec
	nodeMemoryCode          *prometheus.GaugeVec
	nodeMemoryEts           *prometheus.GaugeVec

	viewStaleness *prometheus.GaugeVec

	schedulerJobs *prometheus.GaugeVec
}

func NewExporter(uri string, basicAuth BasicAuth, collectorConfig CollectorConfig, insecure bool) *Exporter {

	return &Exporter{
		client:          NewCouchdbClient(uri, basicAuth, insecure),
		collectorConfig: collectorConfig,

		requestCount: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "exporter",
				Name:      "request_count",
				Help:      "Number of CouchDB requests for this scrape.",
			}),

		up: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "up",
				Help:      "Was the last query of CouchDB stats successful.",
			}),

		databasesTotal: prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "databases_total",
				Help:      "Total number of databases in the cluster",
			}),
		nodeUp: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "httpd",
				Name:      "node_up",
				Help:      "Is the node available.",
			},
			[]string{"node_name"}),
		nodeInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "node_info",
				Help:      "General info about a node.",
			},
			[]string{"node_name", "version", "vendor_name"}),

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
			[]string{"node_name", "metric"}),

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

		dbInfo: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "info",
				Help:      "General info about a database.",
			},
			[]string{"db_name", "disk_format_version", "partitioned"}),
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
		docCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "doc_count",
				Help:      "document count",
			},
			[]string{"db_name"}),
		docDelCount: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "doc_del_count",
				Help:      "deleted document count",
			},
			[]string{"db_name"}),
		compactRunning: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "database",
				Name:      "compact_running",
				Help:      "database compaction running",
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

		couchLog: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "server",
				Name:      "couch_log",
				Help:      "number of messages logged by log level",
			},
			[]string{"level", "node_name"}),

		nodeMemoryOther: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_other",
				Help:      "erlang memory counters - other",
			},
			[]string{"node_name"}),

		nodeMemoryAtom: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_atom",
				Help:      "erlang memory counters - atom",
			},
			[]string{"node_name"}),

		nodeMemoryAtomUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_atom_used",
				Help:      "erlang memory counters - atom_used",
			},
			[]string{"node_name"}),

		nodeMemoryProcesses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_processes",
				Help:      "erlang memory counters - processes",
			},
			[]string{"node_name"}),

		nodeMemoryProcessesUsed: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_processes_used",
				Help:      "erlang memory counters - processes_used",
			},
			[]string{"node_name"}),

		nodeMemoryBinary: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_binary",
				Help:      "erlang memory counters - binary",
			},
			[]string{"node_name"}),

		nodeMemoryCode: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_code",
				Help:      "erlang memory counters - code",
			},
			[]string{"node_name"}),

		nodeMemoryEts: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "erlang",
				Name:      "memory_ets",
				Help:      "erlang memory counters - ets",
			},
			[]string{"node_name"}),

		viewStaleness: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "view",
				Name:      "staleness",
				Help:      "the view's staleness (the view's update_seq compared to the database's update_seq)",
			},
			[]string{"db_name", "design_doc_name", "view_name", "shard_begin", "shard_end"}),

		schedulerJobs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "scheduler",
				Name:      "jobs",
				Help:      "scheduler jobs",
			},
			[]string{"node_name", "job_id", "db_name", "doc_id", "source", "target"}),

		fabricWorker: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "fabric",
				Name:      "worker",
				Help:      "worker metrics",
			},
			[]string{"metric", "node_name"}),

		fabricOpenShard: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "fabric",
				Name:      "open_shard",
				Help:      "open_shard metrics",
			},
			[]string{"metric", "node_name"}),

		fabricReadRepairs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "fabric",
				Name:      "read_repairs",
				Help:      "read repair metrics",
			},
			[]string{"metric", "node_name"}),

		fabricDocUpdate: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "fabric",
				Name:      "doc_update",
				Help:      "doc update metrics",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorChangesReadFailures: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "changes_read_failures",
				Help:      "number of failed replicator changes read failures",
			},
			[]string{"node_name"}),

		couchReplicatorChangesReaderDeaths: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "changes_reader_deaths",
				Help:      "number of failed replicator changes readers",
			},
			[]string{"node_name"}),

		couchReplicatorChangesManagerDeaths: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "changes_manager_deaths",
				Help:      "number of failed replicator changes managers",
			},
			[]string{"node_name"}),

		couchReplicatorChangesQueueDeaths: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "changes_queue_deaths",
				Help:      "number of failed replicator changes work queues",
			},
			[]string{"node_name"}),

		couchReplicatorCheckpoints: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "checkpoints",
				Help:      "replicator checkpoint counters",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorFailedStarts: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "failed_starts",
				Help:      "number of replications that have failed to start",
			},
			[]string{"node_name"}),

		couchReplicatorRequests: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "requests",
				Help:      "number of HTTP requests made by the replicator",
			},
			[]string{"node_name"}),

		couchReplicatorResponses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "responses",
				Help:      "number of HTTP responses by state",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorStreamResponses: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "stream_responses",
				Help:      "number of streaming HTTP responses by state",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorWorkerDeaths: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "worker_deaths",
				Help:      "number of failed replicator workers",
			},
			[]string{"node_name"}),

		couchReplicatorWorkersStarted: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "workers_started",
				Help:      "number of replicator workers started",
			},
			[]string{"node_name"}),

		couchReplicatorClusterIsStable: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "cluster_is_stable",
				Help:      "1 if cluster is stable, 0 if unstable",
			},
			[]string{"node_name"}),

		couchReplicatorDbScans: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "db_scans",
				Help:      "number of times replicator db scans have been started",
			},
			[]string{"node_name"}),

		couchReplicatorDocs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "docs",
				Help:      "replicator metrics shown by type",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorJobs: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "jobs",
				Help:      "replicator jobs shown by type",
			},
			[]string{"metric", "node_name"}),

		couchReplicatorConnection: prometheus.NewGaugeVec(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Subsystem: "replicator",
				Name:      "connections",
				Help:      "replicator connection metrics shown by type",
			},
			[]string{"metric", "node_name"}),
	}
}
