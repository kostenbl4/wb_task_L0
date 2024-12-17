package main

import (
	"fmt"
	"log/slog"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
	"github.com/kostenbl4/wb_task_L0/internal/cache"
	"github.com/kostenbl4/wb_task_L0/internal/database"
	"github.com/kostenbl4/wb_task_L0/internal/env"
	"github.com/kostenbl4/wb_task_L0/internal/kafka"
	"github.com/kostenbl4/wb_task_L0/internal/store"
)

func main() {
	
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	slog.SetDefault(logger)

	err := godotenv.Load()
	if err != nil {
		slog.Error("Error loading .env file, using environment variables instead")
	}
	slog.Debug("Environment variables loaded successfully")

	dbCfg := dbConfig{
		addr: fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
			env.GetString("DB_HOST", "localhost"),
			env.GetString("DB_PORT", "5432"),
			env.GetString("DB_USER", "postgres"),
			env.GetString("DB_PASSWORD", "postgres"),
			env.GetString("DB_NAME", "orders")),
		maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
		maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 30),
		maxIdleTime:  env.GetString("DB_MAX_IDLE_TIME", "15m"),
	}

	kafkaCfg := kafkaConfig{
		addr:  []string{env.GetString("KAFKA_ADDR", "localhost:9092")},
		topic: env.GetString("KAFKA_TOPIC", "orders"),
		group: env.GetString("KAFKA_GROUP", "orders-group"),
	}

	cfg := config{
		addr:  env.GetString("ADDR", ":8080"),
		db:    dbCfg,
		kafka: kafkaCfg,
	}

	slog.Debug("Configuration loaded", "config", cfg)
	
	db, err := database.New(cfg.db.addr, cfg.db.maxOpenConns, cfg.db.maxIdleConns, cfg.db.maxIdleTime)
	if err != nil {
		logger.Error("Failed to connect database: " + err.Error())
		os.Exit(1)
	}
	slog.Info("Connected to the database")

	ttl := time.Duration(env.GetInt("CACHE_TTL_SEC", 10800)) * time.Second
	cache := cache.NewCache[store.Order](ttl)
	slog.Info("Cache initialized", "ttl", ttl)

	store := store.NewStorage(db)
	slog.Info("Storage initialized")

	app := application{
		config: cfg,
		store:  store,
		cache:  cache,
		logger: logger,
	}

	kafkaConsumer, err := kafka.NewKafkaConsumer(
		cfg.kafka.topic,
		cfg.kafka.group,
		cfg.kafka.addr,
		app.GetConsumerFunc(),
	)
	if err != nil {
		logger.Error("Failed to connect kafka: " + err.Error())
		os.Exit(1)
	}
	slog.Info("Kafka consumer connected")

	app.kafkaConsumer = kafkaConsumer
	r := app.mount()
	slog.Info("Router mounted")

	sigs := make(chan os.Signal, 1)
	done := make(chan bool, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigs
		app.logger.Info("Application is shutting down (not gracefully yet)")
		done <- true
	}()

	go app.kafkaConsumer.Start()
	slog.Info("Kafka consumer started")

	go func() {
		if err := app.run(r); err != nil {
			app.logger.Error("Application running error: " + err.Error())
			done <- true
		}
	}()
	slog.Info("Server running")

	<-done
	app.kafkaConsumer.Stop()
	slog.Info("Kafka consumer stopped")
	slog.Info("Application shutdown complete")
}
