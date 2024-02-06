# Change these variables as necessary.
MAIN_PACKAGE_PATH := ./cmd
BINARY_NAME := nextdo-api-go
DB_URL=postgresql://postgres:postgres@localhost:5432/nextdo?sslmode=disable

# ==================================================================================== #
# HELPERS
# ==================================================================================== #

## help: print this help message
.PHONY: help
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'


# ==================================================================================== #
# DATABASE
# ==================================================================================== #

## postgres: run a postgres container
.PHONY: postgres
postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=postgres -d postgres:16-alpine

## postgres/stop: stop the postgres container
.PHONY: postgres/stop
postgres/stop:
	docker stop postgres

## postgres/rm: remove the postgres container
.PHONY: postgres/rm
postgres/rm:
	docker rm postgres

## db/create: create the database in the postgres container
.PHONY: db/create
db/create:
	docker exec -it postgres createdb --username=postgres --owner=postgres nextdo

## db/drop: drop the database in the postgres container
.PHONY: db/drop
db/drop:
	docker exec -it postgres dropdb nextdo

## migrate/up: run all migrations
.PHONY: migrate/up
migrate/up:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

## migrate/up1: run the next migration
.PHONY: migrate/up1
migrate/up1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

## migrate/down: rollback all migrations
.PHONY: migrate/down
migrate/down:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

## migrate/down1: rollback the last migration
.PHONY: migrate/down1
migrate/down1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

## migrate/create: create a new migration
.PHONY: migrate/create
migrate/create:
	$(eval NAME := $(shell read -p "Migration name: " name; echo $$name))
	migrate create -ext sql -dir db/migration -seq $(NAME)


# ==================================================================================== #
# DEVELOPMENT
# ==================================================================================== #

## test: run all tests
.PHONY: test
test:
	go test -v -race -buildvcs ./...

## test/cover: run all tests and display coverage
.PHONY: test/cover
test/cover:
	go test -v -race -buildvcs -coverprofile=/tmp/coverage.out ./...
	go tool cover -html=/tmp/coverage.out

## build: build the application
.PHONY: build
build:
	go build -o=/tmp/bin/${BINARY_NAME} ${MAIN_PACKAGE_PATH}

## build/prod: build the application for production in Linux, Windows and MacOS with stripped binaries
.PHONY: build/prod
build/prod:
	GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/${BINARY_NAME}-linux-amd64 ${MAIN_PACKAGE_PATH}
	GOOS=windows GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/${BINARY_NAME}-windows-amd64.exe ${MAIN_PACKAGE_PATH}
	GOOS=darwin GOARCH=amd64 go build -ldflags="-s -w" -o=./bin/${BINARY_NAME}-darwin-amd64 ${MAIN_PACKAGE_PATH}

## run: run the  application
.PHONY: run
run: build
	/tmp/bin/${BINARY_NAME}

