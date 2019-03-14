.DEFAULT_GOAL := help
FILES_LIST=`find . -iname '*.go' | grep -v 'vendor'`
GO_PKG_LIST=`go list ./... | grep -v 'vendor'`
BINARY_NAME=bryk-did
VERSION_TAG=0.1.0

# Custom compilation tags
LD_FLAGS="\
-X github.com/bryk-io/did-method/client/cli/cmd.coreVersion=$(VERSION_TAG) \
-X github.com/bryk-io/did-method/client/cli/cmd.buildCode=`git log --pretty=format:'%H' -n1` \
-X github.com/bryk-io/did-method/client/cli/cmd.buildTimestamp=`date +'%s'` \
"

test: ## Run all tests excluding the vendor dependencies
	# Formatting
	go vet $(GO_PKG_LIST)
	gofmt -s -w $(FILES_LIST)
	golint -set_exit_status $(GO_PKG_LIST)
	misspell $(FILES_LIST)

	# Static analysis
	ineffassign $(FILES_LIST)
	GO111MODULE=off gosec ./...
	gocyclo -over 15 `find . -iname '*.go' | grep -v 'vendor' | grep -v '_test.go' | grep -v 'pb.go' | grep -v 'pb.gw.go'`
	go-consistent -v ./...

	# Unit tests
	go test -race -cover -v $(GO_PKG_LIST)

build: ## Build for the current architecture in use, intended for devevelopment
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME) github.com/bryk-io/did-method/client/cli
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent github.com/bryk-io/did-method/agent/cli

release: ## Build the binaries for a new release
	@-rm -rf release-$(VERSION_TAG)
	mkdir release-$(VERSION_TAG)
	make build-for os=linux arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=darwin arch=amd64 dest=release-$(VERSION_TAG)/
	make build-for os=windows arch=amd64 suffix=".exe" dest=release-$(VERSION_TAG)/
	make build-for os=windows arch=386 suffix=".exe" dest=release-$(VERSION_TAG)/

build-for: ## Build the availabe binaries for the specified 'os' and 'arch'
	# Build client binary
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags $(LD_FLAGS) \
	-o $(dest)$(BINARY_NAME)_$(VERSION_TAG)_$(os)_$(arch)$(suffix) github.com/bryk-io/did-method/client/cli

	# Build agent binary
	CGO_ENABLED=0 GOOS=$(os) GOARCH=$(arch) \
	go build -v -ldflags $(LD_FLAGS) \
	-o $(dest)$(BINARY_NAME)-agent_$(VERSION_TAG)_$(os)_$(arch)$(suffix) github.com/bryk-io/did-method/agent/cli

install: ## Install the binary to GOPATH and keep cached all compiled artifacts
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME) github.com/bryk-io/did-method/client/cli

clean: ## Download and compile all dependencies and intermediary products
	go mod tidy
	go mod verify

.PHONY: proto
proto: ## Compile protocol buffers and RPC services
	prototool lint
	prototool generate

	# Fix gRPC-Gateway generated code
    # https://github.com/grpc-ecosystem/grpc-gateway/issues/229
	sed -i '' "s/empty.Empty/types.Empty/g" proto/*.pb.gw.go

ca-roots: ## Generate the list of valid CA certificates
	@docker run -dit --rm --name ca-roots debian:stable-slim
	@docker exec --privileged ca-roots sh -c "apt update"
	@docker exec --privileged ca-roots sh -c "apt install -y ca-certificates"
	@docker exec --privileged ca-roots sh -c "cat /etc/ssl/certs/* > /ca-roots.crt"
	@docker cp ca-roots:/ca-roots.crt ca-roots.crt
	@docker stop ca-roots

docker: ## Build docker image
	@-rm bryk-did-agent_$(VERSION_TAG)_linux_amd64 bryk-did_$(VERSION_TAG)_linux_amd64 ca-roots.crt
	@make ca-roots
	@make build-for os=linux arch=amd64
	@-docker rmi $(BINARY_NAME):$(VERSION_TAG)
	@docker build --build-arg VERSION_TAG="$(VERSION_TAG)" --rm -t $(BINARY_NAME):$(VERSION_TAG) .
	@-rm bryk-did-agent_$(VERSION_TAG)_linux_amd64 bryk-did_$(VERSION_TAG)_linux_amd64 ca-roots.crt

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'
