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

# Start podman compose
.PHONY: up
up:
	@echo ">> Starting podman compose"
	podman compose -f ./$(COMPOSE_FILE) up -d --build

# Stop podman compose
.PHONY: down
down:
	@echo ">> Stopping podman compose"
	podman compose -f $(COMPOSE_FILE) down

# Restart podman compose
.PHONY: restart
restart: down up

# Show logs
.PHONY: logs
logs:
	podman compose -f $(COMPOSE_FILE) logs -f

# Clean up all containers & images
.PHONY: clean
clean: down
	@echo ">> Removing dangling images"
	-podman image prune -f

# Stop podman compose
.PHONY: kitexgen
kitexgen:
	@echo ">> Generating a new kitex gen"
	rm -rf kitexgen
	kitex -module github.com/yosuarichel/billing-engine ./conf/.idl/billing_engine/billing_engine_service.thrift