package kafka

import (
	"context"
	"fmt"
	"time"

	"github.com/IBM/sarama"
)

type TopicConfig struct {
	Name string
	Partitions int32
	ReplicationFactor int16
}

type TopicManager struct {
 	adminClient sarama.ClusterAdmin
}

func NewTopicManager(brokers []string) (*TopicManager, error) {
	config := sarama.NewConfig()
	config.Version = sarama.V2_8_0_0
	config.Admin.Timeout = 10 * time.Second
	
	adminClient, err := sarama.NewClusterAdmin(brokers, config)
	if err != nil {
		return nil, fmt.Errorf("Failed to create cluster admin: %v", err)
	}
	
	return &TopicManager{
		adminClient: adminClient,
	}, nil
}

func (tm *TopicManager) CreateTopic(ctx context.Context, cfg TopicConfig) error {
	topicDetail := &sarama.TopicDetail{
		NumPartitions: cfg.Partitions,
		ReplicationFactor: cfg.ReplicationFactor,
	}
	
	err := tm.adminClient.CreateTopic(cfg.Name, topicDetail, false)
	if err != nil {
		if err == sarama.ErrTopicAlreadyExists {
			return nil
		}
		return fmt.Errorf("Failed to create topic %v: %v", cfg.Name, err)
	}
	
	return nil
}

func (tm *TopicManager) DeleteTopic(ctx context.Context, name string) error {
	err := tm.adminClient.DeleteTopic(name)
	if err != nil {
		return fmt.Errorf("Failed to delete topic %v: %v", name, err)
	}
	return nil
}

func (tm *TopicManager) Close() error {
 	return tm.adminClient.Close()
}
