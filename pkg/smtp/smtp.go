package smtp

import (
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"net/smtp"
)

type SMTP struct {
	Selfemail string
	Selfpwd   string
	Subject   string
	SmtpHost  string
	SmtpPort  string
}
type STATOR interface {
	SendVerificationEmail(to string, code string) error
	GenerateVerificationCode(length int) (string, error)
}

func NewSMTP(email, pwd, subject, host, port string) *SMTP {
	return &SMTP{
		Selfemail: email,
		Selfpwd:   pwd,
		Subject:   subject,
		SmtpHost:  host,
		SmtpPort:  port,
	}
}

// 发送验证码邮件
func (s *SMTP) SendVerificationEmail(to string, code string) error {
	from := s.Selfemail
	password := s.Selfpwd
	smtpHost := s.SmtpHost
	smtpPort := s.SmtpPort
	subject := s.Subject
	body := fmt.Sprintf("Your verification code is: %s", code)
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: " + subject + "\n\n" +
		body
	auth := smtp.PlainAuth("", from, password, smtpHost)
	err := smtp.SendMail(smtpHost+":"+smtpPort, auth, from, []string{to}, []byte(msg))
	return err
}

// 生成随机验证码
func (s *SMTP) GenerateVerificationCode(length int) (string, error) {
	buffer := make([]byte, length)
	_, err := rand.Read(buffer)
	if err != nil {
		return "", err
	}

	return base64.RawURLEncoding.EncodeToString(buffer)[:length], nil
}
