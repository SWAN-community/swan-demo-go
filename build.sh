#!/usr/bin/env bash
# Stops the process if something fails
set -xe

# create the application binary that eb uses
GOOS=linux GOARCH=amd64 go build  -ldflags="-s -w" -o ./application ./src/server.go 

zip -r aws-eb-swan-demo.zip application appsettings.json Procfile www

# build the docker container image.
docker build -t swan-community/swan-demo-go .