# Load Tests

This directory contains [k6](https://k6.io/) performance tests for the Go Media Resizer service. The tests are designed to evaluate the performance of different image processing endpoints under various load conditions.

## Overview

The load testing suite includes tests for:
- **Source image serving**: Direct image delivery without processing
- **Image resizing**: Testing resize operations with different dimensions
- **Format conversion**: Testing AVIF and WEBP format conversion
- **Resize + format conversion**: Combined operations
- **CDN-CGI endpoint**: CloudFlare-compatible CDN endpoint testing

## Prerequisites

- [k6](https://k6.io/docs/get-started/installation/) installed
- Go Media Resizer service running
- Test images available in the configured storage backend

## Quick Start

1. **Configure the tests:**
   ```bash
   cp config/config.dist.json config/config.json
   cp config/image-paths.dist.json config/image-paths.json
   ```

2. **Edit configuration files:**
   - Update `config/config.json` with your service URL and test parameters
   - Update `config/image-paths.json` with actual image paths in your storage

3. **Run project endpoint tests:**
   ```bash
   k6 run load-endpoint_project.js
   ```

4. **Run CDN-CGI endpoint tests:**
   ```bash
   k6 run load-endpoint_cdn_cgi.js
   ```

## Configuration

### config/config.json

Main configuration file containing:

- **baseUrl**: Target service URL (default: `http://localhost:8080`)
- **pathPrefix**: Optional path prefix for image URLs
- **scenarios**: Test scenario parameters for each test type

#### Scenario Parameters

Each scenario type (source, format, resize, etc.) can be configured with:

- **rate**: Requests per time unit
- **timeUnit**: Time unit (e.g., "1s", "30s")
- **duration**: Test duration (e.g., "1m")
- **preAllocatedVUs**: Pre-allocated virtual users
- **maxVUs**: Maximum virtual users
- **response_time**: Expected response time threshold (ms)

Example scenario configuration:
```json
{
  "source": {
    "response_time": 200,
    "small": {
      "rate": 5,
      "timeUnit": "1s",
      "preAllocatedVUs": 0,
      "maxVUs": 28,
      "duration": "1m"
    }
  }
}
```

### config/image-paths.json

Contains arrays of test image paths categorized by size:

- **large**: High-resolution images (2400x1800+)
- **medium**: Standard resolution images (800x600 to 1280x960)
- **small**: Low-resolution images (64x64 to 600x400)

## Test Types

### 1. Source Tests (`tests/source.js`)

Tests direct image serving without any processing. Validates:
- HTTP 200 responses
- Correct image content-type headers
- Response times under threshold

**Scenarios:** `source_small_test`, `source_medium_test`, `source_large_test`

### 2. Resize Tests (`tests/resize.js`)

Tests image resizing operations with widths: 400px, 800px, 1200px. Validates:
- Successful resizing operations
- Response time performance by image size
- Image format preservation

**Scenarios:** `resize_400_test`, `resize_800_test`, `resize_1200_test`

### 3. Format Conversion Tests (`tests/format.js`)

Tests format conversion to AVIF and WEBP formats. Validates:
- Format conversion accuracy
- Performance impact of different formats
- Content negotiation via Accept headers

**Scenarios:** `format_avif_test`, `format_webp_test`

### 4. Combined Tests (`tests/resize_format.js`)

Tests simultaneous resizing and format conversion operations.

**Scenarios:** `resizeFormat_small_test`, `resizeFormat_medium_test`, `resizeFormat_large_test`

### 5. CDN-CGI Tests (`tests/cdn_cgi.js`)

Tests CloudFlare-compatible CDN endpoint with the format:
```
/cdn-cgi/image/width=800,format=auto/{source_url}
```

**Scenarios:** `cdnCgi_400_test`, `cdnCgi_800_test`, `cdnCgi_1200_test`

## Environment Variables

Override configuration at runtime:

- **BASE_URL**: Override base service URL
- **PATH_PREFIX**: Override path prefix
- **WIDTH**: Set specific width for resize tests
- **DISTRIBUTION**: Set image size distribution (small/medium/large)

Example:
```bash
BASE_URL=https://your-service.com k6 run load-endpoint_project.js
# or
BASE_URL=https://your-service.com PATH_PREFIX="/my-prefix" k6 run load-endpoint_project.js
```

## Test Execution

### Individual Test Execution

Run specific test types:
```bash
# Source images only
k6 run tests/source.js

# Format conversion only  
k6 run tests/format.js

# Resize operations only
k6 run tests/resize.js
```

### Combined Test Execution

Run multiple test scenarios:
```bash
# All project endpoint tests
k6 run load-endpoint_project.js

# CDN-CGI endpoint tests
k6 run load-endpoint_cdn_cgi.js
```

## Performance Thresholds

Each test defines performance thresholds:

- **HTTP Failure Rate**: < 1% (`http_req_failed: rate<0.01`)
- **Response Times**: 95th percentile under configured thresholds
- **Image Validation**: Content-type and status code checks

Example threshold configuration:
```javascript
'http_req_duration{test_type:source,distribution:small}': ['p(95)<200'],
'http_req_duration{test_type:resize,width:800}': ['p(95)<300'],
'http_req_failed': ['rate<0.01']
```

## Test Results

k6 provides detailed metrics including:

- **Request rate** and **response times**
- **Virtual user** scaling behavior  
- **Threshold pass/fail** status
- **Custom metrics** by test type and parameters

Results are categorized by:
- Test type (source, resize, format, etc.)
- Image size distribution (small/medium/large)
- Processing parameters (width, format)

## Troubleshooting

### Common Issues

1. **Connection refused**: Ensure the service is running on the configured URL
2. **404 errors**: Verify image paths exist in your storage backend
3. **High response times**: Check service resource usage and configuration
4. **Threshold failures**: Adjust scenario parameters or performance expectations

### Debug Mode

Enable verbose output:
```bash
k6 run --verbose load-endpoint_project.js
```

### Health Check

Tests include automatic health checks at startup:
```javascript
const testResponse = http.get(CONFIG.baseUrl + '/health/ping', {timeout: '5s'});
```

## Directory Structure

```
loadtests/
├── README.md                     # This documentation
├── load-endpoint_project.js      # Main project endpoint test runner
├── load-endpoint_cdn_cgi.js      # CDN-CGI endpoint test runner
├── config/
│   ├── config.js                 # Configuration loader
│   ├── config.json               # Runtime configuration
│   ├── config.dist.json          # Configuration template
│   ├── images.js                 # Image path loader
│   ├── image-paths.json          # Runtime image paths
│   └── image-paths.dist.json     # Image paths template
├── common/
│   └── utils.js                  # Shared utilities and image selection
└── tests/
    ├── source.js                 # Source image serving tests
    ├── resize.js                 # Image resizing tests
    ├── format.js                 # Format conversion tests
    ├── resize_format.js          # Combined resize + format tests
    └── cdn_cgi.js               # CDN-CGI endpoint tests
```
