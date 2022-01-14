
# GO QR Generator
**GO QR Generator**, is exactly what it sounds like; a QR code generator written in GoLang. Specifically, this is a small, simple microservice which will generate a QR code upon request - currently this microservice supports two client defined parameters:
* ```data```: (Required) The (URL encoded) string of data that should be encoded in the QR code.
* ```size```: (Optional) The size of the image (default configurable via arguments).

**Example of request**: Spin up the microservice, and go to `http://localhost:8081/?size=350&data=This%20is%20a%20test` 
- (*If you wanted to do this from a command line*, then install `viu` + `curl` and run `curl http://localhost:8081/\?size\=350\&data\=This%20is%20a%20test -s | viu -`!) 

**Public Endpoint**: https://qr-api.ashleyprimo.com/docs

```
Usage of go-qr-generator:
  -address string
    	The address to listen on for HTTP requests. (default "0.0.0.0")
  -defualtQRSize string
    	Default QR Image Size, if unspecified by end user. (default "250")
  -https
    	Enable, or Disable HTTPS
  -log.conn
    	Log connections to API (default true)
  -log.level string
    	The level of logs to log (default "info")
  -port string
    	The port number to listen on for HTTP requests. (default "8080")
  -server_crt string
    	Certificate file (default "server.crt")
  -server_key string
    	Certificate key file. (default "server.key")
```

## Endpoints
* `/` -> API endpoint.
* `/docs` -> Simple API documentation
* `/health` -> Simple Health Check 

## Docker
You can get started with docker quickly, by using the `docker pull ashleyprimo/go-qr-generator:latest` or alternatively you can build your own image using `docker build ./ -t go-qr-generator:latest`

Once you have *pulled* the image, you can now `docker run ashleyprimo/go-qr-generator`. That's it.

## References
* Barcode Library: https://github.com/boombuler/barcode
* Fork from: [echo-golang-qr-generator](https://github.com/kingsleytan/echo-golang-qr-generator)
