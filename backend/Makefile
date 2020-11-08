# Variables
MIGRATIONS_DIR=sql/migrations
MIGRATION_PATH="postgres://postgres:secret@localhost:5432/fruit_store?sslmode=disable"

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor

# Test
.PHONY: test
test:
	go test ./... -cover -coverprofile cover-all.out

# Migration github repo: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate

migrate-new:
	@migrate create -dir $(MIGRATIONS_DIR) -ext sql -format "unix" new_migration

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) up 

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) down
migrate-force:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) force 1604439163 
