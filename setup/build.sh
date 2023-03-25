#!/bin/bash
CGO_ENABLED=0 go build -ldflags="-s -w -extldflags '-static'" -o setup .
CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-s -w -extldflags '-static'" -o setup.exe .
