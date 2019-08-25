WORKING_DIR := $(shell pwd)

.DEFAULT_GOAL := build

.PHONY: build push tests

install-deps:: ## Download and installs dependencies
		@go get ./cmd/bitbadger/...

build:: install-deps ## Build command line binary
		@go build cmd/bitbadger/bitbadger.go

build-static-linux:: install-deps ## Builds a static linux binary
		@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
			go build \
			-o bin/bitbadger \
			-a -ldflags '-extldflags "-static"' \
				cmd/bitbadger/bitbadger.go

tests:: ## Run tests
		@cd test && go test

install:: ## Build and install bitbadger locally
		@cd cmd/bitbadger/ && go install .

# A help target including self-documenting targets (see the awk statement)
define HELP_TEXT
Usage: make [TARGET]... [MAKEVAR1=SOMETHING]...

Available targets:
endef
export HELP_TEXT
help: ## This help target
	@echo
	@echo "$$HELP_TEXT"
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / \
		{printf "\033[36m%-30s\033[0m  %s\n", $$1, $$2}' $(MAKEFILE_LIST)