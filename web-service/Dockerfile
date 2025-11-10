#web-service
FROM golang:1.24.9-bookworm AS builder

WORKDIR /app

# Copy dependency metadata first (leverages Docker layer caching)
COPY go.mod go.sum ./
RUN go mod download

# Copy your actual source
COPY . .

# Compile statically (or semi-statically) for Linux
RUN CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o /webservice

# --- Second stage: minimal runtime ---
FROM debian:bookworm-slim

WORKDIR /app

# Copy binary only
COPY --from=builder /webservice .

EXPOSE 5821
CMD ["./webservice"]
