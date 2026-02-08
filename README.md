# About

- CRM platform / Mailing list manager
- Written in Go
- Main functionalities are mainly controlled through REST-API server

# How to run

## Start docker-compose to start mysql and smtp4dev

- mysql is a database for storing the application data.
- smtp4dev is a SMTP server for testing email. You can access it at http://localhost:4999 to view sent emails.

```sh
docker compose up -d
```

## Database migration

### Install goose cli globally (only once)

```sh
go install github.com/pressly/goose/v3/cmd/goose@latest
```

### Start database migration

```sh
goose -dir migrations mysql "admin:admin@tcp(localhost:3309)/crm-go" up
```

## Environment variables for API server

Please set up `.env` following `.env.example`

## Start API Server

This is REST-API server that you can make a REST http request to send email, manage contact; etc.

```sh
go run cmd/api/main.go
```

## Test with mock data

For local development and testing, you can use this cmd to quickly create mock data. See `cmd/createmock/main.go` for more info.

```sh
go run cmd/createmock/main.go
```

## Start Mock API (OPTIONAL)

Rendering email templates can make external API calls to get data. You can use this Mock API for local development.

```sh
go run cmd/mockapi/main.go
```

## Curl

These are some curl commands for testing.

- Send specific email

```sh
curl -X POST --location 'http://localhost:8080/emails/-100/send'
```
