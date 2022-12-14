## help: print this help message
help:
	@echo 'Usage:'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

## run: run the cmd/app application
run:
	go run ./cmd/app

## swagger: output swagger files
swagger:
	swag init -g ./internal/controller/http/v1/router.go

## build: build the cmd/app application
build:
	@echo 'Building cmd/app...'
	go build -ldflags='-s' -o=./bin/app ./cmd/app
