LIMITER_BIN := "./bin/limiter"
LIMITER_CLI_BIN := "./bin/limiter-cli"

DOCKER_IMG := "limiter:0.1"
COMPOSE_CONFIG := ./deployments/compose.yaml
TEST_COMPOSE_CONFIG := ./deployments/compose.test.yaml
PB_PATH := internal/server/grpc/pb

GIT_HASH := $(shell git log --format="%h" -n 1)
LDFLAGS := -X main.release="develop" -X main.buildDate=$(shell date -u +%Y-%m-%dT%H:%M:%S) -X main.gitHash=$(GIT_HASH)

build:
	go build -v -o $(LIMITER_BIN) -ldflags "$(LDFLAGS)" ./cmd/limiter

build-cli:
	go build -v -o $(LIMITER_CLI_BIN) ./cmd/cli

run: build
	$(LIMITER_BIN) -config ./configs/limiter.yaml

build-img:
	docker build \
		--build-arg=LDFLAGS="$(LDFLAGS)" \
		-t $(DOCKER_IMG) \
		-f build/limiter.Dockerfile .

run-img: build-img
	docker run $(DOCKER_IMG)

generate:
	protoc \
		--proto_path=api \
		--go_out=$(PB_PATH) --go_opt=paths=source_relative \
		--go-grpc_out=$(PB_PATH) --go-grpc_opt=paths=source_relative \
		rate_limiter.proto

up:
	docker compose -f $(COMPOSE_CONFIG) up -d

down:
	docker compose -f $(COMPOSE_CONFIG) down -v

version: build
	$(LIMITER_BIN) version

mock:
	docker run --rm -w /src -v $(shell pwd):/src vektra/mockery

test:
	go test -race -count=10 ./internal/...

integration-tests:
	docker compose -f $(COMPOSE_CONFIG) -f $(TEST_COMPOSE_CONFIG) up -d ;\
	test_status_code=0 ;\
	docker compose -f $(COMPOSE_CONFIG) -f $(TEST_COMPOSE_CONFIG) run integration_tests go test -v /app/tests || test_status_code=$$? ;\
	docker compose -f $(COMPOSE_CONFIG) -f $(TEST_COMPOSE_CONFIG) down -v ;\
	exit $$test_status_code ;

install-lint-deps:
	(which golangci-lint > /dev/null) || curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(shell go env GOPATH)/bin v1.54.2

lint: install-lint-deps
	golangci-lint run ./...

.PHONY: build run build-img run-img generate up down version mock test integration-tests lint
