package initialize

import (
	"flag"

        log "github.com/sirupsen/logrus"
)

var (
        // Webserver Flags
        PortNumber      = flag.String("port", "8080", "The port number to listen on for HTTP requests.")
        Address         = flag.String("address", "0.0.0.0", "The address to listen on for HTTP requests.")

        Https           = flag.Bool("https", false, "Enable, or Disable HTTPS")
        Server_crt      = flag.String("server_crt", "server.crt", "Certificate file")
        Server_key      = flag.String("server_key", "server.key", "Certificate key file.")

        // QR Setup Flags
	QREndpoint      = flag.String("qr.endpoint", "/", "QR API endpoint location")

        DefaultSize     = flag.Int("qr.default.size", 250, "Default QR Image Size, if unspecified by end user.")
	MaxSize         = flag.Int("qr.max.size", 1000, "Maximum QR Image Size")

        // Logging Flags
        LogLevel        = flag.String("log.level", "info", "The level of logs to log")
       	LogConn         = flag.Bool("log.conn", true, "Log connections to API")

        // Documentation Flags
        EnableDocs      = flag.Bool("enable.docs", true, "Enable documentation (/docs) endpoint.")

	// Metrics Flagd
	EnableMetrics   = flag.Bool("enable.metrics", true, "Enable metrics (/metrics) endpoint.")

        // Misc Flags
        VersionFlag     = flag.Bool("v", false, "Outputs package version")
)

func Flags() {
	flag.Parse()

	// Flag Sanity Check
	if *DefaultSize > *MaxSize {
		log.Fatalf("-qr.default.size cannot be bigger than -qr.max.size!")
	}
}
