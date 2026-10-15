// Package mail sends plain-text email over SMTP.
package mail

import (
	"fmt"
	"log"
	"net/smtp"
)

type Config struct {
	Host, Port, User, Pass, From string
}

type Mailer struct{ cfg Config }

func New(cfg Config) *Mailer { return &Mailer{cfg: cfg} }

func (m *Mailer) Send(to, subject, body string) error {
	msg := fmt.Sprintf(
		"From: %s\r\nTo: %s\r\nSubject: %s\r\nContent-Type: text/plain; charset=utf-8\r\n\r\n%s",
		m.cfg.From, to, subject, body)

	var auth smtp.Auth
	if m.cfg.User != "" {
		auth = smtp.PlainAuth("", m.cfg.User, m.cfg.Pass, m.cfg.Host)
	}
	err := smtp.SendMail(m.cfg.Host+":"+m.cfg.Port, auth, m.cfg.From, []string{to}, []byte(msg))
	if err != nil {
		log.Println("mail error:", err)
	}
	return err
}
