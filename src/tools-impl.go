package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BlueMonday/go-scryfall"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

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

func searchCardByText(ctx context.Context, req *mcp.CallToolRequest, args SearchCardByTextArgs) (*mcp.CallToolResult, SearchCardResult, error) {
	if args.Text == "" {
		log.Println("Error: Received request with empty card text.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card text cannot be empty."}},
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

	// oracle: or o: keyword for Scryfall
	searchQuery := fmt.Sprintf(`oracle:"%s"`, args.Text)

	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModeCards,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	log.Printf("Searching Scryfall for card text: %s (Query: %s)", args.Text, searchQuery)
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil {
		log.Printf("Error searching Scryfall for card text %s: %v", args.Text, err)
		if scryfallErr, ok := err.(*scryfall.Error); ok {
			return &mcp.CallToolResult{
				IsError: true,
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Scryfall API error searching for text '%s': %s (Status: %d)", args.Text, scryfallErr.Details, scryfallErr.Status)}},
			}, SearchCardResult{}, nil
		}

		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error searching for card text '%s': %v", args.Text, err)}},
		}, SearchCardResult{}, nil
	}

	if len(result.Cards) == 0 {
		log.Printf("No cards found matching text: %s", args.Text)
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No cards found matching the text '%s'.", args.Text)}},
		}, SearchCardResult{Cards: []scryfall.Card{}}, nil
	}

	log.Printf("Found %d cards matching text: %s", len(result.Cards), args.Text)
	return nil, SearchCardResult{Cards: result.Cards}, nil
}