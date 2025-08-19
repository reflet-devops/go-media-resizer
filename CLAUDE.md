# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

This is a Go-based media resizer service that provides HTTP endpoints for image processing and resizing. It supports multiple image formats (JPEG, PNG, WEBP, AVIF) and offers both direct endpoint access and CDN-CGI compatible API endpoints for image transformation.

## Prerequisites

Before building or running the project, install required system dependencies:

## Architecture

### Core Components

- **CLI Layer** (`cli/`): Cobra-based command line interface with configuration management
- **HTTP Server** (`http/`): Echo-based HTTP server with project-based routing and middleware
- **Media Processing** (`resize/`): Image resizing and format conversion using imaging library
- **Storage Abstraction** (`storage/`): Pluggable storage backends (filesystem, MinIO S3)
- **Cache Purging** (`cache_purge/`): Support for Varnish and Cloudflare cache invalidation
- **Configuration** (`config/`): YAML-based configuration with validation

### Multi-Project Architecture

The server supports multiple projects on different hostnames, each with:
- Custom storage backend configuration
- Configurable URL patterns via regex with named groups
- Independent cache purging configuration
- Per-project headers and file type restrictions

### Key Design Patterns

1. **Factory Pattern**: Used for storage and cache purge implementations
2. **Middleware Pattern**: Domain validation and request processing
3. **Event-Driven**: File change notifications trigger cache purging
4. **Regex-Based Routing**: URL patterns extract resize parameters using named capture groups

### URL Pattern System

Projects use regex patterns with mandatory named groups:
- `source`: File path in storage backend (required)
- `width`, `height`, `quality`: Optional resize parameters
- Pattern validation includes regex tests to ensure correctness

Example regex: `^/test/((?<width>[0-9]{1,4})?(x(?<height>[0-9]{1,4}))?(-(?<quality>[0-9]{1,2}))?\/)?(?<source>.*)`

### Storage Backends

- **Filesystem**: Direct file system access using afero
- **MinIO**: S3-compatible object storage with real-time notifications

### Supported Image Formats

- **Input**: JPEG, PNG, GIF (passthrough), WEBP, AVIF, SVG (passthrough)
- **Resize**: JPEG, PNG with quality control and format conversion
- **Output**: Auto-format detection based on Accept headers or explicit format specification

## Configuration

Configuration uses YAML files with the following structure:
- Global settings (headers, file types, timeouts)
- HTTP server configuration
- CDN-CGI endpoint configuration
- Per-project settings (hostname, storage, endpoints, cache purging)

Reference configuration: `tmp/config.yml`
