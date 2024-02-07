package auth

import (
	"context"
	authv1 "github.com/vaberof/auth-grpc/protos/gen/go/auth"
	"google.golang.org/grpc"
)

type serverAPI struct {
	authv1.UnimplementedAuthServer
}

func Register(gRPC *grpc.Server) {
	authv1.RegisterAuthServer(gRPC, &serverAPI{})
}

func (s *serverAPI) Register(ctx context.Context, req *authv1.RegisterRequest) (*authv1.RegisterResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Login(ctx context.Context, req *authv1.LoginRequest) (*authv1.LoginResponse, error) {
	panic("implement me")
}

func (s *serverAPI) VerifyAccessToken(ctx context.Context, req *authv1.VerifyAccessTokenRequest) (*authv1.VerifyAccessTokenResponse, error) {
	panic("implement me")
}
