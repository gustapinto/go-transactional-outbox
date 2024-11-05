package kafka

import (
	"context"

	"github.com/twmb/franz-go/pkg/kgo"
)

func OpenQueueConnection(seeds []string, topic string) (*kgo.Client, error) {
	kafkaClient, err := kgo.NewClient(kgo.SeedBrokers(seeds...), kgo.ConsumeTopics(topic), kgo.ConsumerGroup("inventory-service"))
	if err != nil {
		return nil, err
	}

	if err := kafkaClient.Ping(context.Background()); err != nil {
		return nil, err
	}

	return kafkaClient, nil
}
