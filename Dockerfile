FROM alpine:3.8

RUN apk --no-cache add ca-certificates && update-ca-certificates

RUN mkdir -p /opt
ADD ./certctl /opt/certctl

ENTRYPOINT ["/opt/certctl"]
