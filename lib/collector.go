package lib

import (
	"fmt"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

var (
	exposedHttpStatusCodes = []string{"200", "201", "202", "301", "304", "400", "401", "403", "404", "405", "409", "412", "500"}
)

// Describe describes all the metrics ever exported by the couchdb exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.up.Desc()

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
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	sendStatus := func() {
		ch <- e.up
	}
	defer sendStatus()

	e.up.Set(0)
	statsByNodeName, err := e.client.getStats()
	if err != nil {
		return fmt.Errorf("Error reading couchdb stats: %v", err)
	}

	e.up.Set(1)

	for name, stats := range statsByNodeName {
		//fmt.Printf("%s -> %v\n", name, stats)
		//glog.Info(fmt.Sprintf("name: %s -> stats: %v\n", name, stats))
		e.authCacheHits.WithLabelValues(name).Set(stats.Couchdb.AuthCacheHits.Current)
		e.authCacheMisses.WithLabelValues(name).Set(stats.Couchdb.AuthCacheMisses.Current)
		e.databaseReads.WithLabelValues(name).Set(stats.Couchdb.DatabaseReads.Current)
		e.databaseWrites.WithLabelValues(name).Set(stats.Couchdb.DatabaseWrites.Current)
		e.openDatabases.WithLabelValues(name).Set(stats.Couchdb.OpenDatabases.Current)
		e.openOsFiles.WithLabelValues(name).Set(stats.Couchdb.OpenOsFiles.Current)
		e.requestTime.WithLabelValues(name).Set(stats.Couchdb.RequestTime.Current)

		for _, code := range exposedHttpStatusCodes {
			if _, ok := stats.HttpdStatusCodes[code]; ok {
				e.httpdStatusCodes.WithLabelValues(code).Set(stats.HttpdStatusCodes[code].Current)
			}
		}

		e.httpdRequestMethods.WithLabelValues("COPY", name).Set(stats.HttpdRequestMethods.COPY.Current)
		e.httpdRequestMethods.WithLabelValues("DELETE", name).Set(stats.HttpdRequestMethods.DELETE.Current)
		e.httpdRequestMethods.WithLabelValues("GET", name).Set(stats.HttpdRequestMethods.GET.Current)
		e.httpdRequestMethods.WithLabelValues("HEAD", name).Set(stats.HttpdRequestMethods.HEAD.Current)
		e.httpdRequestMethods.WithLabelValues("POST", name).Set(stats.HttpdRequestMethods.POST.Current)
		e.httpdRequestMethods.WithLabelValues("PUT", name).Set(stats.HttpdRequestMethods.PUT.Current)

		e.bulkRequests.WithLabelValues(name).Set(stats.Httpd.BulkRequests.Current)
		e.clientsRequestingChanges.WithLabelValues(name).Set(stats.Httpd.ClientsRequestingChanges.Current)
		e.requests.WithLabelValues(name).Set(stats.Httpd.Requests.Current)
		e.temporaryViewReads.WithLabelValues(name).Set(stats.Httpd.TemporaryViewReads.Current)
		e.viewReads.WithLabelValues(name).Set(stats.Httpd.ViewReads.Current)
	}
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
