package recipeservice

import (
	"fmt"
	"net/http"
	"testing"

	"github.com/flucas97/delivery-much-challange/internal/domain/gif"
	"github.com/flucas97/delivery-much-challange/internal/domain/recipe"
	"github.com/flucas97/delivery-much-challange/pkg/services/gifservice"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
	"github.com/jarcoal/httpmock"
	"github.com/stretchr/testify/assert"
)

type gifServiceMock struct {
	GetRandomFn func(tag string) (*gif.Gif, *errortools.APIError)
}

func (sm *gifServiceMock) GetRandom(tag string) (*gif.Gif, *errortools.APIError) {
	return sm.GetRandomFn(tag)
}

type recipeServiceMock struct {
	GetAllFn                 func([]string) ([]recipe.Recipe, *errortools.APIError)
	FetchGifForFn            func([]recipe.Recipe) ([]recipe.Recipe, *errortools.APIError)
	GetGifFn                 func(string) (string, *errortools.APIError)
	ConcatenateIngredientsFn func([]string) string
}

func (rsm *recipeServiceMock) ConcatenateIngredients(ingredients []string) string {
	return rsm.ConcatenateIngredientsFn(ingredients)
}

func (rsm *recipeServiceMock) GetGif(label string) (string, *errortools.APIError) {
	return rsm.GetGifFn(label)
}

func (rsm *recipeServiceMock) GetAll(ingredients []string) ([]recipe.Recipe, *errortools.APIError) {
	return rsm.GetAllFn(ingredients)
}

func (rsm *recipeServiceMock) FetchGifFor(recipes []recipe.Recipe) ([]recipe.Recipe, *errortools.APIError) {
	return rsm.FetchGifForFn(recipes)
}

func TestGetAll(t *testing.T) {
	t.Run("Successfully receive recipes", func(t *testing.T) {
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
		result, err := RecipeService.GetAll(ingredients)

		recipe := result[0]
		assert.Nil(t, err)
		assert.NotNil(t, result)
		assert.Equal(t, "Vegetable-Pasta Oven Omelet", recipe.Title)
		assert.Equal(t, []string{"garlic", "onions", "red pepper", "tomato"}, recipe.Ingredients)
		assert.EqualValues(t, 4, len(recipe.Ingredients))
	})

	t.Run("More than three ingredients", func(t *testing.T) {
		var ingredients = []string{"eggs", "garlic", "onions", "açai"}
		result, err := RecipeService.GetAll(ingredients)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "max of 3 ingredients. recipeservice.GetAll", err.Message)
	})

	t.Run("Error doing request", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()

		RecipeURI = "fake.com.br"
		httpmock.RegisterResponder("GET", RecipeURI+"[]", httpmock.NewStringResponder(http.StatusInternalServerError, ""))

		var ingredients = []string{""}
		result, err := RecipeService.GetAll(ingredients)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error doing request. recipeservice.GetAll", err.Message)
	})

	t.Run("Unmarshal error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", RecipeURI+"eggs,garlic,onions", httpmock.NewStringResponder(http.StatusOK, `{[no_content]}`))

		var ingredients = []string{"eggs", "garlic", "onions"}
		result, err := RecipeService.GetAll(ingredients)

		assert.Nil(t, result)
		assert.NotNil(t, err)
		assert.Equal(t, "error unmarshalling response from client. recipeservice.GetAll", err.Message)
	})
}

func TestGetGif(t *testing.T) {
	t.Run("Successfully receive gif", func(t *testing.T) {
		serviceMock := gifServiceMock{}

		serviceMock.GetRandomFn = func(tag string) (*gif.Gif, *errortools.APIError) {
			mockURL := fmt.Sprintf("%v.gif", tag)

			responseMock := gif.Gif{
				URL: mockURL,
			}
			return &responseMock, nil
		}

		gifservice.GifService = &serviceMock

		result, err := RecipeService.GetGif("dev.test")

		assert.Equal(t, "dev.test.gif", result)
		assert.Nil(t, err)
	})

	t.Run("With error", func(t *testing.T) {
		serviceMock := gifServiceMock{}

		serviceMock.GetRandomFn = func(tag string) (*gif.Gif, *errortools.APIError) {
			responseMock := errortools.APIErrorInterface.NewInternalServerError("error getting Giphy. gifservice.GetRandom")
			return nil, responseMock
		}

		gifservice.GifService = &serviceMock

		result, err := RecipeService.GetGif("dev.test")

		assert.NotNil(t, err)
		assert.Equal(t, "error getting Giphy. gifservice.GetRandom", err.Message)
		assert.Equal(t, "", result)
	})
}

func TestFetchGifFor(t *testing.T) {
	t.Run("With error", func(t *testing.T) {
		example := []recipe.Recipe{
			{
				Title: "",
			},
			{
				Title: "",
			},
		}

		result, err := RecipeService.FetchGifFor(example)

		assert.NotNil(t, err)
		assert.Equal(t, "title cannot be empty. recipeservice.FetchGifFor", err.Message)
		assert.Nil(t, result)
	})
}

func TestConcatenateIngredients(t *testing.T) {
	t.Run("Successfully concatenate", func(t *testing.T) {
		var (
			ingredients = []string{"aveia", "mel", "abobora"}
		)

		result := RecipeService.ConcatenateIngredients(ingredients)
		assert.Equal(t, "aveia,mel,abobora", result)
	})
}
