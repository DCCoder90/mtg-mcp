package main

import (
	"context"
	"log"

	"github.com/BlueMonday/go-scryfall"
)

func main() {
	ctx := context.Background()
	client, err := scryfall.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	sco := scryfall.SearchCardsOptions{
		Unique:              scryfall.UniqueModePrints,
		Order:               scryfall.OrderSet,
		Dir:                 scryfall.DirDesc,
		IncludeExtras:       false,
		IncludeMultilingual: false,
		IncludeVariations:   false,
	}
	result, err := client.SearchCards(ctx, "sliver", sco)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("%s", result.Cards[0].Colors)
}
