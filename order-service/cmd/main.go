package main

import (
	"log"
	"os"

	"github.com/gustapinto/go-transactional-outbox/order-service/internal/repository/postgres"
	"github.com/gustapinto/go-transactional-outbox/order-service/internal/service"
)

var mockedOrders = []struct {
	Title    string
	Product  string
	Quantity int64
	Value    float64
}{
	{Title: "Example order 1", Product: "PRODUCT_1", Quantity: 1, Value: 10.00},
	{Title: "Example order 2", Product: "PRODUCT_1", Quantity: 3, Value: 30.00},
	{Title: "Example order 3", Product: "PRODUCT_1", Quantity: 5, Value: 50.00},
	{Title: "Example order 4", Product: "PRODUCT_2", Quantity: 10, Value: 150.00},
	{Title: "Example order 5", Product: "PRODUCT_3", Quantity: 2, Value: 7.50},
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
		id, err := orderService.Create(mockedOrder.Title, mockedOrder.Product, mockedOrder.Quantity, mockedOrder.Value)
		if err != nil {
			log.Printf("Failed to create order %d, got error: %s", i, err.Error())
			continue
		}

		log.Printf("Order created, ID=%s", id.String())
	}
}
