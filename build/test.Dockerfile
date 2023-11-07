# syntax=docker/dockerfile:1

FROM golang:1.21-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY tests/ tests/
COPY internal/ internal/

CMD ["echo", "test container"]
