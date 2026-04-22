.PHONY: docker-up
docker-up:
	@echo "Starting services..."
	docker-compose up -d
	@echo "Waiting for services to be ready..."
	@sleep 15
	@echo "Services started"
	@echo "Kafka: localhost:9092"
	@echo "Schema Registry: localhost:8081"
	@echo "ClickHouse: localhost:8123"

.PHONY: docker-down
docker-down:
	@echo "Stopping services..."
	docker-compose down

.PHONY: init-schemas
init-schemas:
	@echo "Initializing schemas..."
	go run cmd/producer/main.go

.PHONY: clean
clean: docker-down
	@echo "Cleaning up..."
	docker volume rm movie-analytics_clickhouse_data || true
