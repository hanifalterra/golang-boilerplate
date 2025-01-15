# Golang Boilerplate
A robust, scalable Golang boilerplate for building HTTP services, workers, and cron jobs. This boilerplate follows a clean architecture with a well-defined dependency injection pattern, making it easy to maintain and extend.

## Features
- **Modular structure** for `http`, `worker`, and `cron` services
- **Dependency injection pattern** for managing service initialization
- **Docker Compose** for local development with dependency services
- Follows best practices for organizing code in the `internal` directory
- `.env` file support for environment-specific configurations

## Getting Started

### Prerequisites
- **Golang** installed
- **Docker** & **Docker Compose** installed

### First-Time Setup
1. **Create a `.env` file:**
Copy the provided `.env.example` file and fill in the necessary values:
```bash
cp .env.example .env
```
2. **Start dependency services:**
Run the following command to start the necessary services (e.g., database, message brokers) locally using Docker Compose:
```bash
docker-compose -f docker-compose-development.yml up -d
```
3. **Run the application:**
You can run the application using this method:
```bash
go run cmd/http/main.go
```

## Dependency Injection Pattern
This boilerplate uses a structured dependency injection pattern to ensure maintainability and extensibility. The process follows these steps:
1. **Initialize third-party services** (e.g., database, message broker, cache)
2. **Inject third-party services into the infrastructure layer**
3. **Inject infrastructure dependencies into the service layer**
4. **Inject services into the controller layer**

## Running Services
### HTTP Service
To run the HTTP service, use:
```bash
go run cmd/http/main.go
```
### Worker Service
To run the worker service, use:
```bash
go run cmd/worker/main.go
```
### Cron Service
To run the cron service, use:
```bash
go run cmd/cron/main.go
```