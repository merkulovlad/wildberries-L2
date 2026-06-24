package rabbitmq

import (
	wbrabbit "github.com/wb-go/wbf/rabbitmq"
	"github.com/wb-go/wbf/retry"
)

func NewClient(cfg Config) (*wbrabbit.RabbitClient, error) {
	retryStrategy := retry.Strategy{
		Attempts: defaultRetryAttempts,
		Delay:    defaultRetryDelay,
		Backoff:  defaultRetryBackoff,
	}

	return wbrabbit.NewClient(wbrabbit.ClientConfig{
		URL:            cfg.URL,
		ConnectionName: cfg.ConnectionName,
		ConnectTimeout: cfg.ConnectTimeout,
		Heartbeat:      cfg.Heartbeat,
		ReconnectStrat: retryStrategy,
		ProducingStrat: retryStrategy,
		ConsumingStrat: retryStrategy,
	})
}

func Setup(client *wbrabbit.RabbitClient, cfg Config) error {
	if err := client.DeclareExchange(
		cfg.Exchange,
		cfg.exchangeKind(),
		true,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	return client.DeclareQueue(
		cfg.Queue,
		cfg.Exchange,
		cfg.RoutingKey,
		true,
		false,
		true,
		nil,
	)
}
