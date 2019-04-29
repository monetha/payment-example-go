SHELL := bash
PACKAGE_NAME := github.com/monetha/payment-example
ARTIFACTS_DIR := $(if $(ARTIFACTS_DIR),$(ARTIFACTS_DIR),bin)

PKGS ?= $(shell glide novendor)
PKGS_NO_CMDS ?= $(shell glide novendor | grep -v ./cmd/)
BENCH_FLAGS ?= -benchmem

VERSION := $(if $(TRAVIS_TAG),$(TRAVIS_TAG),$(if $(TRAVIS_BRANCH),$(TRAVIS_BRANCH),development_in_$(shell git rev-parse --abbrev-ref HEAD)))
COMMIT := $(if $(TRAVIS_COMMIT),$(TRAVIS_COMMIT),$(shell git rev-parse HEAD))
BUILD_TIME := $(shell TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ')

CMD_GO_LDFLAGS := '-X "$(PACKAGE_NAME)/cmd.Version=$(VERSION)" -X "$(PACKAGE_NAME)/cmd.BuildTime=$(BUILD_TIME)" -X "$(PACKAGE_NAME)/cmd.GitHash=$(COMMIT)"'

.PHONY: all
all: lint

.PHONY: dependencies
dependencies:
	@echo "Installing Glide and locked dependencies..."
	glide --version || go get -u -f github.com/Masterminds/glide
	glide install
	@echo "Installing goimports..."
	go install ./vendor/golang.org/x/tools/cmd/goimports
	@echo "Installing golint..."
	go install ./vendor/golang.org/x/lint/golint
	@echo "Installing gosimple..."
	go install ./vendor/honnef.co/go/tools/cmd/gosimple
	@echo "Installing unused..."
	go install ./vendor/honnef.co/go/tools/cmd/unused
	@echo "Installing staticcheck..."
	go install ./vendor/honnef.co/go/tools/cmd/staticcheck

.PHONY: lint
lint:
	@echo "Checking formatting..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | grep -v mock) && [ -z "$$gofiles" ] || unformatted=$$(for d in $$gofiles; do goimports -l $$d/*.go; done) && [ -z "$$unformatted" ] || (echo >&2 "Go files must be formatted with goimports. Following files has problem:\n$$unformatted" && false)
	@echo "Checking vet..."
	@go vet $(PKG_FILES)
	@echo "Checking simple..."
	@gosimple $(PKG_FILES)
	@echo "Checking unused..."
	@unused $(PKG_FILES)
	@echo "Checking staticcheck..."
	@staticcheck $(PKG_FILES)
	@echo "Checking lint..."
	@$(foreach dir,$(PKGS),golint $(dir);)

.PHONY: fmt
fmt:
	@echo "Formatting files..."
	@gofiles=$$(go list -f {{.Dir}} $(PKGS) | grep -v mock) && [ -z "$$gofiles" ] || for d in $$gofiles; do goimports -l -w $$d/*.go; done