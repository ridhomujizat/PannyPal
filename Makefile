# Makefile

include .env

DB_URL=postgres://$(DB_USER):$(DB_PASS)@$(DB_HOST):$(DB_PORT)/$(DB_NAME)?sslmode=disable

# Database migration commands using golang-migrate/migrate tool
migrate-down:
	@echo "Running migrate DOWN..."
	migrate -path migrations -database "$(DB_URL)" down

migrate-up:
	@echo "Running migrate UP..."
	migrate -path migrations -database "$(DB_URL)" up

migrate-down-1:
	@echo "Running migrate DOWN 1 step..."
	migrate -path migrations -database "$(DB_URL)" down 1

migrate-up-1:
	@echo "Running migrate UP 1 step..."
	migrate -path migrations -database "$(DB_URL)" up 1

migrate-force:
	@echo "Force migrate to version $(version)..."
	migrate -path migrations -database "$(DB_URL)" force $(version)

migrate-goto:
	@echo "Migrate to version $(version)..."
	migrate -path migrations -database "$(DB_URL)" goto $(version)

migrate-drop:
	@echo "Dropping all migrations..."
	migrate -path migrations -database "$(DB_URL)" drop

migrate-version:
	@echo "Current migration version:"
	migrate -path migrations -database "$(DB_URL)" version

migrate:
	@echo "Creating new migration..."
	migrate create -ext sql -dir migrations -format timestamp $(name)

install-migrate:
	@echo "Installing golang-migrate..."
	go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

# GORM AutoMigrate approach (your existing system)
migrate-gorm:
	@echo "Running GORM AutoMigrate..."
	go run cmd/migrate/main.go

# SQL Migration approach (new system)
migrate-sql-up:
	@echo "Running SQL migrations UP..."
	go run cmd/migrate-sql/main.go -command=up

migrate-sql-down:
	@echo "Running SQL migrations DOWN..."
	go run cmd/migrate-sql/main.go -command=down

migrate-sql-version:
	@echo "Getting SQL migration version..."
	go run cmd/migrate-sql/main.go -command=version

run-app:
	@echo "Running application..."
	go run cmd/api/main.go

build-dashboard:
	@echo "Building dashboard..."
	cd dashboard && npm install && npm run build
	@echo "Dashboard built to dashboard/dist"

run-full: build-dashboard
	@echo "Running full stack (backend + frontend)..."
	go run cmd/api/main.go

build-all: build-dashboard
	@echo "Building Go application..."
	go build -o bin/pannypal cmd/api/main.go
	@echo "Build complete: bin/pannypal"
