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
	docker compose --env-file yours.env up

down:
	docker compose --env-file yours.env down

fulldown:
	docker compose --env-file yours.env down
	docker volume rm gostream_postgres_data
	docker volume rm gostream_minio_data

prune:
	docker volume rm gostream_postgres_data
	docker volume rm gostream_minio_data

unbuild:
	docker rmi gostream-postgres-image
