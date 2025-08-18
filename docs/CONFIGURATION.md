# Configuration Documentation

This documentation describes all configuration possibilities for the Go Media Resizer service.

## General Configuration Structure

The service uses YAML configuration files with the following structure:

```yaml
# HTTP server configuration
http:
  listen: "127.0.0.1:8080"

# Accepted file types (without resizing)
accept_type_files: # Default value
  - "plain"
  - "gif" 
  - "mpeg" 
  - "mp4"
  - "svg"
  - "avif"
  - "webp"

# File types supporting resizing
resize_type_files: # Default value
  - "png"
  - "jpeg"

# Global HTTP headers
headers:
  x-custom: "myvalue"
  cache-control: "max-age=3600"

# HTTP request timeout
request_timeout: "2s"

# CDN-CGI configuration (optional)
resize_cgi:
  enabled: true
  allow_self_domain: true # Check that host in source is the same domain 
  allow_domains: # Check that host in source is present in this list
    - "media.example.com"
  default_resize:
    format: "auto"
  extra_headers:
    x-cgi: "enabled"

# Projects (at least one required)
projects:
  - id: "main"
    hostname: "media.example.com"
    # ... see Projects section
```

## Project Configuration

Each project defines an independent media service with its own hostname:

```yaml
projects:
  - id: "main"                    # Unique project identifier
    hostname: "media.example.com" # Hostname for this project
    prefix_path: "/cdn"           # URL prefix (optional)
    webhook_token: "secret_token" # Bearer token for webhook authentication (optional)
    
    # Storage configuration (required)
    storage:
      type: "minio"  # or "fs"
      config:
        # ... see Storage section
    
    # Cache purging configuration (optional)
    purge_caches:
      - type: "varnish-tag"
        config:
          # ... see Cache Purge section
    
    # Endpoints with regex patterns (at least one required)
    endpoints:
      - regex: '^/prod/(?<source>.*)'
        default_resize:
          format: "auto"
        regex_tests:
          - path: "/prod/image.png"
            result_opts: 
              source: "image.png"
              format: "auto"          # Matches default_resize.format
              origin_format: "png"
      
      - regex: '^/resize/((?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?(-(?<quality>[0-9]{1,2}))?\/)?(?<source>.*)'
        default_resize:
          format: "webp"
          quality: 85
        regex_tests:
          - path: "/resize/500x500-80/product/image.jpg"
            result_opts:
              width: 500
              height: 500  
              quality: 80             # Extracted from regex, overrides default
              source: "product/image.jpg"
              format: "webp"          # Matches default_resize.format
              origin_format: "jpeg"
          - path: "/resize/product/image.jpg"
            result_opts:
              source: "product/image.jpg"
              format: "webp"          # Matches default_resize.format
              quality: 85             # Matches default_resize.quality
              origin_format: "jpeg"
    
    # Accepted file types (inherits from global if not specified)
    #accept_type_files: [] # If defined, it's override global config
    extra_accept_type_files: 
      - "tiff"
    
    # Project-specific headers
    #headers: # If defined, it's override global config
    #  x-project: "main"
    extra_headers:
      x-version: "1.0"
```

### Endpoint Regex Patterns

Regex patterns must contain the following mandatory named groups:

- **`source`** (required): File path in storage backend
- **`width`** (optional): Resize width
- **`height`** (optional): Resize height
- **`quality`** (optional): JPEG quality (1-100)
- **`format`** (optional): Output format

#### Regex Testing

The `regex_tests` array is used to validate that your regex patterns work correctly and extract the expected parameters. Each test case should:

1. **Include all `default_resize` options**: Test results must contain all options defined in `default_resize` for the endpoint
2. **Validate regex behavior**: Ensures the regex captures named groups correctly
3. **Test edge cases**: Include various URL patterns to verify comprehensive coverage

**Important**: The `result_opts` in each test must include all fields from `default_resize`, plus any additional fields extracted by the regex pattern.

