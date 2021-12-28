#!/usr/bin/env bash

go mod tidy
[[ ! -d "out" ]] && mkdir ./out
CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./out/wss-exporter-linux main.go
CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -o ./out/wss-exporter-darwin main.go
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -o ./out/wss-exporter-windows main.go
cp ./out/wss-exporter-linux ./out/wss-exporter
