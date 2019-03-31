package test

import (
	"context"
	"crypto/x509"
	"log"
	"strconv"
	"testing"
	"time"

	hello "github.com/larwef/grpc-test/internal/hello"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

var address = "<yourNLBaddress>:<yourPort>"

var iterations = 10
var insecure = false

func getClient() (hello.HelloServiceClient, func() error) {
	var opts []grpc.DialOption
	if insecure {
		opts = append(opts, grpc.WithInsecure())
	} else {
		pool, err := x509.SystemCertPool()
		if err != nil {
			log.Fatalf("unable to get cert pool: %v", err)
		}
		opts = append(opts, grpc.WithTransportCredentials(credentials.NewClientTLSFromCert(pool, "")))
	}

	conn, err := grpc.Dial(address, opts...)
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
	if err != nil {
		return "", err
	}

	return res.ServerID, nil
}

func Test_SingleCall(t *testing.T) {
	client, close := getClient()
	defer close()

	serverID, err := doCall(client, "TestMessage")
	if err != nil {
		t.Fatalf("Error: %v", err)
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
			t.Fatalf("Error: %v", err)
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
			t.Fatalf("Error: %v", err)
		}
		close()
		servers[serverID]++
	}

	t.Log("Hit the following servers:")
	for k, v := range servers {
		t.Logf("%s: %d", k, v)
	}
}
