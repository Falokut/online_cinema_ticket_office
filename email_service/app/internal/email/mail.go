package email

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"

	"github.com/k3a/html2text"
	"github.com/sirupsen/logrus"
	"gopkg.in/gomail.v2"
)

type MailSender struct {
	logger       *logrus.Logger
	dialler      *gomail.Dialer
	temp         *template.Template
	EmailAddress string
}

type MailSenderConfig struct {
	Password        string `env:"EMAIL_PASSWORD"`
	Port            int    `yaml:"email_port"`
	Host            string `yaml:"email_host"`
	EmailAddress    string `yaml:"email_address"`
	EmailLogin      string `yaml:"email_login"`
	EnableTLS       bool   `yaml:"enable_TLS"`
	TemplatesOrigin string `yaml:"templates_origin"`
}

func NewMailSender(cfg MailSenderConfig, logger *logrus.Logger) *MailSender {
	temp := template.Must(template.ParseGlob(fmt.Sprintf("%s/*.html", cfg.TemplatesOrigin)))

	s := MailSender{logger: logger, EmailAddress: cfg.EmailAddress, temp: temp}

	s.logger.Infoln("Creating mail dialler.")
	s.dialler = gomail.NewDialer(cfg.Host, cfg.Port, cfg.EmailLogin, cfg.Password)
	s.dialler.TLSConfig = &tls.Config{InsecureSkipVerify: true}

	return &s
}

type EmailData struct {
	URL     string
	Subject string
	LinkTTL string
}

func (s *MailSender) SendEmail(data *EmailData, EmailTo string, templateName string) error {
	var body bytes.Buffer

	if err := s.temp.ExecuteTemplate(&body, templateName, &data); err != nil {
		return err
	}

	sender, err := s.dialler.Dial()
	if err != nil {
		return err
	}

	s.logger.Infoln("Creating message.")
	m := gomail.NewMessage()
	m.SetHeader("From", s.EmailAddress)
	m.SetHeader("Subject", data.Subject)
	m.AddAlternative("text/plain", html2text.HTML2Text(body.String()))
	m.SetBody("text/html", body.String())
	defer sender.Close()

	s.logger.Infoln("Sending message.")
	if err := sender.Send(s.EmailAddress, []string{EmailTo}, m); err != nil {
		s.logger.Debug(EmailTo)
		s.logger.Error(err.Error())
		return nil
	}

	s.logger.Infoln("Message sended.")
	return nil
}
