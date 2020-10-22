PROJECTNAME := $(shell basename "$(PWD)")

MYSQL_USER ?= user
MYSQL_PASSWORD ?= password
MYSQL_ADDRESS ?= 127.0.0.1:3306
MYSQL_DATABASE ?= article

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

all:

.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
	curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh| sh -s latest

.PHONY: lint
lint:
	./bin/golangci-lint run \
		--exclude-use-default=false \
		--enable=golint \
		--enable=gocyclo \
		--enable=goconst \
		--enable=unconvert \
		./...

## mockery-prepare: install mockery before generate
.PHONY: mockery-prepare
mockery-prepare:
	@echo "Remove the existing one if exists"
	@rm -rf $(GOPATH)/bin/mockery
	@echo "Installing mockery"
	@go get github.com/vektra/mockery/.../

## mockery-generate: generate all mock
mockery-generate:
	@$(GOPATH)/bin/mockery -all

.PHONY: mysql-test-up
mysql-test-up:
	@docker-compose up -d mysql_test

.PHONY: mysql-down-test
mysql-down-test:
	@docker-compose stop mysql_test

## test-docker: test docker integration
.PHONY: test-docker
test-docker: mysql-test-up
	@go test -v -race ./...
	@make mysql-down-test

## test-unit: run all unit test & lint
.PHONY: test-unit
test-unit:
	@echo "Start: unit test"
	@go test -v -short -race ./...
	@echo "gofmt: start..."
	@gofmt -l -e -d .
	@echo "gofmt: done"
	@echo "golint: start..."
	@golint ./...
	@echo "golint: done"
	@echo "End: unit test"

## go-build: build to compile the project & swagger docs
go-build:
	@echo "Start: build "$(PROJECTNAME)" project"
	@make test-unit
	@make swagger-init
	@make swagger-validate
	@echo "Run: go build -o "$(PROJECTNAME)" app/main.go"
	@go build -o "$(PROJECTNAME)" app/main.go
	@echo "Done: "$(PROJECTNAME)" project has build"

## go-run: run project
go-run:
	@make go-build
	@echo "Run "$(PROJECTNAME)" project"
	./$(PROJECTNAME) http

## docker-build: dockerize the project
.PHONY: docker-build
docker-build:
	@make swagger-init
	@make swagger-validate
	@docker build . -t "$(PROJECTNAME)":latest

## docker-up: run docker compose up
.PHONY: docker-up
docker-up:
	@docker-compose up -d
	@make mysql-down-test

## docker-down: run docker compose down
.PHONY: docker-down
docker-down:
	@docker-compose down

## swagger-init: initialize swagger to folder ./docs
.PHONY: swagger-init
swagger-init:
	@echo "Start: Initialize swagger"
	@swag init -g app/main.go
	@echo "Done: initialize ./docs swagger"

## swagger-validate: validate swagger.yaml in folder ./docs
.PHONY: swagger-validate
swagger-validate:
	@echo "Start: Validate swagger"
	@swagger validate docs/swagger.yaml
	@echo "Done: Validate swagger"

## migrate-prepare: prepare migrate the schema with mysql
.PHONY: migrate-prepare
migrate-prepare:
	@go get -u github.com/golang-migrate/migrate/v4
	@go build -a -o ./bin/migrate -tags 'mysql' github.com/golang-migrate/migrate/v4/cli

## migrate-up: run migration up to latest version
.PHONY: migrate-up
migrate-up:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=internal/database/mysql/migrations up	

## migrate-down: run migration down to oldest version
.PHONY: migrate-down
migrate-down:
	@migrate -database "mysql://$(MYSQL_USER):$(MYSQL_PASSWORD)@tcp($(MYSQL_ADDRESS))/$(MYSQL_DATABASE)" \
	-path=internal/database/mysql/migrations down

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo