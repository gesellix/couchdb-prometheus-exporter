package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"github.com/gesellix/couchdb-cluster-config/pkg"
	"github.com/gesellix/couchdb-prometheus-exporter/lib"
	"github.com/gesellix/couchdb-prometheus-exporter/testutil"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"
	"time"
)

type Handler func(w http.ResponseWriter, r *http.Request)

func BasicAuth(basicAuth lib.BasicAuth, pass Handler) Handler {

	validate := func(basicAuth lib.BasicAuth, username, password string) bool {
		if username == basicAuth.Username && password == basicAuth.Password {
			return true
		}
		return false
	}

	return func(w http.ResponseWriter, r *http.Request) {

		if len(r.Header["Authorization"]) == 0 || len(r.Header["Authorization"][0]) == 0 {
			http.Error(w, "missing Authorization", http.StatusBadRequest)
			return
		}
		auth := strings.SplitN(r.Header["Authorization"][0], " ", 2)

		if len(auth) != 2 || auth[0] != "Basic" {
			http.Error(w, "bad syntax", http.StatusBadRequest)
			return
		}

		payload, _ := base64.StdEncoding.DecodeString(auth[1])
		pair := strings.SplitN(string(payload), ":", 2)

		if len(pair) != 2 || !validate(basicAuth, pair[0], pair[1]) {
			http.Error(w, "authorization failed", http.StatusUnauthorized)
			return
		}

		pass(w, r)
	}
}

func readFile(t *testing.T, filename string) []byte {
	fileContent, err := ioutil.ReadFile(filename)
	if err != nil {
		t.Errorf("Error reading file %s: %v\n", filename, err)
	}
	return fileContent
}

func couchdbResponse(t *testing.T, versionSuffix string) Handler {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/" {
			file := readFile(t, fmt.Sprintf("./testdata/couchdb-%s.json", versionSuffix))
			w.Write([]byte(file))
		} else if r.URL.Path == "/_membership" {
			file := readFile(t, fmt.Sprintf("./testdata/couchdb-membership-response-%s.json", versionSuffix))
			w.Write([]byte(file))
		} else if r.URL.Path == "/_active_tasks" {
			file := readFile(t, fmt.Sprintf("./testdata/active-tasks-%s.json", versionSuffix))
			w.Write([]byte(file))
		} else if r.URL.Path == "/example" {
			file := readFile(t, fmt.Sprintf("./testdata/example-meta-%s.json", versionSuffix))
			w.Write([]byte(file))
		} else if r.URL.Path == "/another-example" {
			file := readFile(t, fmt.Sprintf("./testdata/example-meta-%s.json", versionSuffix))
			w.Write([]byte(file))
		} else {
			file := readFile(t, fmt.Sprintf("./testdata/couchdb-stats-response-%s.json", versionSuffix))
			w.Write([]byte(file))
		}
	}
}

func performCouchdbStatsTest(t *testing.T, couchdbVersion string, expectedMetricsCount int, expectedGetRequestCount float64, expectedDiskSize float64) {
	basicAuth := lib.BasicAuth{Username: "username", Password: "password"}
	handler := http.HandlerFunc(BasicAuth(basicAuth, couchdbResponse(t, couchdbVersion)))
	server := httptest.NewServer(handler)

	e := lib.NewExporter(server.URL, basicAuth, []string{"example", "another-example"}, true)

	ch := make(chan prometheus.Metric)
	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	metricFamilies := testutil.CollectMetrics(ch, false)
	metricsCount := testutil.CountMetrics(metricFamilies)

	if metricsCount < expectedMetricsCount {
		t.Errorf("got less metrics (%d) as expected (%d)", metricsCount, expectedMetricsCount)
	}
	if metricsCount > expectedMetricsCount {
		t.Errorf("got more metrics (%d) as expected (%d)", metricsCount, expectedMetricsCount)
	}

	actualGetRequestCount := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_request_methods", "method", "GET")
	if expectedGetRequestCount != actualGetRequestCount {
		t.Errorf("expected %f GET requests, but got %f instead", expectedGetRequestCount, actualGetRequestCount)
	}

	actualDiskSize := testutil.GetGaugeValue(metricFamilies, "couchdb_database_disk_size", "db_name", "example")
	if expectedDiskSize != actualDiskSize {
		t.Errorf("expected %f disk size, but got %f instead", expectedDiskSize, actualDiskSize)
	}
}

func TestCouchdbStatsV1(t *testing.T) {
	performCouchdbStatsTest(t, "v1", 45, 4711, 12396)
}

func TestCouchdbStatsV2(t *testing.T) {
	performCouchdbStatsTest(t, "v2", 77, 4712, 58570)
}

func TestCouchdbStatsV1Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbAddress := "localhost:4895"
	err := cluster_config.AwaitNodes([]string{dbAddress}, cluster_config.Available)
	if err != nil {
		t.Error(err)
	}
	dbUrl := fmt.Sprintf("http://%s", dbAddress)

	t.Run("node_up", func(t *testing.T) {
		basicAuth := lib.BasicAuth{Username: "root", Password: "a-secret"}
		e := lib.NewExporter(dbUrl, basicAuth, []string{}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		nodeName := "master"
		actualNodeUp := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_node_up", "node_name", nodeName)
		if actualNodeUp != 1 {
			t.Errorf("Expected node '%s' at '%s' to be available.", nodeName, dbUrl)
		}
	})
}

func membership(t *testing.T, basicAuth lib.BasicAuth) (func(address string) (bool, error)) {
	time.Sleep(5 * time.Second)

	return func(address string) (bool, error) {
		dbUrl := fmt.Sprintf("http://%s", address)
		c := lib.NewCouchdbClient(dbUrl, basicAuth, []string{}, true)
		nodeNames, err := c.GetNodeNames()
		if err != nil {
			if err, ok := err.(net.Error); ok && (err.Timeout() || err.Temporary()) {
				return false, nil
			}
			return false, nil
			//return false, err
		}
		log.Println(fmt.Sprintf("%v (%d)", nodeNames, len(nodeNames)))
		return assert.ElementsMatch(t, nodeNames, []string{
			"couchdb@172.16.238.11",
			"couchdb@172.16.238.12",
			"couchdb@172.16.238.13",
		}), nil
	}
}

func TestCouchdbStatsV2Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	// <setup code>
	dbAddress := "localhost:15984"
	err := cluster_config.AwaitNodes([]string{dbAddress}, cluster_config.Available)
	if err != nil {
		t.Error(err)
	}

	dbUrl := fmt.Sprintf("http://%s", dbAddress)
	basicAuth := lib.BasicAuth{Username: "root", Password: "a-secret"}

	err = cluster_config.AwaitNodes([]string{dbAddress}, membership(t, basicAuth))
	if err != nil {
		t.Error(err)
	}

	t.Run("node_up", func(t *testing.T) {
		e := lib.NewExporter(dbUrl, basicAuth, []string{}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		nodeName := "couchdb@172.16.238.11"
		actualNodeUp := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_node_up", "node_name", nodeName)
		if actualNodeUp != 1 {
			t.Errorf("Expected node '%s' at '%s' to be available.", nodeName, dbUrl)
		}
	})

	// <tear-down code>
	// (nothing to do)
}
