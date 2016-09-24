package lib

func (e *Exporter) collectV1(stats Stats, exposedHttpStatusCodes []string) error {
	for name, nodeStats := range stats.StatsByNodeName {
		//fmt.Printf("%s -> %v\n", name, stats)
		//glog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.authCacheHits.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheHits.Current)
		e.authCacheMisses.WithLabelValues(name).Set(nodeStats.Couchdb.AuthCacheMisses.Current)
		e.databaseReads.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseReads.Current)
		e.databaseWrites.WithLabelValues(name).Set(nodeStats.Couchdb.DatabaseWrites.Current)
		e.openDatabases.WithLabelValues(name).Set(nodeStats.Couchdb.OpenDatabases.Current)
		e.openOsFiles.WithLabelValues(name).Set(nodeStats.Couchdb.OpenOsFiles.Current)
		e.requestTime.WithLabelValues(name).Set(nodeStats.Couchdb.RequestTime.Current)

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

	return nil
}
