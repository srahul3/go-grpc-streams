# go-grpc-streams

A Go gRPC server and client implementation with two services: `Foo` and `Bar`. This project demonstrates various gRPC patterns including unary RPCs, server streaming, and client streaming.

## Services

### Foo Service
- **SayHello**: Unary RPC that returns a greeting message
- **StreamNumbers**: Server streaming RPC that streams a sequence of numbers

### Bar Service
- **GetInfo**: Unary RPC that returns information about a query
- **CollectMessages**: Client streaming RPC that collects messages from the client

## Project Structure

```
.
├── proto/
│   ├── services.proto          # Protocol buffer definitions
│   ├── services.pb.go          # Generated Go code (messages)
│   └── services_grpc.pb.go     # Generated Go code (services)
├── server/
│   └── main.go                 # gRPC server implementation
├── client/
│   └── main.go                 # gRPC client implementation
├── go.mod                      # Go module definition
├── go.sum                      # Go dependencies
└── README.md                   # This file
```

## Prerequisites

- Go 1.24 or higher
- Protocol Buffer Compiler (protoc)
- protoc-gen-go plugin
- protoc-gen-go-grpc plugin

## Installation

### 1. Install Go
Download and install Go from [https://golang.org/dl/](https://golang.org/dl/)

### 2. Install Protocol Buffer Compiler

#### On Linux:
```bash
cd /tmp
curl -LO https://github.com/protocolbuffers/protobuf/releases/download/v25.1/protoc-25.1-linux-x86_64.zip
unzip protoc-25.1-linux-x86_64.zip -d $HOME/.local
export PATH="$PATH:$HOME/.local/bin"
```

#### On macOS:
```bash
brew install protobuf
```

#### On Windows:
Download the protoc binary from [https://github.com/protocolbuffers/protobuf/releases](https://github.com/protocolbuffers/protobuf/releases) and add it to your PATH.

### 3. Install Go Plugins for Protocol Buffer Compiler
```bash
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
```

Make sure `$HOME/go/bin` is in your PATH:
```bash
export PATH="$PATH:$HOME/go/bin"
```

## Proto Compilation

To generate Go code from the `.proto` file:

```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/services.proto
```

This will generate two files in the `proto/` directory:
- `services.pb.go`: Contains the message type definitions
- `services_grpc.pb.go`: Contains the service client and server code

## Build

### Build Server
```bash
go build -o server/server ./server
```

### Build Client
```bash
go build -o client/client ./client
```

### Build Both
```bash
go build -o server/server ./server && go build -o client/client ./client
```

## Running the Application

### 1. Start the Server
In one terminal, run:
```bash
./server/server
```

Or using Go:
```bash
go run ./server/main.go
```

The server will start on port `50051` and log:
```
gRPC server is running on port 50051...
Services available: Foo, Bar
```

### 2. Run the Client
In another terminal, run:
```bash
./client/client
```

Or using Go:
```bash
go run ./client/main.go
```

The client will connect to the server and test all four RPC methods:
1. Foo.SayHello (Unary RPC)
2. Foo.StreamNumbers (Server Streaming RPC)
3. Bar.GetInfo (Unary RPC)
4. Bar.CollectMessages (Client Streaming RPC)

## Example Output

### Server Output:
```
gRPC server is running on port 50051...
Services available: Foo, Bar
Foo.SayHello called with name: Alice
Foo.StreamNumbers called with count: 5
Sent number: 1
Sent number: 2
Sent number: 3
Sent number: 4
Sent number: 5
Bar.GetInfo called with query: gRPC
Bar.CollectMessages called
Received message: Hello
Received message: from
Received message: the
Received message: client
Received message: streaming
Received message: RPC
Finished collecting messages: Collected 6 messages
```

### Client Output:
```
Connected to gRPC server at localhost:50051
--------------------------------------------------

1. Testing Foo.SayHello (Unary RPC)...
Response: Hello, Alice!

2. Testing Foo.StreamNumbers (Server Streaming RPC)...
Receiving numbers from server:
  Received: 1
  Received: 2
  Received: 3
  Received: 4
  Received: 5

3. Testing Bar.GetInfo (Unary RPC)...
Response: Information about: gRPC

4. Testing Bar.CollectMessages (Client Streaming RPC)...
Sending messages to server:
  Sending: Hello
  Sending: from
  Sending: the
  Sending: client
  Sending: streaming
  Sending: RPC
Server response: Collected 6 messages (count: 6)

--------------------------------------------------
All tests completed successfully!
```

## Development

### Install Dependencies
```bash
go mod download
```

### Tidy Dependencies
```bash
go mod tidy
```

### Regenerate Proto Files
If you modify `proto/services.proto`, regenerate the Go code:
```bash
protoc --go_out=. --go_opt=paths=source_relative \
       --go-grpc_out=. --go-grpc_opt=paths=source_relative \
       proto/services.proto
```

## License

This project is open source and available under the MIT License.