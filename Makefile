#!/bin/bash
export GOOSE_MIGRATION_DIR=script

run:
	@nodemon --verbose --delay 1 -e go,html,mod --signal SIGTERM --exec 'go run . || exit 1'

build:
	@echo "run go mod tidy"
	@go mod tidy
	@echo "building..."
	@go build -o bin/historia-app .

pull:
	@git pull

update: pull build service-restart

service-start:
	@echo "service starting..."
	@pm2 start historia-app

service-restart:
	@echo "service restarting..."
	@pm2 restart historia-app

migration-create:
	@echo "migration creating..."
	@goose create $(NAME) sql

migration-up:
	@echo "migration up..."
	@GOOSE_DRIVER=mysql goose up