package grpcclient

import (
	"context"
	"fmt"
	pb "github.com/vaberof/auth-grpc/genproto/notification_service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GrpcClient interface {
	NotificationService() pb.NotificationServiceClient
}

type grpcClientImpl struct {
	cfg         *NotificationServiceClientConfig
	connections map[string]interface{}
}

func New(cfg *NotificationServiceClientConfig) (GrpcClient, error) {
	connNotificationService, err := grpc.DialContext(
		context.Background(),
		fmt.Sprintf("%s:%d", cfg.Host, cfg.Port),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("user service dial host=%s port=%d err=%v",
			cfg.Host, cfg.Port, err)
	}

	return &grpcClientImpl{
		cfg: cfg,
		connections: map[string]interface{}{
			"notification_service": pb.NewNotificationServiceClient(connNotificationService),
		},
	}, nil
}

func (g *grpcClientImpl) NotificationService() pb.NotificationServiceClient {
	return g.connections["notification_service"].(pb.NotificationServiceClient)
}
