#!/bin/bash

# Generate Go code from templ files
templ generate ./internal/templates
# generate tailwind css
tailwindcss -i ./static/css/input.css -o ./static/css/output.css
