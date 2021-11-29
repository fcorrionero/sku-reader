setup-env:
	cd etc/devtools/docker && docker-compose up -d

teardown-env:
	cd etc/devtools/docker && docker-compose down

unit-test:
	go test -tags=unit ./...

integration-test: setup-env
	go test -tags=integration ./...
	cd etc/devtools/docker && docker-compose down

run: setup-env
	go run cmd/main.go