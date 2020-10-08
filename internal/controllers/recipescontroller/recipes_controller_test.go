package recipescontroller

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/flucas97/delivery-much-challange/internal/domain/recipe"
	"github.com/flucas97/delivery-much-challange/pkg/services/recipeservice"
	"github.com/flucas97/delivery-much-challange/tools/errortools"
	"github.com/stretchr/testify/assert"

	"github.com/jarcoal/httpmock"
)

type getAllResponse struct {
	Keywords []string        `json:"keywords"`
	Recipes  []recipe.Recipe `json:"recipes"`
}

func TestGetAll(t *testing.T) {
	t.Run("Receive all recipes without error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", recipeservice.RecipeURI+"a,b", httpmock.NewStringResponder(http.StatusOK,
			`{
				"results":[
					{
						   "title":"recipe",
						   "href":"http://mysuperrecipe.com",
						   "ingredients":"d, a, b, c"
					}
				]
			}`))

		var (
			handler = http.HandlerFunc(RecipeController.GetAll)
			req     = httptest.NewRequest(http.MethodGet, "/recipes/?i="+"a,b", nil)
			res     = httptest.NewRecorder()
		)
		var responseAPI = getAllResponse{}

		handler.ServeHTTP(res, req)

		bytes, err := ioutil.ReadAll(res.Body)
		json.Unmarshal(bytes, &responseAPI)

		recipe := responseAPI.Recipes[0]
		assert.Nil(t, err)
		assert.NotNil(t, responseAPI)
		assert.Equal(t, "recipe", recipe.Title)
		assert.Equal(t, []string{"a", "b", "c", "d"}, recipe.Ingredients)
		assert.Equal(t, "http://mysuperrecipe.com", recipe.Link)
		assert.EqualValues(t, 1, len(responseAPI.Recipes))
	})

	t.Run("Passing no ingredients error", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", recipeservice.RecipeURI+"", httpmock.NewStringResponder(http.StatusInternalServerError, ""))

		var (
			handler = http.HandlerFunc(RecipeController.GetAll)
			req     = httptest.NewRequest(http.MethodGet, "/recipes/"+"", nil)
			res     = httptest.NewRecorder()
		)
		var errorAPI = errortools.APIError{}

		handler.ServeHTTP(res, req)

		bytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bytes, &errorAPI)

		assert.Equal(t, "enter at least one ingredient.", errorAPI.Message)
	})

	t.Run("Error on GetAll service response", func(t *testing.T) {
		httpmock.Activate()
		defer httpmock.DeactivateAndReset()
		httpmock.RegisterResponder("GET", recipeservice.RecipeURI+"a,b", httpmock.NewStringResponder(http.StatusInternalServerError, `{[error]}`))

		var (
			handler = http.HandlerFunc(RecipeController.GetAll)
			req     = httptest.NewRequest(http.MethodGet, "/recipes/?i="+"a,b", nil)
			res     = httptest.NewRecorder()
		)

		var errorAPI = errortools.APIError{}

		handler.ServeHTTP(res, req)

		bytes, _ := ioutil.ReadAll(res.Body)
		json.Unmarshal(bytes, &errorAPI)

		assert.Equal(t, "error unmarshalling response from client. recipeservice.GetAll", errorAPI.Message)
	})
}
