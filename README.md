# Real-Time Bus Tracking System Architecture

## Overview
The **real-time bus tracking system** employs an **event-driven microservice architecture** with **choreography**, using **Kafka** for asynchronous event streaming, **Gin** for the HTTP API Gateway, **MySQL** (centralized with master-replica), **Redis** for caching, and **golang-migrate** for migrations. It supports bus registration, trip scheduling, real-time updates (routes, stops, distance/duration, wheelchair accessibility), and route search. The system includes four core microservices (**API Gateway**, **Bus Service**, **Route Service**, **Trip Service**) with potential for up to 12 (e.g., **User Service**, **Notification Service**, **Event Service**). The `helpers` package manages environment variable conversions for database and Kafka retries. The system is structured as a **monorepo**, with each microservice as an individual application.

## Architecture Pattern
- **Primary Pattern**: **Event-Driven Microservices with Shared Database and Choreography**.
  - **Event-Driven**: Services communicate via Kafka events (`BusCreated`, `TripUpdated`, `DriverLocationUpdated`).
  - **Shared Database**: Centralized MySQL (`t_buses`, `t_routes`, etc.) and Redis (`bus:<id>`).
  - **Choreography**: Services react to events independently.
- **Secondary Pattern**: **API Gateway** for client interactions.

## Repository Strategy
- **Monorepo**:
  - All microservices, shared code, and configurations are in a single repository (`real-time-bus-tracking/`).
  - Each microservice is an **individual application** with its own binary (`cmd/<service>/main.go`) and container (`cmd/<service>/Dockerfile`).
  - **Rationale**:
    - Simplifies shared code management (`internal/models/`, `internal/helpers/`, `internal/events/`).
    - Ensures event schema consistency (`internal/events/types/events.go`).
    - Streamlines local development with `config/docker-compose.yml`.
    - Suitable for small to medium teams (~5–10 engineers).
  - **Mitigations**:
    - Selective CI/CD builds for changed services.
    - Minimal `internal/` packages to reduce coupling.
    - Linters to enforce service boundaries.

## Architecture Diagram (Textual)
```
+-------------------+           +-------------------+           +-------------------+
|      Client       |           | API Gateway (Gin) |           |  Microservices    |
| (Mobile/Web App)  |           |                   |           |                   |
+-------------------+           +-------------------+           +-------------------+
|                   |           |                   |           | Bus Service       |
| POST /buses       | --------> | Publish BusCreated| --------> | Route Service     |
| PATCH /buses/:id  |           | Publish BusUpdated| Kafka     | Trip Service      |
| POST /routes      |           | Publish RouteCreated         | User Service      |
| POST /search      |           | Query RouteService| gRPC      | Notification Serv |
| POST /trips       |           | Publish TripCreated          | Analytics Service |
| GET /routes/:id/  |           | Stream TripService| gRPC      | Payment Service   |
| buses             |           |                   |           | Geolocation Serv  |
|                   |           |                   |           | Driver Service    |
|                   |           |                   |           | Config Service    |
|                   |           |                   |           | Event Service     |
| JSON Response     | <-------- | JSON Response     | <-------- |                   |
+-------------------+           +-------------------+           +-------------------+
                                                  |           |
                                                  |           |   +----------------+
                                                  |           |   |   Kafka        |
                                                  |           |   | (bus.created,  |
                                                  |           |   |  trip.updated, |
                                                  | Bus       |-->|  driver.loc)   |
                                                  | Trip      |   +----------------+
                                                  | Notify    |   |
                                                  | Driver    |   +----------------+
                                                  | Event     |-->|   MySQL        |
                                                  |           |   | (Master+2      |
                                                  |           |   | Replicas)      |
                                                  |           |   | (t_buses,      |
                                                  |           |   | t_routes,      |
                                                  |           |   | t_stops,       |
                                                  |           |   | t_trips,       |
                                                  |           |   | t_users,       |
                                                  |           |   | t_notifications|
                                                  |           |   | t_payments,    |
                                                  |           |   | t_configs,     |
                                                  |           |   | t_events,      |
                                                  |           |   | t_drivers)     |
                                                  |           |   |                |
                                                  |           |   |   Redis        |
                                                  |           |-->| (bus:<id>,     |
                                                  |           |   |  geo:<id>,     |
                                                  |           |   |  driver:<id>,  |
                                                  |           |   |  config:<key>) |
                                                  |           |   +----------------+
```

