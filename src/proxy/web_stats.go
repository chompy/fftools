package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var totalRequests = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "http_requests_total",
		Help: "Number of get requests.",
	},
	[]string{"path"},
)

var responseStatus = prometheus.NewCounterVec(
	prometheus.CounterOpts{
		Name: "response_status",
		Help: "Status of HTTP response",
	},
	[]string{"status"},
)

/*var responseSize = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "response_size",
		Help: "Size in bytes of HTTP response",
	},
	[]string{"size"},
)*/

var httpDuration = promauto.NewHistogramVec(prometheus.HistogramOpts{
	Name: "http_response_time_seconds",
	Help: "Duration of HTTP requests.",
}, []string{"path"})

func init() {
	prometheus.Register(totalRequests)
	prometheus.Register(responseStatus)
	//prometheus.Register(responseSize)
	prometheus.Register(httpDuration)
}
