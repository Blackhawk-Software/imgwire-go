SHELL := /bin/sh

NODE ?= yarn
GO ?= go
GOFMT ?= gofmt

.PHONY: help install install-js generate verify-generated test build format format-go format-js release-set clean ci

help:
	@printf "%s\n" \
		"Targets:" \
		"  make install            Install Yarn tooling with frozen lockfile" \
		"  make install-js         Install Yarn tooling with frozen lockfile" \
		"  make generate           Regenerate checked-in OpenAPI and generated client artifacts" \
		"  make verify-generated   Verify checked-in generated artifacts are current" \
		"  make test               Run Go tests" \
		"  make build              Build all Go packages" \
		"  make format             Run Go and repo metadata/doc formatting" \
		"  make format-go          Run gofmt on handwritten Go files" \
		"  make format-js          Run Prettier for repo metadata and docs" \
		"  make release-set VERSION=X.Y.Z  Set the repo tooling version manually" \
		"  make clean              Remove local build artifacts and temporary caches" \
		"  make ci                 Run generation verification, tests, and build"

install: install-js

install-js:
	$(NODE) install --frozen-lockfile

generate:
	$(NODE) generate

verify-generated:
	$(NODE) verify-generated

test:
	$(GO) test ./...

build:
	$(GO) build ./...

format:
	$(MAKE) format-go
	$(MAKE) format-js

format-go:
	find . -path ./generated -prune -o -name '*.go' -print | xargs $(GOFMT) -w

format-js:
	$(NODE) format

release-set:
	@test -n "$(VERSION)" || (echo "VERSION is required. Usage: make release-set VERSION=0.2.0" && exit 1)
	$(NODE) release:set-version $(VERSION)

clean:
	rm -rf /tmp/imgwire-go-buildcache /tmp/imgwire-go-tmp

ci: verify-generated test build
