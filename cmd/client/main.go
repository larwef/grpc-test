package main

import (
	"context"
	"crypto/x509"
	"flag"
	"log"
	"time"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var defaultAddress = "localhost:8080"
var address = flag.String("a", "", "Address to dial.")

var defaultMessage = "Hello from client"
var message = flag.String("m", "", "Message to send.")

var insecure = flag.Bool("i", false, "Use with insecure if set.")

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
	var opts []grpc.DialOption
	if *insecure {
		log.Println("Using WithInsecure.")
		opts = append(opts, grpc.WithInsecure())
	} else {
		pool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("unable to get cert pool: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, "")))
	}

	log.Printf("Dialing: %s", *address)
	conn, err := grpc.Dial(*address, opts...)
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()

	client := hello.NewHelloServiceClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &hello.HelloRequest{
		Message: *message,
	}

	log.Printf("Sending message: %s", *message)
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
