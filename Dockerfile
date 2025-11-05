FROM golang:1.24.2-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git build-base

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o subscription-service ./cmd/app
RUN go build -o migrate ./cmd/migrate

FROM alpine:3.20

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/subscription-service .
COPY --from=builder /app/migrate .
COPY ./config/.env ./config/.env
COPY ./migrations ./migrations

EXPOSE 8080

CMD ["./subscription-service"]
