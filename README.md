# Brewferring - Your Coffee is Brewferring

A web application for ordering coffee through the Terminal.shop API.

## Features

- Browse coffee products
- User authentication via Terminal.shop OAuth
- View and manage orders
- User profile management
- Modern UI with Tailwind CSS and DaisyUI

## Prerequisites

- Go 1.21 or later
- Terminal.shop API credentials

## Setup

1. Clone the repository:
```bash
git clone https://github.com/yourusername/brewferring.git
cd brewferring
```

2. Install dependencies:
```bash
go mod download
```

3. Set up configuration:
```bash
# Copy the template configuration file
cp config.template.yml config.yml

# Edit config.yml with your credentials
nano config.yml
```

4. Run the application:
```bash
go run cmd/server/main.go
```

The application will be available at http://localhost:8080 (or the configured host and port)

## Configuration

The application uses a YAML configuration file (`config.yml`) for its settings. A template file (`config.template.yml`) is provided as an example.

Required configuration:
- `oauth.client_id`: Your Terminal.shop API client ID
- `oauth.client_secret`: Your Terminal.shop API client secret
- `oauth.redirect_uri`: The OAuth callback URL (default: http://localhost:8080/callback)

Optional configuration:
- `server.port`: The port to listen on (default: 8080)
- `server.host`: The host to listen on (default: localhost)

## Environment Variables

- `TERMINAL_CLIENT_ID`: Your Terminal.shop API client ID
- `TERMINAL_CLIENT_SECRET`: Your Terminal.shop API client secret

## Development

The application uses the following technologies:

- Go for the backend
- Templ for HTML templates
- Tailwind CSS and DaisyUI for styling
- HTMX for dynamic interactions
- Terminal.shop SDK for API integration

## Project Structure

```
brewferring/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── config/
│   │   └── config.go
│   ├── handlers/
│   │   └── handlers.go
│   ├── middleware/
│   │   └── auth.go
│   └── templates/
│       ├── home.templ
│       ├── dashboard.templ
│       ├── products.templ
│       ├── profile.templ
│       └── orders.templ
├── static/
│   └── css/
├── config.template.yml
└── config.yml
```

## License

MIT License 