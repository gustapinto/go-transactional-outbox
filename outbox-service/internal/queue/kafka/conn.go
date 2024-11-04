package kafka

import (
	"context"
	"strings"

	"github.com/twmb/franz-go/pkg/kadm"
	"github.com/twmb/franz-go/pkg/kgo"
)

func OpenQueueConnection(seeds []string) (*kgo.Client, error) {
	kafkaClient, err := kgo.NewClient(kgo.SeedBrokers(seeds...))
	if err != nil {
		return nil, err
	}

	if err := kafkaClient.Ping(context.Background()); err != nil {
		return nil, err
	}

	return kafkaClient, nil
}

func InitializeQueue(client *kgo.Client) error {
	partitions := int32(1)
	replicationFactor := int16(len(client.DiscoveredBrokers()))
	configs := map[string]*string{}
	topics := []string{
		"ORDER-CREATED-TOPIC",
	}

	kafkaAdminClient := kadm.NewClient(client)

	for _, topic := range topics {
		_, err := kafkaAdminClient.CreateTopic(
			context.Background(),
			partitions,
			replicationFactor,
			configs,
			topic)
		if err != nil {
			if strings.Contains(err.Error(), "TOPIC_ALREADY_EXISTS") {
				return nil
			}

			return err
		}
	}

	return nil
}
