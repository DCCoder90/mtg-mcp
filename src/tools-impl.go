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

func findRelatedCards(ctx context.Context, req *mcp.CallToolRequest, args FindRelatedCardsArgs) (*mcp.CallToolResult, FindRelatedCardsResult, error) {
	if args.CardName == "" {
		log.Println("Error: Received request with empty card name.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card name cannot be empty."}},
		}, FindRelatedCardsResult{}, nil
	}

	// Defaults
	maxResults := args.MaxResults
	if maxResults <= 0 {
		maxResults = 10
	}

	client, err := scryfall.NewClient()
	if err != nil {
		log.Printf("Error creating Scryfall client: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error initializing card search service: %v", err)}},
		}, FindRelatedCardsResult{}, nil
	}

	searchQuery := fmt.Sprintf(`name:"%s"`, args.CardName)
	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModeCards,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	log.Printf("Searching for main card: %s", args.CardName)
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil || len(result.Cards) == 0 {
		log.Printf("Error finding main card '%s': %v", args.CardName, err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Could not find card '%s'", args.CardName)}},
		}, FindRelatedCardsResult{}, nil
	}

	mainCard := result.Cards[0]
	categories := []RelatedCardCategory{}

	relationTypes := args.RelationType
	if len(relationTypes) == 0 {
		relationTypes = []string{"reprints", "tokens", "mechanics", "same_set"}
	}

	// 1. Find reprints (same oracle_id)
	if contains(relationTypes, "reprints") {
		if category := findReprintCards(ctx, client, mainCard, opts, maxResults); category != nil {
			categories = append(categories, *category)
		}
	}

	// 2. Find tokens created
	if contains(relationTypes, "tokens") {
		if category := findTokenCards(ctx, client, mainCard, maxResults); category != nil {
			categories = append(categories, *category)
		}
	}

	// 3. Find similar mechanics (based on keywords)
	if contains(relationTypes, "mechanics") {
		if category := findMechanicCards(ctx, client, mainCard, opts); category != nil {
			categories = append(categories, *category)
		}
	}

	// Find cards from same set
	if contains(relationTypes, "same_set") && mainCard.Set != "" {
		log.Printf("Searching for cards from set %s", mainCard.Set)
		setQuery := fmt.Sprintf(`set:%s -name:"%s"`, mainCard.Set, mainCard.Name)
		setCards, err := client.SearchCards(ctx, setQuery, opts)
		if err == nil && len(setCards.Cards) > 0 {
			categories = append(categories, RelatedCardCategory{
				CategoryName: fmt.Sprintf("Same Set (%s)", mainCard.SetName),
				Cards:        limitCards(setCards.Cards, maxResults),
				Count:        len(setCards.Cards),
			})
			log.Printf("Found %d cards from same set", len(setCards.Cards))
		}
	}

	log.Printf("Successfully found related cards for '%s' in %d categories", mainCard.Name, len(categories))
	return nil, FindRelatedCardsResult{
		MainCard:   mainCard,
		Categories: categories,
	}, nil
}

func findCardSynergies(ctx context.Context, req *mcp.CallToolRequest, args FindCardSynergiesArgs) (*mcp.CallToolResult, FindCardSynergiesResult, error) {
	if args.CardName == "" {
		log.Println("Error: Received request with empty card name.")
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: Card name cannot be empty."}},
		}, FindCardSynergiesResult{}, nil
	}

	// Set defaults
	maxResults := args.MaxResults
	if maxResults <= 0 {
		maxResults = 15
	}

	client, err := scryfall.NewClient()
	if err != nil {
		log.Printf("Error creating Scryfall client: %v", err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error initializing card search service: %v", err)}},
		}, FindCardSynergiesResult{}, nil
	}

	// Get the main card
	searchQuery := fmt.Sprintf(`name:"%s"`, args.CardName)
	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModeCards,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	log.Printf("Searching for main card: %s", args.CardName)
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil || len(result.Cards) == 0 {
		log.Printf("Error finding main card '%s': %v", args.CardName, err)
		return &mcp.CallToolResult{
			IsError: true,
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Could not find card '%s'", args.CardName)}},
		}, FindCardSynergiesResult{}, nil
	}

	mainCard := result.Cards[0]
	extractedThemes := extractThemesFromCard(mainCard)
	synergies := []SynergyCategory{}

	log.Printf("Extracted themes for %s: %v", mainCard.Name, extractedThemes)

	// Prioritize user defined theme
	searchThemes := extractedThemes
	if args.Theme != "" {
		searchThemes = []string{args.Theme}
		log.Printf("Using user specified theme: %s", args.Theme)
	}

	// Keyword based 
	synergies = findKeywordSynergies(ctx, client, mainCard, opts, synergies)

	// Themebased 
	synergies = findThemeSynergies(ctx, client, mainCard, opts, searchThemes, synergies)

	// Color identity synergies
	synergies = findColorIdentitySynergies(ctx, client, mainCard, opts, synergies)

	if len(synergies) == 0 {
		log.Printf("No synergies found for '%s'", mainCard.Name)
		return &mcp.CallToolResult{
				Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("No clear synergies found for '%s'. Try specifying a specific theme.", mainCard.Name)}},
			}, FindCardSynergiesResult{
				MainCard:        mainCard,
				ExtractedThemes: extractedThemes,
				Synergies:       []SynergyCategory{},
			}, nil
	}

	log.Printf("Successfully found synergy categories")
	return nil, FindCardSynergiesResult{
		MainCard:        mainCard,
		ExtractedThemes: extractedThemes,
		Synergies:       synergies,
	}, nil
}
