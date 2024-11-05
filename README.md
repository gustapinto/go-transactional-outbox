# Go Transactional Outbox

A implementation of the [transactional outbox messaging pattern](https://microservices.io/patterns/data/transactional-outbox.html) using Go, Kafka and Postgres

![](https://raw.githubusercontent.com/gustapinto/go-transactional-outbox/main/docs/go-transactional-outbox-dark.jpg#gh-dark-mode-only)
![](https://raw.githubusercontent.com/gustapinto/go-transactional-outbox/main/docs/go-transactional-outbox-light.jpg#gh-light-mode-only)

## Running the services

### Requirements:

1. Docker
2. Docker Compose

### Instructions:

1. Clone this repository
2. (Optional, only if using Linux) Run the `setup-kafka.sh` script with `./scripts/setup-kafka.sh` to configure `kafka_data/` directory permissions
3. Start the docker containers with `docker compose up`
