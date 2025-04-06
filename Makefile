# Run the app directly (non-containerized)
run:
	go run cmd/api/main.go cmd/api/api.go

# Run migrations locally using golang-migrate CLI
migrate:
	migrate -path migrate/migrations -database \
		"postgres://$(DB_DSN)" up

# Build the app Docker image
build:
	docker build -t gochujang -f Dockerfile .

# Build the migration image
build-migrate:
	docker build -t gochujang-migrate -f Dockerfile.migrate .

# Run local dev environment with PostgreSQL and port 8080 exposed
dev:
	docker-compose --profile dev up --build

# Run production-like environment (no dev db)
prod:
	docker-compose --profile prod up --build

# Tear down containers and volumes
down:
	docker-compose down -v
