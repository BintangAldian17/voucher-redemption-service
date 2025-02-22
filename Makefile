build:
	@go build -o bin/otto_digital_be cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/otto_digital_be 

migration:
	@migrate create -ext sql -dir cmd/migrate/migrations $(filter-out $@,$(MAKECMDGOALS))

migrate-up:
	@go run cmd/migrate/main.go up

migrate-down:
	@go run cmd/migrate/main.go down

migrate-force:
	@go run cmd/migrate/main.go force  $(filter-out $@,$(MAKECMDGOALS))

seed:
	@go run cmd/seeder/main.go