### Resize Options

```yaml
default_resize:
  format: "auto"        # auto, jpeg, png, webp, avif
  width: 800           # Width in pixels
  height: 600          # Height in pixels
  quality: 85          # JPEG quality (1-100)
  fit: "crop"          # Resize method: crop or resize (default)
  
  # Image adjustment parameters
  blur: 2.5            # Blur radius (0 = no blur)
  brightness: 10       # Brightness adjustment (-100 to 100)
  contrast: 15         # Contrast adjustment (-100 to 100)
  saturation: 20       # Saturation adjustment (-100 to 100)
  sharpen: 1.2         # Sharpening amount (0 = no sharpening)
  gamma: 1.2           # Gamma correction (1.0 = no correction)
```

**Supported Formats:**
- `auto`: Automatic detection based on Accept headers
- `jpeg`: JPEG format with quality support
- `png`: PNG format
- `webp`: Modern WebP format (requires libwebp-dev)
- `avif`: AVIF format (requires libaom-dev)

**Resize Methods:**
- `resize` (default): Proportional resizing with Lanczos algorithm
- `crop`: Center crop to exact dimensions

**Image Adjustment Parameters:**
- `blur`: Blur radius in pixels (0 = no blur, typical range: 0.5-5.0)
- `brightness`: Brightness adjustment (-100 to 100, 0 = no change)
- `contrast`: Contrast adjustment (-100 to 100, 0 = no change)
- `saturation`: Saturation adjustment (-100 to 100, 0 = no change)  
- `sharpen`: Sharpening amount (0 = no sharpening, typical range: 0.5-3.0)
- `gamma`: Gamma correction (1.0 = no correction, typical range: 0.5-2.5)

## Storage Configuration

### Filesystem Storage

```yaml
storage:
  type: "fs"
  config:
    prefix_path: "/var/www/media"  # Root path for files
```

### MinIO/S3 Storage

```yaml  
storage:
  type: "minio"
  config:
    endpoint: "s3.amazonaws.com"     # S3/MinIO endpoint
    bucket: "media-bucket"           # Bucket name
    access_key: "AKIAIOSFODNN7EXAMPLE"
    secret_key: "wJalrXUtnFEMI/K7MDENG/bPxRfiCYEXAMPLEKEY"  
    use_ssl: true                    # HTTPS (default: false)
    prefix_path: "images"            # Prefix in bucket (optional)
    
    # Health check configuration (optional)
    health_check_interval: "5s"
    
    # Fallback to another endpoint (optional)
    fallback:
      endpoint: "backup.s3.amazonaws.com"
      bucket: "media-backup"  
      access_key: "BACKUP_ACCESS_KEY"
      secret_key: "BACKUP_SECRET_KEY"
      use_ssl: true
```

**MinIO Features:**
- **Real-time notifications**: Listens for create/delete events
- **Automatic fallback**: Switches to secondary storage on failure
- **Health checks**: Monitors primary storage health

## Cache Purging Configuration

The service supports automatic cache purging when files are modified in storage. There are two purging strategies: **tag-based** (recommended) and **URL-based**.

**Recommendation**: Use tag-based purging (`varnish-tag` or `cloudflare-tag`) as it can purge all cached variants of an image regardless of resize parameters. URL-based purging only works for endpoints without resize parameters.

### Varnish - Tag-based Purge

```yaml
purge_caches:
  - type: "varnish-tag"
    config:
      server: "http://varnish.example.com:6081"
```

**How it works:**
- Uses BAN method with tag header based on path hash
- Generated tag: SHA-256 hash of source path
- Compatible with standard VCL configuration

### Varnish - URL-based Purge

```yaml
purge_caches:
  - type: "varnish-url"  
    config:
      server: "http://varnish.example.com:6081"
```

**How it works:**
- Uses PURGE method on exact URL
- Direct cache purge for specific URL
- **Limitation**: Only works for URLs without resize parameters (i.e., only `source` parameter)
- Does not purge all cached variants when resize options are used

