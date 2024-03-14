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

GOLANG_CROSS_VERSION  ?= v1.21.5
PACKAGE_NAME          := github.com/navikt/knaudit-proxy
release:
	@docker run \
		--rm \
		-v /var/run/docker.sock:/var/run/docker.sock \
		-v `pwd`:/go/src/$(PACKAGE_NAME) \
		-w /go/src/$(PACKAGE_NAME) \
		ghcr.io/goreleaser/goreleaser-cross:${GOLANG_CROSS_VERSION} \
		release --clean --skip validate --skip publish --snapshot --debug --verbose \
		--config /go/src/$(PACKAGE_NAME)/.goreleaser.local.yml

builder:
	docker buildx create --name mybuilder --use --platform linux/amd64,linux/arm64
	docker buildx inspect --bootstrap


PLATFORMS         ?= linux/amd64,linux/arm64
BUILDX_EXTRA_ARGS ?=
image: | release
	docker buildx build --platform $(PLATFORMS) --file Dockerfile -t knaudit-proxy:latest $(BUILDX_EXTRA_ARGS) .

lint: $(GOLANGCILINT)
	$(GOLANGCILINT) run
.PHONY: lint

check: | lint test
.PHONY: check
