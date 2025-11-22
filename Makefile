.PHONY: help build build-backend build-frontend run run-backend run-frontend test test-backend test-frontend docker-build docker-run docker-run-local docker-stop docker-up docker-down docker-logs clean install install-backend install-frontend

# Default target
help:
	@echo "Available targets:"
	@echo "  make install          - Install all dependencies (backend + frontend)"
	@echo "  make build           - Build both backend and frontend"
	@echo "  make run             - Run both backend and frontend in development mode"
	@echo "  make test            - Run all tests"
	@echo "  make docker-build    - Build Docker image"
	@echo "  make docker-run      - Run Docker container (Cloud Run mode, no service account)"
	@echo "  make docker-run-local - Run Docker container with service account file (local dev)"
	@echo "  make docker-stop     - Stop running Docker container"
	@echo "  make docker-up       - Start services with Docker Compose"
	@echo "  make docker-down     - Stop Docker Compose services"
	@echo "  make docker-logs     - View Docker Compose logs"
	@echo "  make clean           - Clean build artifacts"

# Install dependencies
install: install-backend install-frontend

install-backend:
	@echo "Installing backend dependencies..."
	cd backend && go mod download

install-frontend:
	@echo "Installing frontend dependencies..."
	cd frontend && npm install

# Build targets
build: build-backend build-frontend

build-backend:
	@echo "Building backend..."
	cd backend && go build -o server ./cmd/server

build-frontend:
	@echo "Building frontend..."
	cd frontend && npm run build

# Run targets
run-backend:
	@echo "Running backend..."
	cd backend && go run cmd/server/main.go

run-frontend:
	@echo "Running frontend..."
	cd frontend && npm run dev

run:
	@echo "Starting both backend and frontend..."
	@make -j2 run-backend run-frontend

# Test targets
test: test-backend test-frontend

test-backend:
	@echo "Running backend tests..."
	cd backend && go test ./...

test-frontend:
	@echo "Running frontend tests..."
	cd frontend && npm run lint

# Docker targets
docker-build:
	@echo "Building Docker image..."
	docker build -t ai-india-workshop:latest .

