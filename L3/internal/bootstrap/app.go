package bootstrap

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"github.com/merkulovlad/wildberries-L3/notification_service/internal/cache"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/httpapi"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/notification"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/postgres"
	rabbitmqadapter "github.com/merkulovlad/wildberries-L3/notification_service/internal/rabbitmq"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/scheduler"
	"github.com/merkulovlad/wildberries-L3/notification_service/internal/sender"
	pgxdriver "github.com/wb-go/wbf/dbpg/pgx-driver"
	"github.com/wb-go/wbf/logger"
	wbrabbit "github.com/wb-go/wbf/rabbitmq"
	wbredis "github.com/wb-go/wbf/redis"
)

type App struct {
	logger          *log.Logger
	server          *http.Server
	db              *pgxdriver.Postgres
	redis           interface{ Close() error }
	rabbitClient    interface{ Close() error }
	workerScheduler interface{ Shutdown() error }
	cancel          context.CancelFunc
}

func Run() error {
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	app, err := New(ctx, stop)
	if err != nil {
		return err
	}
	app.Run()

	<-ctx.Done()
	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	app.Shutdown(shutdownCtx)
	return nil
}

func New(ctx context.Context, cancel context.CancelFunc) (*App, error) {
	stdLogger := log.New(os.Stdout, "notify ", log.LstdFlags|log.Lmicroseconds)
	appLogger, err := logger.InitLogger(logger.SlogEngine, "notify-service", envOr("APP_ENV", "local"))
	if err != nil {
		return nil, fmt.Errorf("init logger: %w", err)
	}

	db, err := pgxdriver.New(requiredEnv("POSTGRES_DSN"), appLogger)
	if err != nil {
		return nil, fmt.Errorf("connect postgres: %w", err)
	}
	if err := postgres.Setup(ctx, db); err != nil {
		db.Close()
		return nil, fmt.Errorf("setup postgres: %w", err)
	}

	baseRepo := postgres.NewRepository(db)
	redisClient := wbredis.New(envOr("REDIS_ADDR", "localhost:6379"), envOr("REDIS_PASSWORD", ""), envInt("REDIS_DB", 0))
	if err := redisClient.Ping(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("connect redis: %w", err)
	}
	repo := cache.NewRepository(baseRepo, redisClient, envDuration("CACHE_TTL", 5*time.Minute))

	rabbitCfg := rabbitConfig()
	rabbitClient, err := newRabbitClient(ctx, rabbitCfg)
	if err != nil {
		_ = redisClient.Close()
		db.Close()
		return nil, fmt.Errorf("connect rabbitmq: %w", err)
	}
	if err := rabbitmqadapter.Setup(rabbitClient, rabbitCfg); err != nil {
		_ = rabbitClient.Close()
		_ = redisClient.Close()
		db.Close()
		return nil, fmt.Errorf("setup rabbitmq: %w", err)
	}

	publisher := rabbitmqadapter.NewPublisher(rabbitClient, rabbitCfg)
	worker := notification.NewWorker(repo, publisher, envInt("WORKER_LIMIT", 100))
	workerScheduler, err := scheduler.NewWorkerScheduler(worker, envDuration("WORKER_INTERVAL", 10*time.Second))
	if err != nil {
		_ = rabbitClient.Close()
		_ = redisClient.Close()
		db.Close()
		return nil, fmt.Errorf("create worker scheduler: %w", err)
	}
	if err := workerScheduler.Start(ctx); err != nil {
		_ = workerScheduler.Shutdown()
		_ = rabbitClient.Close()
		_ = redisClient.Close()
		db.Close()
		return nil, fmt.Errorf("start worker scheduler: %w", err)
	}

	processor := notification.NewProcessor(repo, sender.NewLogSender(stdLogger))
	consumer := rabbitmqadapter.NewConsumer(rabbitClient, rabbitCfg, processor)
	go func() {
		if err := consumer.Start(ctx); err != nil && !errors.Is(err, context.Canceled) {
			stdLogger.Printf("rabbitmq consumer stopped: %v", err)
			cancel()
		}
	}()

	service := notification.NewService(repo)
	server := &http.Server{
		Addr:              envOr("NOTIFY_ADDR", ":8080"),
		Handler:           httpapi.NewRouter(service),
		ReadHeaderTimeout: 5 * time.Second,
	}

	return &App{
		logger:          stdLogger,
		server:          server,
		db:              db,
		redis:           redisClient,
		rabbitClient:    rabbitClient,
		workerScheduler: workerScheduler,
		cancel:          cancel,
	}, nil
}

func (a *App) Run() {
	go func() {
		a.logger.Printf("notify service listening on %s", a.server.Addr)
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			a.logger.Printf("listen and serve: %v", err)
			a.cancel()
		}
	}()
}

func (a *App) Shutdown(ctx context.Context) {
	if err := a.server.Shutdown(ctx); err != nil {
		a.logger.Printf("shutdown http server: %v", err)
	}
	if err := a.workerScheduler.Shutdown(); err != nil {
		a.logger.Printf("shutdown worker scheduler: %v", err)
	}
	if err := a.rabbitClient.Close(); err != nil {
		a.logger.Printf("close rabbitmq: %v", err)
	}
	if err := a.redis.Close(); err != nil {
		a.logger.Printf("close redis: %v", err)
	}
	a.db.Close()
}

func rabbitConfig() rabbitmqadapter.Config {
	return rabbitmqadapter.Config{
		URL:            envOr("RABBITMQ_URL", "amqp://guest:guest@localhost:5672/"),
		ConnectionName: envOr("RABBITMQ_CONNECTION_NAME", "notify-service"),
		Exchange:       envOr("RABBITMQ_EXCHANGE", "notifications"),
		ExchangeKind:   envOr("RABBITMQ_EXCHANGE_KIND", "direct"),
		Queue:          envOr("RABBITMQ_QUEUE", "notifications.send"),
		RoutingKey:     envOr("RABBITMQ_ROUTING_KEY", "notifications.send"),
		ConsumerTag:    envOr("RABBITMQ_CONSUMER_TAG", "notify-service"),
		Workers:        envInt("RABBITMQ_CONSUMER_WORKERS", 1),
		PrefetchCount:  envInt("RABBITMQ_PREFETCH_COUNT", 1),
	}
}

func newRabbitClient(ctx context.Context, cfg rabbitmqadapter.Config) (*wbrabbit.RabbitClient, error) {
	var lastErr error
	for attempt := 0; attempt < 12; attempt++ {
		client, err := rabbitmqadapter.NewClient(cfg)
		if err == nil {
			return client, nil
		}
		lastErr = err

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(5 * time.Second):
		}
	}
	return nil, lastErr
}

func envOr(key, def string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return def
}

func requiredEnv(key string) string {
	value := os.Getenv(key)
	if value == "" {
		log.Fatalf("%s is required", key)
	}
	return value
}

func envInt(key string, def int) int {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	parsed, err := strconv.Atoi(value)
	if err != nil {
		log.Fatalf("%s must be an integer: %v", key, err)
	}
	return parsed
}

func envDuration(key string, def time.Duration) time.Duration {
	value := os.Getenv(key)
	if value == "" {
		return def
	}
	parsed, err := time.ParseDuration(value)
	if err != nil {
		log.Fatalf("%s must be a duration, for example 10s: %v", key, err)
	}
	return parsed
}
