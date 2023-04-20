run:
	go run cmd/main.go

up-db:
	docker run --name postgres-onela -e POSTGRES_USER=onelab -e POSTGRES_PASSWORD=qwerty -e POSTGRES_DB=onelab_db -d -p5433:5432 --rm postgres

stop-db:
	docker stop postgres-onela

migration-up:
	migrate -path ./storage/postgres/migrations/ -database 'postgres://onelab:qwerty@localhost:5433/onelab_db?sslmode=disable' up

migration-down:
	migrate -path ./storage/postgres/migrations/ -database 'postgres://onelab:qwerty@localhost:5432/onelab_db?sslmode=disable' down
