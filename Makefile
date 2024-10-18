include yours.env

shell = /bin/bash

user = $(shell whoami)
gid = $(shell id -g $(user))

run:
	mkdir -p $(APP_LOGS_PATH)
	go run cmd/main.go

build:
	docker build -f "docker/Dockerfile.postgres" -t "gostream-postgres-image" --build-arg GID=$(gid) .

up:
	mkdir -p $(POSTGRES_LOGS_DIR)
	mkdir -p $(MINIO_LOGS_DIR)
	chmod 777 $(POSTGRES_LOGS_DIR)
	docker compose --env-file yours.env up -d

down:
	docker compose --env-file yours.env down

unbuild:
	docker rmi gostream-postgres-image
