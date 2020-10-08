package recipescontroller

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/flucas97/delivery-much-challange/tools/errortools"

	"github.com/flucas97/delivery-much-challange/internal/domain/recipe"
	"github.com/flucas97/delivery-much-challange/pkg/services/recipeservice"
)

var (
	// RecipeController interface
	RecipeController recipeControllerInterface = &recipeController{}

	rs = recipeservice.RecipeService
)

type recipeController struct{}

type recipeControllerInterface interface {
	GetAll(w http.ResponseWriter, r *http.Request)
}

// GetAll label
func (rc *recipeController) GetAll(w http.ResponseWriter, r *http.Request) {
	var (
		ingredientsParams string
		ingredients       []string
	)

	params := r.URL.Query()["i"]
	if params == nil {
		var e = errortools.APIErrorInterface.NewBadRequestError("enter at least one ingredient.")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(e)
		return
	}

	ingredientsParams = params[0]
	ingredients = strings.Split(ingredientsParams, ",")

	recipes, err := rs.GetAll(ingredients)
	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(err.Status)
		json.NewEncoder(w).Encode(err)
		return
	}

	payload := struct {
		Keywords []string        `json:"keywords"`
		Recipes  []recipe.Recipe `json:"recipes"`
	}{
		ingredients,
		recipes,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(payload)
}
