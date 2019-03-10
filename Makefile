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
	gofmt -s -l $(FILES_LIST)
	golint -set_exit_status $(GO_PKG_LIST)
	misspell $(FILES_LIST)

	# Static analysis
	ineffassign $(FILES_LIST)
	GO111MODULE=off gosec ./...
	gocyclo -over 15 `find . -iname '*.go' | grep -v 'vendor' | grep -v '_test.go' | grep -v 'pb.go' | grep -v 'pb.gw.go'`
	go-consistent -v ./...

	# Unit tests
	go test -race -cover -v $(GO_PKG_LIST)

build: ## Build for the default architecture in use
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-client github.com/bryk-io/did-method/client/cli
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent github.com/bryk-io/did-method/agent/cli

linux: ## Build for 64 bit Linux systems
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-client_linux github.com/bryk-io/did-method/client/cli
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent_linux github.com/bryk-io/did-method/agent/cli

mac: ## Build for 64 bit MacOS systems
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-client_darwin github.com/bryk-io/did-method/client/cli
	CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent_darwin github.com/bryk-io/did-method/agent/cli

windows: ## Build for 32 and 64 bit Windows systems
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-client_windows_64bit.exe github.com/bryk-io/did-method/client/cli
	CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent_windows_64bit.exe github.com/bryk-io/did-method/agent/cli
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-client_windows_32bit.exe github.com/bryk-io/did-method/client/cli
	CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME)-agent_windows_32bit.exe github.com/bryk-io/did-method/agent/cli

install: ## Install the binary to GOPATH and keep cached all compiled artifacts
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME)-client github.com/bryk-io/did-method/client/cli

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
	@-rm bryk-did-agent_linux bryk-did-client_linux ca-roots.crt
	@make ca-roots
	@make linux
	@-docker rmi $(BINARY_NAME):$(VERSION_TAG)
	@docker build --build-arg VERSION="$(VERSION_TAG)" --rm -t $(BINARY_NAME):$(VERSION_TAG) .
	@-rm bryk-did-agent_linux bryk-did-client_linux ca-roots.crt

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'
