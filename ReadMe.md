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
- **Same Artist**: Cards illustrated by the same artist
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

The build process automatically embeds all resource files from `src/res/` into the binary using Go's `embed` package, creating a self-contained executable.

### Dependencies
- github.com/modelcontextprotocol/go-sdk/mcp
- github.com/BlueMonday/go-scryfall
- github.com/google/jsonschema-go/jsonschema

### Schema Handling Notes
The server explicitly defines JSON schemas for certain types returned by the go-scryfall SDK (`scryfall.Date` and various optional slices like `[]scryfall.FrameEffect`, `[]scryfall.Color`, etc.). This is necessary because these types might be marshaled in ways (e.g., date as string, optional slices as null) that differ from the default schema inferred by the MCP SDK, ensuring correct validation of tool output.

## Running the Server

Execute the compiled binary from your terminal.

### Using with Claude Desktop

Update your `claude_desktop_config.json` to include the following under `mcpServers` (adjust your path to the executable):

**Windows:**
```json
"card_server": {
    "command": "C:\\path\\to\\mtg-mcp-windows-amd64.exe",
    "args": []
}
```

**Linux:**
```json
"card_server": {
    "command": "/path/to/mtg-mcp-linux-amd64",
    "args": []
}
```

This configuration file can be found by going to `Settings` -> `Developer` -> `Edit Config` in Claude Desktop.

## Sources

Card effects and keywords are from the [Magic: The Gathering Comprehensive Rules - September 19, 2025](./MagicCompRules%2020250919.pdf)