SCRIPTS_DIR=scripts
TESTDATA_DIR=testdata
GEN_AMOUNT ?= 10

up:
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
	@go run $(SCRIPTS_DIR)/kafka/kafka_producer.go
	@echo "Test data sent to Kafka"

send_testdata:
	@echo "Sending test data to Kafka..."
	@go run $(SCRIPTS_DIR)/kafka/kafka_producer.go
	@echo "Test data sent to Kafka"