include .env

GOLANGCI_LINT_CACHE?=/tmp/gophkeeper-golangci-lint-cache
USER=CURRENT_UID=$$(id -u):0
DOCKER_PROJECT_NAME=gophkeeper


gofmt:
	gofmt -s -w ./
.PHONY: gofmt


containers:
	$(USER) docker-compose --project-name $(DOCKER_PROJECT_NAME) up -d
.PHONY: containers


client: client-run
.PHONY: client


client-run: client-build
	./cmd/client/client
.PHONY: client-run


client-build:
	go build -o ./cmd/client/client ./cmd/client/
.PHONY: client-build


client-build-cover:
	go build -cover -o ./cmd/client/client ./cmd/client/
.PHONY: client-build-cover


test: client-build-cover
	go test ./... -cover -coverprofile=coverage.out && \
	go tool cover -html coverage.out -o test.html && \
	go tool cover -func=coverage.out
.PHONY: test


server: server-run
.PHONY: server


server-run: server-build
	./cmd/server/server
.PHONY: server-run


server-build:
	go build -o ./cmd/server/server ./cmd/server/
.PHONY: server-build


proto:
	protoc --go_out=. --go_opt=paths=source_relative \
  		--go-grpc_out=. --go-grpc_opt=paths=source_relative \
  		internal/proto/server.proto
.PHONY: proto


golangci-lint-run: _golangci-lint-rm-unformatted-report
.PHONY: golangci-lint-run


_golangci-lint-reports-mkdir:
	mkdir -p ./golangci-lint
.PHONY: _golangci-lint-reports-mkdir


_golangci-lint-run: _golangci-lint-reports-mkdir
	-docker run --rm \
    -v $(shell pwd):/app \
    -v $(GOLANGCI_LINT_CACHE):/root/.cache \
    -w /app \
    golangci/golangci-lint:v1.57.2 \
        golangci-lint run \
            -c .golangci.yml \
	> ./golangci-lint/report-unformatted.json
.PHONY: _golangci-lint-run


_golangci-lint-format-report: _golangci-lint-run
	cat ./golangci-lint/report-unformatted.json | jq > ./golangci-lint/report.json
.PHONY: _golangci-lint-format-report


_golangci-lint-rm-unformatted-report: _golangci-lint-format-report
	rm ./golangci-lint/report-unformatted.json
.PHONY: _golangci-lint-rm-unformatted-report


golangci-lint-clean:
	sudo rm -rf ./golangci-lint 
.PHONY: golangci-lint-clean