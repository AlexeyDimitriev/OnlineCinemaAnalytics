package schema

import (
	"context"
	"fmt"
	"log"
	
	"online-cinema-analytics/internal/infrastructure/kafka"
)

type Config struct {
	SchemaRegistryURL string
	KafkaBrokers []string
	TopicName string
	Partitions int32
	ReplicationFactor int16
}

func Initialize(cfg Config) error {
	ctx := context.Background()

	log.Printf("Registering schema in Schema Registry: %v", cfg.SchemaRegistryURL)
	
	schemaClient := NewSchemaRegistryClient(cfg.SchemaRegistryURL)
	
	version, err := schemaClient.RegisterSchema(SchemaSubject, MovieEventAvroSchema)
	if err != nil {
		return fmt.Errorf("Failed to register schema: %v", err)
	}
	
	log.Printf("Schema registered successfully (version %v)", version)

	log.Printf("Creating Kafka topic: %v (partitions=%v, replication=%v)", 
	cfg.TopicName, cfg.Partitions, cfg.ReplicationFactor)
	
	topicManager, err := kafka.NewTopicManager(cfg.KafkaBrokers)
	if err != nil {
		return fmt.Errorf("Failed to create topic manager: %v", err)
	}
	defer topicManager.Close()
	
	topicCfg := kafka.TopicConfig{
		Name: cfg.TopicName,
		Partitions: cfg.Partitions,
		ReplicationFactor: cfg.ReplicationFactor,
	}
	
	err = topicManager.CreateTopic(ctx, topicCfg)
	if err != nil {
		return fmt.Errorf("Failed to create topic: %v", err)
	}
	
	log.Printf("Topic \"%v\" created successfully", cfg.TopicName)
	
	return nil
}
