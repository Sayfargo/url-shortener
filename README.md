# URL Shortener (gRPC)

Simple URL shortener microservice written in Go.  
The service provides an API for creating short links and resolving them back to original URLs via gRPC.

## Features

- Create short URLs
- Resolve original URLs by short code
- gRPC-based communication
- PostgreSQL for persistent storage
- Redis caching layer for faster lookups
- Docker support

## Architecture

The project consists of:

- gRPC server (core service)
- PostgreSQL repository (persistent storage)
- Redis cache (fast access layer)
- API Gateway (HTTP → gRPC)

## Run with Docker

```bash
make serivce-up
