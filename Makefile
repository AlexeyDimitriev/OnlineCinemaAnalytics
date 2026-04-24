.PHONY: help up down clean reset-db reset-kafka test test-gen logs

PRODUCER_URL := http://localhost:8080
CLICKHOUSE_HOST := localhost
CLICKHOUSE_PORT := 9000
CLICKHOUSE_DB := movie_analytics
CLICKHOUSE_USER := analytics
CLICKHOUSE_PASSWORD := analytics

TIMEOUT_INTEGRATION := 60s
TIMEOUT_GENERATOR := 90s
STARTUP_WAIT := 20


help:
	@echo "Available commands:"
	@echo "  make up           - Start all services (infra + producer)"
	@echo "  make down         - Stop all containers"
	@echo "  make clean        - Full reset: stop containers + remove volumes"
	@echo "  make reset-db     - Reset ClickHouse data only"
	@echo "  make reset-kafka  - Reset Kafka topics only"
	@echo "  make test         - Run integration test (HTTP → Kafka → ClickHouse)"
	@echo "  make test-gen     - Run generator test (realistic scenarios)"
	@echo "  make logs         - Show producer logs"

up:
	@echo "Starting infrastructure..."
	docker-compose up -d zookeeper kafka schema-registry clickhouse
	@echo "Waiting for infrastructure to be ready ($(STARTUP_WAIT)s)..."
	@sleep $(STARTUP_WAIT)
	@$(MAKE) _init-infra
	@echo "Building and starting producer..."
	docker-compose up -d --build producer
	@sleep 5
	@$(MAKE) _check-producer
	@echo "All services are up and running"
	@echo "   Producer: $(PRODUCER_URL)"
	@echo "   Kafka: localhost:9092"
	@echo "   Schema Registry: localhost:8081"
	@echo "   ClickHouse: http://localhost:8123"

down:
	@echo "Stopping services..."
	docker-compose down
	@echo "Services stopped"

clean: down
	@echo "Removing volumes (all data will be lost)..."
	-docker volume rm onlinecinemaanalytics_clickhouse_data 2>/dev/null || true
	-docker volume rm onlinecinemaanalytics_kafka_data 2>/dev/null || true
	-docker volume rm onlinecinemaanalytics_zookeeper_data 2>/dev/null || true
	@echo "Clean complete. Run 'make up' to start fresh."

_init-infra:
	@echo "Initializing ClickHouse..."
	@$(MAKE) _wait-clickhouse
	@$(MAKE) _apply-migrations
	@echo "Initializing Kafka..."
	@$(MAKE) _create-kafka-topics

_ensure-infra:
	@curl -sf http://localhost:8123/ping >/dev/null 2>&1 || (echo "✗ ClickHouse not ready" && exit 1)
	@docker exec kafka kafka-broker-api-versions --bootstrap-server localhost:9092 >/dev/null 2>&1 || (echo "✗ Kafka not ready" && exit 1)
	@echo "Infrastructure is ready"

_wait-clickhouse:
	@for i in $$(seq 1 30); do \
		if docker exec clickhouse clickhouse-client -u $(CLICKHOUSE_USER) --password $(CLICKHOUSE_PASSWORD) --query "SELECT 1" >/dev/null 2>&1; then \
		echo "   ClickHouse is ready"; \
		exit 0; \
		fi; \
		echo "   Attempt $$i/30: waiting for ClickHouse..."; \
		sleep 2; \
		done; \
		echo "ClickHouse did not become ready"; \
		exit 1

_apply-migrations:
	@docker exec -i clickhouse clickhouse-client \
		-u $(CLICKHOUSE_USER) \
		--password $(CLICKHOUSE_PASSWORD) \
		--multiquery \
		< migrations/clickhouse/001_init_schemas.sql \
		2>/dev/null || echo "   Migrations may have been applied already"
	@echo "   Migrations applied"

_create-kafka-topics:
	@docker exec kafka kafka-topics \
		--bootstrap-server localhost:9092 \
		--create \
		--topic movie-events \
		--partitions 3 \
		--replication-factor 1 \
		--if-not-exists \
		>/dev/null 2>&1 || true
	@echo "   Kafka topic 'movie-events' ready"

