# syntax=docker/dockerfile:1

ARG GO_VERSION="1.24"
ARG ALPINE_VERSION="3.21"

# Build stage
FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} AS build

RUN apk add --no-cache git ca-certificates tzdata

ARG TARGETOS=linux
ARG TARGETARCH=amd64

WORKDIR /src

COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod \
    go mod download

COPY . .
RUN --mount=type=cache,target=/go/pkg/mod \
    --mount=type=cache,target=/root/.cache/go-build \
    CGO_ENABLED=0 \
    GOOS=${TARGETOS} \
    GOARCH=${TARGETARCH} \
    go build \
    -ldflags="-s -w" \
    -trimpath \
    -o /bin/proompt \
    ./cmd/proompt

# Runtime stage
FROM alpine:${ALPINE_VERSION}

RUN apk add --no-cache \
    ca-certificates \
    tzdata \
    && addgroup -g 1001 -S proompt \
    && adduser -u 1001 -S proompt -G proompt

RUN mkdir -p /app/data /app/config \
    && chown -R proompt:proompt /app

COPY --from=build /bin/proompt /usr/local/bin/proompt
COPY --chown=proompt:proompt internal/db/migrations /app/migrations/
COPY --chown=proompt:proompt proompt.docker.xml /app/config/proompt.xml
COPY --chown=proompt:proompt proompt.example.xml /app/config/

USER proompt
WORKDIR /app
EXPOSE 8080

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD wget --no-verbose --tries=1 --spider http://localhost:8080/api/health || exit 1

CMD ["proompt", "-config", "/app/config/proompt.xml"]