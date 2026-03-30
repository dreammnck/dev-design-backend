# Makefile for Go Backend

# Variables
APP_NAME := backend
MAIN_FILE := main.go
BIN_DIR := bin
EXEC := $(BIN_DIR)/$(APP_NAME)

# Default target
.PHONY: all
all: run

# Run the application
.PHONY: run
run:
	go run $(MAIN_FILE)

# Build the application
.PHONY: build
build:
	@mkdir -p $(BIN_DIR)
	go build -o $(EXEC) $(MAIN_FILE)

# Database operations
.PHONY: db-up
db-up:
	docker-compose up -d db

.PHONY: db-down
db-down:
	docker-compose down

.PHONY: db-logs
db-logs:
	docker-compose logs -f db

# Dependency management
.PHONY: tidy
tidy:
	go mod tidy

# Testing
.PHONY: test
test:
	go test -v -race -cover ./...

# Cleanup
.PHONY: clean
clean:
	rm -rf $(BIN_DIR)

# Cloud Run Deployment
GCP_PROJECT_ID ?= dev-design-491813
GCP_REGION ?= asia-southeast1
GCP_SERVICE_NAME ?= cu-pop-backend
CLOUD_SQL_INSTANCE ?= $(GCP_PROJECT_ID):$(GCP_REGION):dev-design

# Database settings for production
P_DB_USER ?= postgres
P_DB_PASS ?= StrongSecret123!
P_DB_NAME ?= tickets_db
PAYMENT_GATEWAY_URL ?= https://mock-payment-449892262369.asia-southeast1.run.app

.PHONY: deploy
deploy:
	@echo "Deploying to Cloud Run..."
	gcloud run deploy $(GCP_SERVICE_NAME) \
		--source . \
		--project $(GCP_PROJECT_ID) \
		--region $(GCP_REGION) \
		--allow-unauthenticated \
		--port=8080 \
		--add-cloudsql-instances="$(CLOUD_SQL_INSTANCE)" \
		--set-env-vars="PAYMENT_GATEWAY_URL=$(PAYMENT_GATEWAY_URL),DB_HOST=/cloudsql/$(CLOUD_SQL_INSTANCE),DB_USER=$(P_DB_USER),DB_PASSWORD=$(P_DB_PASS),DB_NAME=$(P_DB_NAME),DB_PORT=5432"

# Help
.PHONY: help
help:
	@echo "Available commands:"
	@echo "  make run        - Run the application directly"
	@echo "  make build      - Build the application binary"
	@echo "  make db-up      - Start the database container (PostgreSQL)"
	@echo "  make db-down    - Stop and remove the database container"
	@echo "  make db-logs    - Tail the database container logs"
	@echo "  make tidy       - Tidy Go modules"
	@echo "  make test       - Run all tests"
	@echo "  make clean      - Remove build artifacts"
	@echo "  make deploy     - Deploy the application to Google Cloud Run"
