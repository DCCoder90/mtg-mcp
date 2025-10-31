package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/BlueMonday/go-scryfall"
)

// ThemePattern represents a theme pattern loaded from JSON
type ThemePattern struct {
	Patterns            []string `json:"patterns"`
	SynergyQuery        string   `json:"synergyQuery"`
	SynergyDescription  string   `json:"synergyDescription"`
	SynergyType         string   `json:"synergyType"`
}

// Global cache for resources
var (
	creatureTypesCache  []string
	themePatternsCache  map[string]ThemePattern
	keywordAbilitiesCache []string
)

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}

func limitCards(cards []scryfall.Card, max int) []scryfall.Card {
	if max <= 0 || max > len(cards) {
		return cards
	}
	return cards[:max]
}

func extractKeywordsFromText(text string) []string {
	keywords := []string{}
	keywordList := loadKeywordAbilities()

	textLower := fmt.Sprintf(" %s ", text)
	// Convert to lowercase for case-insensitive comparison
	for i := 0; i < len(textLower); i++ {
		if textLower[i] >= 'A' && textLower[i] <= 'Z' {
			textLower = textLower[:i] + string(textLower[i]+32) + textLower[i+1:]
		}
	}

	for _, kw := range keywordList {
		if len(kw) > 0 {
			kwLower := kw
			// Convert keyword to lowercase
			for i := 0; i < len(kwLower); i++ {
				if kwLower[i] >= 'A' && kwLower[i] <= 'Z' {
					kwLower = kwLower[:i] + string(kwLower[i]+32) + kwLower[i+1:]
				}
			}

			// Check for keyword followed by space or punctuation
			searchStr := fmt.Sprintf(" %s ", kwLower)
			if len(textLower) >= len(searchStr) {
				for i := 0; i <= len(textLower)-len(searchStr); i++ {
					if textLower[i:i+len(searchStr)] == searchStr {
						if !contains(keywords, kwLower) {
							keywords = append(keywords, kwLower)
						}
						break
					}
				}
			}
		}
	}

	return keywords
}

func extractThemesFromCard(card scryfall.Card) []string {
	themes := []string{}
	oracleText := ""
	if card.OracleText != "" {
		oracleText = card.OracleText
	}

	// Extract keywords
	keywords := extractKeywordsFromText(oracleText)
	themes = append(themes, keywords...)

	// Load theme patterns from resource file
	themePatterns := loadThemePatterns()

	for theme, pattern := range themePatterns {
		for _, patternStr := range pattern.Patterns {
			if len(oracleText) > 0 && patternStr != "" {
				regexPattern := "(?i)" + patternStr
				re, err := regexp.Compile(regexPattern)
				if err != nil {
					// If regex compilation fails, fall back to simple substring matching
					log.Printf("Invalid regex pattern '%s' for theme '%s', using substring match: %v", patternStr, theme, err)
					if strings.Contains(strings.ToLower(oracleText), strings.ToLower(patternStr)) {
						if !contains(themes, theme) {
							themes = append(themes, theme)
						}
					}
					continue
				}

				// Check if pattern matches oracle text
				if re.MatchString(oracleText) {
					if !contains(themes, theme) {
						themes = append(themes, theme)
					}
					break // Found a match for this, move on
				}
			}
		}
	}

	return themes
}

// loadCreatureTypes loads creature types from the embedded resource file
func loadCreatureTypes() []string {
	if creatureTypesCache != nil {
		return creatureTypesCache
	}

	data, err := embeddedResources.ReadFile("res/creature-types.txt")
	if err != nil {
		log.Printf("Error loading creature types: %v, using defaults", err)
		return []string{"Human", "Elf", "Goblin", "Zombie", "Vampire", "Soldier", "Wizard", "Knight", "Dragon", "Angel", "Demon", "Spirit", "Merfolk", "Beast"}
	}

	var types []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			types = append(types, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading creature types: %v", err)
		return []string{"Human", "Elf", "Goblin", "Zombie", "Vampire"}
	}

	creatureTypesCache = types
	return types
}

