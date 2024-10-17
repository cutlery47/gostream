include yours.env

shell = /bin/bash

user = $(shell whoami)
gid = $(shell id -g $(user))

postgres_logs_dir = 

build:
	docker build -f "docker/Dockerfile.postgres" -t "gostream-postgres-image" --build-arg GID=$(gid) .

up:
	mkdir -p $(POSTGRES_LOGS_DIR)
	chmod 777 $(POSTGRES_LOGS_DIR)
	docker compose --env-file yours.env up

down:
	docker compose --env-file yours.env down

unbuild:
	docker rmi gostream-postgres-image
