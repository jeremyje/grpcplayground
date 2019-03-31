// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"net"

	pb "github.com/jeremyje/grpcplayground/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	port = ":50051"
)

var (
	privateKeyFlag = flag.String("priv", "secrets/key.pem", "Private Key File")
	publicKeyFlag  = flag.String("pub", "secrets/cert.pem", "Public Key File")
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

// SayHello implements helloworld.GreeterServer
func (s *server) Echo(ctx context.Context, in *pb.EchoRequest) (*pb.EchoResponse, error) {
	log.Printf("Received: %v", in.Text)
	return &pb.EchoResponse{Text: "Hello " + in.Text}, nil
}

func main() {
	flag.Parse()
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	creds, err := credentials.NewServerTLSFromFile(*publicKeyFlag, *privateKeyFlag)
	if err != nil {
		log.Fatalf("failed to create credentials: %v", err)
	}

	s := grpc.NewServer(grpc.Creds(creds))
	s := grpc.NewServer()
	pb.RegisterEchoServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
