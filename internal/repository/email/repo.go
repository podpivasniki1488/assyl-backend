package email

import "context"

type EmailRepo interface {
	SendEmail(ctx context.Context, to []string, subject, body string) error
}
