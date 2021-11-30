directory = $(PWD)
build:
	mkdir -p build
	go mod download
	cd cmd && go build -o ${directory}/build/sku_reader

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
	cd build && ./sku_reader