# go-toggl
#### #### #### #### #### #### ####
DIR := ${CURDIR}
DIR_MIGRATIONS := "${DIR}/internal/data/db/migrations"
DOCKER_NETWORK=host
DOCKER_PROJECT=go-toggl
#### #### #### #### #### #### ####

# Load .env in order to use vars
ifneq (,$(wildcard ./.env))
    include .env
    export
endif

all: setup_project migrate_up

### ### ### ### ### ### ### ### ### ###
# D A T A B A S E  M I G R A T I O N S
### ### ### ### ### ### ### ### ### ###
# source: https://github.com/golang-migrate/migrate/blob/master/database/postgres/TUTORIAL.md
migrate_init:
	docker run -v ${DIR_MIGRATIONS}:/migrations --network ${DOCKER_NETWORK} migrate/migrate -path=/migrations/ \
			create -ext sql -dir migrations/ -seq init_schema

migrate_create:
	docker run -v ${DIR_MIGRATIONS}:/migrations --network ${DOCKER_NETWORK} migrate/migrate -path=/migrations/ \
			-verbose create -ext sql -dir migrations/ -seq create_${NAME}_table

migrate_up:
	docker run -v ${DIR_MIGRATIONS}:/migrations --network ${DOCKER_NETWORK} migrate/migrate -path=/migrations/ \
 			-database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable -verbose up

migrate_down:
	docker run -v ${DIR_MIGRATIONS}:/migrations --network ${DOCKER_NETWORK} migrate/migrate -path=/migrations/ \
			-database postgres://${DB_USER}:${DB_PASSWORD}@${DB_HOST}:${DB_PORT}/${DB_NAME}?sslmode=disable -verbose down

setup_sqlite:
		sqlite3 ${DATABASE_NAME} < ./internal/data/init.sql

### ### ### ### ### ### ### ### ###
### ### ### D O C K E R ### ### ###
### ### ### ### ### ### ### ### ###
setup_project:
	docker-compose -f ./infra/docker/docker-compose.yml up

docker_build_tag_push:
	docker build -t jrmanes/${DOCKER_PROJECT}:latest .
	docker push jrmanes/$(DOCKER_PROJECT):latest