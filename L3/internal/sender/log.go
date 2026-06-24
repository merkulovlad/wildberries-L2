package sender

import (
	"context"
	"log"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
)

type LogSender struct {
	logger *log.Logger
}

func NewLogSender(logger *log.Logger) *LogSender {
	return &LogSender{logger: logger}
}

func (s *LogSender) Send(_ context.Context, n notification.Notification) error {
	s.logger.Printf("send notification id=%s recipient=%s message=%q", n.ID, n.Recipient, n.Message)
	return nil
}
