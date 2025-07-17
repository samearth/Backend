package mailer

import (
	"fmt"
	"github.com/spf13/viper"
	"net/smtp"
)

type Mailer struct {
	From     string
	Password string
	Host     string
	Port     string
}

// NewMailer initializes from env/config
func NewMailer() *Mailer {
	return &Mailer{
		From:     viper.GetString("EMAIL_SENDER"),
		Password: viper.GetString("EMAIL_PASSWORD"),
		Host:     viper.GetString("EMAIL_SMTP_HOST"),
		Port:     viper.GetString("EMAIL_SMTP_PORT"),
	}
}

// Send sends a plain-text email
func (m *Mailer) Send(to string, subject string, body string) error {
	auth := smtp.PlainAuth("", m.From, m.Password, m.Host)

	msg := []byte(fmt.Sprintf("To: %s\r\nSubject: %s\r\n\r\n%s", to, subject, body))

	addr := fmt.Sprintf("%s:%s", m.Host, m.Port)
	return smtp.SendMail(addr, auth, m.From, []string{to}, msg)
}
