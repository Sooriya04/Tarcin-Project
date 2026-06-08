# Backend Setup & Installation Guide

This guide covers everything you need to run the Tarcin gRPC backend locally.

## 1. Prerequisites

You must have **Go** installed on your system.

### Install the Protocol Buffer Compiler (protoc)
To compile `.proto` files into Go code, install the `protoc` binary natively:
```bash
sudo apt-get update
sudo apt-get install -y protobuf-compiler
```

### Install gRPC Go Plugins
Run these commands to install the required Go code generators:
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
export PATH=$PATH:$HOME/go/bin
```

### Install grpcurl (For API testing)
`grpcurl` is the gRPC equivalent of `curl`. It allows you to test your endpoints directly from the command line:
```bash
go install github.com/fullstorydev/grpcurl/cmd/grpcurl@latest
```

## 2. Generating Protobuf Code
If you ever modify the `proto/performance.proto` schema, you must regenerate the Go bindings so the server recognizes the new structures.

Run this exactly from the `backend` directory:
```bash
export PATH=$PATH:$HOME/go/bin
protoc --go_out=. --go-grpc_out=. proto/performance.proto
mv backend/internal/grpc/pb/* internal/grpc/pb/
rm -rf backend
```

## 3. Running the Server

Make sure your PostgreSQL database is running, and your `.env` file is properly configured.

1. Navigate to the backend directory:
   ```bash
   cd backend
   ```
2. Download Go dependencies:
   ```bash
   go mod tidy
   ```
3. Run the server:
   ```bash
   go run ./cmd/server/main.go
   ```
   *Expected Output:* `gRPC Server listening on :50051`
