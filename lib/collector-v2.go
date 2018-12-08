package lib

import (
	"fmt"
	"strconv"
)

func (e *Exporter) collectV2(stats Stats, exposedHttpStatusCodes []string, databases []string) error {
	e.databasesTotal.Set(float64(stats.DatabasesTotal))

	for name, nodeStats := range stats.StatsByNodeName {
		//fmt.Printf("%s -> %v\n", name, stats)
		//glog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.nodeUp.WithLabelValues(name).Set(nodeStats.Up)
		e.nodeInfo.WithLabelValues(name, nodeStats.NodeInfo.Version, nodeStats.NodeInfo.Vendor.Name).Set(1)

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

	for _, dbName := range databases {
		e.diskSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DiskSize)
		e.dataSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DataSize)
		e.docCount.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DocCount)
		e.docDelCount.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DocDelCount)
		e.compactRunning.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].CompactRunning)
		e.diskSizeOverhead.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DiskSizeOverhead)

		for designDoc, view := range stats.DatabaseStatsByDbName[dbName].Views {
			for viewName, updateSeq := range view {
				dbUpdateSeq, err := DecodeUpdateSeq(stats.DatabaseStatsByDbName[dbName].UpdateSeq.String())
				if err != nil {
					return err
				}

				viewUpdateSeq, err := DecodeUpdateSeq(updateSeq)
				if err != nil {
					return err
				}

				for _, viewRangeSeq := range viewUpdateSeq {
					for _, dbRangeSeq := range dbUpdateSeq {
						if viewRangeSeq.Range[0].Cmp(dbRangeSeq.Range[0]) == 0 {
							age := dbRangeSeq.Seq - viewRangeSeq.Seq
							//glog.Infof("dbRangeSeq.Seq %d, viewRangeSeq.Seq %d, age %d", dbRangeSeq.Seq, viewRangeSeq.Seq, age)
							e.viewStaleness.WithLabelValues(
								dbName,
								designDoc,
								viewName,
								viewRangeSeq.Range[0].String(),
								viewRangeSeq.Range[1].String()).Set(float64(age))
						}
					}
				}
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
