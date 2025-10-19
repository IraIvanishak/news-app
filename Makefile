# Makefile for News App

# Environment variables
ENV_FILE := .env
COMPOSE_FILE := docker-compose.yml
COMPOSE_OVERRIDE := docker-compose.override.yml

# Docker Compose Commands
.PHONY: up down rebuild logs test clean

# Start the application in production mode
up:
	docker-compose -f $(COMPOSE_FILE) up -d

# Start the application in development mode
dev:
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_OVERRIDE) up -d

# Stop the application
down:
	docker-compose down

# Rebuild and start containers
rebuild:
	docker-compose -f $(COMPOSE_FILE) up -d --build

# Development rebuild
dev-rebuild:
	docker-compose -f $(COMPOSE_FILE) -f $(COMPOSE_OVERRIDE) up -d --build

# View logs
logs:
	docker-compose logs -f

# Run tests
test:
	go test ./internal/adapter/repositories -v


# Database operations
mongo-cli:
	docker-compose exec news_app_database mongosh

# Go dependencies
deps:
	go mod tidy
	go mod download