## Folder Structure
- **cmd/**: Microservice entry points (`api-gateway`, `bus-service`, `event-service`, etc.), each with `main.go` and `Dockerfile`.
- **internal/**:
  - `models/`: Database initialization (`models.go`) and schemas (`schemas/models.go`).
  - `helpers/`: Environment variable conversions (`env.go`).
  - `events/`: Kafka producer/consumer (`kafka/producer.go`, `kafka/consumer.go`) and event types (`types/events.go`).
  - `grpc/`: Limited gRPC for queries (`bus_service.go`, `trip_service.go`).
  - `http/`: Gin handlers for API Gateway (`bus_handlers.go`).
- **migrations/**: SQL migrations (`001_create_tables.up.sql`, `007_create_events.up.sql`).
- **tests/**: Unit tests (`models_test.go`, `event_service_test.go`).
- **config/**: Deployment files (`docker-compose.yml`, `.env`, `kafka-deployment.yaml`).

## Event-Driven Design
- **Events**:
  - `BusCreated`, `BusUpdated`: Published by `BusService`.
  - `TripCreated`, `TripUpdated`: Published by `TripService`.
  - `DriverLocationUpdated`: Published by `DriverService`.
  - `SendNotification`: Consumed by `NotificationService`.
- **Kafka**: Handles event streaming (`KAFKA_BROKER_ADDR`, `KAFKA_MAX_RETRIES=5`).
- **Event Service**: Centralizes consumer/producer logic (`cmd/event-service/`).

## Database Connection Strategy
- **Shared Database**:
  - **MySQL**: Centralized with master-replica, accessed by services.
  - **Redis**: Caches real-time data (`bus:<id>`, `geo:<id>`).
- **Retry Mechanism**:
  - Environment variables (`MYSQL_MAX_RETRIES`, `KAFKA_RETRY_BACKOFF`) via `helpers`.
  - Retries on transient errors (MySQL 1040, Kafka timeouts).
- **Implementation**: `internal/models/models.go`, `internal/events/kafka/producer.go`.

## Microservices Details
- **API Gateway**: Publishes events for writes, queries services for reads.
- **Bus Service**: Publishes `BusCreated`, `BusUpdated`, consumes `DriverLocationUpdated`.
- **Trip Service**: Publishes `TripCreated`, consumes `DriverLocationUpdated` for streaming.
- **Driver Service**: Publishes `DriverLocationUpdated`.
- **Notification Service**: Consumes `TripUpdated`, `DriverLocationUpdated`.
- **Event Service**: Manages Kafka consumers/producers.

## Scalability and High Availability
- **Kafka**: Partitioned topics for parallel processing.
- **MySQL**: Master-replica with retries.
- **Redis**: Cluster-ready with retries.
- **Microservices**: Independent scaling via Docker/Kubernetes.

## Edge Cases and Breaking Points
- **Edge Cases**:
  - Invalid environment variables: `utils` uses defaults.
  - Event duplication: Kafka consumer idempotency required.
  - Replica lag: `RandomPolicy` falls back to master.
- **Breaking Points**:
  - Kafka outages: Retries exhaust after ~7.5s.
  - Shared database: Limits schema independence.
  - Eventual consistency: Delays in event processing.

## Testing
- **Unit Tests**: `tests/helpers_test.go`, `tests/event_service_test.go`.
- **Integration Tests**: Dockerized MySQL/Redis/Kafka for event scenarios.
- **Coverage**: Target 70% for event, retry, and service logic.