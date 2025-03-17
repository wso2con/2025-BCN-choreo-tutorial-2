# Personal Finance Tracker API

A RESTful API for tracking personal finances, specifically for managing bills and their items.

## Features

- Create, read, update, and delete bills
- Add, modify, and remove items from bills
- Automatic calculation of bill totals based on item prices and quantities
- Support for both MySQL and SQLite databases
- OpenAPI documentation

## Prerequisites

- Go 1.21 or higher
- MySQL (optional)
- SQLite (default)

## Installation

1. Clone the repository

```bash
git clone https://github.com/your-username/accounts.git
cd accounts
```

2. Install dependencies

```bash
go mod download
```

3. Configure the database

Create a `.env` file based on the `.env.example` file:

```bash
cp .env.example .env
```

Edit the `.env` file to configure your database preferences:

```env
# Database configuration
# Options: mysql, sqlite
DB_TYPE=sqlite

# MySQL settings (if DB_TYPE=mysql)
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=accounts

# SQLite settings (if DB_TYPE=sqlite)
DB_PATH=./accounts.db

# Server settings
PORT=8080
```

## Running the API

```bash
go run main.go
```

The API will be available at `http://localhost:8080/api/v1`

Swagger documentation is available at `http://localhost:8080/swagger/`

## API Endpoints

### Bills

- `GET /api/v1/bills` - Get all bills
- `POST /api/v1/bills` - Create a new bill
- `GET /api/v1/bills/{id}` - Get a bill by ID
- `PUT /api/v1/bills/{id}` - Update a bill
- `DELETE /api/v1/bills/{id}` - Delete a bill

## Sample Requests

### Create a bill

```bash
curl -X POST http://localhost:8080/api/v1/bills \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Grocery Shopping",
    "description": "Weekly groceries",
    "due_date": "2023-05-15",
    "paid": false,
    "items": [
      {
        "name": "Milk",
        "description": "1 gallon",
        "amount": 3.99,
        "quantity": 2
      },
      {
        "name": "Bread",
        "description": "Whole wheat",
        "amount": 2.49,
        "quantity": 1
      }
    ]
  }'
```

### Get all bills

```bash
curl -X GET http://localhost:8080/api/v1/bills
```

### Get a bill by ID

```bash
curl -X GET http://localhost:8080/api/v1/bills/1
```

## Database Configuration

The API supports both MySQL and SQLite databases. You can configure which one to use in the `.env` file:

### SQLite (default)

```env
DB_TYPE=sqlite
DB_PATH=./accounts.db
```

### MySQL

```env
DB_TYPE=mysql
DB_HOST=localhost
DB_PORT=3306
DB_USER=root
DB_PASSWORD=your-password
DB_NAME=accounts
```

## Development

### Build

```bash
go build -o accounts
```

### Run tests

```bash
go test ./...
```

## License

This project is licensed under the MIT License - see the LICENSE file for details.
