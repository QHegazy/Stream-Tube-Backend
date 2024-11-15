package auth_grpc

import (
	auth "authService/auth_proto_generated"
	"context"
	"fmt"
)

type server struct {
    auth.UnimplementedAuthServiceServer
}

func (s *server) Register(ctx context.Context, req *auth.RegisterRequest) (*auth.RegisterResponse, error) {
    // Your logic here
    fmt.Println(req.Username, req.Password)
    return &auth.RegisterResponse{Message: "Registration successful"}, nil
}

func (s *server) Login(ctx context.Context, req *auth.LoginRequest) (*auth.LoginResponse, error) {
    // Your logic here
    return &auth.LoginResponse{Session: "generated-session-id"}, nil
}

func (s *server) ValidateSession(ctx context.Context, req *auth.ValidateSessionRequest) (*auth.ValidateSessionResponse, error) {
    // Your logic here
    return &auth.ValidateSessionResponse{IsValid: true}, nil
}
