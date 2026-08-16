# Build stage
FROM golang:1.26 AS builder
WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /app/ordersystem ./cmd/ordersystem

# Runtime stage
FROM scratch
WORKDIR /app

COPY --from=builder /app/ordersystem /app/ordersystem
COPY --from=builder /app/sql /app/sql

EXPOSE 3000 3001 3002

ENTRYPOINT ["/app/ordersystem"]