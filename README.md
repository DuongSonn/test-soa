# README

- [README](#readme)
  - [Tech stack](#tech-stack)
  - [Prerequisites](#prerequisites)
  - [Running the app](#running-the-app)
  - [Structure](#structure)
  - [Testing](#testing)
  - [Question: Log Aggregation Optimization (Focus on efficient log handling)](#question-log-aggregation-optimization-focus-on-efficient-log-handling)

## Tech stack

- [x] Go 1.21.3
- [x] Gin
- [x] GORM
- [x] Sentry
- [x] Redis
- [x] PostgreSQL

## Prerequisites

- [x] Docker (image link: https://hub.docker.com/repository/docker/duongsonn/sondth-test-soa/general)
- [x] Docker Compose
- [x] Git

## Running the app

```bash
docker compose up -d
```
then create database name: "test_soa" then run the cmd `docker compose up -d` again.

## Structure

- app
  - controller: For handling HTTP requests
  - entity: For database entities
  - middleware: For handling middleware
  - model: For business logic
  - repository: For database operations
  - service: For business logic
  - utils: For utility functions
  - helper: For helper functions
- config: For configuration
- package: For external packages

## Testing

- [x] Unit testing

```bash
go test -v ./...
```

## Question: Log Aggregation Optimization (Focus on efficient log handling)

Task: Describe your approach for handling logs efficiently across the system.
Normally, logs are written in multiple places within the system, but we want to
minimize the amount of logging code and centralize log writing in one location

Solution:
- Use a shared logging package across all services
  - This ensures consistency in logging across different services.
  - It minimizes duplicate logging code and simplifies maintenance.
  - The package can standardize log formats (e.g., JSON, structured logging).
- Log to a file first, then use Logstash (from the ELK stack) to collect logs
  - This reduces direct load on the ELK system.
  - Logging to a file first helps prevent log loss if ELK is temporarily unavailable.
  - Logstash can batch process logs, optimizing resource usage.
- Configure ELK by channels, with each service having its own channel
  - This allows easier filtering and management of logs per service in Kibana.
  - It helps with log access control and prevents logs from getting mixed up between services.
- Use OpenTelemetry (Otel) tracing to track the full request cycle
  - Distributed tracing provides full visibility into how a request flows across different services.
  - It makes debugging easier by linking logs with trace IDs.
