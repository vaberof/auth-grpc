package auth

type NotificationService interface {
	SendEmail(to string, emailType string, subject string, body map[string]string) error
}
