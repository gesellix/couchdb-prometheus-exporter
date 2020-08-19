package lib

import (
	"fmt"
	"strconv"
)

func (e *Exporter) collectV1(stats Stats, exposedHttpStatusCodes []string, collectorConfig CollectorConfig) error {
	e.databasesTotal.Set(float64(stats.DatabasesTotal))

	for name, nodeStats := range stats.StatsByNodeName {
		//fmt.Printf("%s -> %v\n", name, stats)
		//klog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.nodeUp.WithLabelValues(name).Set(nodeStats.Up)
		e.nodeInfo.WithLabelValues(name, nodeStats.NodeInfo.Version, nodeStats.NodeInfo.Vendor.Name).Set(1)

		e.authCacheHits.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheHits.Current)
		e.authCacheMisses.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheMisses.Current)
		e.databaseReads.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseReads.Current)
		e.databaseWrites.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseWrites.Current)
		e.openDatabases.WithLabelValues(name).Set(nodeStats.Couchdb.OpenDatabases.Current)
		e.openOsFiles.WithLabelValues(name).Set(nodeStats.Couchdb.OpenOsFiles.Current)
		e.requestTime.WithLabelValues(name, "Current").Set(nodeStats.Couchdb.RequestTime.Current)

		for _, code := range exposedHttpStatusCodes {
			if _, ok := nodeStats.HttpdStatusCodes[code]; ok {
				e.httpdStatusCodes.WithLabelValues(code, name).Set(nodeStats.HttpdStatusCodes[code].Current)
			}
		}

		e.httpdRequestMethods.WithLabelValues("COPY", name).Set(nodeStats.HttpdRequestMethods.COPY.Current)
		e.httpdRequestMethods.WithLabelValues("DELETE", name).Set(nodeStats.HttpdRequestMethods.DELETE.Current)
		e.httpdRequestMethods.WithLabelValues("GET", name).Set(nodeStats.HttpdRequestMethods.GET.Current)
		e.httpdRequestMethods.WithLabelValues("HEAD", name).Set(nodeStats.HttpdRequestMethods.HEAD.Current)
		e.httpdRequestMethods.WithLabelValues("POST", name).Set(nodeStats.HttpdRequestMethods.POST.Current)
		e.httpdRequestMethods.WithLabelValues("PUT", name).Set(nodeStats.HttpdRequestMethods.PUT.Current)

		e.bulkRequests.WithLabelValues(name).Set(nodeStats.Httpd.BulkRequests.Current)
		e.clientsRequestingChanges.WithLabelValues(name).Set(nodeStats.Httpd.ClientsRequestingChanges.Current)
		e.requests.WithLabelValues(name).Set(nodeStats.Httpd.Requests.Current)
		e.temporaryViewReads.WithLabelValues(name).Set(nodeStats.Httpd.TemporaryViewReads.Current)
		e.viewReads.WithLabelValues(name).Set(nodeStats.Httpd.ViewReads.Current)
	}

	for _, dbName := range collectorConfig.ObservedDatabases {
		e.dbInfo.WithLabelValues(
			dbName,
			strconv.FormatFloat(stats.DatabaseStatsByDbName[dbName].DiskFormatVersion, 'G', -1, 32),
			strconv.FormatBool(stats.DatabaseStatsByDbName[dbName].Props.Partitioned),
		).Set(1)
		e.diskSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DiskSize)
		e.dataSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DataSize)
		e.docCount.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DocCount)
		e.docDelCount.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DocDelCount)
		e.compactRunning.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].CompactRunning)
		e.diskSizeOverhead.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DiskSizeOverhead)

		for designDoc, view := range stats.DatabaseStatsByDbName[dbName].Views {
			for viewName, updateSeq := range view {
				dbUpdateSeq, _ := stats.DatabaseStatsByDbName[dbName].UpdateSeq.Int64()
				viewUpdateSeq, _ := strconv.ParseInt(updateSeq, 10, 64)
				age := dbUpdateSeq - viewUpdateSeq
				e.viewStaleness.WithLabelValues(dbName, designDoc, viewName, "0", "1").Set(float64(age))
			}
		}
	}

	activeTasksByNode := make(map[string]ActiveTaskTypes)
	for _, task := range stats.ActiveTasksResponse {
		if task.Type == "replication" {
			e.activeTasksReplicationLastUpdate.WithLabelValues(
				task.Node,
				task.DocId,
				strconv.FormatBool(task.Continuous),
				task.Source,
				task.Target).Set(task.UpdatedOn)
		}

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
