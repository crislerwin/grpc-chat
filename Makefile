.PHONY: proto build run-server run-client clean help

# Generate protobuf files
proto:
	protoc --go_out=. --go_opt=paths=source_relative \
	       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
	       proto/chat/chat.proto

# Build server and client
build: proto
	go build -o server cmd/server/main.go
	go build -o client cmd/client/main.go

# Run server
run-server: build
	./server

# Run client
run-client: build
	./client

# Clean generated files and binaries
clean:
	rm -f server client
	rm -f proto/chat/*.pb.go

# Install dependencies
deps:
	go mod download
	go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
	go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest

# Help
help:
	@echo "Available targets:"
	@echo "  make proto       - Generate protobuf files"
	@echo "  make build       - Build server and client"
	@echo "  make run-server  - Build and run server"
	@echo "  make run-client  - Build and run client"
	@echo "  make clean       - Remove generated files"
	@echo "  make deps        - Install dependencies"
