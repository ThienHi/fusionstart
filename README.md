# Go Project Documentation

## Overview
This project is a Go application designed with modular architecture to ensure scalability and maintainability. It includes features such as HTTP routing, database interaction, and utility functions.

## Build source:

### To build and run the project using Docker Compose:
```
docker compose -f compose.yaml up --build
```

### Build & Run Manually
```
go mod tidy
go run app/main.go
```

### Running Tests
```
go test tests/<file>
```

------------------------------------

# Features

```
* Clean Architecture (Domain-Driven Design)

* Modular and testable codebase

* RESTful API with custom middleware

* Background job processing (workers/)

* Centralized configuration management

* Dockerized for easy deployment
```

------------------------------------

# Structure Architechture:

### **App**
Contains file main

### **Internal**
The `internal` directory contains the core logic of the application. It is organized into modules and utilities.

### **Test**
Unit tests

### Structure
```
fusionstart/
│
├── app/                # Application entry point (main.go)
│
├── internal/           # Core application logic following Clean Architecture
│   ├── configs/        # Environment configuration (e.g., reading .env files or env variables)
│   ├── constants/      # Define constants
│   ├── databases/      # Connect database - manage migrate db
│   ├── dto/            # Define data transfer objects to standard data
│   ├── handlers/       # HTTP handlers (interface with HTTP requests/responses)
│   ├── middleware/     # Custom middleware for HTTP requests (e.g., auth, logging)
│   ├── rabbitmq/       # Setup Queue
│   ├── models/         # Definitions of models/DTOs used across the system
│   ├── repositories/   # Data access layer (e.g., database, cache)
│   ├── routers/        # Route definitions and grouping
│   ├── utils/          # Utility functions and shared helpers
│   └── workers/        # Background workers and asynchronous job processing
│
├── tests/              # Unit tests
│
├── dockerfile          # Docker build instructions
├── compose.yaml        # Docker Compose file to run the entire stack
├── go.mod              # Go module definition and dependencies
└── README.md           # Project documentation
```

------------------------------------

# Environment Variable
```
# PostgreSql
DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=event_booking
DB_SSLMODE=disable

# RabbitMQ
RABBITMQ_HOST=localhost
RABBITMQ_PORT=5672
RABBITMQ_USER=guest
RABBITMQ_PASSWORD=guest
```

------------------------------------

# Contact
For questions, let's send email me `nguyenkimthien2603@gmail.com`