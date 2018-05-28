package lib

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	exposedHttpStatusCodes = []string{"200", "201", "202", "301", "304", "400", "401", "403", "404", "405", "409", "412", "500"}
)

type ActiveTaskTypes struct {
	DatabaseCompaction float64
	ViewCompaction     float64
	Indexer            float64
	Replication        float64
	Sum                float64
}

type ActiveTaskTypesByNodeName map[string]ActiveTaskTypes

// Describe describes all the metrics ever exported by the couchdb exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()
	e.totalDatabases.Describe(ch)
	e.nodeUp.Describe(ch)
	e.nodeInfo.Describe(ch)

	e.authCacheHits.Describe(ch)
	e.authCacheMisses.Describe(ch)
	e.databaseReads.Describe(ch)
	e.databaseWrites.Describe(ch)
	e.openDatabases.Describe(ch)
	e.openOsFiles.Describe(ch)
	e.requestTime.Describe(ch)

	e.httpdStatusCodes.Describe(ch)
	e.httpdRequestMethods.Describe(ch)

	e.bulkRequests.Describe(ch)
	e.clientsRequestingChanges.Describe(ch)
	e.requests.Describe(ch)
	e.temporaryViewReads.Describe(ch)
	e.viewReads.Describe(ch)

	e.diskSize.Describe(ch)
	e.dataSize.Describe(ch)
	e.diskSizeOverhead.Describe(ch)

	e.activeTasks.Describe(ch)
	e.activeTasksDatabaseCompaction.Describe(ch)
	e.activeTasksViewCompaction.Describe(ch)
	e.activeTasksIndexer.Describe(ch)
	e.activeTasksReplication.Describe(ch)
	e.activeTasksReplicationLastUpdate.Describe(ch)
}

func (e *Exporter) resetAllMetrics() {
	metrics := []*prometheus.GaugeVec{
		e.nodeUp,
		e.nodeInfo,

		e.authCacheHits,
		e.authCacheMisses,
		e.databaseReads,
		e.databaseWrites,
		e.openDatabases,
		e.openOsFiles,
		e.requestTime,

		e.httpdStatusCodes,
		e.httpdRequestMethods,

		e.clientsRequestingChanges,
		e.temporaryViewReads,
		e.requests,
		e.bulkRequests,
		e.viewReads,

		e.diskSize,
		e.dataSize,
		e.diskSizeOverhead,

		e.activeTasks,
		e.activeTasksDatabaseCompaction,
		e.activeTasksViewCompaction,
		e.activeTasksIndexer,
		e.activeTasksReplication,
		e.activeTasksReplicationLastUpdate,
	}
	e.resetMetrics(metrics)
}

func (e *Exporter) resetMetrics(metrics []*prometheus.GaugeVec) {
	for _, metricVec := range metrics {
		metricVec.Reset()
	}
}

func (e *Exporter) getMonitoredDatabaseNames(candidates []string) ([]string, error) {
	if len(candidates) == 1 && candidates[0] == "_all_dbs" {
		return e.client.getDatabaseList()
	}
	return candidates, nil
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	e.up.Set(0)
	sendStatus := func() {
		ch <- e.up
	}
	defer sendStatus()

	e.resetAllMetrics()

	var databases, err = e.getMonitoredDatabaseNames(e.databases)
	if err != nil {
		return err
	}

	stats, err := e.client.getStats(databases)
	if err != nil {
		return fmt.Errorf("error collecting couchdb stats: %v", err)
	}
	e.up.Set(1)

	if stats.ApiVersion == "2" {
		e.collectV2(stats, exposedHttpStatusCodes, databases)
	} else {
		e.collectV1(stats, exposedHttpStatusCodes, databases)
	}

	e.totalDatabases.Collect(ch)
	e.nodeUp.Collect(ch)
	e.nodeInfo.Collect(ch)

	e.authCacheHits.Collect(ch)
	e.authCacheMisses.Collect(ch)
	e.databaseReads.Collect(ch)
	e.databaseWrites.Collect(ch)
	e.openDatabases.Collect(ch)
	e.openOsFiles.Collect(ch)
	e.requestTime.Collect(ch)

	e.httpdStatusCodes.Collect(ch)
	e.httpdRequestMethods.Collect(ch)

	e.bulkRequests.Collect(ch)
	e.clientsRequestingChanges.Collect(ch)
	e.requests.Collect(ch)
	e.temporaryViewReads.Collect(ch)
	e.viewReads.Collect(ch)

	e.diskSize.Collect(ch)
	e.dataSize.Collect(ch)
	e.diskSizeOverhead.Collect(ch)

	e.activeTasks.Collect(ch)
	e.activeTasksDatabaseCompaction.Collect(ch)
	e.activeTasksViewCompaction.Collect(ch)
	e.activeTasksIndexer.Collect(ch)
	e.activeTasksReplication.Collect(ch)
	e.activeTasksReplicationLastUpdate.Collect(ch)

	return nil
}

// Collect fetches the stats from configured couchdb location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		glog.Error(fmt.Sprintf("Error collecting stats: %s", err))
	}
	return
}
