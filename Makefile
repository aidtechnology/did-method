.PHONY: *
.DEFAULT_GOAL:=help

# Project setup
BINARY_NAME=didctl
DOCKER_IMAGE=docker.pkg.github.com/bryk-io/did-method/didctl
MAINTAINERS='Ben Cessa <ben@pixative.com>'

# State values
GIT_COMMIT_DATE=$(shell TZ=UTC git log -n1 --pretty=format:'%cd' --date='format-local:%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT_HASH=$(shell git log -n1 --pretty=format:'%H')
GIT_TAG=$(shell git describe --tags --always --abbrev=0 | cut -c 1-8)

# Linker tags
# https://golang.org/cmd/link/
LD_FLAGS += -s -w
LD_FLAGS += -X github.com/bryk-io/did-method/info.CoreVersion=$(GIT_TAG)
LD_FLAGS += -X github.com/bryk-io/did-method/info.BuildCode=$(GIT_COMMIT_HASH)
LD_FLAGS += -X github.com/bryk-io/did-method/info.BuildTimestamp=$(GIT_COMMIT_DATE)

# Proto builder basic setup
proto-builder=docker run --rm -it -v $(shell pwd):/workdir \
docker.pkg.github.com/bryk-io/base-images/buf-builder:0.20.5

## help: Prints this help message
help:
	@echo "Commands available"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | sort

## build: Build for the current architecture in use, intended for development
build:
	go build -v -ldflags '$(LD_FLAGS)' -o $(BINARY_NAME) github.com/bryk-io/did-method/client/cli

## build-for: Build the available binaries for the specified 'os' and 'arch'
build-for:
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags '$(LD_FLAGS)' \
	-o $(BINARY_NAME)_$(os)_$(arch)$(suffix) github.com/bryk-io/did-method/client/cli

## ca-roots: Generate the list of valid CA certificates
ca-roots:
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

## clean: Verify dependencies and remove intermediary products
clean:
	@-rm -rf vendor
	go clean
	go mod tidy
	go mod verify
	go mod vendor

## docker: Build docker image
docker:
	make build-for os=linux arch=amd64
	@-docker rmi $(DOCKER_IMAGE):$(GIT_TAG)
	@docker build \
	"--label=org.opencontainers.image.title=$(BINARY_NAME)" \
	"--label=org.opencontainers.image.authors=$(MAINTAINERS)" \
	"--label=org.opencontainers.image.created=$(GIT_COMMIT_DATE)" \
	"--label=org.opencontainers.image.revision=$(GIT_COMMIT_HASH)" \
	"--label=org.opencontainers.image.version=$(GIT_TAG)" \
	--rm -t $(DOCKER_IMAGE):$(GIT_TAG) .
	@rm $(BINARY_NAME)_linux_amd64

## install: Install the binary to GOPATH and keep cached all compiled artifacts
install:
	@go build -v -ldflags '$(LD_FLAGS)' -i -o ${GOPATH}/bin/$(BINARY_NAME) github.com/bryk-io/did-method/client/cli

## lint: Static analysis
lint:
	# Code
	golangci-lint run -v ./...

	# Helm charts
	helm lint helm/*

## proto: Compile all PB definitions and RPC services
proto:
	# Verify style and consistency
	$(proto-builder) buf check lint --file $(shell echo proto/v1/*.proto | tr ' ' ',')
	@-$(proto-builder) buf check breaking \
	--file $(shell echo proto/v1/*.proto | tr ' ' ',') \
	--against-input proto/v1/image.bin

	# Clean old builds
	@-rm proto/v1/image.bin

	# Build package image
	$(proto-builder) buf image build -o proto/v1/image.bin --file $(shell echo proto/v1/*.proto | tr ' ' ',')

	# Build package code
	$(proto-builder) buf protoc \
	--proto_path=proto \
	--go_out=proto \
	--go-grpc_out=proto \
	--grpc-gateway_out=logtostderr=true:proto \
	--swagger_out=logtostderr=true:proto \
	--govalidators_out=proto \
	proto/v1/*.proto

	# Remove package comment added by the gateway generator to avoid polluting
	# the package documentation.
	@-sed -i '' '/\/\*/,/*\//d' proto/v1/*.pb.gw.go

	# Style adjustments
	gofmt -s -w proto/v1
	goimports -w proto/v1

## release: Prepare assets for a new tagged release
release:
	goreleaser release --skip-validate --skip-publish --rm-dist

## scan: Look for known vulnerabilities in the project dependencies
# https://github.com/sonatype-nexus-community/nancy
scan:
	@nancy -quiet go.sum

## test: Run unit tests excluding the vendor dependencies
test:
	go test -race -v -failfast -coverprofile=coverage.report ./...
	go tool cover -html coverage.report -o coverage.html

## updates: List available updates for direct dependencies
# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
updates:
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null
