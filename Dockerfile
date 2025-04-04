FROM golang:1.24.1-alpine AS builder

WORKDIR /app

RUN apk add --no-cache git gcc musl-dev libc-dev
RUN go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .

ENV CGO_ENABLED=1
RUN sqlc generate && \
    go mod tidy && \
    go mod verify

# Build a statically linked binary
RUN CGO_ENABLED=1 GOOS=linux \
    go build \
    -ldflags="-linkmode external -extldflags -static" \
    -trimpath \
    -o go-rss .

FROM scratch AS production

ENV GIN_MODE=release \
    PORT=8000

COPY --from=builder /app/go-rss /go-rss

EXPOSE 8000

CMD ["/go-rss"]