### Cloudflare - Tag-based Purge

```yaml
purge_caches:
  - type: "cloudflare-tag"
    config:
      zone_id: "d56084c56a852d1ceddaa30f83c72bce"
      
      # Authentication via API Token (recommended)
      auth_token: "YQSn-xWAQiiEh9qM58wZNnyQS7FUdoqGIUAbrh7T"
      
      # OR authentication via Email + Key (legacy)
      auth_email: "user@example.com"  
      auth_key: "c2547eb745079dac9320b638f5e225cf483cc5cfdda41"
```

### Cloudflare - URL-based Purge

```yaml  
purge_caches:
  - type: "cloudflare-url"
    config:
      zone_id: "d56084c56a852d1ceddaa30f83c72bce"
      auth_token: "YQSn-xWAQiiEh9qM58wZNnyQS7FUdoqGIUAbrh7T"
```

**How it works:**
- Uses Cloudflare API to purge specific URLs
- Direct cache purge for exact URL
- **Limitation**: Only works for URLs without resize parameters (i.e., only `source` parameter)
- Does not purge all cached variants when resize options are used

**Cloudflare Notes:**
- Tags are automatically generated by hashing source path
- API Token is recommended for security
- Zone ID available in Cloudflare dashboard

## CDN-CGI Configuration

The service can emulate Cloudflare Image Resizing API:

```yaml
resize_cgi:
  enabled: true                    # Enable CDN-CGI mode
  allow_self_domain: true          # Allow service domain
  allow_domains:                   # Allowed domains for sources
    - "cdn.example.com"
    - "media.example.com"
  default_resize:
    format: "auto"
    quality: 85
  extra_headers:
    x-powered-by: "go-media-resizer"
```

**CDN-CGI Endpoint:**
```
GET /cdn-cgi/image/width=500,height=300,quality=85,format=webp/https://example.com/image.jpg
GET /cdn-cgi/image/width=500,blur=2.0,brightness=10,contrast=15/https://example.com/image.jpg
```

**Supported Parameters:**
- `width`: Width in pixels
- `height`: Height in pixels
- `quality`: JPEG quality (1-100)
- `format`: Output format (auto, jpeg, png, webp, avif)
- `fit`: Resize method (crop or resize)
- `blur`: Blur radius (0 = no blur)
- `brightness`: Brightness adjustment (-100 to 100)
- `contrast`: Contrast adjustment (-100 to 100)
- `saturation`: Saturation adjustment (-100 to 100)
- `sharpen`: Sharpening amount (0 = no sharpening)
- `gamma`: Gamma correction (1.0 = no correction)

## Webhook Configuration

The service provides webhook endpoints for receiving external events that trigger cache purging operations. Each project can have its own webhook configuration with optional authentication.

### Basic Configuration

```yaml
projects:
  - id: "main"
    hostname: "media.example.com"
    prefix_path: "/api"                        # Optional URL prefix
    webhook_token: "your-secret-bearer-token"  # Optional authentication token
    # ... other project settings
```

### Webhook Endpoint

The webhook endpoint is accessible only on the configured project hostname and prefix path:
```
POST https://{hostname}{prefix_path}/webhook
```

Examples:
- Without prefix_path: `POST https://media.example.com/webhook`
- With prefix_path: `POST https://media.example.com/api/webhook`

### Authentication

If `webhook_token` is configured, requests must include the Bearer token in the Authorization header:
```bash
curl -X POST https://media.example.com/api/webhook \
  -H "Authorization: Bearer your-secret-bearer-token" \
  -H "Content-Type: application/json" \
  -d '[{"type": "purge", "path": "images/photo.jpg"}]'
```

If no `webhook_token` is configured, the webhook endpoint accepts requests without authentication.

### Request Format

The webhook expects a JSON array of event objects. Each event must contain the required fields:
For each element `Event` the path must be without `hostname` or `prefix_path`. The system will add this information if is necessary.

