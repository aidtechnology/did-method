FROM scratch

ARG VERSION
LABEL maintainer="Ben Cessa <ben@bryk.io>"
LABEL version=${VERSION}

COPY bryk-did-agent_linux /usr/bin/agent
COPY bryk-did-client_linux /usr/bin/client
COPY ca-roots.crt /etc/ssl/certs/ca-roots.crt

VOLUME ["/tmp", "/etc/bryk-did", "/etc/bryk-did/agent"]

EXPOSE 9090/tcp
