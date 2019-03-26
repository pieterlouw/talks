package main

import (
	"flag"
	"fmt"
	"log"
	"play/cncf/gen"

	"golang.org/x/net/context"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	key := flag.String("k", "", "key ")

	flag.Parse()
	ctx := context.Background()

	// Create the client TLS credentials
	creds, err := credentials.NewClientTLSFromFile("../../certs/demo.crt", "CNCF-gRPC")
	if err != nil {
		log.Fatalf("could not load tls cert: %s", err)
	}

	grpcConn, err := grpc.Dial("localhost:8080",
		//grpc.WithInsecure()
		grpc.WithTransportCredentials(creds),
	)
	if err != nil {
		log.Fatal("grpc.Dial err:", err)
	}

	client := kv.NewKVClient(grpcConn)

	rsp, err := client.Get(ctx, &kv.GetRequest{
		Key: *key,
	})

	if err != nil {
		log.Fatal("client.Get err:", err)
	}

	fmt.Printf("Value set for %s: %s\n", *key, rsp.Value)
}

/*// create gRPC TLS credentials
creds := credentials.NewTLS(&tls.Config{
	InsecureSkipVerify: true, // using self signed certificate for demo, for more secure connections see https://bbengfort.github.io/programmer/2017/03/03/secure-grpc.html
})*/
