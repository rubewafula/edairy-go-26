# syntax=docker/dockerfile:1

FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git ca-certificates

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-s -w" -o /api ./cmd/api

FROM alpine:3.20

RUN apk add --no-cache ca-certificates tzdata wget \
    && addgroup -S edairy && adduser -S edairy -G edairy

WORKDIR /app

COPY --from=builder /api ./api

RUN mkdir -p uploads/members uploads/transporters storage/exports \
    && chown -R edairy:edairy /app

USER edairy

ENV GIN_MODE=release
EXPOSE 8000

HEALTHCHECK --interval=30s --timeout=5s --retries=3 \
    CMD wget -qO- http://127.0.0.1:${PORT:-8000}/health || exit 1

CMD ["./api"]
