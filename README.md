# Swiperproxy-go v1.01

A lightweight, caching forward HTTP proxy written in Go. Inspired by [Swiperproxy](https://github.com/pgodschalk/swiperproxy) by pgodschalk.

**Author:** Saimueli

## Features

- Proxy GET requests via `?url=...`
- In-memory cache with configurable TTL and max entries
- Per-IP rate limiting
- Security headers (nosniff, frame options, XSS protection, referrer policy)
- Request logging
- YAML configuration
- Podman support only (Containerfile + docker-compose.yml)

## Quick start

### 1. Prepare the project
Clone the repository and enter the directory.

### 2. Build the image (Podman)
```bash
podman build -t swiperproxy-go -f Containerfile .
````

> **Note:** If the build fails with a missing `go.sum` entry, run `go mod tidy` locally first to generate the `go.sum` file. This happens because `go.sum` is not included in the repository by default. Once created, rebuild the image.

### 3. Run the container

Bash

```
podman run -d --name swiperproxy -p 8080:8080 swiperproxy-go
```

### 4. Use the proxy

Open your browser or use `curl`:

Bash

```
curl "http://localhost:8080/?url=[https://example.com](https://example.com)"
```

## Testing

- **Caching:** Repeat the same request; the second one will be much faster.
    
- **Rate limiting:** Send more than 10 requests in a minute (configurable) to get a `429 Too Many Requests` response.
    
- **Security headers:** Check `curl -I` output for headers like `X-Content-Type-Options: nosniff`.
    

## Configuration

Edit `configs/config.yaml` before building, or mount a custom config with:

Bash

```
podman run -v $(pwd)/my-config.yaml:/etc/swiperproxy-go/config.yaml:ro -d -p 8080:8080 swiperproxy-go
```

**Default configuration:**

YAML

```
server:
  listen: ":8080"
proxy:
  timeout: 10s
cache:
  ttl: 5m
  max_entries: 500
rate_limit:
  requests: 10
  window: 1m
```

## Project structure

Plaintext

```
swiperproxy-go/
├── cmd/swiperproxy/main.go
├── internal/
│   ├── config/config.go
│   ├── proxy/handler.go
│   ├── proxy/cache.go
│   └── middleware/...
├── configs/config.yaml
├── Containerfile
├── docker-compose.yml
├── go.mod
└── README.md
```

## License

MIT
