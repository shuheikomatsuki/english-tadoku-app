# ==========================================================
# English Tadoku App - Makefile
# ==========================================================

# ----------------------------
# Variables
# ----------------------------
BACKEND_DIR=backend
FRONTEND_DIR=frontend
MIGRATE_CMD=docker compose run --rm migrate
DOCKER_COMPOSE=docker compose

# ----------------------------
# Docker lifecycle
# ----------------------------
.PHONY: up down build logs ps restart

up:
	$(DOCKER_COMPOSE) up -d
	@echo "✅ Containers started."

down:
	$(DOCKER_COMPOSE) down --remove-orphans
	@echo "🧹 Containers stopped and cleaned."

build:
	$(DOCKER_COMPOSE) build --no-cache
	@echo "🔧 Build completed."

logs:
	$(DOCKER_COMPOSE) logs -f

ps:
	$(DOCKER_COMPOSE) ps

restart:
	$(MAKE) down
	$(MAKE) up

# ----------------------------
# Backend commands
# ----------------------------
.PHONY: run-backend fmt lint test-backend

run-backend:
	cd $(BACKEND_DIR) && go run ./cmd/api

fmt:
	cd $(BACKEND_DIR) && go fmt ./...

lint:
	cd $(BACKEND_DIR) && golangci-lint run ./...

test-backend:
	cd $(BACKEND_DIR) && go test ./... -v

# ----------------------------
# Database migration
# ----------------------------
.PHONY: migrate-up migrate-down migrate-status migrate-create

migrate-up:
	$(MIGRATE_CMD) up
	@echo "✅ Database migrated up."

migrate-down:
	$(MIGRATE_CMD) down
	@echo "⚠️ Database rolled back."

migrate-status:
	$(MIGRATE_CMD) status

migrate-create:
	$(MIGRATE_CMD) bash

# ----------------------------
# Frontend commands
# ----------------------------
.PHONY: run-frontend build-frontend lint-frontend

run-frontend:
	cd $(FRONTEND_DIR) && npm run dev

build-frontend:
	cd $(FRONTEND_DIR) && npm run build

lint-frontend:
	cd $(FRONTEND_DIR) && npm run lint

# ----------------------------
# Utilities
# ----------------------------
.PHONY: clean env-check

clean:
	docker system prune -af --volumes
	@echo "🧽 Docker cleaned."

env-check:
	@echo "Current environment variables:"
	env | grep VITE_ || true
