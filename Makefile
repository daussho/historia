#!/bin/bash

run:
	@nodemon -e go --signal SIGTERM --exec 'go run . || exit 1'

build:
	@go build -o bin/app .