#!/bin/sh
set -eux
gofmt -s -w .
reset
#go build -o zoomer -race cmd/zoomer/main.go
go build -o zoomer cmd/zoomer/main.go