.DEFAULT_GOAL := default

protoc:
	@echo "Building protobufs"
	@protoc --go_out=. --go_opt=paths=source_relative pkg/gtfsrt/gtfs-realtime.proto


build:
	@if [ ! -d "./bin" ]; then mkdir bin; fi
	@go build -o bin .

install:
	@go install

tidy:
	@echo "Making mod tidy"
	@go mod tidy

update:
	@echo "Updating..."
	@go get -u ./...
	@go mod tidy

default: tidy build install
