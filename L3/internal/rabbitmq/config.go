package rabbitmq

import "time"

const (
	defaultExchangeKind  = "direct"
	defaultContentType   = "application/json"
	defaultRetryAttempts = 3
	defaultRetryDelay    = time.Second
	defaultRetryBackoff  = 2
)

type Config struct {
	URL            string
	ConnectionName string
	Exchange       string
	ExchangeKind   string
	Queue          string
	RoutingKey     string
	ContentType    string
	ConnectTimeout time.Duration
	Heartbeat      time.Duration
	ConsumerTag    string
	Workers        int
	PrefetchCount  int
	RequeueOnError bool
}

func (c Config) exchangeKind() string {
	if c.ExchangeKind == "" {
		return defaultExchangeKind
	}
	return c.ExchangeKind
}

func (c Config) contentType() string {
	if c.ContentType == "" {
		return defaultContentType
	}
	return c.ContentType
}

func (c Config) consumerTag() string {
	if c.ConsumerTag == "" {
		return "notify-service"
	}
	return c.ConsumerTag
}

func (c Config) consumerWorkers() int {
	if c.Workers <= 0 {
		return 1
	}
	return c.Workers
}

func (c Config) prefetchCount() int {
	if c.PrefetchCount <= 0 {
		return c.consumerWorkers()
	}
	return c.PrefetchCount
}
