FROM scratch

ARG VERSION_TAG
LABEL maintainer="Ben Cessa <ben@bryk.io>"
LABEL version=${VERSION_TAG}

COPY bryk-did-agent_${VERSION_TAG}_linux_amd64 /usr/bin/agent
COPY bryk-did_${VERSION_TAG}_linux_amd64 /usr/bin/client
COPY ca-roots.crt /etc/ssl/certs/ca-roots.crt

VOLUME ["/tmp", "/etc/bryk-did", "/etc/bryk-did/agent"]

EXPOSE 9090/tcp
