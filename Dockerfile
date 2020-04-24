FROM registry.bryk.io/general/shell:0.1.0

ARG VERSION_TAG
LABEL maintainer="Ben Cessa <ben@bryk.io>"
LABEL version=${VERSION_TAG}

COPY didctl_${VERSION_TAG}_linux_amd64 /usr/bin/didctl
COPY ca-roots.crt /etc/ssl/certs/ca-roots.crt

VOLUME ["/etc/didctl"]

EXPOSE 9090/tcp

ENTRYPOINT ["/usr/bin/didctl"]
