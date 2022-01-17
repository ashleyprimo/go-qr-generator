package main

import (
	"fmt"
	"net/http"
	"time"

	log "github.com/sirupsen/logrus"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"github.com/ashleyprimo/go-qr-generator/initialize"
)

var (
	versionExport = promauto.NewGauge(
		prometheus.GaugeOpts{
			Namespace: initialize.MetricNamespace,
			Name:      "version",
			Help:      "current running version",
			ConstLabels: map[string]string{
				"version": initialize.Version,
			},
		},
	)
)

func metrics() {
	// Metrics Endpoint
	if *initialize.EnableMetrics {
		log.Debugf("Metrics Endpoint Enabled")
		if *initialize.MetricServer {
			log.Debugf("Metrics Server Enabled")
			log.Infof("Listening for metrics requests on %s:%s", *initialize.MetricServerHost, *initialize.MetricServerPort)
			go func() {
				mux := http.NewServeMux()
				mux.Handle("/metrics", promhttp.Handler())
				srv := &http.Server{
					Addr:         fmt.Sprintf("%s:%s", *initialize.MetricServerHost, *initialize.MetricServerPort),
					Handler:      mux,
					ReadTimeout:  30 * time.Second,
					WriteTimeout: 30 * time.Second,
					IdleTimeout:  3 * time.Minute,
				}
				log.Fatalf("Failed to start web server: %s", srv.ListenAndServe())
			}()

		} else {
			http.Handle("/metrics", promhttp.Handler())
		}
	}
}
