# ── Build Stage ────────────────────────────────────────────────────
FROM golang:1.23-alpine AS builder

WORKDIR /app

# Cache dependency downloads.
COPY go.mod go.sum ./
RUN go mod download

# Copy source and build a static binary.
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o /ainyx ./cmd/server

# ── Runtime Stage ─────────────────────────────────────────────────
FROM alpine:3.19

RUN apk --no-cache add ca-certificates

WORKDIR /app
COPY --from=builder /ainyx .

EXPOSE 3000

ENTRYPOINT ["./ainyx"]
