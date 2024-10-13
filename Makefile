#!/bin/bash

run:
	@nodemon -e go --signal SIGTERM --exec 'go run . || exit 1'

build:
	@echo "building..."
	@go build -o bin/historia-app .

pull:
	@git pull

update: pull build service-restart

service-start:
	@echo "service starting..."
	@systemctl start historia.service
	
service-enable:
	@echo "service enabling..."
	@systemctl enable historia.service

service-restart:
	@echo "service restarting..."
	@systemctl restart historia.service
