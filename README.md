# About

- CRM platform / Mailing list manager
- Written in Go
- Main functionalities are mainly controlled through REST-API server

# How to run

## Start docker-compose to start mysql and smtp4dev

```bash
docker compose up -d
```

## Database migration

```bash
goose -dir migrations mysql "admin:admin@tcp(localhost:3309)/crm-go" up
```

Access at http://localhost:4999

## Mock API

```bash
go run cmd/mockapi/main.go
```

## API Server

```bash
go run cmd/api/main.go
```

## Test with mock data

See `cmd/createmock/main.go`
