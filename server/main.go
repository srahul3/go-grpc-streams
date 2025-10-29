package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net"

	pb "github.com/srahul3/go-grpc-streams/proto"
	"google.golang.org/grpc"
)

// FooServer implements the Foo service
type FooServer struct {
	pb.UnimplementedFooServer
}

// SayHello implements the SayHello RPC
func (s *FooServer) SayHello(ctx context.Context, req *pb.HelloRequest) (*pb.HelloResponse, error) {
	log.Printf("Foo.SayHello called with name: %s", req.Name)
	return &pb.HelloResponse{
		Message: fmt.Sprintf("Hello, %s!", req.Name),
	}, nil
}

// StreamNumbers implements the StreamNumbers RPC (server streaming)
func (s *FooServer) StreamNumbers(req *pb.NumberRequest, stream pb.Foo_StreamNumbersServer) error {
	log.Printf("Foo.StreamNumbers called with count: %d", req.Count)
	for i := int32(1); i <= req.Count; i++ {
		if err := stream.Send(&pb.NumberResponse{Number: i}); err != nil {
			return err
		}
		log.Printf("Sent number: %d", i)
	}
	return nil
}

// BarServer implements the Bar service
type BarServer struct {
	pb.UnimplementedBarServer
}

// GetInfo implements the GetInfo RPC
func (s *BarServer) GetInfo(ctx context.Context, req *pb.InfoRequest) (*pb.InfoResponse, error) {
	log.Printf("Bar.GetInfo called with query: %s", req.Query)
	return &pb.InfoResponse{
		Info: fmt.Sprintf("Information about: %s", req.Query),
	}, nil
}

// CollectMessages implements the CollectMessages RPC (client streaming)
func (s *BarServer) CollectMessages(stream pb.Bar_CollectMessagesServer) error {
	log.Println("Bar.CollectMessages called")
	var messages []string
	for {
		msg, err := stream.Recv()
		if err == io.EOF {
			// Client has finished sending messages
			result := fmt.Sprintf("Collected %d messages", len(messages))
			log.Printf("Finished collecting messages: %s", result)
			return stream.SendAndClose(&pb.MessageResponse{
				Result: result,
				Count:  int32(len(messages)),
			})
		}
		if err != nil {
			return err
		}
		log.Printf("Received message: %s", msg.Content)
		messages = append(messages, msg.Content)
	}
}

func main() {
	// Create a TCP listener on port 50051
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen: %v", err)
	}

	// Create a gRPC server
	grpcServer := grpc.NewServer()

	// Register the Foo and Bar services
	pb.RegisterFooServer(grpcServer, &FooServer{})
	pb.RegisterBarServer(grpcServer, &BarServer{})

	log.Println("gRPC server is running on port 50051...")
	log.Println("Services available: Foo, Bar")

	// Start serving
	if err := grpcServer.Serve(lis); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
