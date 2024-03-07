GOPATH := $(shell go env GOPATH)
GOBIN  ?= $(GOPATH)/bin # Default GOBIN if not set

# A template function for installing binaries
define install-binary
	 @if ! command -v $(1) &> /dev/null; then \
		  echo "$(1) not found, installing..."; \
		  go install $(2); \
	 fi
endef

GOLANGCILINT         ?= $(shell command -v golangci-lint || echo "$(GOBIN)/golangci-lint")
GOLANGCILINT_VERSION := v1.55.2
GOTEST               ?= $(shell command -v gotest || echo "$(GOBIN)/gotest")
GOTEST_VERSION       := v0.0.6
STATICCHECK          ?= $(shell command -v staticcheck || echo "$(GOBIN)/staticcheck")
STATICCHECK_VERSION  := v0.4.6

$(GOLANGCILINT):
	$(call install-binary,golangci-lint,github.com/golangci/golangci-lint/cmd/golangci-lint@$(GOLANGCILINT_VERSION))

$(GOTEST):
	$(call install-binary,gotest,github.com/rakyll/gotest@$(GOTEST_VERSION))

$(STATICCHECK):
	$(call install-binary,staticcheck,honnef.co/go/tools/cmd/staticcheck@$(STATICCHECK_VERSION))

test: $(GOTEST)
	$(GOTEST) -v ./... -count=1
.PHONY: test

lint: $(GOLANGCILINT)
	$(GOLANGCILINT) run
.PHONY: lint

check: | lint test
.PHONY: check
