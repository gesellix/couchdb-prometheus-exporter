package lib

import (
	"fmt"
	"strconv"
)

func (e *Exporter) collectV2(stats Stats, exposedHttpStatusCodes []string, collectorConfig CollectorConfig) error {
	e.databasesTotal.Set(float64(stats.DatabasesTotal))

	for name, nodeStats := range stats.StatsByNodeName {
		// fmt.Printf("%s -> %v\n", name, stats)
		// klog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.nodeUp.WithLabelValues(name).Set(nodeStats.Up)
		e.nodeInfo.WithLabelValues(name, nodeStats.NodeInfo.Version, nodeStats.NodeInfo.Vendor.Name).Set(1)

		e.authCacheHits.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheHits.Value)
		e.authCacheMisses.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheMisses.Value)
		e.databaseReads.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseReads.Value)
		e.databaseWrites.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseWrites.Value)
		e.openDatabases.WithLabelValues(name).Set(nodeStats.Couchdb.OpenDatabases.Value)
		e.openOsFiles.WithLabelValues(name).Set(nodeStats.Couchdb.OpenOsFiles.Value)

		e.requestTime.WithLabelValues(name, "Min").Set(nodeStats.Couchdb.RequestTime.Value.Min)
		e.requestTime.WithLabelValues(name, "Max").Set(nodeStats.Couchdb.RequestTime.Value.Max)
		e.requestTime.WithLabelValues(name, "ArithmeticMean").Set(nodeStats.Couchdb.RequestTime.Value.ArithmeticMean)
		e.requestTime.WithLabelValues(name, "GeometricMean").Set(nodeStats.Couchdb.RequestTime.Value.GeometricMean)
		e.requestTime.WithLabelValues(name, "HarmonicMean").Set(nodeStats.Couchdb.RequestTime.Value.HarmonicMean)
		e.requestTime.WithLabelValues(name, "Median").Set(nodeStats.Couchdb.RequestTime.Value.Median)
		e.requestTime.WithLabelValues(name, "Variance").Set(nodeStats.Couchdb.RequestTime.Value.Variance)
		e.requestTime.WithLabelValues(name, "StandardDeviation").Set(nodeStats.Couchdb.RequestTime.Value.StandardDeviation)
		e.requestTime.WithLabelValues(name, "Skewness").Set(nodeStats.Couchdb.RequestTime.Value.Skewness)
		e.requestTime.WithLabelValues(name, "Kurtosis").Set(nodeStats.Couchdb.RequestTime.Value.Kurtosis)

		for _, percentile := range nodeStats.Couchdb.RequestTime.Value.Percentile {
			e.requestTime.WithLabelValues(name, fmt.Sprintf("%v", percentile[0])).Set(percentile[1])
		}

		for _, level := range exposedLogLevels {
			e.couchLog.WithLabelValues(level, name).Set(nodeStats.CouchLog.Level[level].Value)
		}

		for _, metric := range exposedWorkerMetrics {
			e.fabricWorker.WithLabelValues(metric, name).Set(nodeStats.Fabric.Worker[metric].Value)
		}

		for _, metric := range exposedOpenShardMetrics {
			e.fabricOpenShard.WithLabelValues(metric, name).Set(nodeStats.Fabric.OpenShard[metric].Value)
		}

		for _, metric := range exposedExitState {
			e.fabricReadRepairs.WithLabelValues(metric, name).Set(nodeStats.Fabric.ReadRepairs[metric].Value)
		}

		for _, metric := range exposedDocUpdateMetrics {
			e.fabricDocUpdate.WithLabelValues(metric, name).Set(nodeStats.Fabric.DocUpdate[metric].Value)
		}

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

		e.couchReplicatorChangesReadFailures.WithLabelValues(name).Set(nodeStats.CouchReplicator.ChangesReadFailures.Value)
		e.couchReplicatorChangesReaderDeaths.WithLabelValues(name).Set(nodeStats.CouchReplicator.ChangesReaderDeaths.Value)
		e.couchReplicatorChangesManagerDeaths.WithLabelValues(name).Set(nodeStats.CouchReplicator.ChangesManagerDeaths.Value)
		e.couchReplicatorChangesQueueDeaths.WithLabelValues(name).Set(nodeStats.CouchReplicator.ChangesQueueDeaths.Value)
		for _, metric := range exposedExitState {
			e.couchReplicatorCheckpoints.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.Checkpoints[metric].Value)
		}
		e.couchReplicatorFailedStarts.WithLabelValues(name).Set(nodeStats.CouchReplicator.FailedStarts.Value)
		e.couchReplicatorRequests.WithLabelValues(name).Set(nodeStats.CouchReplicator.Requests.Value)
		for _, metric := range exposedExitState {
			e.couchReplicatorResponses.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.Responses[metric].Value)
		}
		for _, metric := range exposedExitState {
			e.couchReplicatorStreamResponses.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.StreamResponses[metric].Value)
		}
		e.couchReplicatorWorkerDeaths.WithLabelValues(name).Set(nodeStats.CouchReplicator.WorkerDeaths.Value)
		e.couchReplicatorWorkersStarted.WithLabelValues(name).Set(nodeStats.CouchReplicator.WorkersStarted.Value)
		e.couchReplicatorClusterIsStable.WithLabelValues(name).Set(nodeStats.CouchReplicator.ClusterIsStable.Value)
		e.couchReplicatorDbScans.WithLabelValues(name).Set(nodeStats.CouchReplicator.DbScans.Value)
		for _, metric := range exposedReplicatorDocs {
			e.couchReplicatorDocs.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.Docs[metric].Value)
		}
		for _, metric := range exposedReplicatorJobs {
			e.couchReplicatorJobs.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.Jobs[metric].Value)
		}
		for _, metric := range exposedReplicatorConnection {
			e.couchReplicatorConnection.WithLabelValues(metric, name).Set(nodeStats.CouchReplicator.Connection[metric].Value)
		}
	}

	for _, dbName := range collectorConfig.ObservedDatabases {
		e.dbInfo.WithLabelValues(
			dbName,
			strconv.FormatFloat(stats.DatabaseStatsByDbName[dbName].DiskFormatVersion, 'G', -1, 32),
			strconv.FormatBool(stats.DatabaseStatsByDbName[dbName].Props.Partitioned),
		).Set(1)
		if stats.DatabaseStatsByDbName[dbName].DiskSize == 0 && stats.DatabaseStatsByDbName[dbName].Sizes.File > 0 {
			e.diskSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].Sizes.File)
		} else {
			e.diskSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DiskSize)
		}
		if stats.DatabaseStatsByDbName[dbName].DataSize == 0 && stats.DatabaseStatsByDbName[dbName].Sizes.Active > 0 {
			e.dataSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].Sizes.Active)
		} else {
			e.dataSize.WithLabelValues(dbName).Set(stats.DatabaseStatsByDbName[dbName].DataSize)
		}
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
							// klog.Infof("dbRangeSeq.Seq %d, viewRangeSeq.Seq %d, age %d", dbRangeSeq.Seq, viewRangeSeq.Seq, age)
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

	if collectorConfig.CollectSchedulerJobs {
		for _, job := range stats.SchedulerJobsResponse.Jobs {
			e.schedulerJobs.WithLabelValues(
				job.Node,
				job.ID,
				job.Database,
				job.DocID,
				job.Source,
				job.Target).Set(float64(len(job.History)))
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

	for nodeName, metric := range stats.SystemByNodeName {
		e.nodeMemoryOther.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Other)
		e.nodeMemoryAtom.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Atom)
		e.nodeMemoryAtomUsed.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.AtomUsed)
		e.nodeMemoryProcesses.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Processes)
		e.nodeMemoryProcessesUsed.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.ProcessesUsed)
		e.nodeMemoryBinary.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Binary)
		e.nodeMemoryCode.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Code)
		e.nodeMemoryEts.WithLabelValues(nodeName).Set(metric.MemoryStatsResponse.Ets)
	}

	return nil
}
