package recipe

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIngredientsToSlice(t *testing.T) {
	sr := SearchResult{
		Results: []Recipe{
			{
				Title:       "Bolo de Banana",
				Link:        "https://www.delivery-much-test.com.br/recipes/bolo-de-banana",
				Ingredients: "banana, caramelo, ovo",
			},
		},
	}

	result := sr.IngredientsToSlice()
	assert.Equal(t, []string{"banana", "caramelo", "ovo"}, result[0].Ingredients, "failed to cast ingredients. recipe.IngredientsToSlice")
	assert.Contains(t, result[0].Ingredients, "caramelo", "ovo", "banana")
}
