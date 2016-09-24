package lib

func (e *Exporter) collectV2(stats Stats, exposedHttpStatusCodes []string) error {
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

	return nil
}
