package http

import (
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

var (
	httpRequestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "http_requests_total",
			Help: "How many HTTP requests processed, partitioned by status code and HTTP method.",
		},
		[]string{"path"},
	)

	randomNumber = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: "random_number",
			Help: "the value of random number.",
		}, []string{"key"},
	)

	//https://github.com/kubernetes/kubernetes/blob/master/pkg/volume/util/metrics.go
	storageOperationMetric = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "storage_operation_duration_seconds",
			Help:    "Storage operation duration",
			Buckets: []float64{.1, .25, .5, 1, 2.5, 5, 10, 15, 25, 50},
		},
		[]string{"volume_plugin", "operation_name"},
	)
)

func prometheusRegister() {
	log.WithFields(log.Fields{"name": "httpRequestsTotal"}).Info("prometheus register")
	prometheus.MustRegister(httpRequestsTotal)
	prometheus.MustRegister(randomNumber)
	randomNumber.With(prometheus.Labels{"key": "value"}).Set(0)
	prometheus.MustRegister(storageOperationMetric)
}
