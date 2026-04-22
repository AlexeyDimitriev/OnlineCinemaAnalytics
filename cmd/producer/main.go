package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"
	
	"online-cinema-analytics/internal/infrastructure/schema"
)

func main() {
	cfg := schema.Config{
		SchemaRegistryURL: getEnv("SCHEMA_REGISTRY_URL", "http://localhost:8081"),
		KafkaBrokers: []string{getEnv("KAFKA_BROKER", "localhost:9092")},
		TopicName: schema.TopicName,
		Partitions: 3,
		ReplicationFactor: 1,
	}

	log.Println("Initializing infrastructure...")
	
	if err := schema.Initialize(cfg); err != nil {
		log.Fatalf("Failed to initialize infrastructure: %v", err)
	}
	
	log.Println("Infrastructure initialized successfully")

	ctx, stop := signal.NotifyContext(
		context.Background(), 
		os.Interrupt,
		syscall.SIGTERM,
	)
	defer stop()
	
	<-ctx.Done()
	log.Println("Shutting down...")
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
