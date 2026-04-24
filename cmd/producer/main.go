package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"online-cinema-analytics/internal/application/event_service"
	"online-cinema-analytics/internal/infrastructure/config"
	"online-cinema-analytics/internal/infrastructure/generator"
	"online-cinema-analytics/internal/infrastructure/http"
	"online-cinema-analytics/internal/infrastructure/http/handler"
	"online-cinema-analytics/internal/infrastructure/kafka"
	"online-cinema-analytics/internal/infrastructure/logger"
)

func main() {
	lgr := logger.NewLogger()
	cfg := config.Load()

	prod, err := kafka.NewProducer(cfg.KafkaBrokers, cfg.KafkaTopic)
	if err != nil {
		lgr.Error("Error while creating producer", "Error", err)
	}
	defer prod.Close()

	evtSvc := eventservice.NewEventService(prod)
	evtHandler := handler.NewEventHandler(evtSvc)

	go func() {
		lgr.Info("HTTP server starting...", "Address", cfg.HTTPAddr)
		if err := http.Run(cfg.HTTPAddr, evtHandler); err != nil {
			lgr.Error("HTTP server failed.", "Error", err)
		}
		lgr.Info("HTTP server started.")
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	if cfg.GeneratorEnabled {
		gen := generator.NewEngine(evtSvc, cfg.GeneratorUsers, cfg.GeneratorMovies, cfg.GeneratorInterval)
		lgr.Info("Generator started.")
		go gen.Run(ctx)
	}

	<-ctx.Done()
	lgr.Info("Shutting down...")
}