docker-run: docker-build
	@echo "Running Docker container (Cloud Run mode)..."
	@docker stop ai-india-workshop 2>/dev/null || true
	@docker rm ai-india-workshop 2>/dev/null || true
	@if [ -f .env ]; then \
		export FIRESTORE_SUBCOLLECTION_ID=$$(grep -E '^FIRESTORE_SUBCOLLECTION_ID=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'ai-india-workshop-2024'); \
		export ADMIN_PASSWORD=$$(grep -E '^ADMIN_PASSWORD=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'change-this-password'); \
		export SESSION_SECRET=$$(grep -E '^SESSION_SECRET=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'change-this-secret-min-32-chars'); \
		export GCP_PROJECT_ID=$$(grep -E '^(GCP_PROJECT_ID|GOOGLE_CLOUD_PROJECT|GCLOUD_PROJECT)=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' | head -1); \
		if [ -n "$$GCP_PROJECT_ID" ]; then \
			docker run -d --name ai-india-workshop \
				-p 8080:8080 \
				-e FIRESTORE_SUBCOLLECTION_ID="$$FIRESTORE_SUBCOLLECTION_ID" \
				-e ADMIN_PASSWORD="$$ADMIN_PASSWORD" \
				-e SESSION_SECRET="$$SESSION_SECRET" \
				-e GCP_PROJECT_ID="$$GCP_PROJECT_ID" \
				-e PORT=8080 \
				-e STATIC_DIR=/app/static \
				ai-india-workshop:latest; \
		else \
			docker run -d --name ai-india-workshop \
				-p 8080:8080 \
				-e FIRESTORE_SUBCOLLECTION_ID="$$FIRESTORE_SUBCOLLECTION_ID" \
				-e ADMIN_PASSWORD="$$ADMIN_PASSWORD" \
				-e SESSION_SECRET="$$SESSION_SECRET" \
				-e PORT=8080 \
				-e STATIC_DIR=/app/static \
				ai-india-workshop:latest; \
		fi \
	else \
		docker run -d --name ai-india-workshop \
			-p 8080:8080 \
			-e FIRESTORE_SUBCOLLECTION_ID=ai-india-workshop-2024 \
			-e ADMIN_PASSWORD=change-this-password \
			-e SESSION_SECRET=change-this-secret-min-32-chars \
			-e PORT=8080 \
			-e STATIC_DIR=/app/static \
			ai-india-workshop:latest; \
	fi
	@echo "Container started. Access at http://localhost:8080"
	@echo "View logs with: docker logs -f ai-india-workshop"

docker-run-local: docker-build
	@echo "Running Docker container with service account file (local dev)..."
	@if [ ! -f firebase-service-account.json ]; then \
		echo "Error: firebase-service-account.json not found!"; \
		echo "Please place your Firebase service account JSON file in the project root."; \
		exit 1; \
	fi
	@docker stop ai-india-workshop 2>/dev/null || true
	@docker rm ai-india-workshop 2>/dev/null || true
	@if [ -f .env ]; then \
		export FIRESTORE_SUBCOLLECTION_ID=$$(grep -E '^FIRESTORE_SUBCOLLECTION_ID=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'ai-india-workshop-2024'); \
		export ADMIN_PASSWORD=$$(grep -E '^ADMIN_PASSWORD=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'change-this-password'); \
		export SESSION_SECRET=$$(grep -E '^SESSION_SECRET=' .env | cut -d '=' -f2- | tr -d '"'"'"'"' || echo 'change-this-secret-min-32-chars'); \
		docker run -d --name ai-india-workshop \
			-p 8080:8080 \
			-e FIREBASE_SERVICE_ACCOUNT_PATH=/app/firebase-service-account.json \
			-e FIRESTORE_SUBCOLLECTION_ID="$$FIRESTORE_SUBCOLLECTION_ID" \
			-e ADMIN_PASSWORD="$$ADMIN_PASSWORD" \
			-e SESSION_SECRET="$$SESSION_SECRET" \
			-e PORT=8080 \
			-e STATIC_DIR=/app/static \
			-v $$(pwd)/firebase-service-account.json:/app/firebase-service-account.json:ro \
			ai-india-workshop:latest; \
	else \
		docker run -d --name ai-india-workshop \
			-p 8080:8080 \
			-e FIREBASE_SERVICE_ACCOUNT_PATH=/app/firebase-service-account.json \
			-e FIRESTORE_SUBCOLLECTION_ID=ai-india-workshop-2024 \
			-e ADMIN_PASSWORD=change-this-password \
			-e SESSION_SECRET=change-this-secret-min-32-chars \
			-e PORT=8080 \
			-e STATIC_DIR=/app/static \
			-v $$(pwd)/firebase-service-account.json:/app/firebase-service-account.json:ro \
			ai-india-workshop:latest; \
	fi
	@echo "Container started. Access at http://localhost:8080"
	@echo "View logs with: docker logs -f ai-india-workshop"

docker-stop:
	@echo "Stopping Docker container..."
	@docker stop ai-india-workshop 2>/dev/null || true
	@docker rm ai-india-workshop 2>/dev/null || true
	@echo "Container stopped and removed"

docker-up:
	@echo "Starting Docker Compose services..."
	docker-compose up -d --build

docker-down:
	@echo "Stopping Docker Compose services..."
	docker-compose down

docker-logs:
	@echo "Viewing Docker Compose logs..."
	docker-compose logs -f

docker-clean:
	@echo "Cleaning Docker resources..."
	docker-compose down -v
	docker rmi ai-india-workshop:latest 2>/dev/null || true

# Clean targets
clean:
	@echo "Cleaning build artifacts..."
	rm -f backend/server
	rm -rf frontend/dist
	rm -rf frontend/node_modules/.vite

