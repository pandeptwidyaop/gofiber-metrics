package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

var RequestDuration = promauto.NewHistogramVec(
	prometheus.HistogramOpts{
		Name:    "api_request_duration_seconds",
		Buckets: []float64{.00005, .0005, .005, .01, .025, .05, .1, .25, .5, 1, 2.5},
	},
	[]string{"path"},
)

var RequestTotal = promauto.NewCounterVec(
	prometheus.CounterOpts{
		Name: "api_request_total",
	},
	[]string{"path", "code"},
)
