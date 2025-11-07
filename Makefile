# Variabel
COMPOSE_FILE = docker-compose.yaml
BUILD_SCRIPT = sh ./build.sh

# Default target
.PHONY: all
all: build up

# Build binary / image
.PHONY: build
build:
	@echo ">> Running build script"
	@$(BUILD_SCRIPT)

# Start docker-compose
.PHONY: up
up:
	@echo ">> Starting docker-compose"
	docker-compose -f ./$(COMPOSE_FILE) up --scale billing-engine-http=3 --scale billing-engine-rpc=2 -d --build

# Stop docker-compose
.PHONY: down
down:
	@echo ">> Stopping docker-compose"
	docker-compose -f $(COMPOSE_FILE) down

# Restart docker-compose
.PHONY: restart
restart: down up

# Show logs
.PHONY: logs
logs:
	docker-compose -f $(COMPOSE_FILE) logs -f

# Clean up all containers & images
.PHONY: clean
clean: down
	@echo ">> Removing dangling images"
	-podman image prune -f

# Stop docker-compose
.PHONY: kitexgen
kitexgen:
	@echo ">> Generating a new kitex gen"
	rm -rf kitex_gen
	kitex -module github.com/yosuarichel/billing-engine ./conf/.idl/billing_engine/billing_engine_service.thrift