protoc:
	@echo "Building protobufs"
	@protoc --go_out=. --go_opt=paths=source_relative pkg/gtfsrt/gtfs-realtime.proto

update:
	@echo "Updating data"
	@go get -u ./...