package main

import (
	"context"
	"log"
	"net"
	"os"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

// HelloServer implements the hello service
type HelloServer struct {
	serverID string
}

// SayHello says hello
func (hs *HelloServer) SayHello(ctx context.Context, req *hello.HelloRequest) (*hello.HelloResponse, error) {
	log.Printf("SayHello invoked with message: %q\n", req.Message)

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
	}

	return hs.sayHello(req), nil
}

func (hs *HelloServer) sayHello(req *hello.HelloRequest) *hello.HelloResponse {
	return &hello.HelloResponse{
		ServerID: hs.serverID,
		Response: "Got it!",
	}
}

func main() {
	var serverID string

	if len(os.Args) > 1 {
		serverID = os.Args[1]
	} else {
		log.Fatal("Need to provide a server id")
	}

	log.Printf("Starting server with id %q...\n", serverID)

	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	server := grpc.NewServer()
	hello.RegisterHelloServiceServer(server, &HelloServer{serverID: serverID})

	reflection.Register(server)
	if err := server.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}

	log.Println("Server exited.")
}