// loadThemePatterns loads theme patterns from the embedded resource file
func loadThemePatterns() map[string]ThemePattern {
	if themePatternsCache != nil {
		return themePatternsCache
	}

	data, err := embeddedResources.ReadFile("res/themepatterns.json")
	if err != nil {
		log.Printf("Error loading theme patterns: %v, using defaults", err)
		return getDefaultThemePatterns()
	}

	var patterns map[string]ThemePattern
	if err := json.Unmarshal(data, &patterns); err != nil {
		log.Printf("Error decoding theme patterns: %v", err)
		return getDefaultThemePatterns()
	}

	themePatternsCache = patterns
	return patterns
}

// getDefaultThemePatterns returns default theme patterns as fallback
func getDefaultThemePatterns() map[string]ThemePattern {
	return map[string]ThemePattern{
		"sacrifice": {
			Patterns:           []string{"sacrifice", "dies"},
			SynergyQuery:       "(oracle:\"when a creature dies\" OR oracle:\"whenever you sacrifice\")",
			SynergyDescription: "Cards that benefit from creature sacrifice",
			SynergyType:        "Sacrifice Synergy",
		},
		"tokens": {
			Patterns:           []string{"create.*token", "token"},
			SynergyQuery:       "(oracle:\"create.*token\")",
			SynergyDescription: "Cards that create or benefit from tokens",
			SynergyType:        "Token Synergy",
		},
	}
}

// loadKeywordAbilities loads keyword abilities from the embedded resource file
func loadKeywordAbilities() []string {
	if keywordAbilitiesCache != nil {
		return keywordAbilitiesCache
	}

	data, err := embeddedResources.ReadFile("res/keyword-abilities.txt")
	if err != nil {
		log.Printf("Error loading keyword abilities: %v, using defaults", err)
		return []string{
			"deathtouch", "defender", "double strike", "enchant", "equip", "first strike",
			"flash", "flying", "haste", "hexproof", "indestructible", "intimidate",
			"lifelink", "menace", "protection", "reach", "trample", "vigilance", "ward",
			"cycling", "flashback", "kicker", "madness", "morph", "storm", "convoke",
			"delve", "suspend", "cascade", "miracle", "overload", "prowess",
		}
	}

	var abilities []string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		line := scanner.Text()
		if line != "" {
			abilities = append(abilities, line)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading keyword abilities: %v", err)
		return []string{"flying", "trample", "lifelink", "deathtouch", "haste"}
	}

	keywordAbilitiesCache = abilities
	return abilities
}

// findReprintCards searches for reprints of a card
func findReprintCards(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions, maxResults int) *RelatedCardCategory {
	log.Printf("Searching for reprints of %s (oracle_id: %s)", mainCard.Name, mainCard.OracleID)
	reprintQuery := fmt.Sprintf(`oracle_id:%s`, mainCard.OracleID)
	reprints, err := client.SearchCards(ctx, reprintQuery, opts)
	if err == nil && len(reprints.Cards) > 1 {
		// Filter out the main card itself
		reprintCards := []scryfall.Card{}
		for _, card := range reprints.Cards {
			if card.ID != mainCard.ID {
				reprintCards = append(reprintCards, card)
			}
		}
		if len(reprintCards) > 0 {
			log.Printf("Found %d reprints", len(reprintCards))
			return &RelatedCardCategory{
				CategoryName: "Reprints",
				Cards:        limitCards(reprintCards, maxResults),
				Count:        len(reprintCards),
			}
		}
	}
	return nil
}

