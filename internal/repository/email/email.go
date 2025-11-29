package email

import (
	"context"
	"fmt"

	"go.opentelemetry.io/otel/trace"
	"gopkg.in/mail.v2"
)

type email struct {
	tracer             trace.Tracer
	username, password string
}

func NewEmailRepo(tracer trace.Tracer, username, password string) EmailRepo {
	return &email{
		tracer:   tracer,
		username: username,
		password: password,
	}
}

func (e *email) SendEmail(ctx context.Context, to []string, subject, body string) error {
	ctx, span := e.tracer.Start(ctx, "EmailRepo.SendEmail")
	defer span.End()

	msg := mail.NewMessage()
	msg.SetHeaders(map[string][]string{
		"From": {
			e.username,
		},
		"To": to,
		"Subject": {
			subject,
		},
	})

	msg.SetBody("text/html", fmt.Sprintf("</p>%s</p>", body))

	dialer := mail.NewDialer("smtp.gmail.com", 587, e.username, e.password)

	if err := dialer.DialAndSend(msg); err != nil {
		span.RecordError(err)
		return err
	}

	return nil
}
