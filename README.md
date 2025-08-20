# GO Media Resizer

A high-performance image resizing service written in Go, providing HTTP endpoints for image processing and transformation. Supports multiple image formats (JPEG, PNG, WebP, AVIF) with pluggable storage backends and cache purging integration.

## Features

- **Real-time image resizing** with configurable quality
- **Multiple formats**: JPEG, PNG, WebP, AVIF, GIF, SVG
- **Storage backends**: Local filesystem, MinIO/S3
- **CDN-CGI API** compatible with Cloudflare Image Resizing
- **Cache purging**: Varnish and Cloudflare support
- **Multi-project**: Multiple services on different hostnames
- **Flexible URL patterns**: Regex configuration with named groups
- **Real-time notifications**: File events for automatic cache purging
- **Automatic fallback**: Failover to secondary storage
- **Customizable headers**: Global and per-project configuration

## Quick Start

### Using Docker

```bash
# Start development services (MinIO + Varnish)
docker compose up -d

# First, build the Go binary (required for Docker image)
go build -v .

# Build the Docker image
docker build -t go-media-resizer .

# Run the service
docker run -p 8080:8080 -v ./tmp/config.yml:/config.yml go-media-resizer start -c /config.yml
```

#### Docker environment variables

- **GO_MEDIA_RESIZER_CFG**: The content of variable will be print in configuration path ($GO_MEDIA_RESIZER_CONFIG_PATH) only if not empty.
- **GO_MEDIA_RESIZER_CONFIG_PATH**: Configuration path used to command start. Default: /etc/go-media-resizer/config.yml
- **LOG_LEVEL**: Chose the log level. Default: INFO

### Build from Source

#### System Requirements

##### Debian/Ubuntu:
```bash
sudo apt-get update
sudo apt-get install libaom-dev libwebp-dev
```

##### RHEL/CentOS/Fedora:
```bash
sudo dnf install libaom-devel libwebp-devel
```

##### macOS:
```bash
brew install webp libaom
```

#### Build Process

```bash
# Install Go dependencies
go mod download

# Generate mocks (required for tests)
./bin/mock.sh

# Build
go build -v .

# Test
go test -v ./...

# Run with default configuration
./go-media-resizer start -c tmp/config.yml
```

## Configuration

The service uses YAML files for configuration. See the [complete configuration documentation](docs/CONFIGURATION.md) for all details.

### Minimal Configuration

```yaml
http:
  listen: "127.0.0.1:8080"

projects:
  - id: "main"
    hostname: "media.example.com"
    
    storage:
      type: "fs"
      config:
        prefix_path: "/path/to/images"
    
    endpoints:
      - regex: '^/resize/((?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?\/)?(?<source>.*)'
        default_resize:
          format: "auto"
```

### Usage Examples

```bash
# Original image
curl http://media.example.com/resize/original/product.jpg

# Resize to 500px width
curl http://media.example.com/resize/500/product.jpg

# Resize to 500x300px with quality 80
curl http://media.example.com/resize/500x300-80/product.jpg

# CDN-CGI API (if enabled)
curl http://localhost:8080/cdn-cgi/image/width=500,height=300,quality=85,format=webp/https://example.com/image.jpg
```

## Architecture

### Core Components

- **CLI Layer**: Command-line interface with Cobra
- **HTTP Server**: Echo server with multi-project routing
- **Media Processing**: Image resizing with imaging library
- **Storage Abstraction**: Filesystem and MinIO/S3 backends
- **Cache Purging**: Varnish and Cloudflare support
- **Configuration**: YAML validation with regex testing

### URL Patterns

Projects use regex patterns with mandatory named groups:

- **`source`** (required): File path in storage backend
- **`width`**, **`height`**, **`quality`**: Resize parameters
- **`format`**: Output format (auto, jpeg, png, webp, avif)

### Storage Backends

#### Filesystem
```yaml
storage:
  type: "fs"
  config:
    prefix_path: "/var/www/media"
```

#### MinIO/S3 with Fallback
```yaml
storage:
  type: "minio"
  config:
    endpoint: "s3.amazonaws.com"
    bucket: "media-bucket"
    access_key: "ACCESS_KEY"
    secret_key: "SECRET_KEY"
    use_ssl: true
    
    fallback:
      endpoint: "backup.s3.amazonaws.com"
      bucket: "backup-bucket"
      access_key: "BACKUP_KEY"
      secret_key: "BACKUP_SECRET"
```

### Cache Purging

Integrated support for Varnish and Cloudflare with automatic purging:

```yaml
purge_caches:
  - type: "varnish-tag"
    config:
      server: "http://varnish:6081"
      
  - type: "cloudflare-tag"
    config:
      zone_id: "your-zone-id"
      auth_token: "your-api-token"
```

## Development

### Useful Commands

```bash
# Tests with coverage
go test -v -cover ./...

# Test specific package
go test -v ./resize

# Regenerate mocks
./bin/mock.sh

# Configuration validation
./go-media-resizer validate -c config.yml

# Debug mode
./go-media-resizer start -c config.yml -l DEBUG
```

### Environment Variables

```bash
export GO_MEDIA_RESIZER_HTTP_LISTEN="0.0.0.0:8080"
export GO_MEDIA_RESIZER_REQUEST_TIMEOUT="5s"
```

## Complete Documentation

- **[Configuration](docs/CONFIGURATION.md)** - Complete configuration guide
- **[Resize Options](docs/RESIZE-OPTIONS.md)** - Complete Media Resize Options guide
- **[Architecture](CLAUDE.md)** - Technical details for developers

## Performance

- **Concurrent request** processing
- **Intelligent caching** with automatic purging
- **Optimized formats** (WebP, AVIF) for reduced size
- **Fast fallback** to secondary storage
- **Real-time notifications** for cache synchronization

## License

MIT License - see the LICENSE file for details.
