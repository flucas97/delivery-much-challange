package recipe

import (
	"sort"
	"strings"
)

// Marshaller model as response using slice
type Marshaller struct {
	Title       string   `json:"title"`
	Link        string   `json:"href"`
	Ingredients []string `json:"ingredients"`
	Gif         string   `json:"gif"`
}

// Recipe model from recipepuppy
type Recipe struct {
	Title       string `json:"title"`
	Link        string `json:"href"`
	Ingredients string `json:"ingredients"`
	Gif         string `json:"gif"`
}

// SearchResult from recipepuppy API
type SearchResult struct {
	Results []Recipe `json:"results"`
}

// IngredientsToSlice is responsable to cast ingredients to slice
func (search *SearchResult) IngredientsToSlice() []Marshaller {
	var result []Marshaller

	for _, recipe := range search.Results {
		var actual Marshaller

		ingredients := strings.Split(recipe.Ingredients, ", ")

		actual.Title = recipe.Title
		actual.Link = recipe.Link
		actual.Ingredients = ingredients
		actual.Gif = recipe.Gif

		actual.sortIngredients()

		result = append(result, actual)
	}

	return result
}

func (m *Marshaller) sortIngredients() {
	sort.Strings(m.Ingredients)
}
