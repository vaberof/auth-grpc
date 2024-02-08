package notificationservice

import (
	"context"
	"fmt"
	"github.com/vaberof/auth-grpc/genproto/notification_service"
	"github.com/vaberof/auth-grpc/pkg/grpc/grpcclient"
	"github.com/vaberof/auth-grpc/pkg/logging/logs"
	"log/slog"
)

type NotificationService struct {
	grpcClient grpcclient.GrpcClient

	logger *slog.Logger
}

func New(grpcClient grpcclient.GrpcClient, logs *logs.Logs) *NotificationService {
	logger := logs.WithName("infra.integration.grpc.notificationservice")
	return &NotificationService{grpcClient: grpcClient, logger: logger}
}

func (service *NotificationService) SendEmail(ctx context.Context, to string, emailType string, subject string, body map[string]string) error {
	const operation = "SendEmail"

	log := service.logger.With(
		slog.String("operation", operation),
		slog.String("to_email", to),
		slog.String("email_type", emailType),
		slog.String("subject", subject),
		slog.Any("body", body))

	log.Info("sending email to", to)

	_, err := service.grpcClient.NotificationService().SendEmail(ctx, &notification_service.SendEmailRequest{
		To:      to,
		Subject: subject,
		Body:    body,
		Type:    emailType,
	})
	if err != nil {
		log.Error("failed to send email", err)

		return fmt.Errorf("%s: %w", operation, err)
	}

	log.Info("email has been sent successfully")

	return nil
}
