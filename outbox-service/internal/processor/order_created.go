package processor

import (
	"context"
	"log"

	"github.com/gustapinto/go-transactional-outbox/outbox-service/internal/model"
	"github.com/twmb/franz-go/pkg/kgo"
)

type OrderCreated struct {
	KafkaClient *kgo.Client
}

func (o OrderCreated) Process(ctx context.Context, event model.OutboxEvent) error {
	record := &kgo.Record{
		Value:   event.Data,
		Headers: []kgo.RecordHeader{},
		Topic:   "ORDER-CREATED-TOPIC",
	}

	err := o.KafkaClient.ProduceSync(ctx, record).FirstErr()
	if err != nil {
		return err
	}

	log.Printf("Event %s sent to %s topic", event.String(), record.Topic)

	return nil
}
