FROM golang:1.23-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY shared/ ./shared/
COPY services/inventory-service/ ./services/inventory-service/
WORKDIR /app/services/inventory-service
RUN CGO_ENABLED=0 GOOS=linux go build -o /inventory-service ./cmd/main.go

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/
COPY --from=builder /inventory-service .
EXPOSE 8006
CMD ["./inventory-service"]
