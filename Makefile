directory = $(PWD)
build:
	mkdir -p build
	go mod download
	go mod tidy
	go build -o ${directory}/build/sku-api sku-reader/cmd/api

setup-env:
	cd etc/dev/docker && docker-compose up -d

teardown-env:
	cd etc/dev/docker && docker-compose down

unit-test:
	go test -tags=unit ./...

integration-test: setup-env
	go test -tags=integration ./...
	cd etc/dev/docker && docker-compose down

run: build setup-env
	cd build && ./sku-api