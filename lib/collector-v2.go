package lib

import (
	"fmt"
)

func (e *Exporter) collectV2(stats Stats, exposedHttpStatusCodes []string, databases []string) error {
	for name, nodeStats := range stats.StatsByNodeName {
		//fmt.Printf("%s -> %v\n", name, stats)
		//glog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.authCacheHits.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheHits.Value)
		e.authCacheMisses.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheMisses.Value)
		e.databaseReads.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseReads.Value)
		e.databaseWrites.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseWrites.Value)
		e.openDatabases.WithLabelValues(name).Set(nodeStats.Couchdb.OpenDatabases.Value)
		e.openOsFiles.WithLabelValues(name).Set(nodeStats.Couchdb.OpenOsFiles.Value)
		e.requestTime.WithLabelValues(name).Set(nodeStats.Couchdb.RequestTime.Value.Median)

		for _, code := range exposedHttpStatusCodes {
			if _, ok := nodeStats.Couchdb.HttpdStatusCodes[code]; ok {
				e.httpdStatusCodes.WithLabelValues(code, name).Set(nodeStats.Couchdb.HttpdStatusCodes[code].Value)
			}
		}

		e.httpdRequestMethods.WithLabelValues("COPY", name).Set(nodeStats.Couchdb.HttpdRequestMethods.COPY.Value)
		e.httpdRequestMethods.WithLabelValues("DELETE", name).Set(nodeStats.Couchdb.HttpdRequestMethods.DELETE.Value)
		e.httpdRequestMethods.WithLabelValues("GET", name).Set(nodeStats.Couchdb.HttpdRequestMethods.GET.Value)
		e.httpdRequestMethods.WithLabelValues("HEAD", name).Set(nodeStats.Couchdb.HttpdRequestMethods.HEAD.Value)
		e.httpdRequestMethods.WithLabelValues("POST", name).Set(nodeStats.Couchdb.HttpdRequestMethods.POST.Value)
		e.httpdRequestMethods.WithLabelValues("PUT", name).Set(nodeStats.Couchdb.HttpdRequestMethods.PUT.Value)

		e.bulkRequests.WithLabelValues(name).Set(nodeStats.Couchdb.Httpd.BulkRequests.Value)
		e.clientsRequestingChanges.WithLabelValues(name).Set(nodeStats.Couchdb.Httpd.ClientsRequestingChanges.Value)
		e.requests.WithLabelValues(name).Set(nodeStats.Couchdb.Httpd.Requests.Value)
		e.temporaryViewReads.WithLabelValues(name).Set(nodeStats.Couchdb.Httpd.TemporaryViewReads.Value)
		e.viewReads.WithLabelValues(name).Set(nodeStats.Couchdb.Httpd.ViewReads.Value)
	}

	for name, dbStats := range stats.DatabaseStatsByNodeName {
		for _, dbName := range databases {
			e.diskSize.WithLabelValues(name, dbName).Set(dbStats[dbName].DiskSize)
			e.dataSize.WithLabelValues(name, dbName).Set(dbStats[dbName].DataSize)
			e.diskSizeOverhead.WithLabelValues(name, dbName).Set(dbStats[dbName].DiskSizeOverhead)
		}
	}

	activeTasksByNode := make(map[string]ActiveTaskTypes)
	for _, task := range stats.ActiveTasksResponse {
		if _, ok := activeTasksByNode[task.Node]; !ok {
			activeTasksByNode[task.Node] = ActiveTaskTypes{}
		}
		types := activeTasksByNode[task.Node]

		switch taskType := task.Type; taskType {
		case "database_compaction":
			types.DatabaseCompaction++
			types.Sum++
		case "view_compaction":
			types.ViewCompaction++
			types.Sum++
		case "indexer":
			types.Indexer++
			types.Sum++
		case "replication":
			types.Replication++
			types.Sum++
		default:
			fmt.Printf("unknown task type %s.", taskType)
			types.Sum++
		}
		activeTasksByNode[task.Node] = types
	}
	for nodeName, tasks := range activeTasksByNode {
		e.activeTasks.WithLabelValues(nodeName).Set(tasks.Sum)
		e.activeTasksDatabaseCompaction.WithLabelValues(nodeName).Set(tasks.DatabaseCompaction)
		e.activeTasksViewCompaction.WithLabelValues(nodeName).Set(tasks.ViewCompaction)
		e.activeTasksIndexer.WithLabelValues(nodeName).Set(tasks.Indexer)
		e.activeTasksReplication.WithLabelValues(nodeName).Set(tasks.Replication)
	}

	return nil
}
