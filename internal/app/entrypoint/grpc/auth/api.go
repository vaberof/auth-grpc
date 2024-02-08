package auth

import (
	"context"
	pb "github.com/vaberof/auth-grpc/genproto/auth_service"
	"google.golang.org/grpc"
)

type serverAPI struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
}

func Register(gRPC *grpc.Server, authService AuthService) {
	pb.RegisterAuthServiceServer(gRPC, &serverAPI{authService: authService})
}

func (s *serverAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Verify(ctx context.Context, req *pb.VerifyRequest) (*pb.VerifyResponse, error) {
	panic("implement me")
}

func (s *serverAPI) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	panic("implement me")
}
