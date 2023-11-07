# syntax=docker/dockerfile:1

FROM golang:1.21-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY cmd/limiter/ cmd/limiter/
COPY internal/ internal/

RUN GOOS=linux go build -o /limiter ./cmd/limiter/

FROM alpine:3.18

COPY --from=build /limiter /limiter
COPY ./configs/limiter.prod.yaml /etc/limiter/limiter.prod.yaml

EXPOSE 50051

CMD ["/limiter", "-config", "/etc/limiter/limiter.prod.yaml"]
