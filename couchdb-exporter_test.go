package main

import (
	"encoding/base64"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
	"time"

	"github.com/gesellix/couchdb-cluster-config/v17/pkg"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/stretchr/testify/assert"

	"github.com/gesellix/couchdb-prometheus-exporter/v29/lib"
	"github.com/gesellix/couchdb-prometheus-exporter/v29/testutil"
)

var clusterSetupDelay = 5 * time.Second
var clusterSetupTimeout = 40 * time.Second

type Handler func(w http.ResponseWriter, r *http.Request)

func BasicAuthHandler(basicAuth lib.BasicAuth, pass Handler) Handler {

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
		var err error
		if r.URL.Path == "/" {
			response := readFile(t, fmt.Sprintf("./testdata/couchdb-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.Path == "/_all_dbs" {
			response := readFile(t, "./testdata/all-dbs.json")
			_, err = w.Write(response)
		} else if r.URL.Path == "/_membership" {
			response := readFile(t, fmt.Sprintf("./testdata/couchdb-membership-response-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.Path == "/_active_tasks" {
			response := readFile(t, fmt.Sprintf("./testdata/active-tasks-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.Path == "/_scheduler/jobs" {
			response := readFile(t, fmt.Sprintf("./testdata/scheduler-jobs-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.Path == "/example" {
			response := readFile(t, fmt.Sprintf("./testdata/example-meta-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.Path == "/another-example" {
			response := readFile(t, fmt.Sprintf("./testdata/example-meta-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.String() == "/example/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"&include_docs=true" {
			response := readFile(t, "./testdata/example-all-design-docs.json")
			_, err = w.Write(response)
		} else if r.URL.String() == "/another-example/_all_docs?startkey=\"_design/\"&endkey=\"_design0\"&include_docs=true" {
			response := readFile(t, "./testdata/example-all-design-docs.json")
			_, err = w.Write(response)
		} else if r.URL.String() == "/example/_design/views/_view/by_id?stale=ok&update=false&stable=true&update_seq=true&include_docs=false&limit=0" {
			response := readFile(t, fmt.Sprintf("./testdata/example-view-stale-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.String() == "/another-example/_design/views/_view/by_id?stale=ok&update=false&stable=true&update_seq=true&include_docs=false&limit=0" {
			response := readFile(t, fmt.Sprintf("./testdata/example-view-stale-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.String() == "/example/_design/views/_view/by_id?stale=ok&update=false&stable=true&update_seq=true&include_docs=false&limit=0" {
			response := readFile(t, fmt.Sprintf("./testdata/example-view-stale-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else if r.URL.String() == "/another-example/_design/views/_view/by_id?stale=ok&update=false&stable=true&update_seq=true&include_docs=false&limit=0" {
			response := readFile(t, fmt.Sprintf("./testdata/example-view-stale-%s.json", versionSuffix))
			_, err = w.Write(response)
		} else {
			response := readFile(t, fmt.Sprintf("./testdata/couchdb-stats-response-%s.json", versionSuffix))
			_, err = w.Write(response)
		}
		if err != nil {
			t.Error(err)
		}
	}
}

func performCouchdbStatsTest(t *testing.T, couchdbVersion string, expectedMetricsCount int, expectedGetRequestCount float64, expectedDiskSize float64, expectedRequestCount float64) {
	basicAuth := lib.BasicAuth{Username: "username", Password: "password"}
	handler := http.HandlerFunc(BasicAuthHandler(basicAuth, couchdbResponse(t, couchdbVersion)))
	server := httptest.NewServer(handler)

	e := lib.NewExporter(server.URL, basicAuth, lib.CollectorConfig{
		Databases:            []string{"example", "another-example"},
		CollectViews:         true,
		CollectSchedulerJobs: true,
	}, true)

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

	actualGetRequestCount, err := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_request_methods", "method", "GET")
	if err != nil {
		t.Error(err)
	}
	if expectedGetRequestCount != actualGetRequestCount {
		t.Errorf("expected %f GET requests, but got %f instead", expectedGetRequestCount, actualGetRequestCount)
	}

	actualDiskSize, err := testutil.GetGaugeValue(metricFamilies, "couchdb_database_disk_size", "db_name", "example")
	if err != nil {
		t.Error(err)
	}
	if expectedDiskSize != actualDiskSize {
		t.Errorf("expected %f disk size, but got %f instead", expectedDiskSize, actualDiskSize)
	}

	actualRequestCount, err := testutil.GetGaugeValue(metricFamilies, "couchdb_exporter_request_count", "", "")
	if err != nil {
		t.Error(err)
	}
	if expectedRequestCount != actualRequestCount {
		t.Errorf("expected %f request count, but got %f instead", expectedRequestCount, actualRequestCount)
	}
}

func TestCouchdbStatsV1(t *testing.T) {
	performCouchdbStatsTest(t, "v1", 58, 4711, 12396, 11)
}

func TestCouchdbStatsV2(t *testing.T) {
	performCouchdbStatsTest(t, "v2", 306, 4712, 58570, 17)
}

func TestCouchdbStatsV2Prerelease(t *testing.T) {
	performCouchdbStatsTest(t, "v2-pre", 294, 4712, 58570, 17)
}

func TestCouchdbStatsV1Integration(t *testing.T) {
	if testing.Short() {
		t.Skip("skipping integration test")
	}

	dbAddress := "localhost:4895"
	err := cluster_config.AwaitNodes([]string{dbAddress}, clusterSetupDelay, clusterSetupTimeout, cluster_config.Available)
	if err != nil {
		t.Error(err)
	}
	dbUrl := fmt.Sprintf("http://%s", dbAddress)
	basicAuth := lib.BasicAuth{Username: "root", Password: "a-secret"}

	client := lib.NewCouchdbClient(dbUrl, basicAuth, true)
	databases := []string{"v1_testdb1", "v1_test/db2"}
	for _, db := range databases {
		_, err = client.Request("PUT", fmt.Sprintf("%s/%s", client.BaseUri, url.QueryEscape(db)), nil)
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("node_up", func(t *testing.T) {
		e := lib.NewExporter(dbUrl, basicAuth, lib.CollectorConfig{
			Databases:    []string{},
			CollectViews: true,
		}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		nodeName := "master"
		actualNodeUp, err := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_node_up", "node_name", nodeName)
		if err != nil {
			t.Error(err)
		}
		if actualNodeUp != 1 {
			t.Errorf("Expected node '%s' at '%s' to be available.", nodeName, dbUrl)
		}
	})

	t.Run("_all_dbs", func(t *testing.T) {
		e := lib.NewExporter(dbUrl, basicAuth, lib.CollectorConfig{
			Databases:    []string{"_all_dbs"},
			CollectViews: true,
		}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		for _, db := range databases {
			databaseDataSize, err := testutil.GetGaugeValue(metricFamilies, "couchdb_database_data_size", "db_name", db)
			if err != nil {
				log.Println(err)
				t.Errorf("Expected stats to be collected for database '%s'.", db)
			}
			if databaseDataSize != 0 {
				t.Errorf("Expected database data size to be 0 for database '%s'.", db)
			}
		}
	})

	for _, db := range databases {
		_, err = client.Request("DELETE", fmt.Sprintf("%s/%s", client.BaseUri, url.QueryEscape(db)), nil)
		if err != nil {
			t.Error(err)
		}
	}
}

func awaitMembership(t *testing.T, basicAuth lib.BasicAuth) func(address string) (bool, error) {
	time.Sleep(5 * time.Second)

	return func(address string) (bool, error) {
		dbUrl := fmt.Sprintf("http://%s", address)
		c := lib.NewCouchdbClient(dbUrl, basicAuth, true)
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
	err := cluster_config.AwaitNodes([]string{dbAddress}, clusterSetupDelay, clusterSetupTimeout, cluster_config.Available)
	if err != nil {
		t.Error(err)
	}

	dbUrl := fmt.Sprintf("http://%s", dbAddress)
	basicAuth := lib.BasicAuth{Username: "root", Password: "a-secret"}

	err = cluster_config.AwaitNodes([]string{dbAddress}, clusterSetupDelay, clusterSetupTimeout, awaitMembership(t, basicAuth))
	if err != nil {
		t.Error(err)
	}

	client := lib.NewCouchdbClient(dbUrl, basicAuth, true)
	databases := []string{"v2_testdb1", "v2_test/db2"}
	for _, db := range databases {
		_, err = client.Request("PUT", fmt.Sprintf("%s/%s", client.BaseUri, url.QueryEscape(db)), nil)
		if err != nil {
			t.Error(err)
		}
	}

	t.Run("node_up", func(t *testing.T) {
		e := lib.NewExporter(dbUrl, basicAuth, lib.CollectorConfig{
			Databases:    []string{},
			CollectViews: true,
		}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		nodeName := "couchdb@172.16.238.11"
		actualNodeUp, err := testutil.GetGaugeValue(metricFamilies, "couchdb_httpd_node_up", "node_name", nodeName)
		if err != nil {
			t.Error(err)
		}
		if actualNodeUp != 1 {
			t.Errorf("Expected node '%s' at '%s' to be available.", nodeName, dbUrl)
		}
	})

	t.Run("_all_dbs", func(t *testing.T) {
		e := lib.NewExporter(dbUrl, basicAuth, lib.CollectorConfig{
			Databases:    []string{"_all_dbs"},
			CollectViews: true,
		}, true)

		ch := make(chan prometheus.Metric)
		go func() {
			defer close(ch)
			e.Collect(ch)
		}()

		metricFamilies := testutil.CollectMetrics(ch, false)

		for _, db := range databases {
			databaseDataSize, err := testutil.GetGaugeValue(metricFamilies, "couchdb_database_data_size", "db_name", db)
			if err != nil {
				log.Println(err)
				t.Errorf("Expected stats to be collected for database '%s'.", db)
			}
			if databaseDataSize != 0 {
				t.Errorf("Expected database data size to be 0 for database '%s'.", db)
			}
		}
	})

	for _, db := range databases {
		_, err = client.Request("DELETE", fmt.Sprintf("%s/%s", client.BaseUri, url.QueryEscape(db)), nil)
		if err != nil {
			t.Error(err)
		}
	}
}
