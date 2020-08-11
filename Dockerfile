FROM docker.pkg.github.com/bryk-io/base-images/shell:0.1.0

COPY didctl_linux_amd64 /usr/bin/didctl

VOLUME ["/etc/didctl"]

EXPOSE 9090/tcp

ENTRYPOINT ["/usr/bin/didctl"]
