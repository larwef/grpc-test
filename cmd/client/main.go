package main

import (
	"context"
	"flag"
	"log"
	"time"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
)

var defaultAddress = "localhost:8080"
var address = flag.String("a", "", "Address to dial")

var defaultMessage = "Hello from client"
var message = flag.String("m", "", "Message to send")

func main() {
	startProgram := time.Now()
	flag.Parse()

	log.Println("Starting client...")

	if *address == "" {
		address = &defaultAddress
	}

	if *message == "" {
		message = &defaultMessage
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := &hello.HelloRequest{
		Message: *message,
	}

	res, err := client.SayHello(ctx, req)
	if err != nil {
		log.Fatalf("Error calling SayHello: %v", err)
	} else {
		log.Printf("Response from server %q: %s\n", res.ServerID, res.Response)
	}

	now := time.Now()
	log.Printf("Call took: %v\n", now.Sub(startProgram))
	log.Printf("Program took: %v\n", now.Sub(startProgram))

	log.Println("Client exited.")
}
