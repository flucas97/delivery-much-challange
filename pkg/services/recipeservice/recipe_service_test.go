package recipeservice

import (
	"fmt"
	"net/http"
	"os"
	"testing"

	"github.com/flucas97/delivery-much-challange/internal/domain/recipe"
	"github.com/flucas97/delivery-much-challange/pkg/services/gifservice"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

var rs recipeService

func TestGetRandomByTagRestClientNoError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", RecipeURI+"eggs,garlic,onions", httpmock.NewStringResponder(http.StatusOK,
		`{
			"results":[
				{
					   "title":"Vegetable-Pasta Oven Omelet",
					   "href":"http:\/\/find.myrecipes.com\/recipes\/recipefinder.dyn?action=displayRecipe&recipe_id=520763",
					   "ingredients":"tomato, onions, red pepper, garlic"
				}
			]
		}`))

	var ingredients = []string{"eggs", "garlic", "onions"}
	result, err := rs.GetAll(ingredients)

	assert.Nil(t, err)
	assert.NotNil(t, result)
	assert.Equal(t, "Vegetable-Pasta Oven Omelet", result[0].Title)
	assert.Equal(t, []string{"garlic", "onions", "red pepper", "tomato"}, result[0].Ingredients)
}

func TestGetRandomByTagRestClientWithFourIngredients(t *testing.T) {
	var ingredients = []string{"eggs", "garlic", "onions", "açai"}
	result, err := rs.GetAll(ingredients)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "max of 3 ingredients. recipeservice.GetAll", err.Message)
}

func TestGetRandomByTagRestClientErrorDoingRequest(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()

	RecipeURI = "fake.com.br"
	httpmock.RegisterResponder("GET", RecipeURI+"[]", httpmock.NewStringResponder(http.StatusInternalServerError, ""))

	var ingredients = []string{""}
	result, err := rs.GetAll(ingredients)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "error doing request. recipeservice.GetAll", err.Message)
}

func TestGetRandomByTagRestClientUnmarshallingError(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", RecipeURI+"eggs,garlic,onions", httpmock.NewStringResponder(http.StatusOK, `{[no_content]}`))

	var ingredients = []string{"eggs", "garlic", "onions"}
	result, err := rs.GetAll(ingredients)

	assert.Nil(t, result)
	assert.NotNil(t, err)
	assert.Equal(t, "error unmarshalling response from client. recipeservice.GetAll", err.Message)
}

func TestConcatenateIngredients(t *testing.T) {
	var (
		ingredients = []string{"aveia", "mel", "abobora"}
	)

	result := rs.concatenateIngredients(ingredients)
	assert.Equal(t, "aveia,mel,abobora", result)
}

func TestGetGifToRecipe(t *testing.T) {
	httpmock.Activate()
	defer httpmock.DeactivateAndReset()
	httpmock.RegisterResponder("GET", fmt.Sprintf(gifservice.GiphyURI, "dev", os.Getenv("GIPHY_API_KEY")), httpmock.NewStringResponder(http.StatusOK,
		`{
			"title": "dev",
			"images": {
				"original": {
					"url": "https://delivery-much-challange-dev.gif"
				}	
			}
		}`,
	))

	recipe := recipe.Recipe{
		Title: "dev",
	}

	err := rs.getGifToRecipe(&recipe)

	assert.Nil(t, err)
	assert.NotNil(t, recipe.Gif)
	assert.Equal(t, "https://delivery-much-challange-dev.gif", recipe.Gif)
}
