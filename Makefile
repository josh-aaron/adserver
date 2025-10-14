build:
	go build -o bin/adserver cmd/api/*.go

run: build
	./bin/adserver

test:
	go test -v ./...