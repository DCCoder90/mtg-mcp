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

## Development
### Prerequisites

* Go programming language (version 1.24 or later recommended for `go-sdk`)

### Building

1.  Navigate to the directory containing `main.go`.
2.  Ensure dependencies are downloaded:
    ```bash
    go mod tidy
    ```
3.  Build the executable:
    ```bash
    go build .
    ```


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
```json
        "card_server": {
            "command": "c:\\Users\\awesomeguy\\Desktop\\mtgmcp\\main.exe",
            "args": [
            ]
        }
```

This file can be found by going to `Settings` -> `Developer` -> `Edit Config`

## Sources

Card effects and keywords are from the [Magic: The Gathering Comprehensive Rules - September 19, 2025](./MagicCompRules%2020250919.pdf)