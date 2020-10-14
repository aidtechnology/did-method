FROM ghcr.io/bryk-io/shell:0.1.0

EXPOSE 9090/tcp

VOLUME ["/etc/didctl"]

COPY didctl /usr/bin/didctl
ENTRYPOINT ["/usr/bin/didctl"]
