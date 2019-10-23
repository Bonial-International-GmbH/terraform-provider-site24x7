.DEFAULT_GOAL := help

GOLANGCI_LINT_VERSION ?= v1.19.1
BINARY := terraform-provider-site24x7

.PHONY: help
help:
	@grep -E '^[a-zA-Z-]+:.*?## .*$$' Makefile | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "[32m%-12s[0m %s\n", $$1, $$2}'

.PHONY: build
build: ## build terraform-provider-site24x7
	go build \
		-ldflags "-s -w" \
		-o $(BINARY) \
		main.go


.PHONY: install
install: build ## install terraform-provider-site24x7
	mkdir -p $(HOME)/.terraform.d/plugins
	cp $(BINARY) $(HOME)/.terraform.d/plugins

.PHONY: test
test: ## run tests
	go test -race -tags="$(TAGS)" $$(go list ./... | grep -v /vendor/)

.PHONY: coverage
coverage: ## generate code coverage
	scripts/coverage

.PHONY: lint
lint: ## run golangci-lint
	command -v golangci-lint > /dev/null 2>&1 || \
	  curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin $(GOLANGCI_LINT_VERSION)
	golangci-lint run
