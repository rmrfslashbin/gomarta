.DEFAULT_GOAL := default

protoc:
	# https://developers.google.com/transit/gtfs-realtime/
	# https://developers.google.com/static/transit/gtfs-realtime/gtfs-realtime.proto
	# Add to proto file:
	#   option go_package = "github.com/rmrfslashbin/gomarta/pkg/gtfsrt";
	@echo "Building protobufs"
	@curl -s -o pkg/gtfsrt/gtfs-realtime.proto  https://developers.google.com/static/transit/gtfs-realtime/gtfs-realtime.proto 
	@protoc --go_opt=Mpkg/gtfsrt/gtfs-realtime.proto=github.com/rmrfslashbin/gomarta/pkg/gtfsrt --go_out=. --go_opt=paths=source_relative pkg/gtfsrt/gtfs-realtime.proto


build: protoc
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
