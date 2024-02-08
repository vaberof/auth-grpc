package auth

import "context"

type NotificationService interface {
	SendEmail(ctx context.Context, to string, emailType string, subject string, body map[string]string) error
}
