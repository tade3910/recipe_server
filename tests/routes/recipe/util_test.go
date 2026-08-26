package recipe_test

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	"gorm.io/gorm"
)

func setupRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	handler := recipe.NewRecipeHandler(db)
	r.Mount("/recipes", handler.Routes())
	return r
}
