package main

import (
	"log/slog"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/kostenbl4/wb_task_L0/internal/cache"
	"github.com/kostenbl4/wb_task_L0/internal/kafka"
	"github.com/kostenbl4/wb_task_L0/internal/store"
)

type application struct {
	config        config
	store         *store.Storage
	cache         cache.Cache[store.Order]
	logger        *slog.Logger
	kafkaConsumer *kafka.KafkaConsumer
}

type config struct {
	addr  string
	db    dbConfig
	kafka kafkaConfig
}

type dbConfig struct {
	addr         string
	maxOpenConns int
	maxIdleConns int
	maxIdleTime  string
}

type kafkaConfig struct {
	addr  []string
	topic string
	group string
}

func (app *application) mount() http.Handler {
	r := chi.NewRouter()

	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(middleware.Timeout(60 * time.Second))

	r.Get("/health", app.healthCheckHandler)

	r.Get("/", app.HandleIndex)

	r.Route("/orders", func(r chi.Router) {
		r.Get("/{id}", app.GetOrderById)
		r.Post("/", app.CreateOrder)
		r.Delete("/{id}", app.DeleteOrder)
	})

	return r
}

func (app *application) run(r http.Handler) error {
	app.logger.Info("Starting server", "address", app.config.addr)

	srv := &http.Server{
		Addr:         app.config.addr,
		Handler:      r,
		WriteTimeout: time.Second * 30,
		ReadTimeout:  time.Second * 10,
		IdleTimeout:  time.Minute,
	}

	return srv.ListenAndServe()
}
