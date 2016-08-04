FROM busybox:ubuntu-14.04

RUN mkdir -p /opt
ADD ./certctl /opt/certctl

ENTRYPOINT ["/opt/certctl"]
