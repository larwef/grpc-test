package test

import (
	"context"
	"log"
	"testing"
	"time"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
)

var address = "<yourNLBaddress>:<yourPort>"

func getClient() (hello.HelloServiceClient, func() error) {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}

	client := hello.NewHelloServiceClient(conn)

	return client, conn.Close
}

func doCall(client hello.HelloServiceClient, message string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &hello.HelloRequest{
		Message: message,
	}

	res, err := client.SayHello(ctx, req)
	return res.ServerID, err
}
func Test_SingleCall(t *testing.T) {
	client, close := getClient()
	defer close()

	serverID, err := doCall(client, "TestMessage")
	if err != nil {
		t.Errorf("Error: %v", err)
	}

	t.Logf("Successfull call to server %s", serverID)
}
