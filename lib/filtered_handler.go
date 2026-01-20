package lib

import (
	"fmt"
	"log/slog"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// CollectorGroup defines groups of metrics that can be collected independently
type CollectorGroup string

const (
	// CollectorGroupStandard includes basic stats without per-database metrics
	CollectorGroupStandard CollectorGroup = "standard"
	// CollectorGroupDatabases includes per-database metrics (heavy operation)
	CollectorGroupDatabases CollectorGroup = "databases"
	// CollectorGroupViews includes view staleness metrics (heavy operation)
	CollectorGroupViews CollectorGroup = "views"
	// CollectorGroupScheduler includes scheduler jobs metrics
	CollectorGroupScheduler CollectorGroup = "scheduler"
)

// FilteredExporter wraps the standard Exporter and provides methods to
// selectively register metrics based on collector groups
type FilteredExporter struct {
	*Exporter
}

// NewFilteredExporter creates a new FilteredExporter
func NewFilteredExporter(uri string, localOnly bool, basicAuth BasicAuth, collectorConfig CollectorConfig, insecure bool) *FilteredExporter {
	// Create the base exporter but don't start auto-scraping
	// since we'll be using per-request registries
	baseExporter := &Exporter{
		client:          NewCouchdbClient(uri, localOnly, basicAuth, insecure),
		collectorConfig: collectorConfig,
		requestCount:    createRequestCountMetric(),
		up:              createUpMetric(),
		databasesTotal:  createDatabasesTotalMetric(),
		nodeUp:          createNodeUpMetric(),
		nodeInfo:        createNodeInfoMetric(),

		authCacheHits:   createAuthCacheHitsMetric(),
		authCacheMisses: createAuthCacheMissesMetric(),
		databaseReads:   createDatabaseReadsMetric(),
		databaseWrites:  createDatabaseWritesMetric(),
		openDatabases:   createOpenDatabasesMetric(),
		openOsFiles:     createOpenOsFilesMetric(),
		requestTime:     createRequestTimeMetric(),

		httpdStatusCodes:    createHttpdStatusCodesMetric(),
		httpdRequestMethods: createHttpdRequestMethodsMetric(),

		clientsRequestingChanges: createClientsRequestingChangesMetric(),
		temporaryViewReads:       createTemporaryViewReadsMetric(),
		requests:                 createRequestsMetric(),
		bulkRequests:             createBulkRequestsMetric(),
		viewReads:                createViewReadsMetric(),

		dbInfo:           createDbInfoMetric(),
		diskSize:         createDiskSizeMetric(),
		dataSize:         createDataSizeMetric(),
		docCount:         createDocCountMetric(),
		docDelCount:      createDocDelCountMetric(),
		compactRunning:   createCompactRunningMetric(),
		diskSizeOverhead: createDiskSizeOverheadMetric(),

		activeTasks:                          createActiveTasksMetric(),
		activeTasksDatabaseCompaction:        createActiveTasksDatabaseCompactionMetric(),
		activeTasksViewCompaction:            createActiveTasksViewCompactionMetric(),
		activeTasksIndexer:                   createActiveTasksIndexerMetric(),
		activeTasksReplication:               createActiveTasksReplicationMetric(),
		activeTasksReplicationLastUpdate:     createActiveTasksReplicationLastUpdateMetric(),
		activeTasksReplicationChangesPending: createActiveTasksReplicationChangesPendingMetric(),

		couchLog: createCouchLogMetric(),

		fabricWorker:      createFabricWorkerMetric(),
		fabricOpenShard:   createFabricOpenShardMetric(),
		fabricReadRepairs: createFabricReadRepairsMetric(),
		fabricDocUpdate:   createFabricDocUpdateMetric(),

		couchReplicatorChangesReadFailures:  createCouchReplicatorChangesReadFailuresMetric(),
		couchReplicatorChangesReaderDeaths:  createCouchReplicatorChangesReaderDeathsMetric(),
		couchReplicatorChangesManagerDeaths: createCouchReplicatorChangesManagerDeathsMetric(),
		couchReplicatorChangesQueueDeaths:   createCouchReplicatorChangesQueueDeathsMetric(),
		couchReplicatorCheckpoints:          createCouchReplicatorCheckpointsMetric(),
		couchReplicatorFailedStarts:         createCouchReplicatorFailedStartsMetric(),
		couchReplicatorRequests:             createCouchReplicatorRequestsMetric(),
		couchReplicatorResponses:            createCouchReplicatorResponsesMetric(),
		couchReplicatorStreamResponses:      createCouchReplicatorStreamResponsesMetric(),
		couchReplicatorWorkerDeaths:         createCouchReplicatorWorkerDeathsMetric(),
		couchReplicatorWorkersStarted:       createCouchReplicatorWorkersStartedMetric(),
		couchReplicatorClusterIsStable:      createCouchReplicatorClusterIsStableMetric(),
		couchReplicatorDbScans:              createCouchReplicatorDbScansMetric(),
		couchReplicatorDocs:                 createCouchReplicatorDocsMetric(),
		couchReplicatorJobs:                 createCouchReplicatorJobsMetric(),
		couchReplicatorConnection:           createCouchReplicatorConnectionMetric(),

		nodeMemoryOther:         createNodeMemoryOtherMetric(),
		nodeMemoryAtom:          createNodeMemoryAtomMetric(),
		nodeMemoryAtomUsed:      createNodeMemoryAtomUsedMetric(),
		nodeMemoryProcesses:     createNodeMemoryProcessesMetric(),
		nodeMemoryProcessesUsed: createNodeMemoryProcessesUsedMetric(),
		nodeMemoryBinary:        createNodeMemoryBinaryMetric(),
		nodeMemoryCode:          createNodeMemoryCodeMetric(),
		nodeMemoryEts:           createNodeMemoryEtsMetric(),

		mangoUnindexedQueries:   createMangoUnindexedQueriesMetric(),
		mangoInvalidIndexes:     createMangoInvalidIndexesMetric(),
		mangoTooManyDocs:        createMangoTooManyDocsMetric(),
		mangoDocsExamined:       createMangoDocsExaminedMetric(),
		mangoQuorumDocsExamined: createMangoQuorumDocsExaminedMetric(),
		mangoResultsReturned:    createMangoResultsReturnedMetric(),
		mangoQueryTime:          createMangoQueryTimeMetric(),
		mangoEvaluateSelectors:  createMangoEvaluateSelectorsMetric(),

		viewStaleness: createViewStalenessMetric(),
		schedulerJobs: createSchedulerJobsMetric(),
	}

	return &FilteredExporter{Exporter: baseExporter}
}

// RegisterStandardMetrics registers the lightweight standard metrics
func (e *FilteredExporter) RegisterStandardMetrics(registry *prometheus.Registry) {
	// Exporter meta-metrics
	registry.MustRegister(e.requestCount)
	registry.MustRegister(e.up)
	registry.MustRegister(e.databasesTotal)
	registry.MustRegister(e.nodeUp)
	registry.MustRegister(e.nodeInfo)

	// HTTP stats
	registry.MustRegister(e.authCacheHits)
	registry.MustRegister(e.authCacheMisses)
	registry.MustRegister(e.databaseReads)
	registry.MustRegister(e.databaseWrites)
	registry.MustRegister(e.openDatabases)
	registry.MustRegister(e.openOsFiles)
	registry.MustRegister(e.requestTime)
	registry.MustRegister(e.httpdStatusCodes)
	registry.MustRegister(e.httpdRequestMethods)
	registry.MustRegister(e.bulkRequests)
	registry.MustRegister(e.clientsRequestingChanges)
	registry.MustRegister(e.requests)
	registry.MustRegister(e.temporaryViewReads)
	registry.MustRegister(e.viewReads)

	// Active tasks
	registry.MustRegister(e.activeTasks)
	registry.MustRegister(e.activeTasksDatabaseCompaction)
	registry.MustRegister(e.activeTasksViewCompaction)
	registry.MustRegister(e.activeTasksIndexer)
	registry.MustRegister(e.activeTasksReplication)
	registry.MustRegister(e.activeTasksReplicationLastUpdate)
	registry.MustRegister(e.activeTasksReplicationChangesPending)

	// Logs
	registry.MustRegister(e.couchLog)

	// Fabric metrics
	registry.MustRegister(e.fabricWorker)
	registry.MustRegister(e.fabricOpenShard)
	registry.MustRegister(e.fabricReadRepairs)
	registry.MustRegister(e.fabricDocUpdate)

	// Replicator metrics
	registry.MustRegister(e.couchReplicatorChangesReadFailures)
	registry.MustRegister(e.couchReplicatorChangesReaderDeaths)
	registry.MustRegister(e.couchReplicatorChangesManagerDeaths)
	registry.MustRegister(e.couchReplicatorChangesQueueDeaths)
	registry.MustRegister(e.couchReplicatorCheckpoints)
	registry.MustRegister(e.couchReplicatorFailedStarts)
	registry.MustRegister(e.couchReplicatorRequests)
	registry.MustRegister(e.couchReplicatorResponses)
	registry.MustRegister(e.couchReplicatorStreamResponses)
	registry.MustRegister(e.couchReplicatorWorkerDeaths)
	registry.MustRegister(e.couchReplicatorWorkersStarted)
	registry.MustRegister(e.couchReplicatorClusterIsStable)
	registry.MustRegister(e.couchReplicatorDbScans)
	registry.MustRegister(e.couchReplicatorDocs)
	registry.MustRegister(e.couchReplicatorJobs)
	registry.MustRegister(e.couchReplicatorConnection)

	// Memory metrics
	registry.MustRegister(e.nodeMemoryOther)
	registry.MustRegister(e.nodeMemoryAtom)
	registry.MustRegister(e.nodeMemoryAtomUsed)
	registry.MustRegister(e.nodeMemoryProcesses)
	registry.MustRegister(e.nodeMemoryProcessesUsed)
	registry.MustRegister(e.nodeMemoryBinary)
	registry.MustRegister(e.nodeMemoryCode)
	registry.MustRegister(e.nodeMemoryEts)

	// Mango metrics
	registry.MustRegister(e.mangoUnindexedQueries)
	registry.MustRegister(e.mangoInvalidIndexes)
	registry.MustRegister(e.mangoTooManyDocs)
	registry.MustRegister(e.mangoDocsExamined)
	registry.MustRegister(e.mangoQuorumDocsExamined)
	registry.MustRegister(e.mangoResultsReturned)
	registry.MustRegister(e.mangoQueryTime)
	registry.MustRegister(e.mangoEvaluateSelectors)
}

// RegisterAllDbsMetrics registers per-database metrics (heavy operation)
func (e *FilteredExporter) RegisterAllDbsMetrics(registry *prometheus.Registry) {
	registry.MustRegister(e.dbInfo)
	registry.MustRegister(e.diskSize)
	registry.MustRegister(e.dataSize)
	registry.MustRegister(e.docCount)
	registry.MustRegister(e.docDelCount)
	registry.MustRegister(e.compactRunning)
	registry.MustRegister(e.diskSizeOverhead)
}

// RegisterViewsMetrics registers view staleness metrics (heavy operation)
func (e *FilteredExporter) RegisterViewsMetrics(registry *prometheus.Registry) {
	registry.MustRegister(e.viewStaleness)
}

// RegisterSchedulerMetrics registers scheduler jobs metrics
func (e *FilteredExporter) RegisterSchedulerMetrics(registry *prometheus.Registry) {
	registry.MustRegister(e.schedulerJobs)
}

// CreateFilteredHandler returns an HTTP handler that supports collect[] parameters
// for selective metric collection
func CreateFilteredHandler(exporter *FilteredExporter) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// Parse collect[] parameters from query string
		collectParams := r.URL.Query()["collect[]"]
		
		// Create a new registry for this specific request
		registry := prometheus.NewRegistry()
		
		// Determine which collectors to enable
		groups := parseCollectorGroups(collectParams)
		
		// Log the requested collector groups
		if len(groups) > 0 {
			slog.Info(fmt.Sprintf("Scrape requested with filters: %v", groups))
		} else {
			slog.Debug("Scrape requested with default (standard) collectors")
		}
		
		// Register collectors based on requested groups
		// If no groups specified, use default (standard only)
		if len(groups) == 0 {
			exporter.RegisterStandardMetrics(registry)
		} else {
			for group := range groups {
				switch group {
				case CollectorGroupStandard:
					exporter.RegisterStandardMetrics(registry)
				case CollectorGroupDatabases:
					exporter.RegisterAllDbsMetrics(registry)
				case CollectorGroupViews:
					exporter.RegisterViewsMetrics(registry)
				case CollectorGroupScheduler:
					exporter.RegisterSchedulerMetrics(registry)
				default:
					slog.Warn(fmt.Sprintf("Unknown collector group: %s", group))
				}
			}
		}
		
		// Trigger a scrape to populate the metrics
		// The metrics are already registered, now we need to collect data
		// Lock the mutex to prevent concurrent scrape operations
		exporter.Exporter.mutex.Lock()
		err := exporter.Exporter.scrape()
		exporter.Exporter.mutex.Unlock()
		if err != nil {
			slog.Warn("Error during scrape", "error", err)
		}
		
		// Create a handler for this specific registry
		handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
			ErrorLog:      slog.NewLogLogger(slog.Default().Handler(), slog.LevelError),
			ErrorHandling: promhttp.ContinueOnError,
		})
		
		// Serve the metrics
		handler.ServeHTTP(w, r)
	}
}

