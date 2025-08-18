#!/usr/bin/env bash

# generate swagger documentation
echo "Generate swagger documentation..."
swag fmt
swag init -ot go,yaml -g cmd/main.go -o api/swagger
