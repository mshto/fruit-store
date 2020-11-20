# Variables
MIGRATIONS_DIR=sql/migrations
# MIGRATION_PATH="postgres://postgres:secret@localhost:5432/fruit_store?sslmode=disable"
MIGRATION_PATH="postgres://tlizipypsyzdzk:722e33386708aa46f59834a0c627cfbbc5c28b30ec98f24c27cb13216c7b73aa@ec2-34-237-236-32.compute-1.amazonaws.com:5432/d56t4fqipu28od"

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor

# Test
.PHONY: test
test:
	go test `go list ./... | grep -v mock` -cover -coverprofile cover-all.out

.PHONY: cover-html
cover-html:
	go tool cover -html=cover-all.out

.PHONY: build
build:
	# go build -o bin/go-getting-started -v .
	GOOS=linux GOARCH=amd64 go build -o fruit-store -v .

# Migration github repo: https://github.com/golang-migrate/migrate/tree/master/cmd/migrate
migrate-new:
	@migrate create -dir $(MIGRATIONS_DIR) -ext sql -format "unix" new_migration

migrate-up:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) up 

migrate-down:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) down
migrate-force:
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) force 1605191158 
