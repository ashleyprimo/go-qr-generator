package main

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"

        "github.com/ashleyprimo/go-qr-generator/initialize"
)

var (
        versionExport = promauto.NewGauge(
                prometheus.GaugeOpts{
                        Namespace: initialize.MetricNamespace,
                        Name: "version",
                        Help: "current running version",
			ConstLabels: map[string]string{
				"version": initialize.Version,
			},
                },
        )
)

