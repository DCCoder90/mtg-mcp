package main

import (
	"log"
	"reflect"

	"github.com/BlueMonday/go-scryfall"
	"github.com/google/jsonschema-go/jsonschema"
)

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
				Items:       &jsonschema.Schema{},
			},
		},
	}

	typeSchemas := map[reflect.Type]*jsonschema.Schema{
		reflect.TypeOf(scryfall.Date{}):          customDateSchema,
		reflect.TypeOf([]scryfall.FrameEffect{}): customNilArraySchema,
		reflect.TypeOf([]scryfall.CardFace{}):    customNilArraySchema,
		reflect.TypeOf([]scryfall.Color{}):       customNilArraySchema,
		reflect.TypeOf([]scryfall.RelatedCard{}): customNilArraySchema,
	}

	schema, err := jsonschema.For[SearchCardResult](&jsonschema.ForOptions{
		TypeSchemas: typeSchemas,
	})
	if err != nil {
		log.Fatalf("Failed to generate output schema: %v", err)
	}

	outputSchema = schema
	log.Println("Common Scryfall output schema generated.")

	relatedSchema, err := jsonschema.For[FindRelatedCardsResult](&jsonschema.ForOptions{
		TypeSchemas: typeSchemas,
	})

	if err != nil {
		log.Fatalf("Failed to generate related cards schema: %v", err)
	}

	relatedCardsSchema = relatedSchema
	log.Println("Related cards output schema generated.")

	synergiesSchemaGen, err := jsonschema.For[FindCardSynergiesResult](&jsonschema.ForOptions{
		TypeSchemas: typeSchemas,
	})

	if err != nil {
		log.Fatalf("Failed to generate synrgies schema: %v", err)
	}

	synergiesSchema = synergiesSchemaGen
	log.Println("Card synergies output schema generated.")
}