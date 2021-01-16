#!/usr/bin/env bash
set -e

echo "#################### migrating local db"
go run cmd/migrator/main.go

echo "#################### downloading CompileDaemon"
GO111MODULE=off go get github.com/githubnemo/CompileDaemon

echo "#################### starting deamon"
CompileDaemon --build="go build -o main cmd/portdomainservice/main.go" --command=./main
