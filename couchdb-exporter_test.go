package main

import (
	"encoding/base64"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"sort"
	"strings"
	"testing"

	"github.com/gesellix/couchdb-exporter/lib"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const ()

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

func CouchdbStatsResponse(exampleFile []byte) handler {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(exampleFile))
	}
}

func TestCouchdbStats(t *testing.T) {
	expectedMetricsCount := 32
	exampleFile, err := ioutil.ReadFile("./couchdb-stats-example.json")
	if err != nil {
		t.Error("File error: %v\n", err)
	}

	basicAuth := lib.BasicAuth{Username: "username", Password: "password"}
	handler := http.HandlerFunc(BasicAuth(basicAuth, CouchdbStatsResponse(exampleFile)))
	server := httptest.NewServer(handler)

	e := lib.NewExporter(server.URL, basicAuth)
	ch := make(chan prometheus.Metric)

	go func() {
		defer close(ch)
		e.Collect(ch)
	}()

	metricStrings := []string{}
	for metric := range ch {
		dtoMetric := &dto.Metric{}
		metric.Write(dtoMetric)
		metricStrings = append(metricStrings, proto.CompactTextString(dtoMetric))
	}
	sort.Strings(metricStrings)
	//	fmt.Println(metricStrings)

	if len(metricStrings) < expectedMetricsCount {
		t.Errorf("got less metrics (%d) as expected (%d)", len(metricStrings), expectedMetricsCount)
	}
	if len(metricStrings) > expectedMetricsCount {
		t.Errorf("got more metrics (%d) as expected (%d)", len(metricStrings), expectedMetricsCount)
	}
}
