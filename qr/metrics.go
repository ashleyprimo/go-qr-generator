package qr

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

        "github.com/ashleyprimo/go-qr-generator/initialize"
)

var (
        numRequests = promauto.NewCounter(
                prometheus.CounterOpts{
                        Namespace: initialize.MetricNamespace,
                        Name: "requests_total",
                        Help: "Total number of requests to QR API endpoint.",
                },
        )

	requests = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: initialize.MetricNamespace,
			Name: "requests",
			Help: "Completed requests to QR API endpoint.",
		},
                []string{"type"},
	)
)
