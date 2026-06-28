# URL Shortener

A simple URL shortener microservice written in Go.

The service exposes a gRPC API for creating short URLs and resolving them back to their original addresses. An HTTP API Gateway is also provided as a thin adapter over the gRPC service.

## Features

- Create short URLs
- Resolve original URLs by short code
- gRPC API
- PostgreSQL for persistent storage
- Redis caching
- Docker & Docker Compose support

## Architecture

The project consists of:

- gRPC service (core business logic)
- PostgreSQL
- Redis

## Run

```bash
make service-up
```

The service will start all required containers using Docker Compose.