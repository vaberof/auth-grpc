package auth

import (
	"context"
	pb "github.com/vaberof/auth-grpc/genproto/auth_service"
	"github.com/vaberof/auth-grpc/pkg/domain"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type serverAPI struct {
	pb.UnimplementedAuthServiceServer
	authService AuthService
}

func Register(gRPC *grpc.Server, authService AuthService) {
	pb.RegisterAuthServiceServer(gRPC, &serverAPI{authService: authService})
}

func (s *serverAPI) Register(ctx context.Context, req *pb.RegisterRequest) (*emptypb.Empty, error) {
	err := s.authService.Register(context.Background(), domain.Email(req.Email), domain.Password(req.Password))
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	panic("implement me")
}

func (s *serverAPI) Verify(ctx context.Context, req *pb.VerifyRequest) (*emptypb.Empty, error) {
	err := s.authService.Verify(context.Background(), domain.Email(req.Email), domain.Code(req.Code))
	if err != nil {
		return &emptypb.Empty{}, status.Errorf(codes.Internal, "Internal server error: %v", err)
	}
	return &emptypb.Empty{}, nil
}

func (s *serverAPI) VerifyToken(ctx context.Context, req *pb.VerifyTokenRequest) (*pb.VerifyTokenResponse, error) {
	panic("implement me")
}
