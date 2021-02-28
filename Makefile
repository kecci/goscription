PROJECTNAME := $(shell go run app/main.go project)
VERSION := $(shell go run app/main.go version)
RELEASENAME := $(PROJECTNAME)_$(VERSION)

MYSQL_DB ?= $(shell go run app/main.go mysql)
MYSQL_DIR ?= database/migrations

ifndef $(GOPATH)
    GOPATH=$(shell go env GOPATH)
    export GOPATH
endif

all:

.PHONY: lint-prepare
lint-prepare:
	@echo "Preparing Linter"
ifeq (,$(shell which golangci-lint))
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin latest
else
	@echo "golangci-lint already installed"
endif

## lint, l: run linter to check quality code
.PHONY: lint l
lint l:
	@make lint-prepare
	golangci-lint run -v \
		--exclude-use-default=false \
		--disable-all \
		--enable=goimports \
		--enable=gocyclo \
		--enable=nakedret \
		--enable=golint \
		--enable=gosimple \
		--enable=goconst \
		--enable=misspell \
		--enable=unconvert \
		--enable=varcheck \
		--enable=unused \
		--enable=deadcode \
		--enable=unparam \
		--enable=ineffassign \
		--enable=gocritic \
		--enable=prealloc \
		--enable=scopelint \
		--enable=staticcheck \
		--enable=gosec \
		./...

.PHONY: mockery-prepare
mockery-prepare:
	@echo "prepare mockery"
ifeq (,$(shell which mockery))
	@echo "Installing mockery"
	@go get github.com/vektra/mockery/v2/
else
	@echo "mockery already installed"
endif

## mockery-generate: generate all mock
.PHONY: mockery-generate
mockery-generate:
	@make mockery-prepare
	@echo "start mockery all"
	@mockery --all

## test-unit, t: run all unit test
.PHONY: test-unit t
test-unit t:
	@echo "Start: unit test"
	@go test -v -short -race ./...
	@echo "End: unit test"

## go-build, b: build to compile the project & swagger docs
.PHONY: go-build b
go-build b:
	@echo "Start: build "$(PROJECTNAME)" project"
	@make swagger-init
	@make swagger-validate
	@echo "Run: go build -o bin/"$(PROJECTNAME)" app/main.go"
	@go build -o bin/"$(PROJECTNAME)" app/main.go
	@echo "Done: "$(PROJECTNAME)" project has build"

## go-run, r: run http rest api project
.PHONY: go-run r
go-run r:
	@echo "Run "$(PROJECTNAME)" project"
	@make swagger-init
	@make swagger-validate
	@go run app/main.go

## worker, w: run worker for consumer
.PHONY: worker w
worker w:
	@echo "Run "$(PROJECTNAME)" worker"
	@go run app/main.go worker

## docker-build: dockerize the project
.PHONY: docker-build
docker-build:
	@docker build . -f build/builder/Dockerfile -t "$(PROJECTNAME)":latest

## docker-up: run docker compose up
.PHONY: docker-up
docker-up:
	@docker-compose up -d

## docker-down: run docker compose down
.PHONY: docker-down
docker-down:
	@docker-compose down

.PHONY: swagger-init-prepare
swagger-init-prepare:
	@echo "Prepare swag"
ifeq (,$(shell which swag))
	go get github.com/swaggo/swag/cmd/swag
else
	@echo "swag already installed"
endif

.PHONY: swagger-validate-prepare
swagger-validate-prepare:
	@echo "Prepare go-swagger"
ifeq (,$(shell which swagger))
	go get github.com/go-swagger/go-swagger/cmd/swagger
else
	@echo "go-swagger already installed"
endif

## swagger-init: initialize swagger to folder ./api/docs
.PHONY: swagger-init
swagger-init:
	@make swagger-init-prepare
	@echo "Start: Initialize swagger"
	swag init -g app/main.go -o ./api/docs
	@echo "Done: initialize ./api/docs swagger"

## swagger-validate: validate swagger.yaml in folder ./api/docs
.PHONY: swagger-validate
swagger-validate:
	@make swagger-validate-prepare
	@echo "Start: Validate swagger"
	swagger validate api/docs/swagger.yaml
	@echo "Done: Validate swagger"

.PHONY: migrate-prepare
migrate-prepare:
	@echo "prepare golang migrate using mysql"
ifeq (,$(wildcard ./bin/migrate))
	@go build -a -o ./bin/migrate -tags 'mysql' github.com/golang-migrate/migrate/v4/cli
else
	@echo "golang-migrate already installed"
endif	

## migrate-create [file_name]: create migration file
.PHONY: migrate-create
migrate-create:
	@make migrate-prepare
	./bin/migrate create -ext sql -dir $(MYSQL_DIR) -seq $(filter-out $@,$(MAKECMDGOALS))

## migrate-up [N]: Apply all or N up migrations
.PHONY: migrate-up
migrate-up:
	@make migrate-prepare
	./bin/migrate -database "$(MYSQL_DB)" -path=$(MYSQL_DIR) up $(filter-out $@,$(MAKECMDGOALS))

## migrate-down [N]: Apply all or N down migrations
.PHONY: migrate-down
migrate-down:
	@make migrate-prepare
	./bin/migrate -database "$(MYSQL_DB)" -path=$(MYSQL_DIR) down $(filter-out $@,$(MAKECMDGOALS))

## release: create release app with version
.PHONY: release
release:
	@echo "Start: build release "$(RELEASENAME)""
	@make lint
	@make test-unit
	@make swagger-init
	@make swagger-validate

	@rm -rf release
	@mkdir release

	@echo "Build release "$(RELEASENAME)"_linux_amd64"
	@GOOS=linux GOARCH=amd64 go build -o release/"$(RELEASENAME)"_linux_amd64 app/main.go

	@echo "Build release "$(RELEASENAME)"_darwin_amd64"
	@GOOS=darwin GOARCH=amd64 go build -o release/"$(RELEASENAME)"_darwin_amd64 app/main.go

	@echo "Done: build release "$(RELEASENAME)""

.PHONY: help
all: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo