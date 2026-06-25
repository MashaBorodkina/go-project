include .env
export


MIGRATIONS_PATH = migrations

fmt:
	$(HOME)/go/bin/golangci-lint fmt

lint:
	$(HOME)/go/bin/golangci-lint run --config .golangci.yml

run:
	go run ./cmd/server

migrate-up-local:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up

migrate-up:
	$(HOME)/go/bin/migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" up	

migrate-down-local:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-down:
	$(HOME)/go/bin/migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" down 1

migrate-force-local:
	migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)

migrate-force:
	$(HOME)/go/bin/migrate -path $(MIGRATIONS_PATH) -database "$(DATABASE_URL)" force $(VERSION)

	