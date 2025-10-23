package main

import (
	"context"
	"fmt"
	"log"

	"github.com/BlueMonday/go-scryfall" // Import the new SDK
)

func main() {
	log.Println("Searching for white 'Sliver' cards with CMC 3 using go-scryfall...")

	client, err := scryfall.NewClient()
	if err != nil {
		log.Fatalf("Failed to create Scryfall client: %v", err)
	}

	searchQuery := "name:Sliver color=w cmc=3"

	opts := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModePrints,
		IncludeMultilingual: false,
		IncludeExtras:       false,
		IncludeVariations:   false,
	}

	ctx := context.Background()
	result, err := client.SearchCards(ctx, searchQuery, opts)
	if err != nil {
		if scryfallErr, ok := err.(*scryfall.Error); ok {
			log.Fatalf("Scryfall API error: %s (Status: %d)", scryfallErr.Details, scryfallErr.Status)
		}
		log.Fatalf("Failed to search cards: %v", err)
	}

	if len(result.Cards) == 0 {
		log.Printf("No cards found matching the query: %q", searchQuery)
		return
	}

	log.Printf("Found %d cards matching %q (Total found: %d):\n", len(result.Cards), searchQuery, result.TotalCards)
	for _, card := range result.Cards {
		fmt.Printf("--------------------------------\n")
		fmt.Printf("Name:       %s\n", card.Name)
		fmt.Printf("Mana Cost:  %s\n", card.ManaCost)
		fmt.Printf("CMC:        %.1f\n", card.CMC)
		fmt.Printf("Colors:     %v\n", card.Colors)
		fmt.Printf("Color ID:   %v\n", card.ColorIdentity)
		fmt.Printf("Set:        %s (%s)\n", card.SetName, card.Set)
		fmt.Printf("Type:       %s\n", card.TypeLine)
		fmt.Printf("Effects:	%s\n", card.OracleText)
	}
}
