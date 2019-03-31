// Package main implements a server for Greeter service.
package main

import (
	"context"
	"flag"
	"log"
	"net"

	"strings"

	pb "github.com/jeremyje/grpcplayground/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"google.golang.org/grpc/status"
)

const (
	port = ":50051"
)

var (
	privateKeyFlag = flag.String("priv", "secrets/key.pem", "Private Key File")
	publicKeyFlag  = flag.String("pub", "secrets/cert.pem", "Public Key File")
	tokenFlag      = flag.String("token", "secret", "Token")
)

// server is used to implement helloworld.GreeterServer.
type server struct{}

func AuthUnaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		err := authenticateRequest(ctx)
		if err != nil {
			return nil, err
		}
		return req, nil
	}
}

func purgeHeader(ctx context.Context, header string) context.Context {
	md, _ := metadata.FromIncomingContext(ctx)
	mdCopy := md.Copy()
	mdCopy[header] = nil
	return metadata.NewIncomingContext(ctx, mdCopy)
}

func tryTokenAuth(ctx context.Context) (context.Context, error) {
	auth, err := extractHeader(ctx, "authorization")
	if err != nil {
		return ctx, err
	}

	const prefix = "Bearer "
	if !strings.HasPrefix(auth, prefix) {
		return ctx, status.Error(codes.Unauthenticated, `missing "Bearer " prefix in "Authorization" header`)
	}

	if strings.TrimPrefix(auth, prefix) != *tokenFlag {
		return ctx, status.Error(codes.Unauthenticated, "invalid token")
	}

	// Remove token from headers from here on
	return purgeHeader(ctx, "authorization"), nil
}

func extractHeader(ctx context.Context, header string) (string, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no headers in request")
	}

	authHeaders, ok := md[header]
	if !ok {
		return "", status.Error(codes.Unauthenticated, "no header in request")
	}

	if len(authHeaders) != 1 {
		return "", status.Error(codes.Unauthenticated, "more than 1 header in request")
	}

	return authHeaders[0], nil
}

func authenticateRequest(ctx context.Context) error {
	p, ok := peer.FromContext(ctx)
	if !ok {
		return status.Error(codes.Unauthenticated, "no peer found")
	}

	tlsAuth, ok := p.AuthInfo.(credentials.TLSInfo)
	if !ok {
		return status.Error(codes.Unauthenticated, "unexpected peer transport credentials")
	}

	log.Printf("peer: %v", p)
	log.Printf("tlsAuth: %v", tlsAuth)
	log.Printf("tlsAuth.State.VerifiedChains: %v", tlsAuth.State.VerifiedChains)
	log.Printf("tlsAuth.State.PeerCertificates: %v", tlsAuth.State.PeerCertificates)

	authedCtx, err := tryTokenAuth(ctx)
	if err != nil {
		return err
	}
	log.Printf("%v\n", authedCtx)

	/*
		if len(tlsAuth.State.VerifiedChains) == 0 || len(tlsAuth.State.VerifiedChains[0]) == 0 {
			return status.Error(codes.Unauthenticated, "could not verify peer certificate")
		}

		// Check subject common name against configured username
		if tlsAuth.State.VerifiedChains[0][0].Subject.CommonName != "username" {
			return status.Error(codes.Unauthenticated, "invalid subject common name")
		}
	*/
	return nil
}

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

	s := grpc.NewServer(grpc.Creds(creds), grpc.UnaryInterceptor(AuthUnaryServerInterceptor()))
	pb.RegisterEchoServiceServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
