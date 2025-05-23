include .env

GOLANGCI_LINT_CACHE?=/tmp/gophkeeper-golangci-lint-cache
USER=CURRENT_UID=$$(id -u):0
DOCKER_PROJECT_NAME=gophkeeper
DATABASE_DSN="postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable"
SRV_LD_FLAGS=-ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit='test' -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'"
CLT_LD_FLAGS=-ldflags "-X main.buildVersion=v1.0.1 -X main.buildCommit='test' -X 'main.buildDate=$(shell date +'%Y/%m/%d %H:%M:%S')'"


gofmt:
	gofmt -s -w ./
.PHONY: gofmt


image-server:
	$(USER) docker build \
	--build-arg SERVER_BUILD_VERSION=$(SERVER_BUILD_VERSION) \
	--build-arg SERVER_BUILD_COMMIT=$(SERVER_BUILD_COMMIT) \
	-t "gophkeeper-server" ./
.PHONY: image-server


containers: image-server
	$(USER) docker-compose --project-name $(DOCKER_PROJECT_NAME) up -d
.PHONY: containers


client: client-run
.PHONY: client


client-run: client-build
	./cmd/client/client
.PHONY: client-run


client-build:
	go build \
	${CLT_LD_FLAGS} \
	-o ./cmd/client/client ./cmd/client/
.PHONY: client-build


client-cross-build: client-amd64-windows client-amd64-linux client-amd64-darwin
.PHONY: client-cross-build


client-amd64-windows:
	GOOS=windows GOARCH=amd64 go build \
	${CLT_LD_FLAGS} \
	-o ./bin/client-amd64-windows.exe ./cmd/client/
.PHONY: client-amd64-windows


client-amd64-linux:
	GOOS=linux GOARCH=amd64 go build \
	${CLT_LD_FLAGS} \
	-o ./bin/client-amd64-linux ./cmd/client/
.PHONY: client-amd64-linux


client-amd64-darwin:
	GOOS=darwin GOARCH=amd64 go build \
	${CLT_LD_FLAGS} \
	-o ./bin/client-amd64-darwin ./cmd/client/
.PHONY: client-amd64-darwin


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
	go build \
	${SRV_LD_FLAGS} \
	-o ./cmd/server/server ./cmd/server/
.PHONY: server-build


migrate-up:
	migrate -path ./cmd/server/db/migrations -database ${DATABASE_DSN} up
.PHONY: migrate-up


migrate-down:
	migrate -path ./cmd/server/db/migrations -database ${DATABASE_DSN} down
.PHONY: migrate-down


migrate-create:
	migrate create -ext sql -dir ./cmd/server/db/migrations $(name)
.PHONY: migrate-create


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