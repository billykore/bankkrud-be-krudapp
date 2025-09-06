#!/usr/bin/env bash

echo "Run unit tests with coverage..."
go test -v ./internal/usecase/... -coverprofile=coverage.out -covermode=atomic

echo -e "\nOpen coverage report in browser..."
go tool cover -html=coverage.out
echo "Done."
