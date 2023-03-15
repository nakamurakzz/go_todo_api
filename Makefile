.PHONY: help build build-local up down logs ps test
.DEFAULT_GOAL := help

DOCKER_TAG := latest
build: ## Build the docker image
	docker build -t docker.io/library/gotodo:$(DOCKER_TAG) --target deploy .

build-local: ## Build the docker image locally
	docker compose build --no-cache

up: ## Start the docker containers
	docker compose up -d

down: ## Stop the docker containers
	docker compose down

logs: ## Show the logs
	docker compose logs -f

ps: ## Show the running containers
	docker compose ps

test: ## Run the tests
	go test -race -shuffle=on ./...
