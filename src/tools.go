package main

import (
	"log"

	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var outputSchema *jsonschema.Schema
var relatedCardsSchema *jsonschema.Schema
var synergiesSchema *jsonschema.Schema

func registerSearchByNameTool(server *mcp.Server) {
	searchTool := &mcp.Tool{
		Name:         "search_card_by_name",
		Description:  "Searches Scryfall for MTG card details by the card's exact name.",
		OutputSchema: outputSchema,
	}

	mcp.AddTool(server, searchTool, searchCardByName)

	log.Println("Tool 'search_card_by_name' registered.")
}

func registerSearchByTextTool(server *mcp.Server) {
	searchTool := &mcp.Tool{
		Name:         "search_card_by_text",
		Description:  "Searches Scryfall for MTG card details by the card's oracle text.",
		OutputSchema: outputSchema,
	}

	mcp.AddTool(server, searchTool, searchCardByName)

	log.Println("Tool 'search_card_by_text' registered.")
}

func registerSearchByColorTool(server *mcp.Server) {
	searchTool := &mcp.Tool{
		Name:         "search_card_by_color",
		Description:  "Searches Scryfall for MTG card details by the card's colors.",
		OutputSchema: outputSchema,
	}

	mcp.AddTool(server, searchTool, searchCardByColor)

	log.Println("Tool 'search_card_by_color' registered.")
}

func registerFindRelatedCardsTool(server *mcp.Server) {
	relatedCardsTool := &mcp.Tool{
		Name:         "find_related_cards",
		Description:  "Find cards related to a given card, including reprints, tokens created, cards with similar mechanics, or from the same set.",
		OutputSchema: relatedCardsSchema,
	}

	mcp.AddTool(server, relatedCardsTool, findRelatedCards)

	log.Println("Tool 'find_related_cards' registered.")
}

func registerFindCardSynergiesTool(server *mcp.Server) {
	synergiesTool := &mcp.Tool{
		Name:         "find_card_synergies",
		Description:  "Find cards that synergize with provided card based on keywords, themes, and mechanics. The user can also specify a theme to focus the search.",
		OutputSchema: synergiesSchema,
	}

	mcp.AddTool(server, synergiesTool, findCardSynergies)

	log.Println("Tool 'find_card_synergies' registered.")
}

func registerTools(server *mcp.Server) {
	registerSearchByTextTool(server)
	registerSearchByNameTool(server)
	registerSearchByColorTool(server)
	registerFindRelatedCardsTool(server)
	registerFindCardSynergiesTool(server)
}