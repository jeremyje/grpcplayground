package main

import (
	"context"
	"flag"
	"log"
	"os"
	"time"

	pb "github.com/jeremyje/grpcplayground/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

const (
	defaultName = "world"
)

var (
	publicKeyFlag = flag.String("pub", "secrets/cert.pem", "Public Key File")
	domainFlag    = flag.String("domain", "github.com", "Domain")
	addressFlag   = flag.String("address", "localhost:50051", "Server Address")
)

func main() {
	// Create tls based credential.
	creds, err := credentials.NewClientTLSFromFile(*publicKeyFlag, "github.com")
	if err != nil {
		log.Fatalf("failed to load credentials: %v", err)
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(*addressFlag, grpc.WithTransportCredentials(creds))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewEchoServiceClient(conn)

	// Contact the server and print out its response.
	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.Echo(ctx, &pb.EchoRequest{Text: name})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.Text)
}
