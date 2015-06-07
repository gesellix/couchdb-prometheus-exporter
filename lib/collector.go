package lib

import (
	"encoding/json"
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

	ch <- e.authCacheHits.Desc()
	ch <- e.authCacheMisses.Desc()
	ch <- e.databaseReads.Desc()
	ch <- e.databaseWrites.Desc()
	ch <- e.openDatabases.Desc()
	ch <- e.openOsFiles.Desc()
	ch <- e.requestTime.Desc()

	e.httpdStatusCodes.Describe(ch)
	e.httpdRequestMethods.Describe(ch)

	ch <- e.bulkRequests.Desc()
	ch <- e.clientsRequestingChanges.Desc()
	ch <- e.requests.Desc()
	ch <- e.temporaryViewReads.Desc()
	ch <- e.viewReads.Desc()
}

func (e *Exporter) collect(ch chan<- prometheus.Metric) error {
	sendStatus := func() {
		ch <- e.up
	}
	defer sendStatus()

	e.up.Set(0)
	data, err := e.client.getStats()
	if err != nil {
		return fmt.Errorf("Error reading couchdb stats: %v", err)
	}

	e.up.Set(1)

	var stats StatsResponse
	err = json.Unmarshal(data, &stats)
	if err != nil {
		return fmt.Errorf("error unmarshalling stats: %v", err)
	}
	//	glog.Info(fmt.Sprintf("stats: %v\n", stats))

	e.authCacheHits.Set(stats.Couchdb.AuthCacheHits.Current)
	ch <- e.authCacheHits
	e.authCacheMisses.Set(stats.Couchdb.AuthCacheMisses.Current)
	ch <- e.authCacheMisses
	e.databaseReads.Set(stats.Couchdb.DatabaseReads.Current)
	ch <- e.databaseReads
	e.databaseWrites.Set(stats.Couchdb.DatabaseWrites.Current)
	ch <- e.databaseWrites
	e.openDatabases.Set(stats.Couchdb.OpenDatabases.Current)
	ch <- e.openDatabases
	e.openOsFiles.Set(stats.Couchdb.OpenOsFiles.Current)
	ch <- e.openOsFiles
	e.requestTime.Set(stats.Couchdb.RequestTime.Current)
	ch <- e.requestTime

	for _, code := range exposedHttpStatusCodes {
		if _, ok := stats.HttpdStatusCodes[code]; ok {
			e.httpdStatusCodes.WithLabelValues(code).Set(stats.HttpdStatusCodes[code].Current)
		}
	}
	e.httpdStatusCodes.Collect(ch)

	e.httpdRequestMethods.WithLabelValues("COPY").Set(stats.HttpdRequestMethods.COPY.Current)
	e.httpdRequestMethods.WithLabelValues("DELETE").Set(stats.HttpdRequestMethods.DELETE.Current)
	e.httpdRequestMethods.WithLabelValues("GET").Set(stats.HttpdRequestMethods.GET.Current)
	e.httpdRequestMethods.WithLabelValues("HEAD").Set(stats.HttpdRequestMethods.HEAD.Current)
	e.httpdRequestMethods.WithLabelValues("POST").Set(stats.HttpdRequestMethods.POST.Current)
	e.httpdRequestMethods.WithLabelValues("PUT").Set(stats.HttpdRequestMethods.PUT.Current)
	e.httpdRequestMethods.Collect(ch)

	e.bulkRequests.Set(stats.Httpd.BulkRequests.Current)
	ch <- e.bulkRequests
	e.clientsRequestingChanges.Set(stats.Httpd.ClientsRequestingChanges.Current)
	ch <- e.clientsRequestingChanges
	e.requests.Set(stats.Httpd.Requests.Current)
	ch <- e.requests
	e.temporaryViewReads.Set(stats.Httpd.TemporaryViewReads.Current)
	ch <- e.temporaryViewReads
	e.viewReads.Set(stats.Httpd.ViewReads.Current)
	ch <- e.viewReads

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
