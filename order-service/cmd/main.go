package main

import (
	"log"
	"os"

	"github.com/gustapinto/go-transactional-outbox/order-service/internal/model"
	"github.com/gustapinto/go-transactional-outbox/order-service/internal/repository/postgres"
	"github.com/gustapinto/go-transactional-outbox/order-service/internal/service"
)

var mockedOrders = []model.CreateOrderPayload{
	{
		Title: "Example order 1",
		Value: 10.00,
	},
	{
		Title: "Example order 2",
		Value: 20.50,
	},
	{
		Title: "Example order 3",
		Value: 25.00,
	},
	{
		Title: "Example order 4",
		Value: -25.00,
	},
	{
		Title: "Example order 5",
		Value: 5.75,
	},
}

var (
	PostgresDSN = os.Getenv("POSTGRES_DSN")
)

func main() {
	log.Println("Starting Order Service")

	db, err := postgres.OpenDatabaseConnection(PostgresDSN)
	if err != nil {
		log.Fatalln(err)
	}
	defer db.Close()

	if err := postgres.InitializeDatabase(db); err != nil {
		log.Fatalln(err)
	}

	orderRepository := postgres.Order{DB: db}
	orderService := service.Order{OrderRepository: orderRepository}

	for i, mockedOrder := range mockedOrders {
		id, err := orderService.Create(mockedOrder.Title, mockedOrder.Value)
		if err != nil {
			log.Printf("Failed to create order %d, got error: %s", i, err.Error())
			continue
		}

		log.Printf("Order created, ID=%s", id.String())
	}
}
