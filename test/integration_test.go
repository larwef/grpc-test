package test

import (
	"context"
	"log"
	"strconv"
	"testing"
	"time"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
)

var address = "<yourNLBaddress>:<yourPort>"

var iterations = 100

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

func Test_MultipleCallsOneConnection(t *testing.T) {
	client, close := getClient()
	defer close()

	servers := make(map[string]int)
	for i := 0; i < iterations; i++ {
		serverID, err := doCall(client, "TestMessage "+strconv.Itoa(i))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		servers[serverID]++
	}

	t.Log("Hit the following servers:")
	for k, v := range servers {
		t.Logf("%s: %d", k, v)
	}
}

func Test_MultipleCallsMultipleConnections(t *testing.T) {
	servers := make(map[string]int)
	for i := 0; i < iterations; i++ {
		client, close := getClient()

		serverID, err := doCall(client, "TestMessage "+strconv.Itoa(i))
		if err != nil {
			t.Errorf("Error: %v", err)
		}
		close()
		servers[serverID]++
	}

	t.Log("Hit the following servers:")
	for k, v := range servers {
		t.Logf("%s: %d", k, v)
	}
}
