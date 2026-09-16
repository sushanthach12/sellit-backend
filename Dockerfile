# ---- Build stage ----
FROM golang:1.27.1-alpine3.24 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -tags netgo -ldflags '-s -w' -o bin/api ./cmd/api

# ---- Run stage ----
FROM alpine:3.24

WORKDIR /app

RUN apk add --no-cache ca-certificates

COPY --from=builder /app/bin/api .

EXPOSE 8000

CMD ["./api"]