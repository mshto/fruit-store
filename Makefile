# Variables
MIGRATIONS_DIR=sql/migrations
MIGRATION_PATH="postgres://postgres:secret@localhost:5432/fruit_store?sslmode=disable"
# MIGRATION_PATH="postgres://tlizipypsyzdzk:722e33386708aa46f59834a0c627cfbbc5c28b30ec98f24c27cb13216c7b73aa@ec2-34-237-236-32.compute-1.amazonaws.com:5432/d56t4fqipu28od"

.PHONY: dependencies
dependencies:
	echo "Installing dependencies"
	go mod vendor

# Test
.PHONY: test
test:
	go test ./... -cover -coverprofile cover-all.out

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
	@migrate -path $(MIGRATIONS_DIR) -database $(MIGRATION_PATH) force 1604439163 

# GO_BUILD_ENV := CGO_ENABLED=0 GOOS=linux GOARCH=amd64
# DOCKER_BUILD=$(shell pwd)/.docker_build
# DOCKER_CMD=$(DOCKER_BUILD)/go-getting-started

# $(DOCKER_CMD): clean
# 	mkdir -p $(DOCKER_BUILD)
# 	$(GO_BUILD_ENV) go build -v -o $(DOCKER_CMD) .

# clean:
# 	rm -rf $(DOCKER_BUILD)

# heroku: $(DOCKER_CMD)
# 	heroku container:push web