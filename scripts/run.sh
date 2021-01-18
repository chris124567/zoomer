#!/bin/sh
set -eux
gofmt -s -w .
reset
go run cmd/zoomer/main.go $@