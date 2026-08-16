package service

import (
	"MyMessenger/services/email/internal/config"
	"context"
	"log"

	"gopkg.in/gomail.v2"
)

type TmpEmailSender struct {
	conf *config.EmailConfig
	dial *gomail.Dialer
}

func NewEmailSender(conf *config.EmailConfig) *TmpEmailSender {
	d := gomail.NewDialer(conf.SmtpHost, conf.SmtpPort, conf.FromEmail, conf.Password)
	return &TmpEmailSender{dial: d, conf: conf}
}

func (s *TmpEmailSender) SendVerificationCode(ctx context.Context, email, code string) {
	log.Printf("Email: %q, Code: %s", email, code)

	m := gomail.NewMessage()
	m.SetHeader("From", s.conf.FromEmail)
	m.SetHeader("To", email)
	m.SetHeader("Subject", "Yout verification code: "+code)
	m.SetBody("text/plain", "For Signalline registration")

	err := s.dial.DialAndSend(m)
	if err != nil {
		log.Printf("SendVerificationCode: Trouble with DialAndSend, err: %q", err)
		return
	}

	log.Printf("Email sended")
}
