setup-kafka:
	mkdir kafka_data && sudo chown -R 1001:1001 kafka_data

build-order-service:
	cd order-service && go build -o bin/order-service ./cmd

build-message-relay-service:
	cd message-relay-service && go build -o bin/message-relay-service ./cmd

build-inventory-service:
	cd inventory-service && go build -o bin/inventory-service ./cmd