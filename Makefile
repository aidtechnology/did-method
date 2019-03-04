.DEFAULT_GOAL := help
FILES_LIST=`find . -iname '*.go' | grep -v 'vendor'`
GO_PKG_LIST=`go list ./... | grep -v 'vendor'`
BINARY_NAME=bryk-id
VERSION_TAG=0.1.0

# Custom compilation tags
LD_FLAGS="\
-X github.com/bryk-io/id/client/cli/cmd.coreVersion=$(VERSION_TAG) \
-X github.com/bryk-io/id/client/cli/cmd.buildCode=`git log --pretty=format:'%H' -n1` \
-X github.com/bryk-io/id/client/cli/cmd.buildTimestamp=`date +'%s'` \
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
	go build -v -ldflags $(LD_FLAGS) -o $(BINARY_NAME) github.com/bryk-io/id/client/cli

install: ## Install the binary to GOPATH and keep cached all compiled artifacts
	@go build -v -ldflags $(LD_FLAGS) -i -o ${GOPATH}/bin/$(BINARY_NAME) github.com/bryk-io/id/client/cli

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

help: ## Display available make targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[33m%-16s\033[0m %s\n", $$1, $$2}'
