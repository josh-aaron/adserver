build:
	go build -o bin/adserver cmd/api/*.go

run: build
	./bin/adserver

test:
	go test -v ./...

buildWindows:
	cd cmd/api; \
	go build -o ../../bin/adserver.exe

runWindows: buildWindows
	./bin/adserver.exe