package docs

import (
	"fmt"
	"strconv"
        "net/http"

        "github.com/ashleyprimo/go-qr-generator/initialize"
)

func Landing(w http.ResponseWriter, r *http.Request) {
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
    <p>Please Note: The requested size must be no more than %[4]sx%[4]s pixels<p>
    <p>Example Request: <code>http://%[2]s%[3]s?size=%[5]s&data=This%%20is%%20an%%20example</code></p>
    <p>This above example, will generate a QR code with the text "This is an example", with a size of %[5]sx%[5]s pixels- it's that easy.</p>
		<img src="http://%[2]s%[3]s?size=%[5]s&data=This%%20is%%20an%%20example.">
	</body>
</html>
		`
	)
	w.Write([]byte(fmt.Sprintf(documentationPage, initialize.ApplicationName, r.Host, *initialize.QREndpoint, strconv.Itoa(*initialize.MaxSize), strconv.Itoa(*initialize.DefaultSize))))
}

