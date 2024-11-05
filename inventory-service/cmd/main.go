package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/gustapinto/go-transactional-outbox/inventory-service/internal/model"
	"github.com/gustapinto/go-transactional-outbox/inventory-service/internal/queue/kafka"
	"github.com/gustapinto/go-transactional-outbox/inventory-service/internal/repository/postgres"
	"github.com/gustapinto/go-transactional-outbox/inventory-service/internal/service"
)

var (
	PostgresDSN = os.Getenv("POSTGRES_DSN")
	KafkaSeeds  = strings.Split(os.Getenv("KAFKA_SEEDS"), ",")
)

func main() {
	log.Println("Starting Inventory Service")

	db, err := postgres.OpenDatabaseConnection(PostgresDSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := postgres.InitializeDatabase(db); err != nil {
		log.Fatalln(err)
	}

	kafkaClient, err := kafka.OpenQueueConnection(KafkaSeeds, "ORDER-CREATED-TOPIC")
	if err != nil {
		log.Fatalln(err)
	}
	defer kafkaClient.Close()

	inventoryRepository := postgres.Inventory{DB: db}
	inventoryService := service.Inventory{InventoryRepository: inventoryRepository}

	for {
		time.Sleep(1 * time.Second)

		fetches := kafkaClient.PollFetches(context.Background())
		if errs := fetches.Errors(); len(errs) > 0 {
			log.Fatalln(errs)
		}

		iter := fetches.RecordIter()
		for !iter.Done() {
			record := iter.Next()

			fmt.Printf("Processing KafkaMessage[Key=%s]", string(record.Key))

			var event model.OrderCreatedEvent
			if err := json.Unmarshal(record.Value, &event); err != nil {
				fmt.Printf("Failed to process KafkaMessage[Key=%s], got error %s", string(record.Key), err.Error())
				continue
			}

			if err := inventoryService.ProcessOrder(event); err != nil {
				fmt.Printf("Failed to process KafkaMessage[Key=%s], got service error %s", string(record.Key), err.Error())
				continue
			}
		}
	}
}
