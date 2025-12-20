# About

CRM platform written in Go

# How to run

## Mock SMTP server

```bash
docker compose up -d
```

Access at http://localhost:4999

## Mock API

```bash
go run cmd/mockapi/main.go
```

## Send email

### Send email

```bash
go run cmd/sendemail/main.go -debug=true -smtp-host=localhost -smtp-port=25 -from-addr=from@test.com -to-addr=to@test.com -subject='Subject test' -api-url="http://localhost:8080/" -body-template='Hello {{.name}}. You are from {{ .extraData.location }} !'
```

### Check usage

```bash
go run cmd/sendemail/main.go -h
```
