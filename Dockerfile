# Build stage
FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.24.0

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN sqlc generate && \
    go mod tidy && \
    go mod verify

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 \
    go build \
    -ldflags="-s -w -X main.Version=$(git describe --tags --always)" \
    -trimpath \
    -o api-golang .

FROM gcr.io/distroless/static-debian12:nonroot AS production

ENV GIN_MODE=release \
    PORT=8000

COPY --from=builder --chown=nonroot:nonroot /app/api-golang /api-golang

HEALTHCHECK --interval=30s --timeout=3s \
    CMD ["/api-golang", "healthz"]

USER nonroot:nonroot
EXPOSE $PORT
ENTRYPOINT ["/api-golang"]