// findTokenCards searches for tokens created by a card
func findTokenCards(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, maxResults int) *RelatedCardCategory {
	if mainCard.AllParts == nil {
		return nil
	}

	log.Printf("Searching for tokens created by %s", mainCard.Name)
	tokenCards := []scryfall.Card{}
	for _, part := range mainCard.AllParts {
		if part.Component == "token" {
			token, err := client.GetCard(ctx, part.ID)
			if err == nil {
				tokenCards = append(tokenCards, token)
			}
		}
	}
	if len(tokenCards) > 0 {
		log.Printf("Found %d tokens", len(tokenCards))
		return &RelatedCardCategory{
			CategoryName: "Tokens Created",
			Cards:        limitCards(tokenCards, maxResults),
			Count:        len(tokenCards),
		}
	}
	return nil
}

// findMechanicCards searches for cards with similar mechanics
func findMechanicCards(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions) *RelatedCardCategory {
	log.Printf("Searching for cards with similar mechanics to %s", mainCard.Name)
	keywords := extractKeywordsFromText(mainCard.OracleText)
	if len(keywords) > 0 {
		// Use the first 2 keywords for search to keep results relevant
		searchKeywords := keywords
		if len(searchKeywords) > 2 {
			searchKeywords = searchKeywords[:2]
		}

		for _, kw := range searchKeywords {
			mechanicQuery := fmt.Sprintf(`oracle:"%s" -name:"%s"`, kw, mainCard.Name)
			mechanics, err := client.SearchCards(ctx, mechanicQuery, opts)
			if err == nil && len(mechanics.Cards) > 0 {
				log.Printf("Found %d cards with '%s' mechanic", len(mechanics.Cards), kw)
				return &RelatedCardCategory{
					CategoryName: fmt.Sprintf("Similar Mechanics (%s)", kw),
					Cards:        limitCards(mechanics.Cards, 5),
					Count:        len(mechanics.Cards),
				}
			}
		}
	}
	return nil
}

// findArtistCards searches for cards by the same artist
func findArtistCards(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions, maxResults int) *RelatedCardCategory {
	if mainCard.Artist == nil || *mainCard.Artist == "" {
		return nil
	}

	log.Printf("Searching for cards by artist %s", *mainCard.Artist)
	artistQuery := fmt.Sprintf(`artist:"%s" -name:"%s"`, *mainCard.Artist, mainCard.Name)
	artistCards, err := client.SearchCards(ctx, artistQuery, opts)
	if err == nil && len(artistCards.Cards) > 0 {
		log.Printf("Found %d cards by same artist", len(artistCards.Cards))
		return &RelatedCardCategory{
			CategoryName: fmt.Sprintf("Same Artist (%s)", *mainCard.Artist),
			Cards:        limitCards(artistCards.Cards, maxResults),
			Count:        len(artistCards.Cards),
		}
	}
	return nil
}

// findKeywordSynergies searches for cards with shared keywords
func findKeywordSynergies(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions, synergies []SynergyCategory) []SynergyCategory {
	keywords := extractKeywordsFromText(mainCard.OracleText)
	if len(keywords) > 0 {
		for _, keyword := range keywords {
			if len(synergies) >= 3 { // Limit to avoid too many categories
				break
			}
			keywordQuery := fmt.Sprintf(`oracle:"%s" -name:"%s"`, keyword, mainCard.Name)
			log.Printf("Searching for keyword synergy: %s", keyword)
			keywordCards, err := client.SearchCards(ctx, keywordQuery, opts)
			if err == nil && len(keywordCards.Cards) > 0 {
				synergies = append(synergies, SynergyCategory{
					SynergyType: "Keyword Synergy",
					Description: fmt.Sprintf("Cards that share the '%s' keyword ability", keyword),
					Cards:       limitCards(keywordCards.Cards, 5),
					Count:       len(keywordCards.Cards),
				})
				log.Printf("Found %d cards with '%s' keyword", len(keywordCards.Cards), keyword)
			}
		}
	}
	return synergies
}

