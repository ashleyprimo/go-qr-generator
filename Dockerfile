FROM golang:1.17.5-alpine3.15 as builder
ADD . /go/qr-generator/
WORKDIR /go/qr-generator
RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /go/bin/go-qr-generator


FROM alpine:latest
ENV CMD_FLAGS ""

WORKDIR /app
COPY --from=builder /go/bin/go-qr-generator .

CMD ./go-qr-generator $CMD_FLAGS
EXPOSE 8080
