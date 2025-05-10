# URL Shortener Implementation

This is a simple URL shortener service built with Go and Gin framework.

## Features

- REST API to shorten URLs
- Returns the same short URL for identical long URLs
- Tracks domain metrics
- Provides redirection from short URL to original URL
- Provides a metrics API to show top 3 most shortened domains

## Running the Application

### Using Go

```bash
# Install dependencies
go mod download

# Run the application
go run main.go
```

The server will start on port 8080.

### Using Docker

```bash
# Build the Docker image
docker build -t url-shortener .

# Run the container
docker run -p 8080:8080 url-shortener
```

## API Endpoints

1. **Shorten URL**
   - `POST /api/shorten`
   - Request Body: `{"url": "https://example.com/long/url"}`
   - Response: `{"short_url": "http://localhost:8080/r/AbCdEfGh"}`

2. **Redirect**
   - `GET /r/:shortURL`
   - Redirects to the original URL

3. **Top Domains Metrics**
   - `GET /api/metrics/top-domains`
   - Returns top 3 domains being shortened
   - Response: `{"domains": [{"domain": "example", "count": 5}, ...]}`

## Testing

```bash
go test ./...
```

## Architecture

- The application uses an in-memory storage for URL mappings
- The Gin framework handles HTTP routing
- A service layer contains the business logic
- Handler layer manages HTTP requests and responses
