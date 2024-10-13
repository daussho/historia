#!/bin/bash

run:
	@nodemon -e go --signal SIGTERM --exec 'go run . || exit 1'

build:
	@go build -o bin/historia-app .

update: build service-restart

service-start:
	@systemctl start historia.service
	
service-enable:
	@systemctl enable historia.service

service-restart:
	@systemctl restart historia.service
