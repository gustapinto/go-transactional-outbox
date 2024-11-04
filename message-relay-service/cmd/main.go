package main

import (
	"log"
	"os"
	"strings"
	"time"

	"github.com/gustapinto/go-transactional-outbox/message-relay-service/internal/queue/kafka"
	"github.com/gustapinto/go-transactional-outbox/message-relay-service/internal/repository/postgres"
	"github.com/gustapinto/go-transactional-outbox/message-relay-service/internal/service"
)

var (
	PostgresDSN = os.Getenv("POSTGRES_DSN")
	KafkaSeeds  = strings.Split(os.Getenv("KAFKA_SEEDS"), ",")
)

func main() {
	log.Println("Starting Outbox Service")

	db, err := postgres.OpenDatabaseConnection(PostgresDSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	kafkaClient, err := kafka.OpenQueueConnection(KafkaSeeds)
	if err != nil {
		log.Fatalln(err)
	}
	defer kafkaClient.Close()

	if err := kafka.InitializeQueue(kafkaClient); err != nil {
		log.Fatalln(err)
	}

	outboxRepository := postgres.Outbox{DB: db}
	outboxService := service.Outbox{OutboxRepository: outboxRepository}
	orderCreatedEventProcessor := kafka.OrderCreatedProcessor{KafkaClient: kafkaClient}

	processorMapping := map[string]service.OutboxEventProcessor{
		"ORDER_CREATED": orderCreatedEventProcessor,
	}

	for {
		time.Sleep(5 * time.Second)

		err := outboxService.GetAndProcessNonProcessedOutboxEvents(processorMapping)
		if err != nil {
			log.Printf("Error while processing events: %s", err.Error())
			continue
		}
	}
}
