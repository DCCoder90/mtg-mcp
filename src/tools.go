package main

import (
	"log"
	"reflect"

	"github.com/BlueMonday/go-scryfall"
	"github.com/google/jsonschema-go/jsonschema"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var outputSchema *jsonschema.Schema

func init() {
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
				Description: "A list of related items (colors, faces, etc.).",
				Items: &jsonschema.Schema{},
			},
		},
	}

	schema, err := jsonschema.For[SearchCardResult](&jsonschema.ForOptions{
		TypeSchemas: map[reflect.Type]*jsonschema.Schema{
			reflect.TypeOf(scryfall.Date{}):        customDateSchema,
			reflect.TypeOf([]scryfall.FrameEffect{}): customNilArraySchema,
			reflect.TypeOf([]scryfall.CardFace{}):    customNilArraySchema,
			reflect.TypeOf([]scryfall.Color{}):       customNilArraySchema,
			reflect.TypeOf([]scryfall.RelatedCard{}): customNilArraySchema,
		},
	})
	if err != nil {
		log.Fatalf("Failed to generate output schema: %v", err)
	}

	outputSchema = schema
	log.Println("Common Scryfall output schema generated.")
}


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

func registerTools(server *mcp.Server){
	registerSearchByTextTool(server)
	registerSearchByNameTool(server)
	registerSearchByColorTool(server)
}