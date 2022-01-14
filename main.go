package main

import (
	"flag"
	"bytes"
	"fmt"
	"image/png"
	"net"
	"net/url"
	"strconv"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

const (
	version     = "1.0.0"
)

var (
	// Webserver Flags
	portNumber  = flag.String("port", "8080", "The port number to listen on for HTTP requests.")
	address     = flag.String("address", "0.0.0.0", "The address to listen on for HTTP requests.")

	https       = flag.Bool("https", false, "Enable, or Disable HTTPS")
	server_crt  = flag.String("server_crt", "server.crt", "Certificate file")
	server_key  = flag.String("server_key", "server.key", "Certificate key file.")

	// QR Setup Flags
	defaultSize = flag.String("defualtQRSize", "250", "Default QR Image Size, if unspecified by end user.")

	// Logging Flags
	logLevel    = flag.String("log.level", "info", "The level of logs to log")
	logConn     = flag.Bool("log.conn", true, "Log connections to API")

	// Misc Flags
	versionFlag = flag.Bool("v", false, "Outputs package version")

	// Useful Vars	
        applicationName = os.Args[0]

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

func InternalServerError(w http.ResponseWriter, r *http.Request, l string) {
	log.Errorf(l)
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Sorry, we were unable to serve your request due an internal server error!"))
	return
}

func BadReqeustError(w http.ResponseWriter, r *http.Request, l string) {
	// This is a user error, log as debug.
	log.Debugf(l)
	
	// Write Output
        w.WriteHeader(http.StatusBadRequest)
	w.Write([]byte(fmt.Sprintf("Invalid user input; please refer to %s/docs", r.Host)))
	return
}

func logConnection (r *http.Request) {
        if *logConn {
                source, _, _ := net.SplitHostPort(r.RemoteAddr)
                log.Infof("Request for %s from %s", r.RequestURI, source)
        }
}

func qrEngine(w http.ResponseWriter, r *http.Request) {
	logConnection(r)

	// Collect Parameter(s)
	parameters := r.URL.Query()

	// dataString (data to be encoded)
	if _, present := parameters["data"]; present == false {
		BadReqeustError(w, r, "No 'data' string provided.")
		return
	}
	dataString := parameters["data"][0]
	log.Debugf("Data String: %s", dataString)

	// sizeString (size of QR code)
	var sizeString string
        if _, present := parameters["size"]; present == false {
                log.Debugf("No 'size' string provided; will default to %s", *defaultSize)
		sizeString = *defaultSize
        } else {
		sizeString = parameters["size"][0]
	}

        log.Debugf("Size String: %s", sizeString)

	// Convert sizeString to sizeInt
	sizeInt, err := strconv.Atoi(sizeString)
	if err != nil {
                InternalServerError(w, r, "Unable to convert sizeString to sizeInt")
	}

	//
	dataStringUnescaped, err := url.QueryUnescape(dataString)
	if err != nil {
                InternalServerError(w, r, "Unable to 'Unescape' data string provided...")
	} else {
	        log.Debugf("Data String Unescaped: %s", dataStringUnescaped,)
	}

	// Generate QR
	code, err := qr.Encode(dataStringUnescaped, qr.L, qr.Auto)
        if err != nil {
	        InternalServerError(w, r, "Failed to generate QR code")
        } else {
                log.Debugf("Generated QR Code: %s", code,)
        }

	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, sizeInt, sizeInt)
        if err != nil {
                InternalServerError(w, r, "Unable to scale QR code.")
        } else {
                log.Debugf("Generated QR Code: %s", code,)
        }

	// Encode PNG
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, code); err != nil {
                InternalServerError(w, r, "Unable to encode PNG from code buffer.")
	}

	// Output QR
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

        if _, err := w.Write(buffer.Bytes()); err != nil {
                InternalServerError(w, r, "Unable to write/output QR code.")
        }

	return
}


func documentation(w http.ResponseWriter, r *http.Request) {
        logConnection(r)

	const (
		documentationPage = `
<html>
	<head>
		<title>%[1]s API documentation</title>
		
		<meta charset="utf-8">
		<meta content="width=device-width, initial-scale=1" name="viewport">

		<style>
			body {
				background-color: #5c5b5b;

				color: #ffffff;
				font-family: monospace;

				padding: 5%%;
			}

			h1 {
				text-align: center;
				text-transform: capitalize;
			}

			.img-center {
				margin-left: auto;
				margin-right: auto;
				display: block;
			}
		</style>
	</head>
	<body>
		<h1><a href="https://github.com/ashleyprimo/go-qr-generator">%[1]s</a> API documentation</h1>
		<h2>How to make a request</h2>
    <p>It's very simple, there are two possible parameters currently, <code>size</code> (which is optional) and <code>data</code> (which is mandatory)!<p>
    <p>Example Request: <code>http://%[2]s/?size=200&data=This%%20is%%20an%%20example</code></p>
    <p>This above example, will generate a QR code with the text "This is an example", with a size of 200x200 pixels- it's that easy.</p>
		<img src="http://%[2]s/?size=200&data=This%%20is%%20an%%20example.">
	</body>
</html>
		`
	)
	w.Write([]byte(fmt.Sprintf(documentationPage, applicationName, r.Host)))
}

func health(w http.ResponseWriter, r *http.Request) {
        logConnection(r)
        w.WriteHeader(http.StatusOK)
        w.Write([]byte("Ping."))
}

func main() {
	flag.Parse()

	if *versionFlag {
      		fmt.Printf("%s v%s", applicationName, version)
      		os.Exit(0)
    	}

	loglevel(*logLevel)

	http.HandleFunc("/", qrEngine)
        http.HandleFunc("/docs", documentation)
        http.HandleFunc("/health", health)


	log.Infof("Listening for requests on %s:%s", *address, *portNumber)

	if *https {
		log.Fatalf("Failed to start web server with TLS: %s", http.ListenAndServeTLS(fmt.Sprintf("%s:%s", *address, *portNumber), *server_crt, *server_key, nil))
	} else {
		log.Fatalf("Failed to start web server: %s", http.ListenAndServe(fmt.Sprintf("%s:%s", *address, *portNumber), nil))
	}
}
