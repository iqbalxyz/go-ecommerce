# Go E-Commerce RESTful API

A modular, robust, and high-performance RESTful API for an e-commerce platform built with **Go**, **Fiber v3**, and **GORM** (MySQL). Designed with a clean layered architecture (Handler &rarr; Service &rarr; Repository), strict input validation, JWT-based authentication with Role-Based Access Control (RBAC), and ACID-compliant order transactions.

---

## 🛠 Tech Stack
![Go](https://img.shields.io/badge/Go-1.23+-00ADD8?logo=go&logoColor=white)
![Fiber](https://img.shields.io/badge/Fiber-v3-00ACD7)
![MySQL](https://img.shields.io/badge/MySQL-8.0-4479A1?logo=mysql&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green)

| Category | Technology | Description |
| :--- | :--- | :--- |
| **Language** | [Go](https://go.dev/)  v1.27+ | Fast, statically typed compiled programming language |
| **Web Framework** | [Fiber v3](https://github.com/gofiber/fiber/v3) (`v3.5.0`) | Express-inspired, high-performance web framework built on Fasthttp |
| **ORM** | [GORM](https://gorm.io/) (`v1.31.2`) | Feature-rich Object Relational Mapping library for Go |
| **Database Driver** | [GORM MySQL Driver](https://github.com/go-gorm/mysql) (`v1.6.0`) | MySQL driver for GORM |
| **Database** | [MySQL 8.0+](https://www.mysql.com/) | Relational database management system |
| **Authentication** | [golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) (`v5.3.1`) | JSON Web Token (JWT) implementation for Go |
| **Password Hashing**| [golang.org/x/crypto](https://pkg.go.dev/golang.org/x/crypto) | Industry-standard Bcrypt hashing algorithm |
| **Validation** | [go-playground/validator/v10](https://github.com/go-playground/validator) | Struct and field validation |
| **Environment** | [joho/godotenv](https://github.com/joho/godotenv) (`v1.5.1`) | Environment variable loader from `.env` |
| **Database Migrations** | [golang-migrate](https://github.com/golang-migrate/migrate) | CLI-driven database migration tooling |
| **Automation** | GNU Make | Build and task runner via `Makefile` |

---

## ✨ Features

- **Authentication & Authorization**:
  - User registration with field validation (unique email, name, secure password).
  - Secure password hashing using `bcrypt`.
  - JWT token generation (24-hour expiration) on successful login.
  - Role-Based Access Control (RBAC) with `user` and `admin` roles.
  - Authenticated user profile retrieval (`/api/v1/me`).

- **Product Catalog Management**:
  - Public product catalog browsing with pagination (`page`, `limit`).
  - Retrieve single product details by ID.
  - Admin-only product creation, updates, and soft/hard deletion.
  - SKU uniqueness and stock quantity tracking.

- **Shopping Cart**:
  - Persistent shopping cart per authenticated user.
  - Add items to cart with real-time stock availability and active status checks.
  - Update item quantities with validation.
  - Remove items from cart.
  - Dynamic calculation of subtotal, total cart price, and total item count.

- **Orders & Checkout**:
  - Atomic checkout process using database transactions.
  - Stock validation and automatic inventory deduction.
  - Product price and name snapshotting at the time of purchase to preserve order integrity against future product changes.
  - Cart clearing upon successful checkout.
  - Paginated list of user order history.
  - Itemized order detail retrieval with order status tracking.

- **Architecture & Quality**:
  - **Clean Layered Architecture**: Clear separation of concerns (Handler &rarr; Service &rarr; Repository).
  - **DTO Pattern**: Separation of internal models from API requests and responses.
  - **Standardized API Response**: Uniform JSON response envelope (`success`, `data`, `meta`, `error`, `errors`).
  - **Robust Middleware**: Global panic recovery, HTTP request logging, and CORS handling.

---

## 📁 Folder Structure

```text
go-ecommerce/
├── bin/                       # Compiled application binaries
├── cmd/
│   └── api/
│       └── main.go            # Application entry point & dependency wiring
├── internal/
│   ├── config/                # Environment variables and config loader
│   ├── database/              # MySQL connection pool setup using GORM
│   ├── dto/                   # Data Transfer Objects for requests and responses
│   │   ├── auth_dto.go
│   │   ├── cart_dto.go
│   │   ├── order_dto.go
│   │   ├── product_dto.go
│   │   └── user_dto.go
│   ├── errors/                # Centralized application domain error definitions
│   ├── handler/               # HTTP controllers / Fiber route handlers
│   │   ├── cart_handler.go
│   │   ├── order_handler.go
│   │   ├── product_handler.go
│   │   └── user_handler.go
│   ├── middleware/            # Custom HTTP middleware (Auth, Role, Logger, etc.)
│   │   ├── auth.go
│   │   └── role.go
│   ├── models/                # GORM database entity models
│   │   ├── cart.go
│   │   ├── order.go
│   │   ├── product.go
│   │   └── user.go
│   ├── repository/            # Data access layer interfacing with MySQL
│   │   ├── cart_repository.go
│   │   ├── order_repository.go
│   │   ├── product_repository.go
│   │   └── user_repository.go
│   ├── router/                # Route registrations and endpoint grouping
│   │   └── router.go
│   ├── service/               # Core business logic layer
│   │   ├── cart_service.go
│   │   ├── order_service.go
│   │   ├── product_service.go
│   │   └── user_service.go
│   └── utils/                 # Helpers (JWT, Bcrypt, Validator, JSON Envelope)
│       ├── jwt.go
│       ├── password.go
│       ├── response.go
│       └── validator.go
├── migrations/                # SQL migration files for golang-migrate
├── .env.example               # Example environment variables template
├── .gitignore                 # Git ignore configuration
├── go.mod                     # Go module definitions
├── go.sum                     # Go checksum records
├── Makefile                   # Automation commands for build, run, and migrate
└── README.md                  # Project documentation
```

---

## ⚙️ Environment Variables

Create a `.env` file in the root directory by copying the sample configuration:

```bash
cp .env.example .env
```

| Variable | Required | Default / Example | Description |
| :--- | :---: | :--- | :--- |
| `APP_PORT` | **Yes** | `3000` | Port number for the Fiber HTTP server |
| `APP_ENV` | No | `development` | Environment mode (`development`, `production`, `test`) |
| `DB_HOST` | **Yes** | `127.0.0.1` | MySQL database host address |
| `DB_PORT` | No | `3306` | MySQL database port |
| `DB_USER` | No | `root` | MySQL username |
| `DB_PASSWORD` | No | `""` | MySQL password |
| `DB_NAME` | **Yes** | `ecommerce_db` | MySQL database name |
| `DB_URL` | **Yes** (for migrations) | `mysql://root:secret@tcp(127.0.0.1:3306)/ecommerce_db` | Full DSN URL used by `golang-migrate` |
| `JWT_SECRET` | **Yes** | *(Min 32 characters)* | Secret key used to sign and verify JWT tokens |

> [!TIP]
> You can quickly generate a secure 32-character random string for `JWT_SECRET` using OpenSSL:
> ```bash
> openssl rand -base64 32
> ```

---

## 🚀 Setup and Run Instructions

### 1. Prerequisites
Ensure you have the following installed on your machine:
- **Go**: Version `1.27+` installed & added to `PATH`
- **MySQL**: Running instance (v8.0+)
- **golang-migrate**: Migration CLI tool ([Installation Guide](https://github.com/golang-migrate/migrate/tree/master/cmd/migrate))
- **Make**: (Optional, standard on Linux/macOS or via MinGW/Choco on Windows)

### 2. Clone and Install Dependencies
```bash
git clone https://github.com/iqbalxyz/go-ecommerce.git
cd go-ecommerce
go mod download
```

### 3. Setup the Database
Create the database in your MySQL instance:
```sql
CREATE DATABASE ecommerce_db CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;
```

### 4. Configure Environment Variables
Copy `.env.example` to `.env` and fill in your database credentials and a 32+ character `JWT_SECRET`:
```bash
cp .env.example .env
```

### 5. Run Database Migrations
Apply all migration files to set up the database tables:
```bash
make migrate-up
```
*Or using the migrate CLI directly:*
```bash
migrate -path ./migrations -database "mysql://user:password@tcp(127.0.0.1:3306)/ecommerce_db" up
```

### 6. Run the Application
Start the development server:
```bash
make run
```
*Or run directly with Go:*
```bash
go run ./cmd/api
```
The server will start listening at: `http://localhost:3000`

### 7. Build for Production
To compile a standalone binary to the `./bin` directory:
```bash
make build
```
Run the generated binary:
- **Windows**: `.\bin\api.exe`
- **Linux/macOS**: `./bin/api`

### 8. Run Unit Tests
```bash
make test
```

### Handy Makefile Commands
| Command | Description |
| :--- | :--- |
| `make run` | Run the application with `go run ./cmd/api` |
| `make build` | Compile the application into `./bin/api` |
| `make test` | Run all test suites with verbose output (`go test ./... -v`) |
| `make tidy` | Run `go mod tidy` to clean dependencies |
| `make migrate-up` | Apply all pending migrations |
| `make migrate-down` | Rollback the latest migration step |
| `make migrate-reset` | Rollback all migrations and re-apply from scratch |
| `make migrate-version`| Show current database migration version |
| `make migrate-create NAME=<name>` | Generate new `.up.sql` and `.down.sql` migration pair |

---

## 📡 List of Endpoints

All routes are prefixed with `/api/v1`.

### Response Format
All responses adhere to a consistent JSON structure:
- **Standard Success**: `{"success": true, "data": ...}`
- **Paginated Success**: `{"success": true, "data": [...], "meta": {"page": 1, "limit": 10, "total": 45, "total_pages": 5}}`
- **Error**: `{"success": false, "error": "error message"}`
- **Validation Error**: `{"success": false, "error": "validation failed", "errors": {"field": "error reason"}}`

### Endpoints Table

| Category | Method | Endpoint | Access | Description | Request Body / Query |
| :--- | :---: | :--- | :---: | :--- | :--- |
| **Auth** | `POST` | `/api/v1/auth/register` | Public | Register a new user | `{ "name": "John Doe", "email": "john@example.com", "password": "password123" }` |
| **Auth** | `POST` | `/api/v1/auth/login` | Public | Authenticate user & get JWT token | `{ "email": "john@example.com", "password": "password123" }` |
| **User** | `GET` | `/api/v1/me` | User / Admin | Get current authenticated user profile | *Bearer Token required* |
| **Products** | `GET` | `/api/v1/products` | Public | List products (paginated) | Query: `?page=1&limit=10` |
| **Products** | `GET` | `/api/v1/products/:id` | Public | Get product details by ID | Path param: `:id` |
| **Products** | `POST` | `/api/v1/products` | Admin | Create a new product | `{ "name": "Keyboard", "description": "Mechanical", "price": 750000, "stock": 25, "sku": "KB-001", "is_active": true }` |
| **Products** | `PUT` | `/api/v1/products/:id` | Admin | Update existing product | `{ "name": "Keyboard Pro", "price": 850000, "stock": 20 }` |
| **Products** | `DELETE`| `/api/v1/products/:id` | Admin | Delete product | Path param: `:id` |
| **Cart** | `GET` | `/api/v1/cart` | User / Admin | View current user's shopping cart | *Bearer Token required* |
| **Cart** | `POST` | `/api/v1/cart/items` | User / Admin | Add item to shopping cart | `{ "product_id": 1, "quantity": 2 }` |
| **Cart** | `PUT` | `/api/v1/cart/items/:id` | User / Admin | Update quantity of a cart item | `{ "quantity": 5 }` |
| **Cart** | `DELETE`| `/api/v1/cart/items/:id` | User / Admin | Remove item from cart | Path param: `:id` (cart item ID) |
| **Orders** | `POST` | `/api/v1/orders` | User / Admin | Checkout cart into a new order | *Bearer Token required* (creates order from active cart) |
| **Orders** | `GET` | `/api/v1/orders` | User / Admin | List user's orders (paginated) | Query: `?page=1&limit=10` |
| **Orders** | `GET` | `/api/v1/orders/:id` | User / Admin | Get specific order with item snapshot | Path param: `:id` |

> [!NOTE]
> To access protected endpoints, supply the JWT token in the `Authorization` header:
> ```http
> Authorization: Bearer <your_jwt_token>
> ```

---

## 🗄 Screenshots

### 1. REST API Clients Interaction
![Request/Response](assets/image.png)

### 2. Server Log
![alt text](image-1.png)

---

## 📄 License

This project is licensed under the [MIT License](https://opensource.org/licenses/MIT). You are free to use, modify, and distribute this software in accordance with the license conditions.
