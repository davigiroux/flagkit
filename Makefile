.PHONY: api dashboard docker-up docker-down test migrate-up migrate-down

docker-up:
	docker compose up -d

docker-down:
	docker compose down

api:
	cd api && go run ./cmd/server

test:
	cd api && go test ./...

migrate-up:
	cd api && go run ./cmd/server -migrate-only

dashboard:
	pnpm --filter dashboard dev
