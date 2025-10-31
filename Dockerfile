## Build container
FROM golang:1.24-alpine AS builder

RUN apk add --no-cache git ca-certificates
WORKDIR /build
COPY src/go.mod src/go.sum ./
RUN go mod download
COPY src/ ./
RUN GOOS=linux GOARCH=amd64 go build \
    -ldflags="-s -w" \
    -o mtg-mcp .

## Runtime container
FROM alpine:latest
RUN apk --no-cache add ca-certificates
RUN addgroup -g 1000 mtgmcp && \
    adduser -D -u 1000 -G mtgmcp mtgmcp

WORKDIR /app
COPY --from=builder /build/mtg-mcp /app/mtg-mcp
RUN chown -R mtgmcp:mtgmcp /app
USER mtgmcp

ENV MCP_SERVER_NAME="scryfall-card-search-server" \
    MCP_SERVER_VERSION="v1.0.0" \
    MCP_LOG_TO_FILE="false"

ENV MCP_SSE_HOST="0.0.0.0" \
    MCP_SSE_PORT="3000" \
    MCP_SSE_PATH="/sse"

EXPOSE 3000

HEALTHCHECK --interval=30s --timeout=3s --start-period=5s --retries=3 \
    CMD [ -x /app/mtg-mcp ] || exit 1

ENTRYPOINT ["/app/mtg-mcp"]
