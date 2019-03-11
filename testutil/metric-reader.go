package testutil

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
	"log"
	"regexp"
)

func CollectMetrics(ch chan prometheus.Metric, debugMetrics bool) map[string]*dto.MetricFamily {
	metricFamiliesByName := make(map[string]*dto.MetricFamily, 1000)
	r, _ := regexp.Compile("fqName: ?\"([^\"]+)\"")

	var metrics []*dto.Metric
	for metric := range ch {
		dtoMetric := &dto.Metric{}
		metric.Write(dtoMetric)
		metrics = append(metrics, dtoMetric)

		desc := metric.Desc()
		name := r.FindStringSubmatch(desc.String())[1]

		metricFamily, ok := metricFamiliesByName[name]
		if ok {
			// TODO validity checks?
		} else {
			metricFamily = &dto.MetricFamily{}
			metricFamily.Name = proto.String(name)
			metricFamiliesByName[name] = metricFamily
		}
		metricFamily.Metric = append(metricFamily.Metric, dtoMetric)

		if debugMetrics {
			//log.Println(fmt.Sprintf("%s", name))
			log.Println(fmt.Sprintf("%s\n%s\n*********", name, proto.MarshalTextString(dtoMetric)))
		}
	}
	return metricFamiliesByName
}

func GetGaugeValue(metricFamilies map[string]*dto.MetricFamily, metricDesc string, labelName string, labelValue string) (float64, error) {
	//func getGaugeValue(metrics []*dto.Metric, metricDesc string, labelName string, labelValue string) float64 {
	for desc, metrics := range metricFamilies {
		if metricDesc == "" || desc == metricDesc {
			for _, metric := range metrics.Metric {
				if len(metric.Label) == 0 && labelName == "" && labelValue == "" {
					return *metric.Gauge.Value, nil
				}
				for _, label := range metric.Label {
					if *label.Name == labelName && *label.Value == labelValue {
						return *metric.Gauge.Value, nil
					}
				}
			}
		}
	}
	return 0, fmt.Errorf("no gauge found")
}

func CountMetrics(metricFamilies map[string]*dto.MetricFamily) int {
	count := 0
	for _, metrics := range metricFamilies {
		count += len(metrics.Metric)
	}
	return count
}
