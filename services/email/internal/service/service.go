package service

import "context"

type EmailService interface {
	SendVerificationCode(ctx context.Context, email, code string)
}
