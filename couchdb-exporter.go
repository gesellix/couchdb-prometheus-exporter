package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/golang/glog"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	namespace = "couchdb"
)

var (
	listeningAddress = flag.String("telemetry.address", "localhost:9984", "Address on which to expose metrics.")
	metricsEndpoint = flag.String("telemetry.endpoint", "/metrics", "Path under which to expose metrics.")
	couchdbURI = flag.String("couchdb.uri", "http://localhost:5984", "URI to the CouchDB instance")
	insecure = flag.Bool("insecure", true, "Ignore server certificate if using https")
)

// Exporter collects couchdb stats from the given URI and exports them using
// the prometheus metrics package.
type Exporter struct {
	URI    string
	client *http.Client
	mutex  sync.RWMutex

	up prometheus.Gauge
}

// NewExporter returns an initialized Exporter.
func NewExporter(uri string) *Exporter {
	couchDbClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{InsecureSkipVerify: *insecure},
		},
	}
	return &Exporter{
		URI:    fmt.Sprintf("%s/_stats", uri),
		client: couchDbClient,
		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      "up",
			Help:      "Was the last query of CouchDB stats successful.",
		}),
	}
}

// Describe describes all the metrics ever exported by the couchdb exporter. It
// implements prometheus.Collector.
func (e *Exporter) Describe(ch chan <- *prometheus.Desc) {
	ch <- e.up.Desc()
}

func (e *Exporter) collect(ch chan <- prometheus.Metric) error {
	sendStatus := func() {
		ch <- e.up
	}
	defer sendStatus()

	e.up.Set(0)
	resp, err := e.client.Get(e.URI)
	if err != nil {
		return fmt.Errorf("Error reading couchdb stats: %v", err)
	}

	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		if err != nil {
			data = []byte(err.Error())
		}
		return fmt.Errorf("Status %s (%d): %s", resp.Status, resp.StatusCode, data)
	}

	e.up.Set(1)

	var stats StatsResponse
	err = json.Unmarshal(data, &stats)
	//	glog.Info(fmt.Sprintf("stats: %v\n", stats))

	return nil
}

// Collect fetches the stats from configured couchdb location and delivers them
// as Prometheus metrics. It implements prometheus.Collector.
func (e *Exporter) Collect(ch chan <- prometheus.Metric) {
	e.mutex.Lock() // To protect metrics from concurrent collects.
	defer e.mutex.Unlock()
	if err := e.collect(ch); err != nil {
		glog.Error(fmt.Sprintf("Error collecting stats: %s", err))
	}
	return
}

func main() {
	flag.Parse()

	exporter := NewExporter(*couchdbURI)
	prometheus.MustRegister(exporter)

	log.Printf("Starting exporter at %s to read from CouchDB at %s", *listeningAddress, *couchdbURI)
	http.Handle(*metricsEndpoint, prometheus.Handler())
	log.Fatal(http.ListenAndServe(*listeningAddress, nil))
}
