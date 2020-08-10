FROM docker.pkg.github.com/bryk-io/base-images/shell:0.1.0

ARG VERSION_TAG
LABEL maintainer="Ben Cessa <ben@bryk.io>"
LABEL version=${VERSION_TAG}

COPY didctl_*_linux_amd64 /usr/bin/didctl

# Use the CA roots already included in the base shell image.
# COPY ca-roots.crt /etc/ssl/certs/ca-roots.crt

VOLUME ["/etc/didctl"]

EXPOSE 9090/tcp

ENTRYPOINT ["/usr/bin/didctl"]
