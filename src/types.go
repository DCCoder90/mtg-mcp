package main

import "github.com/BlueMonday/go-scryfall"

type SearchCardArgs struct {
	Name string `json:"name" jsonschema:"the name of the Magic: The Gathering card"`
}

type SearchCardResult struct {
	Cards []scryfall.Card `json:"cards" jsonschema:"list of cards found matching the name"`
}

type SearchCardByTextArgs struct {
	Text string `json:"text" jsonschema:"The Oracle text to search for on the card."`
}

type SearchCardByColorArgs struct {
	Color string `json:"color" jsonschema:"The card color(s) to search for. Use W, U, B, R, G. (e.g., 'W', 'UB', 'M' for multicolor, 'C' for colorless)."`
}

type FindRelatedCardsArgs struct {
	CardName     string   `json:"card_name" jsonschema:"required,The name of the card to find relationships for"`
	RelationType []string `json:"relation_type,omitempty" jsonschema:"Types of relationships to find. Options: reprints, tokens, mechanics, same_artist, same_set. If empty, returns all types."`
	MaxResults   int      `json:"max_results,omitempty" jsonschema:"Maximum number of results per category (default: 10)"`
}

type RelatedCardCategory struct {
	CategoryName string          `json:"category_name" jsonschema:"The type of relationship (e.g., 'Reprints', 'Tokens Created', 'Similar Mechanics')"`
	Cards        []scryfall.Card `json:"cards" jsonschema:"List of related cards in this category"`
	Count        int             `json:"count" jsonschema:"Number of cards found in this category"`
}

type FindRelatedCardsResult struct {
	MainCard   scryfall.Card         `json:"main_card" jsonschema:"The original card being queried"`
	Categories []RelatedCardCategory `json:"categories" jsonschema:"Categories of related cards"`
}

type FindCardSynergiesArgs struct {
	CardName   string `json:"card_name" jsonschema:"required,The name of the card to find synergies for"`
	Theme      string `json:"theme,omitempty" jsonschema:"Optional theme or strategy to focus on (e.g., 'sacrifice', 'tokens', 'graveyard', 'counters')"`
	MaxResults int    `json:"max_results,omitempty" jsonschema:"Maximum number of synergistic cards to return (default: 15)"`
}

type SynergyCategory struct {
	SynergyType string          `json:"synergy_type" jsonschema:"Type of synergy (e.g., 'Keyword Synergy', 'Mechanic Synergy', 'Thematic Synergy')"`
	Description string          `json:"description" jsonschema:"Explanattion of why these cards synergize"`
	Cards       []scryfall.Card `json:"cards" jsonschema:"Cards that synergize with the main card"`
	Count       int             `json:"count" jsonschema:"Number of cards in this synergy category"`
}

type FindCardSynergiesResult struct {
	MainCard        scryfall.Card     `json:"main_card" jsonschema:"The original card being analyzed"`
	ExtractedThemes []string          `json:"extracted_themes" jsonschema:"Themes and mechanics identified from the card"`
	Synergies       []SynergyCategory `json:"synergies" jsonschema:"Categories of synrgistic cards"`
}
