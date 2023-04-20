FROM golang:1.19.8-buster as builder

ARG goproxy

COPY . /go/src/github.com/xchenhao/data-transport-service
RUN cd /go/src/github.com/xchenhao/data-transport-service \
    && GOPROXY=${goproxy} GOOS=linux go install -v

FROM debian:buster-slim

COPY --from=builder /usr/share/zoneinfo /usr/share/zoneinfo
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /etc/passwd /etc/
COPY --from=builder /go/bin/data-transport-service /bin/

WORKDIR /

USER nobody

ENTRYPOINT ["/bin/data-transport-service", "-config"]

CMD ["/etc/config.yml"]