```json
[
  {
    "type": "purge",
    "path": "images/photo.jpg"
  },
  {
    "type": "purge", 
    "path": "images/old-image.png"
  }
]
```

**Supported Event Types:**
- `purge`: Triggers cache purging for the specified file path

### Response Codes

- **202 Accepted**: Events received and queued successfully
- **400 Bad Request**: Invalid JSON format or validation failure
- **401 Unauthorized**: Missing or invalid Bearer token (when webhook_token is configured)

### Event Processing

Events received via webhook are processed by the same system that handles storage notifications, triggering configured cache purging operations for the affected files.

**Note**: The webhook provides a way for external systems to notify the media resizer of file changes when direct storage notifications (like MinIO events) are not available.

## Environment Variables

The service supports configuration via environment variables with `GO_MEDIA_RESIZER_` prefix:

```bash
export GO_MEDIA_RESIZER_HTTP_LISTEN="0.0.0.0:8080"
export GO_MEDIA_RESIZER_REQUEST_TIMEOUT="5s"  
export GO_MEDIA_RESIZER_PROJECTS_0_STORAGE_CONFIG_ENDPOINT="minio.example.com:9000"
```

## Complete Configuration Example

```yaml
http:
  listen: "0.0.0.0:8080"

headers:
  cache-control: "public, max-age=31536000"
  x-powered-by: "go-media-resizer"

request_timeout: "10s"

accept_type_files:
  - "plain"
  - "gif"
  - "mpeg"
  - "mp4"
  - "svg"
  - "avif"
  - "webp"

resize_type_files:
  - "png"
  - "jpeg"

resize_cgi:
  enabled: true
  allow_self_domain: true
  allow_domains:
    - "media.example.com"

projects:
  - id: "production"
    hostname: "media.example.com"
    prefix_path: "/v1"
    
    storage:
      type: "minio"
      config:
        endpoint: "s3.amazonaws.com"
        bucket: "production-media"
        access_key: "${AWS_ACCESS_KEY_ID}"
        secret_key: "${AWS_SECRET_ACCESS_KEY}"
        use_ssl: true
        prefix_path: "images"
        
        fallback:
          endpoint: "backup.s3.example.com"
          bucket: "backup-media"
          access_key: "${BACKUP_ACCESS_KEY}"
          secret_key: "${BACKUP_SECRET_KEY}"
          use_ssl: true
    
    purge_caches:
      - type: "cloudflare-tag"
        config:
          zone_id: "d56084c56a852d1ceddaa30f83c72bce"
          auth_token: "${CLOUDFLARE_API_TOKEN}"
      
      - type: "varnish-tag"
        config:
          server: "http://varnish.internal:6081"
    
    endpoints:
      - regex: '^/original/(?<source>.*)'
        default_resize:
          format: "auto"
      
      - regex: '^/thumb/((?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?(-(?<quality>[0-9]{1,2}))?(-(?<format>webp|avif|jpeg|png))?\/)?(?<source>.*)'
        default_resize:
          format: "webp"
          quality: 85
          fit: "crop"
        regex_tests:
          - path: "/thumb/300x200-90-avif/products/image.jpg"
            result_opts:
              width: 300
              height: 200
              quality: 90             # Extracted from regex, overrides default
              format: "avif"          # Extracted from regex, overrides default
              fit: "crop"             # Matches default_resize.fit
              source: "products/image.jpg"
              origin_format: "jpeg"
          - path: "/thumb/products/image.jpg"
            result_opts:
              source: "products/image.jpg"
              format: "webp"          # Matches default_resize.format
              quality: 85             # Matches default_resize.quality
              fit: "crop"             # Matches default_resize.fit
              origin_format: "jpeg"
    
    headers:
      x-project: "production"
      x-environment: "prod"
      
    extra_headers:
      x-version: "1.2.0"
```

This configuration serves images from S3 with resizing, Cloudflare and Varnish caching, and automatic fallback on failures.

