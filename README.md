# Internal Transfer System

A Go-based financial transaction system that facilitates internal transfers between accounts using PostgreSQL as the database.

## Features

- ✅ Account creation with initial balance
- ✅ Account balance queries
- ✅ Internal transfers between accounts
- ✅ Transaction logging and status tracking
- ✅ PostgreSQL database with proper indexing
- ✅ RESTful HTTP API with JSON responses
- ✅ Data integrity with database transactions
- ✅ Comprehensive error handling
- ✅ Graceful server shutdown

## Prerequisites

- Go 1.21 or higher
- Docker and Docker Compose
- PostgreSQL (or use provided Docker setup)
- curl and jq (for testing)

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
- Start the HTTP server on port 8080

## API Endpoints

### Base URL
```
http://localhost:8080
```

### 1. Create Account

**POST** `/accounts`

Creates a new account with an initial balance.

**Request Body:**
```json
{
  "account_id": 123,
  "initial_balance": "100.23344"
}
```

**Success Response:**
- Status: `201 Created`
- Body: Empty

**Error Responses:**
- `400 Bad Request` - Invalid request format or business logic error
- `500 Internal Server Error` - Database or server error

### 2. Get Account Balance

**GET** `/accounts/{account_id}`

Retrieves the account balance for the specified account.

**Success Response:**
- Status: `200 OK`
- Body:
```json
{
  "account_id": 123,
  "balance": "100.23344"
}
```

**Error Responses:**
- `404 Not Found` - Account does not exist
- `400 Bad Request` - Invalid account ID format
- `500 Internal Server Error` - Database or server error

### 3. Create Transaction

**POST** `/transactions`

Creates a transaction to transfer money between accounts.

**Request Body:**
```json
{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "100.12345"
}
```

**Success Response:**
- Status: `201 Created`
- Body: Empty

**Error Responses:**
- `400 Bad Request` - Invalid request format, insufficient balance, or business logic error
- `500 Internal Server Error` - Database or server error

### 4. Health Check

**GET** `/health`

Returns the health status of the service.

**Success Response:**
- Status: `200 OK`
- Body:
```json
{
  "status": "healthy",
  "service": "internal-transfer-system"
}
```

## Testing the API

### Using the Test Script

A test script is provided to verify all API endpoints:

```bash
chmod +x test_api.sh
./test_api.sh
```

### Manual Testing Examples

#### 1. Create accounts:
```bash
curl -X POST "http://localhost:8080/accounts" \
  -H "Content-Type: application/json" \
  -d '{"account_id": 123, "initial_balance": "100.50"}'

curl -X POST "http://localhost:8080/accounts" \
  -H "Content-Type: application/json" \
  -d '{"account_id": 456, "initial_balance": "200.75"}'
```

#### 2. Check balances:
```bash
curl "http://localhost:8080/accounts/123"
curl "http://localhost:8080/accounts/456"
```

#### 3. Create transaction:
```bash
curl -X POST "http://localhost:8080/transactions" \
  -H "Content-Type: application/json" \
  -d '{"source_account_id": 123, "destination_account_id": 456, "amount": "25.25"}'
```

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
- `PORT` (default: 8080)

## Architecture

The application follows a clean architecture pattern:

### Layers:
1. **Handler Layer** - HTTP request/response handling
2. **Service Layer** - Business logic and validation
3. **Repository Layer** - Database operations
4. **Model Layer** - Data structures and DTOs

### Key Features:
- **Data Integrity** - Uses database transactions with row-level locking
- **Error Handling** - Comprehensive error handling with appropriate HTTP status codes
- **Validation** - Input validation and business rule enforcement
- **Logging** - Structured logging for debugging and monitoring
- **Graceful Shutdown** - Proper server shutdown handling

## Project Structure

```
internal-transfer-system/
├── cmd/
│   └── main.go                 # Application entry point
├── internal/
│   ├── database/
│   │   ├── connection.go       # Database connection management
│   │   └── schema.go          # Database schema and migrations
│   ├── model/
│   │   ├── account.go         # Account model and DTOs
│   │   └── transaction.go     # Transaction model and DTOs
│   ├── repository/
│   │   ├── account_repository.go    # Account data access
│   │   └── transaction_repository.go # Transaction data access
│   ├── service/
│   │   ├── account_service.go      # Account business logic
│   │   └── transaction_service.go  # Transaction business logic
│   ├── handler/
│   │   ├── account_handler.go      # Account HTTP handlers
│   │   └── transaction_handler.go  # Transaction HTTP handlers
│   ├── router/
│   │   └── router.go              # HTTP router setup
│   └── utils/
│       └── string_utils.go        # Utility functions
├── docker-compose.yml         # PostgreSQL setup
├── go.mod                     # Go module dependencies
├── test_api.sh               # API test script
└── README.md                 # This file
```

## Data Integrity & Consistency

The system ensures data integrity through:

1. **Database Transactions** - All account balance updates are wrapped in database transactions
2. **Row-Level Locking** - Uses `FOR UPDATE` to prevent concurrent balance modifications
3. **Validation** - Comprehensive input validation and business rule enforcement
4. **Atomic Operations** - Either all operations in a transaction succeed or all fail
5. **Referential Integrity** - Foreign key constraints ensure data consistency

## Error Handling

The system provides comprehensive error handling:

- **Input Validation** - Validates all input parameters
- **Business Rules** - Enforces business logic (e.g., sufficient balance, positive amounts)
- **Database Errors** - Handles database connection and query errors
- **HTTP Status Codes** - Returns appropriate HTTP status codes for different error types
- **Error Messages** - Provides clear error messages for debugging

## Assumptions

- All accounts use the same currency
- Account IDs are provided by the client and must be positive integers
- Initial balances and transaction amounts must be non-negative
- The system is designed for internal transfers only
- Authentication and authorization are not implemented (as per requirements)
- High precision decimal arithmetic is used for financial calculations

## Production Considerations

For production deployment, consider:

1. **Security** - Add authentication and authorization
2. **Monitoring** - Add metrics and health checks
3. **Scaling** - Consider database connection pooling and horizontal scaling
4. **Logging** - Implement structured logging with log levels
5. **Configuration** - Use configuration management for different environments
6. **Testing** - Add comprehensive unit and integration tests
7. **CI/CD** - Implement automated testing and deployment pipelines