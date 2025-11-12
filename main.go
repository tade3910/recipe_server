package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/databse"
	"github.com/tade3910/recipe_server/pkg/routes"
	"github.com/tade3910/recipe_server/pkg/services"
)

func main() {
	// Load environment variables
	loadedEnvs := util.LoadEnvs()
	if loadedEnvs.Port == "" {
		log.Fatal("Could not read port from .env file")
	} else if loadedEnvs.DbUrl == "" {
		log.Fatal("Could not read dbUrl from .env file")
	}

	db := databse.Init()
	r := chi.NewRouter()

	recipeHandler := routes.NewRecipeHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Post("/", recipeHandler.CreateRecipe)       // e.g. POST /recipe
		r.Get("/{id}", recipeHandler.GetRecipe)       // e.g. GET /recipe/123
		r.Put("/{id}", recipeHandler.UpdateRecipe)    // e.g. PUT /recipe/123
		r.Delete("/{id}", recipeHandler.DeleteRecipe) // e.g. DELETE /recips/123
	})

	recipesHandler := routes.NewRecipesHandler(db)
	r.Route("/recipes", func(r chi.Router) {
		r.Get("/{id}", recipesHandler.GetRecipes) // e.g. GET /recipe/123?page=1&limit=10
	})

	scraperHandler := routes.NewscraperHandler(db, &services.RecipeScraper{})
	r.Route("/scrape", func(r chi.Router) {
		r.Post("/", scraperHandler.ParseRecipe)
	})

	addr := ":" + loadedEnvs.Port
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
