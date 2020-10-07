package router

import (
	"github.com/flucas97/delivery-much-challange/internal/controllers/recipescontroller"
)

// Routes avaliables
func Routes() {
	router.HandleFunc("/recipes/", recipescontroller.RecipeController.GetAll)
}
