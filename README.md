# Snippetbox

<div align="center">

![Go](https://img.shields.io/badge/Go-1.24%2B-00ADD8?style=for-the-badge&logo=go&logoColor=white)
![MySQL](https://img.shields.io/badge/MySQL-8.0%2B-4479A1?style=for-the-badge&logo=mysql&logoColor=white)
![HTML](https://img.shields.io/badge/HTML5-E34F26?style=for-the-badge&logo=html5&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=for-the-badge&logo=docker&logoColor=white)
![License](https://img.shields.io/badge/License-MIT-green?style=for-the-badge)
![Status](https://img.shields.io/badge/Status-Production--ready-success?style=for-the-badge)

<p align="center">
  A modern Go web application for creating, sharing, and managing code snippets.
</p>

</div>

## About this repository

Snippetbox is a practical web application built with Go and MySQL for storing and displaying reusable code snippets. The project follows the classic MVC-style structure and demonstrates secure user authentication, database-backed persistence, form validation, and production-ready web patterns.

It is designed as a clean example of building a full-featured web app with Go: server-side rendering, session management, CSRF protection, user accounts, and a structured template layer.

## Tech stack

- Go 1.24+
- MySQL 8+
- HTML templates
- `httprouter` for routing
- `scs` session management
- `bcrypt` password hashing
- Docker support for local setup

## Features

- Create and view snippets with expiration dates
- Latest snippets homepage with dynamic listing
- User sign up and login flows
- Secure password hashing and session authentication
- CSRF protection and request validation
- Template-driven UI with reusable layout partials
- MySQL-backed persistence with automatic timestamps

## Project structure

```text
snippetbox/
├── cmd/
│   └── web/                  # HTTP server, routes, handlers, middleware
├── internal/
│   ├── models/               # Domain models and DB access logic
│   ├── validator/            # Input validation helpers
│   └── assert/              # Test utility helpers
├── ui/
│   ├── html/                 # Page templates and shared layout
│   └── static/              # CSS, JS, and static assets
├── tls/                      # Local TLS certificate files
├── Dockerfile                # Container configuration
├── docker-compose.yml        # Local Docker orchestration
├── go.mod                    # Go module definition
├── go.sum                    # Dependency lock file
├── .env                      # Local environment variables
├── LICENSE                   # Project license
├── README.md                 # Project documentation
└── sql/
    └── schema.sql            # Database schema
```

## Prerequisites

- Go 1.24 or newer
- MySQL 8 or newer
- Docker (optional, for local containerized setup)

## Quick start

1. Create a MySQL database and user.
2. Configure the database connection in the environment or pass it at runtime.
3. Run the application:

```bash
go run ./cmd/web -dsn="user:password@tcp(localhost:3306)/snippetbox?parseTime=true"
```

Then open:

```text
http://localhost:4000
```

## Environment

Set the database connection using `DB_DSN`:

```bash
export DB_DSN="user:password@tcp(localhost:3306)/snippetbox?parseTime=true"
```

You can also run in debug mode:

```bash
go run ./cmd/web -debug
```

## Application highlights

- Secure session-based authentication
- Clean Go MVC patterns and reusable service layers
- Simple, maintainable template rendering
- Practical foundation for a production-ready snippet platform

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.
