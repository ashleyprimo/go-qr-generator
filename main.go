package main

import (
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"

        "github.com/ashleyprimo/go-qr-generator/initialize"
	"github.com/ashleyprimo/go-qr-generator/qr"
        "github.com/ashleyprimo/go-qr-generator/documentation"

        "github.com/prometheus/client_golang/prometheus/promhttp"
)

func loglevel(opt string) {
	switch opt {
	case "error":
		log.SetLevel(log.ErrorLevel)
	case "warn":
		log.SetLevel(log.WarnLevel)
	case "info":
		log.SetLevel(log.InfoLevel)
	case "debug":
		log.SetLevel(log.DebugLevel)
	default:
		log.Warnln("Unrecognized log level, will default to `info` log level")
	}
}

func health(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Ping."))
}

func main() {
        initialize.Flags()

	if *initialize.VersionFlag {
                fmt.Printf("%s v%s", initialize.ApplicationName, initialize.Version)
      		os.Exit(0)
    	}

	loglevel(*initialize.LogLevel)

	// QR Engine API Endpoint
	http.HandleFunc(*initialize.QREndpoint, qr.Engine)

	// Documentation Endpoint
	if *initialize.EnableDocs {
                log.Debugf("Documentation Endpoint Enabled")
	        http.HandleFunc("/docs", docs.Landing)
	}

	// Health Check Endpoint
        http.HandleFunc("/health", health)

	// Metrics Endpoint
	if *initialize.EnableMetrics {
                log.Debugf("Metrics Endpoint Enabled")
		http.Handle("/metrics", promhttp.Handler())
	}

	log.Infof("Listening for requests on %s:%s", *initialize.Address, *initialize.PortNumber)

	if *initialize.Https {
		log.Fatalf("Failed to start web server with TLS: %s", http.ListenAndServeTLS(fmt.Sprintf("%s:%s", *initialize.Address, *initialize.PortNumber), *initialize.Server_crt, *initialize.Server_key, nil))
	} else {
		log.Fatalf("Failed to start web server: %s", http.ListenAndServe(fmt.Sprintf("%s:%s", *initialize.Address, *initialize.PortNumber), nil))
	}
}
