package qr

import (
	"fmt"
	"net"
	//	"flag"
	"bytes"
	"image/png"
	"net/http"
	"net/url"
	"strconv"

	"github.com/ashleyprimo/go-qr-generator/initialize"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
	log "github.com/sirupsen/logrus"
)

func InternalServerError(w http.ResponseWriter, r *http.Request, l string) {
	requests.WithLabelValues("5xx").Inc()

	// Log Error (server side)
	log.Errorf(l)

	// Write Output
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte("Sorry, we were unable to serve your request due an internal server error!"))
}

func BadReqeustError(w http.ResponseWriter, r *http.Request, l string) {
	requests.WithLabelValues("4xx").Inc()

	// This is a user error/failure, log as debug.
	log.Debugf(l)

	// Write Output
	w.WriteHeader(http.StatusBadRequest)

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte("Bad Request!"))
	if *initialize.EnableDocs {
		w.Write([]byte(fmt.Sprintf(" Maybe try looking at our documentation at %[1]s/docs", r.Host)))
	}
}

func logConnection(r *http.Request) {
	if *initialize.LogConn {
		source, _, _ := net.SplitHostPort(r.RemoteAddr)
		log.Infof("Request for %s from %s", r.RequestURI, source)
	}
}

func Engine(w http.ResponseWriter, r *http.Request) {
	initialize.Flags()

	numRequests.Inc()
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
		defaultSize := strconv.Itoa(*initialize.DefaultSize)
		log.Debugf("No 'size' string provided; will default to %s", defaultSize)
		sizeString = defaultSize
	} else {
		sizeString = parameters["size"][0]
	}

	log.Debugf("Size String: %s", sizeString)

	// Convert sizeString to sizeInt
	sizeInt, err := strconv.Atoi(sizeString)
	if err != nil {
		InternalServerError(w, r, "Unable to convert sizeString to sizeInt")
		return
	}

	// Ensure requested size is within set limits
	if sizeInt > *initialize.MaxSize {
		BadReqeustError(w, r, "Requested image size exceeds maximum set")
		return
	}

	//Unescape Data Input
	dataStringUnescaped, err := url.QueryUnescape(dataString)
	if err != nil {
		InternalServerError(w, r, "Unable to 'Unescape' data string provided...")
		return
	} else {
		log.Debugf("Data String Unescaped: %s", dataStringUnescaped)
	}

	// Generate QR
	code, err := qr.Encode(dataStringUnescaped, qr.L, qr.Auto)
	if err != nil {
		InternalServerError(w, r, "Failed to generate QR code")
		return
	} else {
		log.Debugf("Generated QR Code: %s", code)
	}

	// Ensure QR scale is possible
	codeBounds := code.Bounds()
	codeHeight := codeBounds.Max.Y - codeBounds.Min.Y
	if sizeInt < codeHeight {
		BadReqeustError(w, r, "Requested image size is smaller than actual QR minimum size")
		return
	}

	// Scale the barcode to the appropriate size
	code, err = barcode.Scale(code, sizeInt, sizeInt)
	if err != nil {
		InternalServerError(w, r, "Unable to scale QR code.")
		return
	} else {
		log.Debugf("Generated QR Code: %s", code)
	}

	// Encode PNG
	buffer := new(bytes.Buffer)
	if err := png.Encode(buffer, code); err != nil {
		InternalServerError(w, r, "Unable to encode PNG from code buffer.")
		return
	}

	// Output QR
	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Content-Length", strconv.Itoa(len(buffer.Bytes())))

	if _, err := w.Write(buffer.Bytes()); err != nil {
		InternalServerError(w, r, "Unable to write/output QR code.")
		return
	} else {
		// Log Success
		requests.WithLabelValues("2xx").Inc()
		log.Debugf("Successfully wrote QR code")
	}

	return
}
