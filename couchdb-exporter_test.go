package main

import (
	"encoding/base64"
	"fmt"
	"github.com/gesellix/couchdb-exporter/lib"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"
)

type handler func(w http.ResponseWriter, r *http.Request)

func BasicAuth(basicAuth lib.BasicAuth, pass handler) handler {

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

func couchdbResponse(t *testing.T, versionSuffix string) handler {
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

func getGaugeValue(metrics []*dto.Metric, labelName string, labelValue string) float64 {
	for _, metric := range metrics {
		for _, label := range metric.Label {
			if *label.Name == labelName && *label.Value == labelValue {
				return *metric.Gauge.Value
			}
		}
	}
	return 0
}

func printMetrics(metrics []*dto.Metric) {
	metricStrings := []string{}
	for _, metric := range metrics {
		metricStrings = append(metricStrings, proto.MarshalTextString(metric))
	}

	sort.Strings(metricStrings)
	fmt.Println(metricStrings)
}

func performCouchdbStatsTest(t *testing.T, couchdbVersion string, expectedMetricsCount int, expectedGetRequestCount float64, expectedDiskSize float64) {
	basicAuth := lib.BasicAuth{Username: "username", Password: "password"}
	handler := http.HandlerFunc(BasicAuth(basicAuth, couchdbResponse(t, couchdbVersion)))
	server := httptest.NewServer(handler)

	e := lib.NewExporter(server.URL, basicAuth, []string{"example", "another-example"})
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	metrics := []*dto.Metric{}
	for metric := range ch {
		dtoMetric := &dto.Metric{}
		metric.Write(dtoMetric)
		metrics = append(metrics, dtoMetric)
	}
	//printMetrics(metrics)
	actualGetRequestCount := getGaugeValue(metrics, "method", "GET")
	actualDiskSize := getGaugeValue(metrics, "db_name", "example")

	if len(metrics) < expectedMetricsCount {
		t.Errorf("got less metrics (%d) as expected (%d)", len(metrics), expectedMetricsCount)
	}
	if len(metrics) > expectedMetricsCount {
		t.Errorf("got more metrics (%d) as expected (%d)", len(metrics), expectedMetricsCount)
	}
	if expectedGetRequestCount != actualGetRequestCount {
		t.Errorf("expected %f GET requests, but got %f instead", expectedGetRequestCount, actualGetRequestCount)
	}
	if expectedDiskSize != actualDiskSize {
		t.Errorf("expected %f disk size, but got %f instead", expectedDiskSize, actualDiskSize)
	}
}

func TestCouchdbStatsV1(t *testing.T) {
	performCouchdbStatsTest(t, "v1", 43, 4711, 12396)
}

func TestCouchdbStatsV2(t *testing.T) {
	performCouchdbStatsTest(t, "v2", 74, 4712, 58570)
}
