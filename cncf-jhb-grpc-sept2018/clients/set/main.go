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

	key := flag.String("k", "", "key")
	value := flag.String("v", "", "value")

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

	_, err = client.Set(ctx, &kv.SetRequest{
		Key:   *key,
		Value: *value,
	})

	if err != nil {
		log.Fatal("kv.Set err:", err)
	}

	fmt.Printf("KV Set successful! key=%s value=%s\n", *key, *value)

}
