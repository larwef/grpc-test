package main

import (
	"context"
	"io"
	"log"
	"net"
	"net/http"
	"os"

	"github.com/google/uuid"
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

func httpHandler(res http.ResponseWriter, req *http.Request) {
	io.WriteString(res, "Ok")
}

func main() {
	port, exists := os.LookupEnv("port")
	if !exists {
		log.Fatal("Need to provide a port via the 'port' enviroment variable")
	}

	healthPort, exists := os.LookupEnv("healthPort")
	if !exists {
		log.Fatal("Need to provide a healthPort via the 'healthPort' enviroment variable")
	}

	http.HandleFunc("/", httpHandler)
	go func() {
		log.Println("Starting health route on :" + healthPort + "/health")
		if err := http.ListenAndServe(":"+healthPort, nil); err != nil {
			log.Fatalf("Error starting health route: %v", err)
		}
	}()

	serverID := uuid.New().String()

	log.Printf("Starting server with id %q...\n", serverID)

	listener, err := net.Listen("tcp", ":"+port)
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
