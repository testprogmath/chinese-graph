# Chinese Learning Graph - Backend Makefile

.PHONY: help install start stop restart seed test lint clean dev prod logs check

help: ## Show this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-15s\033[0m %s\n", $$1, $$2}'

install: ## Install dependencies
	@echo "📦 Installing Go dependencies..."
	go mod download
	go mod tidy
	@echo "✅ Dependencies installed"

start: ## Start all services (Neo4j, Redis, Backend)
	@echo "🚀 Starting infrastructure..."
	docker-compose up -d
	@echo "⏳ Waiting for Neo4j to be ready..."
	@sleep 10
	@echo "🔥 Starting backend server..."
	go run cmd/server/main.go

stop: ## Stop all services
	@echo "🛑 Stopping services..."
	docker-compose down
	@echo "✅ Services stopped"

restart: stop start ## Restart all services

seed: ## Load sample data
	@echo "🌱 Loading sample data..."
	go run scripts/seed_data.go
	@echo "✅ Sample data loaded"

test: ## Run tests
	@echo "🧪 Running tests..."
	go test -v -race ./...

lint: ## Run linters
	@echo "🔍 Running linters..."
	@if ! command -v golangci-lint > /dev/null; then \
		echo "Installing golangci-lint..."; \
		curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin; \
	fi
	golangci-lint run ./...

lint-fix: ## Fix linting issues
	@echo "🔧 Fixing linting issues..."
	golangci-lint run --fix ./...
	go fmt ./...
	go mod tidy

fmt: ## Format code
	@echo "✨ Formatting code..."
	go fmt ./...
	gofmt -s -w .

coverage: ## Run tests with coverage
	@echo "📊 Running tests with coverage..."
	go test -v -race -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html
	@echo "📊 Coverage report: coverage.html"

clean: ## Clean up everything
	@echo "🧹 Cleaning up..."
	docker-compose down -v
	rm -rf coverage.* server
	@echo "✅ Cleaned"

dev: ## Run in development mode with hot reload
	@echo "🔥 Starting in dev mode..."
	@if ! command -v air > /dev/null; then \
		echo "Installing air for hot reload..."; \
		go install github.com/air-verse/air@latest; \
	fi
	air

prod: ## Build for production
	@echo "📦 Building for production..."
	CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o server cmd/server/main.go
	@echo "✅ Built: ./server"

logs: ## Show Docker logs
	docker-compose logs -f

neo4j-console: ## Open Neo4j console
	@echo "🔗 Opening Neo4j console..."
	@echo "URL: http://localhost:7474"
	@echo "Username: neo4j"
	@echo "Password: chinesegraph123"
	@open http://localhost:7474 2>/dev/null || xdg-open http://localhost:7474 2>/dev/null || echo "Please open http://localhost:7474"

graphql: ## Open GraphQL playground
	@echo "🎮 Opening GraphQL playground..."
	@open http://localhost:8080 2>/dev/null || xdg-open http://localhost:8080 2>/dev/null || echo "Please open http://localhost:8080"

check: ## Check if services are running
	@echo "🔍 Checking services..."
	@docker-compose ps
	@curl -s http://localhost:8080/query > /dev/null && echo "✅ Backend is running" || echo "❌ Backend is not running"
	@docker-compose exec -T neo4j cypher-shell -u neo4j -p chinesegraph123 "RETURN 1" > /dev/null 2>&1 && echo "✅ Neo4j is running" || echo "❌ Neo4j is not running"
	@docker-compose exec -T redis redis-cli ping > /dev/null 2>&1 && echo "✅ Redis is running" || echo "❌ Redis is not running"

generate: ## Generate GraphQL code
	@echo "🔧 Generating GraphQL code..."
	go run github.com/99designs/gqlgen generate
	@echo "✅ Generated"

migrate: ## Run database migrations
	@echo "🗄️ Running migrations..."
	go run scripts/migrate.go
	@echo "✅ Migrations complete"

# Development shortcuts
s: start
st: stop
r: restart
t: test
l: logs
c: clean