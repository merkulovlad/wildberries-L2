package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	wbrabbit "github.com/wb-go/wbf/rabbitmq"
)

var _ notification.Publisher = (*Publisher)(nil)

type Publisher struct {
	publisher  *wbrabbit.Publisher
	routingKey string
}

func NewPublisher(client *wbrabbit.RabbitClient, cfg Config) *Publisher {
	return &Publisher{
		publisher:  wbrabbit.NewPublisher(client, cfg.Exchange, cfg.contentType()),
		routingKey: cfg.RoutingKey,
	}
}

func (p *Publisher) Publish(ctx context.Context, notificationID string) error {
	body, err := json.Marshal(message{NotificationID: notificationID})
	if err != nil {
		return err
	}

	return p.publisher.Publish(ctx, body, p.routingKey)
}

type message struct {
	NotificationID string `json:"notification_id"`
}
