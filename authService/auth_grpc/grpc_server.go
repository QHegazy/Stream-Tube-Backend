package auth_grpc

import (
	"authService/auth_proto_generated"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
)


func AuthGrpcServer() {
    lis, err := net.Listen("tcp", ":5050")
    if err != nil {
        log.Fatalf("failed to listen: %v", err)
    }

    s := grpc.NewServer()
    auth_proto_generated.RegisterAuthServiceServer(s, &server{})

    fmt.Println("Server is running on port :50051...")
    if err := s.Serve(lis); err != nil {
        log.Fatalf("failed to serve: %v", err)
    }
}
