SCRIPTS_DIR=scripts
TESTDATA_DIR=testdata
GEN_AMOUNT ?= 10

include .env
export $(shell sed 's/=.*//' .env)

.PHONY: up down gen_json send_random_testdata send_testdata

up:
	@echo "Starting containers with config: $(CONFIG_PATH)"
	@docker-compose up -d

down:
	@docker-compose down

gen_json:
	@rm -rf $(TESTDATA_DIR)
	@echo "Generating test JSON files..."
	@go run $(SCRIPTS_DIR)/gen/generate_orders.go $(GEN_AMOUNT)
	@echo "Test data generated in $(TESTDATA_DIR)/ directory"

send_random_testdata: gen_json
	@echo "Sending test data to Kafka..."
	@go run $(SCRIPTS_DIR)/kafka/kafka_producer.go CONFIG_PATH=$(CONFIG_PATH)
	@echo "Test data sent to Kafka"

send_testdata:
	@echo "Sending test data to Kafka..."
	@go run $(SCRIPTS_DIR)/kafka/kafka_producer.go CONFIG_PATH=$(CONFIG_PATH)
	@echo "Test data sent to Kafka"