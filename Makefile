# Variables
MIGRATIONS_DIR=sql/migrations
MIGRATION_PATH="postgres://postgres:secret@localhost:5432/fruit_store?sslmode=disable"

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor

.PHONY: test
test:
	go test `go list ./... | grep -v mock` -cover -coverprofile cover-all.out

.PHONY: cover-html
cover-html:
	go tool cover -html=cover-all.out

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: build
build:
	GOOS=linux GOARCH=amd64 go build -o fruit-store -v .

# Migration github repo: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
migrate-new:
	@migrate create -dir $(MIGRATIONS_DIR) -ext sql -format "unix" new_migration

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) up 

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) down

migrate-force:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) force n 
