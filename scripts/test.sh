#!/usr/bin/env bash

echo "Run unit tests with coverage..."
go test -v ./internal/usecase/... -coverprofile=coverage.out -covermode=atomic
