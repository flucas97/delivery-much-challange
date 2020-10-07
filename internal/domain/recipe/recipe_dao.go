package recipe

import (
	"sort"
	"strings"
)

// Recipe model as response, using slice of ingredients
type Recipe struct {
	Title       string   `json:"title"`
	Link        string   `json:"href"`
	Ingredients []string `json:"ingredients"`
	Gif         string   `json:"gif"`
}

// FromRecipepuppy Recipe model
type FromRecipepuppy struct {
	Title       string `json:"title"`
	Link        string `json:"href"`
	Ingredients string `json:"ingredients"`
}

// SearchResult from recipepuppy API
type SearchResult struct {
	Results []FromRecipepuppy `json:"results"`
}

// IngredientsToSortedSlice is responsable to cast ingredients to slice
func (search *SearchResult) IngredientsToSortedSlice() []Recipe {
	var result []Recipe

	for _, recipe := range search.Results {
		var actual Recipe

		ingredients := strings.Split(recipe.Ingredients, ", ")

		actual.Title = recipe.Title
		actual.Link = recipe.Link
		actual.Ingredients = ingredients

		actual.sortIngredients()

		result = append(result, actual)
	}

	return result
}

func (m *Recipe) sortIngredients() {
	sort.Strings(m.Ingredients)
}
