.PHONY: *
.DEFAULT_GOAL:=help

# Project setup variables
BINARY_NAME=didctl
DOCKER_IMAGE=docker.pkg.github.com/bryk-io/did-method/didctl
VERSION_TAG=0.6.0

# Linker tags
# https://golang.org/cmd/link/
LD_FLAGS += -s -w
LD_FLAGS += -X github.com/bryk-io/did-method/info.CoreVersion=$(VERSION_TAG)
LD_FLAGS += -X github.com/bryk-io/did-method/info.BuildTimestamp=$(shell date +'%s')
LD_FLAGS += -X github.com/bryk-io/did-method/info.BuildCode=$(shell git log --pretty=format:'%H' -n1)

# Proto builder basic setup
proto-builder=docker run --rm -it -v $(shell pwd):/workdir \
docker.pkg.github.com/bryk-io/base-images/buf-builder:0.20.5

## help: Prints this help message
help:
	@echo "Commands available"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /' | sort

## lint: Static analysis
lint:
	# Code
	golangci-lint run -v ./...

	# Helm charts
	helm lint helm/*

## test: Run unit tests excluding the vendor dependencies
test:
	go test -race -v -failfast -coverprofile=coverage.report ./...
	go tool cover -html coverage.report -o coverage.html

## updates: List available updates for direct dependencies
# https://github.com/golang/go/wiki/Modules#how-to-upgrade-and-downgrade-dependencies
updates:
	@go list -u -f '{{if (and (not (or .Main .Indirect)) .Update)}}{{.Path}}: {{.Version}} -> {{.Update.Version}}{{end}}' -m all 2> /dev/null

## scan: Look for known vulnerabilities in the project dependencies
# https://github.com/sonatype-nexus-community/nancy
scan:
	@nancy -quiet go.sum

## release: Prepare assets for a new tagged release
release:
	@-rm -rf release-$(VERSION_TAG)
	mkdir release-$(VERSION_TAG)
	make build-for os=linux arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=darwin arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=windows arch=amd64 suffix=".exe" dest=release-$(VERSION_TAG)/

## build: Build for the current architecture in use, intended for development
build:
	go build -v -ldflags '$(LD_FLAGS)' -o $(BINARY_NAME) github.com/bryk-io/did-method/client/cli

## build-for: Build the available binaries for the specified 'os' and 'arch'
build-for:
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags '$(LD_FLAGS)' \
	-o $(dest)$(BINARY_NAME)_$(VERSION_TAG)_$(os)_$(arch)$(suffix) github.com/bryk-io/did-method/client/cli
	chmod +x $(dest)$(BINARY_NAME)_$(VERSION_TAG)_$(os)_$(arch)$(suffix)

## install: Install the binary to GOPATH and keep cached all compiled artifacts
install:
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME) github.com/bryk-io/did-method/client/cli

## clean: Verify dependencies and remove intermediary products
clean:
	@-rm -rf vendor
	go clean
	go mod tidy
	go mod verify
	go mod vendor

## ca-roots: Generate the list of valid CA certificates
ca-roots:
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

## docker: Build docker image
docker:
	@make build-for os=linux arch=amd64
	@-docker rmi $(DOCKER_IMAGE):$(VERSION_TAG)
	@docker build --build-arg VERSION_TAG="$(VERSION_TAG)" --rm -t $(DOCKER_IMAGE):$(VERSION_TAG) .
	@-rm $(BINARY_NAME)_$(VERSION_TAG)_linux_amd64

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
