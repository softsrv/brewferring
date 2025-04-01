#!/bin/bash

# Generate templates
templ generate ./internal/templates/...

# Format code
go fmt ./...

# Run linter
go vet ./...

# Generate Go code from templ files
templ generate ./internal/templates
# generate tailwind css
tailwindcss -i ./static/css/input.css -o ./static/css/output.css
