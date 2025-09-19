#!/usr/bin/env bash

echo -e "Open coverage report in browser..."
go tool cover -html=coverage.out
echo "Done."
