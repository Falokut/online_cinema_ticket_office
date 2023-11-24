package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/Falokut/online_cinema_ticket_office/email_service/internal/utils"
	"github.com/Falokut/online_cinema_ticket_office/email_service/pkg/logging"
	"github.com/segmentio/kafka-go"
)

type MailWorkerConfig struct {
	MailTypes          []string          `yaml:"mail_types"`
	MailSubjectsByType map[string]string `yaml:"mail_subjects_by_type"`
	TemplatesNames     map[string]string `yaml:"templates_names"`
	MaxWorkersCount    int               `yaml:"max_mail_workers_count" env:"MAX_MAIL_WORKERS_COUNT"`
}

type MailWorker struct {
	mailSender  *MailSender
	logger      logging.Logger
	cfg         MailWorkerConfig
	kafkaReader *kafka.Reader

	shutDownCtx context.CancelFunc
}

func NewMailWorker(mailSender *MailSender, logger logging.Logger, cfg MailWorkerConfig, kafkaReader *kafka.Reader) *MailWorker {
	w := MailWorker{mailSender: mailSender, logger: logger, cfg: cfg, kafkaReader: kafkaReader}
	return &w
}

func (w *MailWorker) Run() {
	ctx, cancel := context.WithCancel(context.Background())
	w.shutDownCtx = cancel
	for {
		mailData, m, err := w.getMailFromQueue(ctx)
		if err != nil {
			w.logger.Error(err)
			continue
		}
		err = w.routine(ctx, mailData, m)
		if err != nil {
			w.logger.Error(err)
		}
		select {
		case <-ctx.Done():
			return
		default:
			continue
		}
	}
}

func (w *MailWorker) ShutDown() {
	w.logger.Infoln("Mail worker shutting down")
	w.shutDownCtx()
}

type queueData struct {
	EmailTo  string  `json:"email"`
	URL      string  `json:"url"`
	MailType string  `json:"mail_type"`
	LinkTTL  float64 `json:"link_TTL"`
}

func (w *MailWorker) routine(ctx context.Context, mailData queueData, m kafka.Message) error {
	Expired := time.Since(m.Time) > time.Second*time.Duration(mailData.LinkTTL)
	if Expired {
		w.logger.Debugf("Message expired, message sended: %s. %s since message sended. linkTTL: %s",
			m.Time, time.Since(m.Time), time.Duration(mailData.LinkTTL))
		return nil
	}

	templateName, templateOk := w.cfg.TemplatesNames[mailData.MailType]
	if !templateOk {
		errorMessage := fmt.Sprintf("Can't find template name for mail type: %s. Skipping message", mailData.MailType)
		w.logger.Warningf(errorMessage)
		return errors.New(errorMessage)
	}

	subject, subjectOk := w.cfg.MailSubjectsByType[mailData.MailType]
	if !subjectOk {
		errorMessage := fmt.Sprintf("Can't find subject name mail type: %s. Skipping message", mailData.MailType)
		w.logger.Warningf(errorMessage)
		return errors.New(errorMessage)
	}

	LinkTTL, err := utils.ResolveTime(mailData.LinkTTL)
	if err != nil {
		errorMessage := fmt.Sprintf("Can't parse link ttl, unsupported seconds amount %f", mailData.LinkTTL)
		w.logger.Error(errorMessage)
		return errors.New(errorMessage)
	}

	data := &EmailData{URL: mailData.URL, Subject: subject, LinkTTL: LinkTTL}
	err = w.mailSender.SendEmail(data, mailData.EmailTo, templateName)
	if err != nil {
		w.logger.Error(err)
		return err
	}

	return nil
}

func (w *MailWorker) getMailFromQueue(ctx context.Context) (queueData, kafka.Message, error) {
	m, err := w.kafkaReader.ReadMessage(ctx)
	if err != nil {
		return queueData{}, kafka.Message{}, err
	}

	var data queueData
	err = json.Unmarshal(m.Value, &data)
	if err != nil {
		return queueData{}, kafka.Message{}, err
	}

	return data, m, nil
}
