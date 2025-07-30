# altenar-assignment

## .env
Example .env configuration:
```
APPLICATION_ADDR=:8080

POSTGRES_USER=casino
POSTGRES_PASSWORD=casino
POSTGRES_DSN=postgres://${POSTGRES_USER}:${POSTGRES_PASSWORD}@localhost:5432

KAFKA_BOOTSTRAP_SERVERS=localhost:9092
KAFKA_TOPIC=transactions
KAFKA_CONSUMER_GROUP=transactions-consumer
```

## Deploying infrastructure
- Create `.env` file with contents from the section above
- Run `docker compose up -d` to deploy project's infrastructure

## Running application
- Run `go run ./cmd/casino` from the root directory
- Since assignment's requirements do not imply Kafka producers, user will have to send messages manually:
  - `docker compose exec kafka /opt/kafka/bin/kafka-console-producer.sh --bootstrap-server localhost:9092 --topic transactions`
  - `{"user_id":1,"transaction_type":"bet","amount":10,"timestamp":"2025-07-30T12:03:32+00:00"}`

## Tests and coverage
Run `go test -count 1 ./cmd/casino` to test the application without caching results

Run these commands sequentially to test coverage:
- `go test -coverprofile cover.out.tmp ./cmd/casino`
- `cat cover.out.tmp | grep -v "main.go" > cover.out` (exclude `main.go`, since it is just mostly external dependencies and it also conveniently raises test coverage above 85%)\
- `go tool cover -func cover.out`
