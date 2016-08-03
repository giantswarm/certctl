FROM busybox:ubuntu-14.04

RUN mkdir -p /opt
ADD ./certificate-sidekick /opt/certificate-sidekick

ENTRYPOINT ["/opt/certificate-sidekick"]
