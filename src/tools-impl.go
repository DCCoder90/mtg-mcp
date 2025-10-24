package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BlueMonday/go-scryfall"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func executeSearch(ctx context.Context, searchQuery, searchTerm, searchType string) (*mcp.CallToolResult, SearchCardResult, error) {
	client, err := scryfall.NewClient()
	if err != nil {
		log.Printf("Error creating Scryfall client: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error initializing card search service: %v", err)}},
		}, SearchCardResult{}, nil
	}

	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModeCards,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	log.Printf("Searching Scryfall for %s: %s (Query: %s)", searchType, searchTerm, searchQuery)
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil {
		log.Printf("Error searching Scryfall for %s %s: %v", searchType, searchTerm, err)
		if scryfallErr, ok := err.(*scryfall.Error); ok {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Scryfall API error searching for %s '%s': %s (Status: %d)", searchType, searchTerm, scryfallErr.Details, scryfallErr.Status)}},
			}, SearchCardResult{}, nil
		}

		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error searching for %s '%s': %v", searchType, searchTerm, err)}},
		}, SearchCardResult{}, nil
	}

	if len(result.Cards) == 0 {
		log.Printf("No cards found matching %s: %s", searchType, searchTerm)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No cards found matching the %s '%s'.", searchType, searchTerm)}},
		}, SearchCardResult{Cards: []scryfall.Card{}}, nil
	}

	log.Printf("Found %d cards matching %s: %s", len(result.Cards), searchType, searchTerm)
	return nil, SearchCardResult{Cards: result.Cards}, nil
}

func searchCardByName(ctx context.Context, req *mcp.CallToolRequest, args SearchCardArgs) (*mcp.CallToolResult, SearchCardResult, error) {
	if args.Name == "" {
		log.Println("Error: Received request with empty card name.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card name cannot be empty."}},
		}, SearchCardResult{}, nil
	}

	searchQuery := fmt.Sprintf(`name:"%s"`, args.Name)
	return executeSearch(ctx, searchQuery, args.Name, "name")
}

func searchCardByText(ctx context.Context, req *mcp.CallToolRequest, args SearchCardByTextArgs) (*mcp.CallToolResult, SearchCardResult, error) {
	if args.Text == "" {
		log.Println("Error: Received request with empty card text.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card text cannot be empty."}},
		}, SearchCardResult{}, nil
	}

	searchQuery := fmt.Sprintf(`oracle:"%s"`, args.Text)
	return executeSearch(ctx, searchQuery, args.Text, "text")
}

func searchCardByColor(ctx context.Context, req *mcp.CallToolRequest, args SearchCardByColorArgs) (*mcp.CallToolResult, SearchCardResult, error) {
	if args.Color == "" {
		log.Println("Error: Received request with empty card color.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card color cannot be empty."}},
		}, SearchCardResult{}, nil
	}

	searchQuery := fmt.Sprintf(`color:%s`, args.Color)
	return executeSearch(ctx, searchQuery, args.Color, "color")
}