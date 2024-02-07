package auth

import (
	"context"
	pb "github.com/vaberof/auth-grpc/genproto/auth_service"
	"google.golang.org/grpc"
)

type serverAPI struct {
	pb.UnimplementedAuthServiceServer
}

func Register(gRPC *grpc.Server) {
	pb.RegisterAuthServiceServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	panic("implement me")
}
