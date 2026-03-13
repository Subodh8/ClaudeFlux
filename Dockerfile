# ---- Build Stage ----
FROM golang:1.23-alpine AS builder

RUN apk add --no-cache git gcc musl-dev sqlite-dev

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=1 go build \
    -ldflags "-X main.version=$(git describe --tags --always --dirty 2>/dev/null || echo dev) -s -w" \
    -o /bin/claudeflux \
    ./cmd/claudeflux

# ---- Runtime Stage ----
FROM alpine:3.19

RUN apk add --no-cache git sqlite-libs ca-certificates

# Create non-root user
RUN addgroup -S claudeflux && adduser -S claudeflux -G claudeflux

COPY --from=builder /bin/claudeflux /usr/local/bin/claudeflux

# Default state directory
RUN mkdir -p /data/state && chown claudeflux:claudeflux /data/state
VOLUME ["/data/state"]

USER claudeflux
WORKDIR /workspace

EXPOSE 7070 7071

ENTRYPOINT ["claudeflux"]
CMD ["--help"]
