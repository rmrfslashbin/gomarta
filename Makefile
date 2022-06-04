.DEFAULT_GOAL := default

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