_check-producer:
	@for i in $$(seq 1 30); do \
		if curl -sf $(PRODUCER_URL)/health >/dev/null 2>&1; then \
		echo "   Producer is healthy"; \
		exit 0; \
		fi; \
		echo "   Attempt $$i/30: waiting for producer..."; \
		sleep 2; \
		done; \
		echo "Producer did not become healthy"; \
		exit 1

reset-db:
	@echo "Resetting ClickHouse..."
	docker-compose down clickhouse
	-docker volume rm onlinecinemaanalytics_clickhouse_data 2>/dev/null || true
	docker-compose up -d clickhouse
	@echo "Waiting for ClickHouse initialization..."
	@sleep $(STARTUP_WAIT)
	@$(MAKE) _apply-migrations
	@echo "ClickHouse reset complete"

reset-kafka:
	@echo "Resetting Kafka topic..."
	@docker exec kafka kafka-topics \
		--bootstrap-server localhost:9092 \
		--delete \
		--topic movie-events \
		>/dev/null 2>&1 || true
	@$(MAKE) _create-kafka-topics
	@echo "Kafka topic reset complete"

test: _ensure-up
	@echo "Running integration test (pipeline)..."
	@PRODUCER_URL=$(PRODUCER_URL) \
	CLICKHOUSE_HOST=$(CLICKHOUSE_HOST) \
	CLICKHOUSE_PORT=$(CLICKHOUSE_PORT) \
	CLICKHOUSE_DB=$(CLICKHOUSE_DB) \
	CLICKHOUSE_USER=$(CLICKHOUSE_USER) \
	CLICKHOUSE_PASSWORD=$(CLICKHOUSE_PASSWORD) \
	go test -v -timeout $(TIMEOUT_INTEGRATION) ./internal/test/integration -run TestFullPipeline

test-gen: _ensure-infra
	@echo "Running generator test..."
	@echo "   Preparing producer with generator enabled..."
	@docker rm -f producer 2>/dev/null || true
	@docker-compose build producer >/dev/null
	@docker-compose run --rm -d \
		-e GENERATOR_ENABLED=true \
		-e GENERATOR_USERS=1 \
		-e GENERATOR_MOVIES=1 \
		-e GENERATOR_INTERVAL=2s \
		--name producer \
	producer
	@echo "   Waiting 30s for event generation..."
	@sleep 30
	@echo "   Running Go test..."
	@PRODUCER_URL=$(PRODUCER_URL) \
		CLICKHOUSE_HOST=$(CLICKHOUSE_HOST) \
		CLICKHOUSE_PORT=$(CLICKHOUSE_PORT) \
		CLICKHOUSE_DB=$(CLICKHOUSE_DB) \
		CLICKHOUSE_USER=$(CLICKHOUSE_USER) \
		CLICKHOUSE_PASSWORD=$(CLICKHOUSE_PASSWORD) \
		go test -v -timeout $(TIMEOUT_GENERATOR) ./internal/test/integration -run TestGeneratorFlow; \
		EXIT_CODE=$$?; \
		echo "   Restoring producer to normal mode..."; \
		docker rm -f producer 2>/dev/null || true; \
		docker-compose up -d --build producer; \
	exit $$EXIT_CODE

_ensure-up:
	@curl -sf $(PRODUCER_URL)/health >/dev/null 2>&1 || (echo "✗ Producer not running. Run 'make up' first." && exit 1)
	@docker exec clickhouse clickhouse-client -u $(CLICKHOUSE_USER) --password $(CLICKHOUSE_PASSWORD) --query "SELECT 1" >/dev/null 2>&1 || (echo "✗ ClickHouse not ready. Run 'make up' first." && exit 1)
	@echo "Services are ready"

logs:
	@docker logs -f producer

status:
	@docker-compose ps
	@echo ""
	@echo "Tables in ClickHouse:"
	@docker exec clickhouse clickhouse-client -u $(CLICKHOUSE_USER) --password $(CLICKHOUSE_PASSWORD) --query "SHOW TABLES" 2>/dev/null || echo "   (ClickHouse not ready)"
	@echo ""
	@echo "Kafka topics:"
	@docker exec kafka kafka-topics --bootstrap-server localhost:9092 --list 2>/dev/null || echo "   (Kafka not ready)"
