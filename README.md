# Go Microservices with Kafka

A robust microservices architecture implementation using Go and Apache Kafka for event-driven communication between services.

## 🏗️ Architecture Overview

![High Level Architecture](docs/images/diagram_high_level.png)

This project consists of five microservices:

- **Orchestra Service**: Orchestrates workflows and manages the communication between services
- **Order Service**: Handles order processing and management
- **Payment Service**: Manages payment processing
- **Product Service**: Handles product inventory and details
- **User Service**: Manages user authentication and profiles

### Event-Driven Architecture
The services communicate asynchronously using Apache Kafka as the message broker, enabling loose coupling and high scalability.

## 🔄 Order Creation Flow

![Order Creation Flow](docs/images/flow_create_order.png)

## 📊 Orchestra Service ERD

![Orchestra Service ERD](docs/images/orchestra_svc_erd.png)

## 🚀 Technologies

- **Go**: Primary programming language
- **Apache Kafka**: Event streaming platform
- **PostgreSQL**: Primary database
- **Docker**: Containerization
- **Sarama**: Kafka client for Go
- **SQLC**: SQL query generator for Go
- **Gin**: HTTP web framework

## 📋 Prerequisites

- Go 1.19 or higher
- Docker and Docker Compose
- PostgreSQL
- Apache Kafka
- Make

## 🛠️ Setup and Installation

1. Clone the repository
```bash
git clone <repository-url>
cd go-microservices-kafka
```

2. Start the infrastructure services
```bash
docker-compose up -d
```

3. Create the databases for each service
```bash
make -C orchestra-svc createdb
make -C order-svc createdb
make -C payment-svc createdb
make -C product-svc createdb
make -C user-svc createdb
```

4. Run database migrations
```bash
make -C orchestra-svc migrateup
make -C order-svc migrateup
make -C payment-svc migrateup
make -C product-svc migrateup
make -C user-svc migrateup
```

## 🏃‍♂️ Running the Services

Each service can be started individually:

```bash
# Start Orchestra Service (Port 8082)
cd orchestra-svc && go run main.go

# Start Order Service (Port 8081)
cd order-svc && go run main.go

# Start Payment Service (Port 8083)
cd payment-svc && go run main.go

# Start Product Service (Port 8084)
cd product-svc && go run main.go

# Start User Service (Port 8085)
cd user-svc && go run main.go
```

## 🔄 Workflow Example

1. Client creates an order
2. Order Service processes the request and publishes an event
3. Orchestra Service receives the event and orchestrates the workflow:
   - Validates user details
   - Checks product availability
   - Processes payment
   - Updates order status

## 📁 Project Structure

```plaintext
.
├── orchestra-svc/     # Workflow orchestration service
├── order-svc/        # Order management service
├── payment-svc/      # Payment processing service
├── product-svc/      # Product management service
├── user-svc/         # User management service
└── docker-compose.yml
```

## 🛠️ Development

### Generate SQLC Code
```bash
make -C <service-name> sqlc
```

### Create New Migration
```bash
make -C <service-name> migratecreate name=migration_name
```

### Generate Mocks
```bash
make -C <service-name> mockdb
```

## 🔐 Environment Variables

Each service requires its own `.env` file. Example configuration:

```plaintext
PORT=8082
DB_DRIVER=postgres
DB_SOURCE=postgresql://root:root@localhost:5432/db_name?sslmode=disable
KAFKA_BROKER=localhost:29092
GROUP_ID=service-group
```

## 📚 API Documentation

API documentation is available at:
- Orchestra Service: `http://localhost:8082/swagger/index.html`
- Order Service: `http://localhost:8081/swagger/index.html`
- Payment Service: `http://localhost:8083/swagger/index.html`
- Product Service: `http://localhost:8084/swagger/index.html`
- User Service: `http://localhost:8085/swagger/index.html`

## 🤝 Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## 📄 License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.
```

This README.md follows industry best practices by including:

1. Clear project description and architecture overview
2. Technology stack details
3. Setup and installation instructions
4. Running instructions
5. Development guidelines
6. Project structure
7. Environment configuration
8. API documentation references
9. Contributing guidelines
10. License information

The structure is clean, well-organized, and provides all necessary information for developers to get started with the project. The use of emojis makes it more visually appealing and helps in quick section identification.