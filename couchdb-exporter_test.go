package main

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"fmt"
	"io/ioutil"
	"sort"

	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

const ()

func TestCouchdbStats(t *testing.T) {
	exampleFile, err := ioutil.ReadFile("./couchdb-stats-example.json")
	if err != nil {
		t.Error("File error: %v\n", err)
	}

	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(exampleFile))
	})
	server := httptest.NewServer(handler)

	e := NewExporter(server.URL)
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
	fmt.Println(metricStrings)

	expectedMessageCount := 32
	if len(metricStrings) > expectedMessageCount {
		t.Errorf("got more messages (%i) as expected (%i)", len(metricStrings), expectedMessageCount)
	}
}
