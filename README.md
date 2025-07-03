# Internal Transfer System

A Go-based financial transaction system that facilitates internal transfers between accounts using PostgreSQL as the database.

## Features

- Account management with balance tracking
- Internal transfers between accounts
- Transaction logging and status tracking
- PostgreSQL database with proper indexing
- RESTful HTTP API (coming in next iteration)

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (or use provided Docker setup)

## Installation and Setup

### 1. Clone the Repository

```bash
git clone <repository-url>
cd internal-transfer-system
```

### 2. Start PostgreSQL Database

Using Docker Compose (recommended):

```bash
docker-compose up -d postgres
```

This will start a PostgreSQL container with:
- Database: `internal_transfer`
- Username: `postgres`
- Password: `postgres`
- Port: `5432`

### 3. Install Dependencies

```bash
go mod download
```

### 4. Run the Application

```bash
go run cmd/main.go
```

The application will:
- Connect to the PostgreSQL database
- Create necessary tables (`accounts` and `transactions`)
- Set up proper indexes and triggers

## Database Schema

### Accounts Table
- `account_id` (BIGINT, Primary Key)
- `balance` (DECIMAL(20,8))
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

### Transactions Table
- `transaction_id` (BIGSERIAL, Primary Key)
- `source_account_id` (BIGINT, Foreign Key)
- `destination_account_id` (BIGINT, Foreign Key)
- `amount` (DECIMAL(20,8))
- `status` (VARCHAR(20))
- `created_at` (TIMESTAMP)
- `updated_at` (TIMESTAMP)

## Environment Variables

The application supports the following environment variables:

- `DB_HOST` (default: localhost)
- `DB_PORT` (default: 5432)
- `DB_USER` (default: postgres)
- `DB_PASSWORD` (default: postgres)
- `DB_NAME` (default: internal_transfer)
- `DB_SSL_MODE` (default: disable)

## Development Status

This is the first checkpoint of the development process. The following features are implemented:

✅ Database connection and configuration
✅ Database schema with proper relationships
✅ Account and Transaction models
✅ Database table creation and indexing
✅ Docker setup for PostgreSQL

**Next Iteration:**
- HTTP API endpoints
- Account creation and querying
- Transaction processing
- Error handling and validation
- Unit tests

## Project Structure

```
internal-transfer-system/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   ├── connection.go       # Database connection management
│   │   └── schema.go          # Database schema and migrations
│   └── model/
│       ├── account.go         # Account model and DTOs
│       └── transaction.go     # Transaction model and DTOs
├── docker-compose.yml         # PostgreSQL setup
├── go.mod                     # Go module dependencies
└── README.md                  # This file
```

## Assumptions

- All accounts use the same currency
- Account IDs are provided by the client
- Initial implementation focuses on core functionality
- Authentication and authorization will be added in future iterations