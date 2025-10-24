package main

import "github.com/BlueMonday/go-scryfall"

type SearchCardArgs struct {
	Name string `json:"name" jsonschema:"the name of the Magic: The Gathering card"`
}

type SearchCardResult struct {
	Cards []scryfall.Card `json:"cards" jsonschema:"list of cards found matching the name"`
}

type SearchCardByTextArgs struct {
	Text string `json:"text" description:"The Oracle text to search for on the card."`
}