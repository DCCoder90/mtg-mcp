package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"reflect"

	"github.com/BlueMonday/go-scryfall"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

type SearchCardArgs struct {
	Name string `json:"name" jsonschema:"the name of the Magic: The Gathering card"`
}

type SearchCardResult struct {
	Cards []scryfall.Card `json:"cards" jsonschema:"list of cards found matching the name"`
}

func searchCardByName(ctx context.Context, req *mcp.CallToolRequest, args SearchCardArgs) (*mcp.CallToolResult, SearchCardResult, error) {
	if args.Name == "" {
		log.Println("Error: Received request with empty card name.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card name cannot be empty."}},
		}, SearchCardResult{}, nil
	}

	client, err := scryfall.NewClient()
	if err != nil {
		log.Printf("Error creating Scryfall client: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error initializing card search service: %v", err)}},
		}, SearchCardResult{}, nil
	}

	searchQuery := fmt.Sprintf(`name:"%s"`, args.Name)

	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModeCards,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	log.Printf("Searching Scryfall for card: %s (Query: %s)", args.Name, searchQuery)
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil {
		log.Printf("Error searching Scryfall for card %s: %v", args.Name, err)
		if scryfallErr, ok := err.(*scryfall.Error); ok {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Scryfall API error searching for '%s': %s (Status: %d)", args.Name, scryfallErr.Details, scryfallErr.Status)}},
			}, SearchCardResult{}, nil
		}

		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error searching for card '%s': %v", args.Name, err)}},
		}, SearchCardResult{}, nil
	}

	if len(result.Cards) == 0 {
		log.Printf("No cards found matching: %s", args.Name)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No cards found matching the name '%s'.", args.Name)}},
		}, SearchCardResult{Cards: []scryfall.Card{}}, nil
	}

	log.Printf("Found %d cards matching: %s", len(result.Cards), args.Name)
	return nil, SearchCardResult{Cards: result.Cards}, nil
}

func main() {
	logFileName := "mcp-server.log"
	logFile, err := os.OpenFile(logFileName, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		log.Printf("Failed to open log file %s: %v. Logging to stderr only.", logFileName, err)
	} else {
		defer logFile.Close()
		mw := io.MultiWriter(os.Stderr, logFile)
		log.SetOutput(mw)
		log.Printf("Logging initialized. Outputting to stderr and %s", logFileName)
	}

	//Define new servr
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "scryfall-card-search-server",
		Version: "v1.0.0"}, nil)

	customDateSchema := &jsonschema.Schema{
		Type:        "string",
		Format:      "date",
		Description: "The date the card was released (YYYY-MM-DD).",
	}

	//Ensure array is present even if result is nil
	customNilArraySchema := &jsonschema.Schema{
		OneOf: []*jsonschema.Schema{
			{Type: "null"},
			{
				Type:        "array",
				Description: "The frame effects applied to the card.",
				Items:       &jsonschema.Schema{Type: "string"},
			},
		},
	}

	outputSchema, err := jsonschema.For[SearchCardResult](&jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeOf(scryfall.Date{}):          customDateSchema,
			reflect.TypeOf([]scryfall.FrameEffect{}): customNilArraySchema,
			reflect.TypeOf([]scryfall.CardFace{}):    customNilArraySchema,
			reflect.TypeOf([]scryfall.Color{}):       customNilArraySchema,
			reflect.TypeOf([]scryfall.RelatedCard{}): customNilArraySchema,
		},
	})
	if err != nil {
		log.Fatalf("Failed to generate output schema: %v", err)
	}

	searchTool := &mcp.Tool{
		Name:         "search_card_by_name",
		Description:  "Searches Scryfall for MTG card details by the card's exact name.",
		OutputSchema: outputSchema,
	}

	mcp.AddTool(server, searchTool, searchCardByName)

	log.Println("Starting MCP server for MTG card search...")
	log.Println("Tool 'search_card_by_name' registered.")

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}

	log.Println("Server stopped.")
}
