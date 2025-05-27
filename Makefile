run:
	templ generate
	go fmt ./...
	go vet ./...
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css
	go run main.go

templ:
	templ generate

tw:
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css

build:
	templ generate
	go fmt ./...
	go vet ./...
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css
	go build -o bin/brewferring

build-linux:
	templ generate
	go fmt ./...
	go vet ./...
	npx @tailwindcss/cli -i ./static/css/input.css -o ./static/css/output.css
	GOOS=linux GOARCH=amd64 go build -o bin/brewferring-linux-amd64
