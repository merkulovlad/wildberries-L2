package rabbitmq

import (
	"context"
	"encoding/json"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	"github.com/rabbitmq/amqp091-go"
	wbrabbit "github.com/wb-go/wbf/rabbitmq"
)

func NewConsumer(client *wbrabbit.RabbitClient, cfg Config, processor *notification.Processor) *wbrabbit.Consumer {
	return wbrabbit.NewConsumer(
		client,
		wbrabbit.ConsumerConfig{
			Queue:         cfg.Queue,
			ConsumerTag:   cfg.consumerTag(),
			AutoAck:       false,
			Workers:       cfg.consumerWorkers(),
			PrefetchCount: cfg.prefetchCount(),
			Nack:          wbrabbit.NackConfig{Requeue: cfg.RequeueOnError},
		},
		func(ctx context.Context, d amqp091.Delivery) error {
			var msg message
			if err := json.Unmarshal(d.Body, &msg); err != nil {
				return err
			}
			return processor.Process(ctx, msg.NotificationID)
		},
	)
}
