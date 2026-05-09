# Make
export

SHELL := /bin/bash -o errexit -o nounset -o pipefail

MAKEFLAGS += --warn-undefined-variables
MAKEFLAGS += --no-builtin-rules

VERBOSE ?= false
ifeq (${VERBOSE}, false)
	# --silent drops the need to prepend `@` to suppress command output
	MAKEFLAGS += --silent
endif

# Variables
GOBASE       ?= $(shell pwd)
GOBIN        ?= ${GOBASE}/bin
COVER_MIN    ?= 80
COVERPROFILE ?= coverage.out
COVER_PACKAGES ?= ./pkg/...

GOFUMPT_VERSION           ?= v0.5.0
GOIMPORTS_REVISER_VERSION ?= v3.4.1
GOLANGCI_LINT_VERSION     ?= v2.4.0

# Ensure that we use vendored binaries before consulting the system.
PATH := ${GOBIN}:${PATH}

# Applications
GO ?= go

GOLANGCI_LINT     ?= ${GOBIN}/golangci-lint
GOFUMPT           ?= ${GOBIN}/gofumpt
GOIMPORTS_REVISER ?= ${GOBIN}/goimports-reviser

# Dependencies
.PHONY: depend
depend: ## Update project dependencies
	$(GO) mod tidy

.PHONY: $(GOFUMPT)
$(GOFUMPT):
	$(GO) install mvdan.cc/gofumpt@${GOFUMPT_VERSION}

.PHONY: $(GOIMPORTS_REVISER)
$(GOIMPORTS_REVISER):
	$(GO) install github.com/incu6us/goimports-reviser/v3@${GOIMPORTS_REVISER_VERSION}

.PHONY: $(GOLANGCI_LINT)
$(GOLANGCI_LINT):
	$(GO) install github.com/golangci/golangci-lint/v2/cmd/golangci-lint@${GOLANGCI_LINT_VERSION}

# Helpers
.PHONY: fmt
fmt: $(GOIMPORTS_REVISER) $(GOFUMPT) ## Format source files
	find . -type f -name '*.go' -not -path "./vendor/*" | \
		xargs -I {} $(GOIMPORTS_REVISER) -company-prefixes="github.com/adambrett/" -project-name="github.com/adambrett/go-fyne" {}
	# In some cases you need to run gofumpt twice to resolve all formatting issues as one simplification
	# can allow another one, but gofumpt is not smart enough to apply both at the same time.
	find . -type f -name '*.go' -not -path "./vendor/*" | xargs $(GOFUMPT) -w
	find . -type f -name '*.go' -not -path "./vendor/*" | xargs $(GOFUMPT) -w

# Linting
.PHONY: lint
lint: $(GOLANGCI_LINT) ## Run the linter
	$(GOLANGCI_LINT) run ./...

# Testing
.PHONY: test
test: ## Run tests
	$(GO) test -race ./...

.PHONY: coverage
coverage: ## Run coverage and require total coverage to meet COVER_MIN
	$(GO) test -race -coverprofile=${COVERPROFILE} ${COVER_PACKAGES}
	$(GO) tool cover -func=${COVERPROFILE} | awk -v min="${COVER_MIN}" '{ print } /^total:/ { pct = $$3; sub(/%/, "", pct); if (pct < min) { printf "coverage %.1f%% is below %.1f%%\n", pct, min; exit 1 } }'

# Examples
.PHONY: run-%
run-%: ## Run an example, such as make run-launcher
	$(GO) run ./examples/$*

# Cleaning
.PHONY: clean
clean: ## Clean generated local artifacts
	rm -f ${COVERPROFILE}

.PHONY: clean-go
clean-go: ## Clean Go build artifacts
	$(GO) clean -modcache

.PHONY: clean-all
clean-all: clean clean-go ## Clean all build artifacts

# Make Helpers
.PHONY: help
help: ## Print this help message
	grep -E '^[/a-zA-Z_-]+:.*?## .*$$' ${MAKEFILE_LIST} | sort | awk 'BEGIN {FS = ":|##"}; {printf "%-20s\033[36m%-20s \033[0m %s\n", $$1, $$2, $$4}'

print-%: ## Print the value of a variable
	echo $* = $($*)
