package main

import (
	"context"
	"io"
	"log"
	"time"

	pb "github.com/srahul3/go-grpc-streams/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

func main() {
	// Connect to the gRPC server
	conn, err := grpc.NewClient("localhost:50051", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create clients for both services
	fooClient := pb.NewFooClient(conn)
	barClient := pb.NewBarClient(conn)

	log.Println("Connected to gRPC server at localhost:50051")
	log.Println("--------------------------------------------------")

	// Test Foo service - SayHello (Unary RPC)
	log.Println("\n1. Testing Foo.SayHello (Unary RPC)...")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	
	helloResp, err := fooClient.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})
	if err != nil {
		log.Fatalf("Foo.SayHello failed: %v", err)
	}
	log.Printf("Response: %s", helloResp.Message)

	// Test Foo service - StreamNumbers (Server streaming RPC)
	log.Println("\n2. Testing Foo.StreamNumbers (Server Streaming RPC)...")
	ctx2, cancel2 := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel2()
	
	stream, err := fooClient.StreamNumbers(ctx2, &pb.NumberRequest{Count: 5})
	if err != nil {
		log.Fatalf("Foo.StreamNumbers failed: %v", err)
	}
	
	log.Println("Receiving numbers from server:")
	for {
		numResp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatalf("Failed to receive number: %v", err)
		}
		log.Printf("  Received: %d", numResp.Number)
	}

	// Test Bar service - GetInfo (Unary RPC)
	log.Println("\n3. Testing Bar.GetInfo (Unary RPC)...")
	ctx3, cancel3 := context.WithTimeout(context.Background(), time.Second)
	defer cancel3()
	
	infoResp, err := barClient.GetInfo(ctx3, &pb.InfoRequest{Query: "gRPC"})
	if err != nil {
		log.Fatalf("Bar.GetInfo failed: %v", err)
	}
	log.Printf("Response: %s", infoResp.Info)

	// Test Bar service - CollectMessages (Client streaming RPC)
	log.Println("\n4. Testing Bar.CollectMessages (Client Streaming RPC)...")
	ctx4, cancel4 := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel4()
	
	msgStream, err := barClient.CollectMessages(ctx4)
	if err != nil {
		log.Fatalf("Bar.CollectMessages failed: %v", err)
	}

	messages := []string{"Hello", "from", "the", "client", "streaming", "RPC"}
	log.Println("Sending messages to server:")
	for _, msg := range messages {
		log.Printf("  Sending: %s", msg)
		if err := msgStream.Send(&pb.MessageRequest{Content: msg}); err != nil {
			log.Fatalf("Failed to send message: %v", err)
		}
	}

	msgResp, err := msgStream.CloseAndRecv()
	if err != nil {
		log.Fatalf("Failed to receive response: %v", err)
	}
	log.Printf("Server response: %s (count: %d)", msgResp.Result, msgResp.Count)

	log.Println("\n--------------------------------------------------")
	log.Println("All tests completed successfully!")
}
