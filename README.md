# Mail Scheduler

A Go-based REST API service for scheduling and sending emails using MailerSend. This application allows users to register, login, create greetings, and schedule emails to be sent at specific times.

## Features

- User authentication (JWT-based)
- Greeting management (CRUD operations)
- Email scheduling
- Background email worker for processing scheduled emails
- PostgreSQL database integration
- Environment-based configuration

## Tech Stack

- Go 1.23.2
- Chi Router
- PostgreSQL (via pgx)
- JWT for authentication
- MailerSend for email delivery
- Environment variables management

## Prerequisites

- Go 1.23.2 or higher
- PostgreSQL database
- MailerSend API key
- Environment variables configured (see Configuration section)

## Installation

1. Clone the repository:

```bash
git clone https://github.com/jaykapade/mail-scheduler.git
cd mail-scheduler
```

2. Install dependencies:

```bash
go mod download
```

3. Set up your environment variables in a `.env` file (see Configuration section)

4. Run database migrations:

```bash
# Install goose if you haven't already
go install github.com/pressly/goose/v3/cmd/goose@latest

# Run migrations
goose -dir migrations postgres "user=postgres password=postgres dbname=mail_scheduler sslmode=disable" up
```

5. Start the server:

```bash
go run cmd/main.go
```

## Configuration

Create a `.env` file in the root directory with the following variables:

```
# Database configuration
DB_HOST=
DB_PORT=
DB_USER=
DB_PASSWORD=
DB_NAME=

# JWT configuration
JWT_SECRET=

# MailerSend configuration
MAILERSEND_API_KEY=
```

## API Endpoints

### Authentication

- `POST /register` - Register a new user
- `POST /login` - Login and get JWT token

### Greetings

- `GET /greetings` - List all greetings
- `POST /greetings` - Create a new greeting
- `PUT /greetings/{id}` - Update a greeting
- `DELETE /greetings/{id}` - Delete a greeting

### Schedules

- `GET /schedules` - List all schedules
- `POST /schedules` - Create a new schedule

## Development

The project follows a modular structure:

- `cmd/` - Main application entry point
- `internal/` - Internal packages
  - `auth/` - Authentication related code
  - `db/` - Database operations
  - `greetings/` - Greeting management
  - `schedules/` - Schedule management
- `migrations/` - Database migrations
