default: build

init:
	@cp .env.template .env

build:
	@docker-compose build

devshell:
	@docker-compose run --rm --service-ports service sh

up:
	@docker-compose run --rm --service-ports service

down:
	@docker-compose down --remove-orphans

t:
	@go test -cover ./...