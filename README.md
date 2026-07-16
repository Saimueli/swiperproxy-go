# SwiperProxy Go

A modern, secure, and high-performance HTTP proxy server with intelligent caching, rate limiting, and comprehensive security features. Built with Go 1.21+.

## ✨ Features

- **HTTP/HTTPS Proxy** – Full HTTP/1.1 and HTTP/2 support
    
- **Intelligent Caching** – In-memory LRU cache with configurable TTL
    
- **Rate Limiting** – Protect against DoS attacks with configurable limits
    
- **Security Headers** – Comprehensive security headers (HSTS, XSS, CSP, etc.)
    
- **Structured Logging** – JSON format logs for easy integration
    
- **Docker Ready** – Containerized with Docker/Podman support
    
- **Zero Dependencies** – Single binary, no external runtime required
    

## 🚀 Quick Start

### Using Docker

Bash

```
# Pull and run
docker run -d -p 8080:8080 --name swiperproxy yourusername/swiperproxy-go

# Test it
curl http://localhost:8080/get
```

### Using Binary

Bash

```
# Download the latest release
wget https://github.com/yourusername/swiperproxy-go/releases/latest/swiperproxy-linux-amd64
chmod +x swiperproxy-linux-amd64

# Run with config
./swiperproxy-linux-amd64 -config configs/config.yaml
```

### From Source

Bash

```
# Clone the repository
git clone https://github.com/yourusername/swiperproxy-go.git
cd swiperproxy-go

# Build
go build -o swiperproxy ./cmd/swiperproxy

# Run
./swiperproxy -config configs/config.yaml
```

## ⚙️ Configuration

Edit `configs/config.yaml` to customize:

YAML

```
server:
  host: "0.0.0.0"
  port: "8080"
  read_timeout: 10s
  write_timeout: 10s
  idle_timeout: 120s

proxy:
  target: "https://httpbin.org"  # Your target server
  timeout: 30s
  max_redirects: 5

cache:
  enabled: true
  type: "memory"  # memory or redis
  max_size: 100   # MB
  ttl: 300        # seconds

rate_limit:
  enabled: false
  requests_per_second: 100
  burst: 50

security:
  enable_headers: true
  enable_cors: false
  allowed_origins: ["*"]
```

## 📊 Usage Examples

### Basic Proxy

Bash

```
# Forward requests to target server
curl http://localhost:8080/get
curl http://localhost:8080/post -X POST -d "data=test"
```

### Cache Testing

Bash

```
# First request - cache MISS
curl -v http://localhost:8080/get | grep "X-Cache"

# Second request - cache HIT (faster response)
curl -v http://localhost:8080/get | grep "X-Cache"
```

### Rate Limiting

Bash

```
# Enable rate limiting in config, then test
for i in {1..150}; do
  curl -s -o /dev/null -w "%{http_code}\n" http://localhost:8080/get
done
```

## 🐳 Docker Build

### Using Docker

Bash

```
# Build image
docker build -t swiperproxy-go .

# Run container
docker run -d -p 8080:8080 --name swiperproxy swiperproxy-go
```

### Using Podman

Bash

```
# Build image
podman build -f Containerfile -t swiperproxy-go .

# Run container
podman run -d -p 8080:8080 --name swiperproxy swiperproxy-go
```

### Docker Compose

Bash

```
docker-compose up -d
```

## 🧪 Testing

Bash

```
# Run all tests
go test ./...

# Run with coverage
go test -cover ./...

# Run specific tests
go test ./internal/proxy/...
```

## 📁 Project Structure

Plaintext

```
swiperproxy-go/
├── cmd/
│   └── swiperproxy/
│       └── main.go          # Application entry point
├── internal/
│   ├── config/
│   │   └── config.go        # Configuration management
│   ├── proxy/
│   │   ├── handler.go       # Proxy handler
│   │   └── cache.go         # Cache implementation
│   └── middleware/
│       ├── logger.go        # Logging middleware
│       ├── ratelimit.go     # Rate limiting
│       └── security.go      # Security headers
├── configs/
│   └── config.yaml          # Default configuration
├── Containerfile            # Docker/Podman build
├── docker-compose.yml       # Docker Compose setup
├── go.mod                   # Go module
└── README.md                # This file
```

## 🔒 Security Features

> [!info] Security Highlights
> 
> - **Security Headers:** `X-Content-Type-Options`, `X-Frame-Options`, `X-XSS-Protection`, `Permissions-Policy`
>     
> - **Rate Limiting:** Token bucket algorithm with configurable limits
>     
> - **CORS Support:** Configurable cross-origin resource sharing
>     
> - **TLS/SSL Ready:** Built-in HTTPS support (configure with custom certificates)
>     

## 📊 Logging

JSON structured logs for easy integration with log aggregation tools:

JSON

```
{"time":"2026-07-14T09:27:18Z","level":"INFO","msg":"server starting","addr":"0.0.0.0:8080"}
{"time":"2026-07-14T09:28:49Z","level":"INFO","msg":"request","method":"GET","path":"/get","status":200,"duration":"815ms","ip":"10.0.2.100"}
```

## 🛠️ Development

### Prerequisites

- Go 1.21 or higher
    
- Make (optional)
    
- Docker/Podman (optional)
    

### Build

Bash

```
# Build binary
go build -o swiperproxy ./cmd/swiperproxy

# Build with optimizations
go build -ldflags="-s -w" -o swiperproxy ./cmd/swiperproxy

# Cross-platform build
GOOS=linux GOARCH=amd64 go build -o swiperproxy-linux-amd64 ./cmd/swiperproxy
GOOS=windows GOARCH=amd64 go build -o swiperproxy-windows-amd64.exe ./cmd/swiperproxy
GOOS=darwin GOARCH=amd64 go build -o swiperproxy-darwin-amd64 ./cmd/swiperproxy
```

### Testing

Bash

```
# Unit tests
go test ./...

# Integration tests
go test -tags=integration ./...

# Benchmarks
go test -bench=. ./...
```

## 🤝 Contributing

1. Fork the repository
    
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
    
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
    
4. Push to the branch (`git push origin feature/amazing-feature`)
    
5. Open a Pull Request
    

## 📄 License

This project is licensed under the MIT License - see the LICENSE file for details.

## 🙏 Acknowledgments

- Originally inspired by SwiperProxy by pgodschalk
    
- Built with Go's standard library
    
- Community contributions and feedback
    

_Made with ❤️ and Go_