// findThemeSynergies searches for theme-based synergies
func findThemeSynergies(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions, searchThemes []string, synergies []SynergyCategory) []SynergyCategory {
	themePatterns := loadThemePatterns()
	creatureTypes := loadCreatureTypes()
	themesSearched := 0

	for _, theme := range searchThemes {
		if themesSearched >= 2 { // Limit theme searches
			break
		}
		if pattern, ok := themePatterns[theme]; ok {
			if pattern.SynergyQuery == "" && theme == "tribal" {
				// Handle tribal synergy specially
				if mainCard.TypeLine != "" {
					for _, cType := range creatureTypes {
						typeLower := mainCard.TypeLine
						ctLower := cType
						found := false
						for i := 0; i <= len(typeLower)-len(ctLower); i++ {
							match := true
							for j := 0; j < len(ctLower); j++ {
								c1 := typeLower[i+j]
								c2 := ctLower[j]
								if c1 >= 'A' && c1 <= 'Z' {
									c1 = c1 + 32
								}
								if c2 >= 'A' && c2 <= 'Z' {
									c2 = c2 + 32
								}
								if c1 != c2 {
									match = false
									break
								}
							}
							if match {
								found = true
								break
							}
						}
						if found {
							tribalQuery := fmt.Sprintf(`type:%s -name:"%s"`, cType, mainCard.Name)
							log.Printf("Searching for tribal synergy: %s", cType)
							tribalCards, err := client.SearchCards(ctx, tribalQuery, opts)
							if err == nil && len(tribalCards.Cards) > 0 {
								synergies = append(synergies, SynergyCategory{
									SynergyType: pattern.SynergyType,
									Description: fmt.Sprintf("Cards that share the %s creature type", cType),
									Cards:       limitCards(tribalCards.Cards, 5),
									Count:       len(tribalCards.Cards),
								})
								log.Printf("Found %d cards with %s type", len(tribalCards.Cards), cType)
								themesSearched++
								break
							}
						}
					}
				}
			} else if pattern.SynergyQuery != "" {
				themeQuery := fmt.Sprintf(`%s -name:"%s"`, pattern.SynergyQuery, mainCard.Name)
				log.Printf("Searching for theme synergy: %s", theme)
				themeCards, err := client.SearchCards(ctx, themeQuery, opts)
				if err == nil && len(themeCards.Cards) > 0 {
					synergies = append(synergies, SynergyCategory{
						SynergyType: pattern.SynergyType,
						Description: pattern.SynergyDescription,
						Cards:       limitCards(themeCards.Cards, 5),
						Count:       len(themeCards.Cards),
					})
					log.Printf("Found %d cards for %s theme", len(themeCards.Cards), theme)
					themesSearched++
				}
			}
		}
	}
	return synergies
}

// findColorIdentitySynergies searches for cards with matching color identity
func findColorIdentitySynergies(ctx context.Context, client *scryfall.Client, mainCard scryfall.Card, opts scryfall.SearchCardsOptions, synergies []SynergyCategory) []SynergyCategory {
	if mainCard.Colors != nil && len(mainCard.Colors) > 0 && len(synergies) < 4 {
		colorStr := ""
		for _, color := range mainCard.Colors {
			colorStr += string(color)
		}
		if colorStr != "" {
			colorQuery := fmt.Sprintf(`color:%s -name:"%s"`, colorStr, mainCard.Name)
			log.Printf("Searching for color identity synergy: %s", colorStr)
			colorCards, err := client.SearchCards(ctx, colorQuery, opts)
			if err == nil && len(colorCards.Cards) > 0 {
				synergies = append(synergies, SynergyCategory{
					SynergyType: "Color Identity Synergy",
					Description: fmt.Sprintf("Cards that share the same color identity"),
					Cards:       limitCards(colorCards.Cards, 5),
					Count:       len(colorCards.Cards),
				})
				log.Printf("Found %d cards with matching colors", len(colorCards.Cards))
			}
		}
	}
	return synergies
}

