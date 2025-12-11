# gRPC Chat

A simple real-time chat application using bidirectional gRPC streaming in Go.

## Quick Start

```bash
# Start the server
make run-server

# Start clients (in separate terminals)
make run-client
```

## Project Structure

```
.
├── proto/chat/          # Protocol buffers definitions
├── cmd/
│   ├── server/          # Chat server
│   └── client/          # Chat client
└── Makefile
```

## How it works

The server maintains connections to multiple clients using gRPC bidirectional streaming. When a client sends a message, it's broadcast to all connected clients in real-time.

## Building

```bash
make build    # Build server and client
make proto    # Regenerate proto files
make clean    # Clean up
```

## Sharing gRPC contracts between microservices

Don't copy proto files around. Instead:

1. **Create a separate repo for your protos** (recommended)
   ```
   proto-contracts/
   ├── chat/v1/chat.proto
   ├── user/v1/user.proto
   └── go.mod
   ```

2. **Import it as a module in your services**
   ```go
   // go.mod
   require github.com/yourorg/proto-contracts v1.0.0
   
   // main.go
   import chatv1 "github.com/yourorg/proto-contracts/chat/v1"
   ```

3. **Version with git tags**
   ```bash
   git tag v1.0.0
   go get github.com/yourorg/proto-contracts@v1.0.0
   ```

This keeps a single source of truth and makes updates easy across all your services.
