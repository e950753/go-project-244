COVERAGE_MIN ?= 80
COVERPKG := code,code/diff,code/parsers,code/formatters

.PHONY: build test lint test-coverage

build:
	go build -o bin/gendiff ./cmd/gendiff

test:
	go test ./...

lint:
	golangci-lint run ./...

test-coverage:
	go test -coverpkg=$(COVERPKG) -coverprofile=coverage.out ./...
	@total=$$(go tool cover -func=coverage.out | awk '/^total:/ {gsub("%","",$$3); print $$3}'); \
	echo "total coverage: $$total% (min $(COVERAGE_MIN)%)"; \
	awk -v t="$$total" -v m="$(COVERAGE_MIN)" 'BEGIN { exit (t+0 < m+0) }' \
	  || { echo "coverage is below $(COVERAGE_MIN)%"; exit 1; }
