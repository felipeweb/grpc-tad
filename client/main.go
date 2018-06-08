package main

import (
	"context"
	"log"

	pb "github.com/felipeweb/grpc-tad/proto/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
)

func main() {
	creds, err := credentials.NewClientTLSFromFile("server/server-cert.pem", "")
	if err != nil {
		log.Fatalf("cert load error: %s", err)
	}

	conn, err := grpc.Dial("localhost:6000", grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("Failed to start gRPC connection: %v", err)
	}
	defer conn.Close() // nolint

	client := pb.NewSimpleServerClient(conn)

	md := metadata.Pairs("authorization", "tad")
	ctx := metadata.NewOutgoingContext(context.Background(), md)

	_, err = client.CreateUser(ctx, &pb.CreateUserRequest{User: &pb.User{Username: "pery", Role: "joker"}})
	if err != nil {
		log.Fatalf("Failed to create user: %v", err)
	}
	log.Println("Created user!")

	resp, err := client.GetUser(ctx, &pb.GetUserRequest{Username: "pery"})
	if err != nil {
		log.Fatalf("Failed to get created user: %v", err)
	}
	log.Printf("User exists: %v\n", resp)
}
