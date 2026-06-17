# Shinkyu Shotokan Karate

A Go web application for the Shinkyu Shotokan karate dojo in South San Francisco.

## Prerequisites

- [Go](https://go.dev/) 1.25+
- [Docker](https://www.docker.com/) + [Docker Compose](https://docs.docker.com/compose/)
- [git](https://git-scm.com/)

## Local Setup

### 1. Clone and install dependencies

```
git clone https://github.com/pappygobappy/shinkyuShotokanWebsite.git
cd shinkyuShotokanWebsite
go mod download
```

### 2. Start the database

```
docker compose up -d
```

This spins up a PostgreSQL instance using the credentials from `.env`. The database will be available on `localhost:5432`.

### 3. Configure environment

Copy `.env.example` and adjust values:

```
cp .env.example .env
```

Required variables:

| Variable | Default | Description |
|---|---|---|
| `DB_USERNAME` | `citizix_user` | PostgreSQL username |
| `DB_PASSWORD` | `S3cret` | PostgreSQL password |
| `DB_HOST` | `localhost` | Database host |
| `DB_PORT` | `5432` | Database port |
| `DB_NAME` | `citizix_db` | Database name |
| `PORT` | `8080` | HTTP server port |
| `HMAC_SECRET` | *(64-char hex)* | JWT signing key |
| `UPLOAD_DIR` | `./upload` | File upload directory |
| `SMTP_*` | *(Gmail)* | Email sending credentials |

### 4. Run the app

```
go run main.go
```

The server starts on port **8080** (or whatever `PORT` is set to). Database migrations and seed data are applied automatically on startup via `initializers/syncDb.go`.

## Stopping the database

```
docker compose down
```

Use `-v` to also remove persisted data (resets the database):

```
docker compose down -v
```
