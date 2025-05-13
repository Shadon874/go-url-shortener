# Go URL Shortener

A simple URL shortener service built in Go.

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/Shadon874/go-url-shortener.git
   cd url-shortener
   ```

2. Install Go dependencies:
   ```bash
   go mod tidy
   ```

3. Run the application:
   ```bash
   go run cmd/server/main.go
   ```

4. Visit `http://localhost:8080` to access the URL shortener.

## API

### POST `/api/shorten`
- **Request body**: JSON containing a URL to shorten:
  ```json
  {
    "url": "https://example.com"
  }
  ```

- **Response**: JSON containing the shortened URL:
  ```json
  {
    "shortened_url": "http://localhost:8080/abc123"
  }
  ```

## License
MIT License
