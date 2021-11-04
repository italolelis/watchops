NO_COLOR=\033[0m
OK_COLOR=\033[32;01m
ERROR_COLOR=\033[31;01m
WARN_COLOR=\033[33;01m
SERVICE_NAME=fourkeys

.PHONY: all test build
all: test build

setup: migrate kinesis

build: private-repo
	@echo "$(OK_COLOR)==> Building $(SERVICE_NAME)... $(NO_COLOR)"
	@CGO_ENABLED=0 go build -ldflags "-s -w" -ldflags "-X main.version=${VERSION}" -o "dist/$(SERVICE_NAME)" github.com/italolelis/$(SERVICE_NAME)/cmd/$(SERVICE_NAME)

test: private-repo
	@echo "$(OK_COLOR)==> Running tests$(NO_COLOR)"
	@go test -v -failfast -cover ./...

migrate: tools.migrate
	@./bin/migrate.linux-amd64 -path="configs/migrations/postgres" -database="postgres://fourkeys:qwerty123@localhost:5432/fourkeys-db?sslmode=disable" up

compose:
	@echo "$(OK_COLOR)==> Bringing containers up for $(SERVICE_NAME)... $(NO_COLOR)"
	@docker-compose up -d

kinesis:
	@echo "$(OK_COLOR)==> Creating Kinesis streams... $(NO_COLOR)"
	@aws --endpoint-url http://127.0.0.1:4566 kinesis create-stream --stream-name fourkeys_github --shard-count 1
	@aws --endpoint-url http://127.0.0.1:4566 kinesis create-stream --stream-name fourkeys_opsgenie --shard-count 1

#---------------
#-- tools
#---------------

.PHONY: tools tools.migrate
tools: tools.migrate

tools.migrate:
	@command -v ./bin/migrate.linux-amd64 >/dev/null ; if [ $$? -ne 0 ]; then \
		echo "--> installing migrate"; \
		curl -L https://github.com/golang-migrate/migrate/releases/download/v4.4.0/migrate.linux-amd64.tar.gz | tar xvz; \
		mkdir -p bin; \
		mv migrate.linux-amd64 ./bin; \
	fi
