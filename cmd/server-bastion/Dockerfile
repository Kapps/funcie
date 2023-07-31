# syntax=docker/dockerfile:1

FROM golang:1.20-alpine AS base

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

RUN go mod download

COPY . .

RUN GOOS=linux go build -o /app/server-bastion ./cmd/server-bastion/main.go

FROM alpine:latest AS final

RUN apk --no-cache add ca-certificates

COPY --from=base /app /app

CMD ["/app/server-bastion"]