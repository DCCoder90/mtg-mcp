# MTGMCP

This is a simple MCP server for the [Scryfall](https://scryfall.com/docs/api) to assist with card data and deck building.  This allows clients such as [Claude Desktop](https://claude.com) to search for Magic: The Gathering card details using the Scryfall API.

## Description

This server exposes a number of MCP tools to assist with MTG card search and deck building.  The tools and their purposes are listed below.

## Tools

### `search_card_by_name`

This tool takes a card name as input, queries the Scryfall API for cards matching that exact name, and returns the details of the found card(s) in a structured format.

### `search_card_by_color`

This tool takes a card color as input, queries the Scryfall API for cards with that color, and returns the details of the found card(s) in a structured format.

### `search_card_by_text`

This tool takes a card's text as input, queries the Scryfall API for cards with that text or similar text, and returns the details of the found card(s) in a structured format.

### `find_related_cards`

This tool finds cards related to a specified card through various relationship types:
- **Reprints**: Other printings of the same card 
- **Tokens**: Token cards created by the card
- **Mechanics**: Cards sharing similar keyword abilities
- **Same Set**: Other cards from the same expansion

### `find_card_synergies`

This advanced tool analyzes a card and finds synergistic cards for deck building:
- **Keyword Synergies**: Cards sharing keyword abilities (flying, trample, etc.)
- **Theme Synergies**: Cards fitting strategic themes (sacrifice, tokens, graveyard, counters, etc.)
- **Tribal Synergies**: Cards sharing creature types
- **Color Identity**: Cards matching color requirements

The tool automatically extracts themes from card text.

## Installation

### Download Pre-built Binaries

Pre-built binaries are available for Windows, Linux on the [Releases page](https://github.com/DCCoder90/mtg-mcp/releases).

#### Available Platforms:
- **Windows**: `mtg-mcp-windows-amd64.zip` (x64) / `mtg-mcp-windows-arm64.zip` (ARM64)
- **Linux**: `mtg-mcp-linux-amd64.tar.gz` (x64) / `mtg-mcp-linux-arm64.tar.gz` (ARM64)

#### Installation Steps:

**Windows:**
1. Download the appropriate `.zip` file for your architecture
2. Extract the archive to your desired location

**Linux:**
1. Download the appropriate `.tar.gz` file for your architecture
2. Extract: `tar -xzf mtg-mcp-*.tar.gz`
3. Make executable: `chmod +x mtg-mcp-*`

## Development
### Prerequisites

* Go programming language (version 1.24 or later recommended for `go-sdk`)

### Building

1.  Navigate to the `src/` directory:
    ```bash
    cd src
    ```
2.  Ensure dependencies are downloaded:
    ```bash
    go mod tidy
    ```
3.  Build the executable:
    ```bash
    go build .
    ```

## Running the Server

> **Warning**
> This server does not have built-in authentication. As such, it should not be used in any sensitive environments or exposed to the public internet.

The server supports two [communication transports](https://modelcontextprotocol.io/specification/2025-06-18/basic/transports):

  - **STDIO**: For local execution and SSH-tunneled remote access
  - **SSE**: For remote HTTP/HTTPS access over the internet

### Local Execution (STDIO Mode)

Execute the compiled binary from your terminal:

```bash
.\mtg-mcp-windows-amd64.exe
```

### Remote HTTP Access

For remote access over the internet, run with SSE transport:

```bash
# Linux
MCP_TRANSPORT=sse MCP_SSE_PORT=3000 ./mtg-mcp-linux-amd64

# Windows
$env:MCP_TRANSPORT="sse"; $env:MCP_SSE_PORT="3000"; .\mtg-mcp-windows-amd64.exe

# Docker
docker run --rm -p 3000:3000 \
  -e MCP_TRANSPORT=sse \
  -e MCP_SSE_PORT=3000 \
  mtg-mcp:latest

```

The server will start on `http://0.0.0.0:3000/sse` and be accessible from anywhere.

### Using with Claude Desktop

Update your `claude_desktop_config.json` to include the following under `mcpServers`:

#### Docker

```
    "mtg-mcp":{
      "command": "docker",
      "args": [
        "run",
        "-i",
        "dccoder/mtg-mcp:latest",
	"-e MCP_TRANSPORT=stdio"
      ]
    }
```

#### Local STDIO

**Windows:**
```json
"mtg-mcp": {
    "command": "C:\\path\\to\\mtg-mcp-windows-amd64.exe",
    "args": []
}
```

**Linux:**
```json
"mtg-mcp": {
    "command": "/path/to/mtg-mcp-linux-amd64",
    "args": []
}
```

This configuration file can be found by going to `Settings` -\> `Developer` -\> `Edit Config` in Claude Desktop.

## Environment Variables

This server supports the following environment variables:

| Variable | Default | Description |
|----------|---------|-------------|
| `MCP_SERVER_NAME` | `scryfall-card-search-server` | Server identification name |
| `MCP_SERVER_VERSION` | `v1.0.0` | Server version string |
| `MCP_LOG_TO_FILE` | `true` | Enable/disable file logging |
| `MCP_LOG_FILE` | `mcp-server.log` | Log file path (when logging enabled) |
| `MCP_TRANSPORT` | `stdio` | Transport type: `stdio` or `sse` |
| `MCP_SSE_HOST` | `0.0.0.0` | SSE server bind address (SSE mode only) |
| `MCP_SSE_PORT` | `3000` | SSE server port (SSE mode only) |
| `MCP_SSE_PATH` | `/sse` | SSE endpoint path (SSE mode only) |
| `MCP_SSL_CERT_FILE` | `nil` | Path to TLS certificate file (for https) |
| `MCP_SSL_KEY_FILE` | `nil` | Path to TLS certificate key (for https) |

**Example with environment variables:**

```json
{
  "mcpServers": {
    "card_server": {
      "command": "/path/to/mtg-mcp-linux-amd64",
      "args": [],
      "env": {
        "MCP_LOG_TO_FILE": "false",
        "MCP_SERVER_NAME": "my-mtg-server"
      }
    }
  }
}
```

## Sources

Card effects and keywords are from the [Magic: The Gathering Comprehensive Rules - September 19, 2025](./MagicCompRules%2020250919.pdf)