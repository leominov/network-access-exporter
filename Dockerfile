FROM golang:1.14 as builder
WORKDIR /go/src/github.com/leominov/network-access-exporter
COPY . .
RUN make build

FROM scratch
COPY --from=builder /go/src/github.com/leominov/network-access-exporter/network-access-exporter /go/bin/network-access-exporter
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
ENTRYPOINT ["/go/bin/network-access-exporter"]