// parseCollectorGroups converts the collect[] query parameters into a set of CollectorGroups
func parseCollectorGroups(params []string) map[CollectorGroup]struct{} {
	groups := make(map[CollectorGroup]struct{})
	
	for _, param := range params {
		param = strings.TrimSpace(strings.ToLower(param))
		switch param {
		case "standard":
			groups[CollectorGroupStandard] = struct{}{}
		case "databases":
			groups[CollectorGroupDatabases] = struct{}{}
		case "views":
			groups[CollectorGroupViews] = struct{}{}
		case "scheduler":
			groups[CollectorGroupScheduler] = struct{}{}
		case "":
			// Ignore empty parameters
		default:
			slog.Warn(fmt.Sprintf("Unknown collector parameter: %s", param))
		}
	}
	
	return groups
}

// Helper functions to create metrics (to avoid duplication and allow reuse)
// These replace the inline metric creation in NewExporter

func createRequestCountMetric() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "exporter",
		Name:      "request_count",
		Help:      "Number of CouchDB requests for this scrape.",
	})
}

func createUpMetric() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "up",
		Help:      "Was the last query of CouchDB stats successful.",
	})
}

func createDatabasesTotalMetric() prometheus.Gauge {
	return prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "databases_total",
		Help:      "Total number of databases in the cluster",
	})
}

func createNodeUpMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "node_up",
		Help:      "Is the node available.",
	}, []string{"node_name"})
}

func createNodeInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "node_info",
		Help:      "General info about a node.",
	}, []string{"node_name", "version", "vendor_name"})
}

func createAuthCacheHitsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "auth_cache_hits",
		Help:      "number of authentication cache hits",
	}, []string{"node_name"})
}

func createAuthCacheMissesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "auth_cache_misses",
		Help:      "number of authentication cache misses",
	}, []string{"node_name"})
}

func createDatabaseReadsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "database_reads",
		Help:      "number of times a document was read from a database",
	}, []string{"node_name"})
}

func createDatabaseWritesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "database_writes",
		Help:      "number of times a database was changed",
	}, []string{"node_name"})
}

func createOpenDatabasesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "open_databases",
		Help:      "number of open databases",
	}, []string{"node_name"})
}

func createOpenOsFilesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "open_os_files",
		Help:      "number of file descriptors CouchDB has open",
	}, []string{"node_name"})
}

func createRequestTimeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "request_time",
		Help:      "length of a request inside CouchDB without MochiWeb",
	}, []string{"node_name", "metric"})
}

func createHttpdStatusCodesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "status_codes",
		Help:      "number of HTTP responses by status code",
	}, []string{"code", "node_name"})
}

func createHttpdRequestMethodsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "request_methods",
		Help:      "number of HTTP requests by method",
	}, []string{"method", "node_name"})
}

func createBulkRequestsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "bulk_requests",
		Help:      "number of bulk requests",
	}, []string{"node_name"})
}

func createClientsRequestingChangesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "clients_requesting_changes",
		Help:      "number of clients for continuous _changes",
	}, []string{"node_name"})
}

func createRequestsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "requests",
		Help:      "number of HTTP requests",
	}, []string{"node_name"})
}

func createTemporaryViewReadsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "temporary_view_reads",
		Help:      "number of temporary view reads",
	}, []string{"node_name"})
}

func createViewReadsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "httpd",
		Name:      "view_reads",
		Help:      "number of view reads",
	}, []string{"node_name"})
}

func createDbInfoMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "info",
		Help:      "General info about a database.",
	}, []string{"db_name", "disk_format_version", "partitioned"})
}

func createDiskSizeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "disk_size",
		Help:      "disk size",
	}, []string{"db_name"})
}

func createDataSizeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "data_size",
		Help:      "data size",
	}, []string{"db_name"})
}

func createDocCountMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "doc_count",
		Help:      "document count",
	}, []string{"db_name"})
}

func createDocDelCountMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "doc_del_count",
		Help:      "deleted document count",
	}, []string{"db_name"})
}

func createCompactRunningMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "compact_running",
		Help:      "database compaction running",
	}, []string{"db_name"})
}

func createDiskSizeOverheadMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "database",
		Name:      "overhead",
		Help:      "disk size overhead",
	}, []string{"db_name"})
}

func createActiveTasksMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks",
		Help:      "active tasks",
	}, []string{"node_name"})
}

func createActiveTasksDatabaseCompactionMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_database_compaction",
		Help:      "active tasks database compaction",
	}, []string{"node_name"})
}

func createActiveTasksViewCompactionMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_view_compaction",
		Help:      "active tasks view compaction",
	}, []string{"node_name"})
}

func createActiveTasksIndexerMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_indexer",
		Help:      "active tasks indexer",
	}, []string{"node_name"})
}

func createActiveTasksReplicationMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_replication",
		Help:      "active tasks replication",
	}, []string{"node_name"})
}

func createActiveTasksReplicationLastUpdateMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_replication_updated_on",
		Help:      "active tasks replication updated on",
	}, []string{"node_name", "doc_id", "continuous", "source", "target"})
}

func createActiveTasksReplicationChangesPendingMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "active_tasks_replication_changes_pending",
		Help:      "active tasks replication changes pending ",
	}, []string{"node_name", "doc_id", "continuous", "source", "target"})
}

func createCouchLogMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "server",
		Name:      "couch_log",
		Help:      "number of messages logged by log level",
	}, []string{"level", "node_name"})
}

func createFabricWorkerMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "fabric",
		Name:      "worker",
		Help:      "worker metrics",
	}, []string{"metric", "node_name"})
}

func createFabricOpenShardMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "fabric",
		Name:      "open_shard",
		Help:      "open_shard metrics",
	}, []string{"metric", "node_name"})
}

func createFabricReadRepairsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "fabric",
		Name:      "read_repairs",
		Help:      "read repair metrics",
	}, []string{"metric", "node_name"})
}

func createFabricDocUpdateMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "fabric",
		Name:      "doc_update",
		Help:      "doc update metrics",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorChangesReadFailuresMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "changes_read_failures",
		Help:      "number of failed replicator changes read failures",
	}, []string{"node_name"})
}

func createCouchReplicatorChangesReaderDeathsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "changes_reader_deaths",
		Help:      "number of failed replicator changes readers",
	}, []string{"node_name"})
}

func createCouchReplicatorChangesManagerDeathsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "changes_manager_deaths",
		Help:      "number of failed replicator changes managers",
	}, []string{"node_name"})
}

func createCouchReplicatorChangesQueueDeathsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "changes_queue_deaths",
		Help:      "number of failed replicator changes work queues",
	}, []string{"node_name"})
}

func createCouchReplicatorCheckpointsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "checkpoints",
		Help:      "replicator checkpoint counters",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorFailedStartsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "failed_starts",
		Help:      "number of replications that have failed to start",
	}, []string{"node_name"})
}

func createCouchReplicatorRequestsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "requests",
		Help:      "number of HTTP requests made by the replicator",
	}, []string{"node_name"})
}

func createCouchReplicatorResponsesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "responses",
		Help:      "number of HTTP responses by state",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorStreamResponsesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "stream_responses",
		Help:      "number of streaming HTTP responses by state",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorWorkerDeathsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "worker_deaths",
		Help:      "number of failed replicator workers",
	}, []string{"node_name"})
}

func createCouchReplicatorWorkersStartedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "workers_started",
		Help:      "number of replicator workers started",
	}, []string{"node_name"})
}

func createCouchReplicatorClusterIsStableMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "cluster_is_stable",
		Help:      "1 if cluster is stable, 0 if unstable",
	}, []string{"node_name"})
}

func createCouchReplicatorDbScansMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "db_scans",
		Help:      "number of times replicator db scans have been started",
	}, []string{"node_name"})
}

func createCouchReplicatorDocsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "docs",
		Help:      "replicator metrics shown by type",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorJobsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "jobs",
		Help:      "replicator jobs shown by type",
	}, []string{"metric", "node_name"})
}

func createCouchReplicatorConnectionMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "replicator",
		Name:      "connections",
		Help:      "replicator connection metrics shown by type",
	}, []string{"metric", "node_name"})
}

func createNodeMemoryOtherMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_other",
		Help:      "erlang memory counters - other",
	}, []string{"node_name"})
}

func createNodeMemoryAtomMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_atom",
		Help:      "erlang memory counters - atom",
	}, []string{"node_name"})
}

func createNodeMemoryAtomUsedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_atom_used",
		Help:      "erlang memory counters - atom_used",
	}, []string{"node_name"})
}

func createNodeMemoryProcessesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_processes",
		Help:      "erlang memory counters - processes",
	}, []string{"node_name"})
}

func createNodeMemoryProcessesUsedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_processes_used",
		Help:      "erlang memory counters - processes_used",
	}, []string{"node_name"})
}

func createNodeMemoryBinaryMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_binary",
		Help:      "erlang memory counters - binary",
	}, []string{"node_name"})
}

func createNodeMemoryCodeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_code",
		Help:      "erlang memory counters - code",
	}, []string{"node_name"})
}

func createNodeMemoryEtsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "erlang",
		Name:      "memory_ets",
		Help:      "erlang memory counters - ets",
	}, []string{"node_name"})
}

func createMangoUnindexedQueriesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "unindexed_queries",
		Help:      "number of mango queries that could not use an index",
	}, []string{"node_name"})
}

func createMangoInvalidIndexesMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "query_invalid_index",
		Help:      "number of mango queries that generated an invalid index warning",
	}, []string{"node_name"})
}

func createMangoTooManyDocsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "too_many_docs_scanned",
		Help:      "number of mango queries that generated an index scan warning",
	}, []string{"node_name"})
}

func createMangoDocsExaminedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "docs_examined",
		Help:      "number of documents examined by mango queries coordinated by this node",
	}, []string{"node_name"})
}

func createMangoQuorumDocsExaminedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "quorum_docs_examined",
		Help:      "number of documents examined by mango queries, using cluster quorum",
	}, []string{"node_name"})
}

func createMangoResultsReturnedMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "results_returned",
		Help:      "number of rows returned by mango queries",
	}, []string{"node_name"})
}

func createMangoQueryTimeMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "query_time",
		Help:      "length of time processing a mango query",
	}, []string{"node_name", "metric"})
}

func createMangoEvaluateSelectorsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "mango",
		Name:      "evaluate_selector",
		Help:      "number of mango selector evaluations",
	}, []string{"node_name"})
}

func createViewStalenessMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "view",
		Name:      "staleness",
		Help:      "the view's staleness (the view's update_seq compared to the database's update_seq)",
	}, []string{"db_name", "design_doc_name", "view_name", "shard_begin", "shard_end"})
}

func createSchedulerJobsMetric() *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace: namespace,
		Subsystem: "scheduler",
		Name:      "jobs",
		Help:      "scheduler jobs",
	}, []string{"node_name", "job_id", "db_name", "doc_id", "source", "target"})
}
