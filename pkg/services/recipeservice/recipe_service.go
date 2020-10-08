package recipeservice

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/flucas97/delivery-much-challange/internal/domain/recipe"
	"github.com/flucas97/delivery-much-challange/pkg/services/gifservice"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
)

var (
	// RecipeService interface with other layers
	RecipeService recipeServiceInterface = &recipeService{}

	// RecipeURI recipepuppy API endpoint
	RecipeURI = "http://www.recipepuppy.com/api/?i="
)

type recipeServiceInterface interface {
	GetAll([]string) ([]recipe.Recipe, *errortools.APIError)
	FetchGifFor([]recipe.Recipe) ([]recipe.Recipe, *errortools.APIError)
	GetGif(string) (string, *errortools.APIError)
	ConcatenateIngredients([]string) string
}

type recipeService struct{}

// GetAll method responsible for seeking recipes based on ingredients
func (rs *recipeService) GetAll(ingredients []string) ([]recipe.Recipe, *errortools.APIError) {
	if len(ingredients) > 3 {
		return nil, errortools.APIErrorInterface.NewBadRequestError("max of 3 ingredients. recipeservice.GetAll")
	}

	ingredientsConcatenated := rs.ConcatenateIngredients(ingredients)

	request, err := http.NewRequest(http.MethodGet, fmt.Sprintf("%s%s", RecipeURI, ingredientsConcatenated), nil)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error mounting request. recipeservice.GetAll")
	}

	client := &http.Client{}

	resp, err := client.Do(request)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error doing request. recipeservice.GetAll")
	}
	defer resp.Body.Close()

	bytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error reading body. recipeservice.GetAll")
	}

	var sr recipe.SearchResult

	if err := json.Unmarshal(bytes, &sr); err != nil {
		return nil, errortools.APIErrorInterface.NewInternalServerError("error unmarshalling response from client. recipeservice.GetAll")
	}

	recipesWithoutGif := sr.IngredientsToSortedSlice()

	result, e := rs.FetchGifFor(recipesWithoutGif)
	if e != nil {
		return nil, e
	}

	return result, nil
}

// ConcatenateIngredients to join each ingredient, sep by comma
func (rs *recipeService) ConcatenateIngredients(ingredients []string) string {
	return strings.Join(ingredients, ",")
}

// GetGif for a specif label/tag
func (rs *recipeService) GetGif(label string) (string, *errortools.APIError) {
	gif, err := gifservice.GifService.GetRandom(label)
	if err != nil {
		return "", err
	}

	return gif.URL, nil
}

// FetchGifFor fills the Gif attribute for each Recipe
func (rs *recipeService) FetchGifFor(recipes []recipe.Recipe) ([]recipe.Recipe, *errortools.APIError) {
	recipesSize := len(recipes) - 1

	if recipesSize != 0 {
		var atIndex int

		for atIndex <= recipesSize {
			var recipeTitle = recipes[atIndex].Title

			if recipeTitle == "" {
				return nil, errortools.APIErrorInterface.NewInternalServerError("title cannot be empty. recipeservice.FetchGifFor")
			}

			gifURL, err := rs.GetGif(recipeTitle)
			if err != nil {
				return nil, err
			}

			recipes[atIndex].Gif = gifURL
			atIndex++
		}
	}

	return recipes, nil
}
