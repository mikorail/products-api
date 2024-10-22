# Products API

This is a RESTful API built with **Golang** using the **Gin** framework. The API allows you to manage customers, products, orders, and categories. It features full CRUD (Create, Read, Update, Delete) operations for each resource and is integrated with **PostgreSQL** for database storage and **Redis** for caching. The repository includes Swagger documentation for API reference as well as DDL (Data Definition Language) and DML (Data Manipulation Language) scripts for setting up the database.

## Features

- CRUD operations for:
  - Customers
  - Products
  - Orders
  - Categories
- **PostgreSQL** as the relational database
- **Redis** for caching
- **Swagger** documentation for easy API interaction
- **Gin** web framework for routing and handling HTTP requests

## Table of Contents

- [Getting Started](#getting-started)
- [API Documentation](#api-documentation)
- [Project Structure](#project-structure)
- [Database Setup](#database-setup)
  - [DDL](#ddl)
  - [DML](#dml)
- [Environment Variables](#environment-variables)
- [Running the Application](#running-the-application)
- [Contributing](#contributing)
- [License](#license)

## Getting Started

### Prerequisites

- **Go** (version 1.18+ recommended)
- **PostgreSQL** (version 12+ recommended)
- **Redis** (version 6+ recommended)
- **Swagger** (optional but recommended for interactive API docs)

### Installing Dependencies

Run the following command to install the Go modules:

```bash
go mod tidy
```

## API Documentation

This project includes **Swagger** for easy API documentation and testing. Once the application is running, you can access the Swagger UI at:

```
http://localhost:<port>/swagger/index.html
```

## Project Structure

```bash
├── controllers      # Contains API route handlers (Customers, Orders, Products, Categories)
├── models           # Contains data models that map to PostgreSQL tables
├── services         # Business logic and interaction with the database
├── helpers          # Utility functions (response handling, validation, etc.)
├── database         # Database connection setup (PostgreSQL and Redis)
├── migrations       # Contains DDL and DML scripts for setting up the database
├── docs             # Swagger documentation setup
└── main.go          # Entry point for the application
```

## Database Setup

This project uses PostgreSQL for data persistence. The `migrations/` folder contains scripts for creating and populating the database tables.

### DDL

The **DDL** script (Data Definition Language) creates the necessary tables for the API, including:

- `customers`
- `products`
- `orders`
- `categories`

To execute the DDL script, run the following command in your PostgreSQL instance:

```sql
psql -U <username> -d <dbname> -f migrations/ddl.sql
```

### DML

The **DML** script (Data Manipulation Language) inserts sample data into the tables for testing purposes.

```sql
psql -U <username> -d <dbname> -f migrations/dml.sql
```

## Environment Variables

The application uses environment variables for configuration. You can set them manually or use a `.env` file.

Here’s a list of required environment variables and their defaults:

```env
DB_HOST=localhost
DB_PORT=5432
DB_USER=your_db_user
DB_PASSWORD=your_db_password
DB_NAME=your_db_name
DB_SSLMODE=disable

REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
```

## Running the Application

1. Set up PostgreSQL and Redis locally or in Docker.
2. Configure your `.env` file with the appropriate credentials.
3. Run the following command to start the server:

```bash
go run main.go
```

The application should be running on `http://localhost:8080` by default.

### Running Tests

You can run the tests using the command:

```bash
go test ./...
```

## Contributing

Contributions are welcome! Please submit a pull request or open an issue to improve the project.

## License

This project is licensed under the MIT License.

---

Feel free to modify the `README.md` to fit your specific project needs.