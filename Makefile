run:
	go run cmd/main.go

start-db:
	docker run --name postgres-onelab -e POSTGRES_USER=onelab -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=onelab_db -d -p5432:5432 --rm postgres

stop-db:
	docker stop postgres-onelab

migration-up:
	migrate -path ./internal/storage/postgres/migrations/ -database 'postgres://onelab:qwerty@localhost:5432/onelab_db?sslmode=disable' up

migration-down:
	migrate -path ./internal/storage/postgres/migrations/ -database 'postgres://onelab:qwerty@localhost:5432/onelab_db?sslmode=disable' down

.PHONY: run