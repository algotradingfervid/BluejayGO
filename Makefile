.PHONY: help run build dev migrate-up migrate-down migrate-create sqlc seed test clean deploy deploy-build deploy-upload deploy-restart

help:
	@echo "BlueJay CMS - Available commands:"
	@echo "  make run           - Run the server"
	@echo "  make build         - Build the server binary"
	@echo "  make dev           - Run with hot-reload (air)"
	@echo "  make migrate-up    - Run all migrations"
	@echo "  make migrate-down  - Rollback all migrations"
	@echo "  make sqlc          - Generate sqlc code"
	@echo "  make seed          - Seed database with sample data"
	@echo "  make test          - Run tests"
	@echo "  make clean         - Clean build artifacts"

run:
	go run cmd/server/main.go

build:
	go build -o bin/bluejay-cms cmd/server/main.go

dev:
	air

migrate-up:
	migrate -path db/migrations -database "sqlite://bluejay.db" up

migrate-down:
	migrate -path db/migrations -database "sqlite://bluejay.db" down

sqlc:
	sqlc generate

seed:
	sqlite3 bluejay.db < seed.sql

test:
	go test -v ./...

clean:
	rm -rf bin/
	rm -f bluejay.db
	rm -f bluejay.db-shm
	rm -f bluejay.db-wal

deploy: deploy-build deploy-upload deploy-restart

deploy-build:
	GOOS=linux GOARCH=amd64 go build -o bluejay-cms cmd/server/main.go

deploy-upload:
	scp bluejay-cms user@yourserver:/var/www/bluejay-cms/
	scp deploy/Caddyfile user@yourserver:/etc/caddy/Caddyfile
	scp deploy/bluejay-cms.service user@yourserver:/etc/systemd/system/

deploy-restart:
	ssh user@yourserver 'sudo systemctl daemon-reload && sudo systemctl restart bluejay-cms && sudo systemctl reload caddy'
