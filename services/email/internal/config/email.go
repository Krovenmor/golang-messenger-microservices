package config

import "MyMessenger/pkg/config"

type EmailConfig struct {
	Password  string
	FromEmail string

	SmtpHost string
	SmtpPort int
}

func GetEmailConfig() (*EmailConfig, error) {
	r := config.NewConfigReader()

	conf := EmailConfig{
		Password:  r.GetString("PASSWORD"),
		FromEmail: r.GetString("FROM"),

		SmtpHost: r.GetString("SMTP_HOST"),
		SmtpPort: r.GetInt("SMTP_PORT"),
	}

	return &conf, r.Err()
}
