FROM golang:1.22-alpine AS builder

WORKDIR /app
COPY . .

RUN go mod download
RUN go build -o server main.go handlers.go prometheus.go rateLimits.go storage.go

FROM alpine
WORKDIR /app
COPY --from=builder /app/server .

EXPOSE 8080
CMD ["./server"]
