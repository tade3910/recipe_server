package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"

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
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   loadedEnvs.AllowedOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
		AllowCredentials: true,
		MaxAge:           300, // 5 minutes
	}))

	authHandler := routes.NewAuthHandler(db)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", authHandler.SignUp)
		r.Post("/signin", authHandler.SignIn)
		r.Post("/refresh", authHandler.SignIn)
		r.Group(func(r chi.Router) {
			r.Use(authHandler.AuthMiddleware)
			r.Post("/logout", authHandler.Logout)
		})
	})

	recipeHandler := routes.NewRecipeHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Post("/", recipeHandler.CreateRecipe)
		r.Get("/{id}", recipeHandler.GetRecipe)
		r.Put("/{id}", recipeHandler.UpdateRecipe)
		r.Delete("/{id}", recipeHandler.DeleteRecipe)
	})

	recipesHandler := routes.NewRecipesHandler(db)
	r.Route("/recipes", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Get("/", recipesHandler.GetRecipes)
	})

	scraperHandler := routes.NewscraperHandler(db, &services.RecipeScraper{})
	r.Route("/scrape", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Post("/", scraperHandler.ParseRecipe)
	})

	addr := ":" + loadedEnvs.Port
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	// Log all registered routes
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Registered Route: [%s] %s\n", method, route)
		return nil
	})
	log.Fatal(http.ListenAndServe(addr, r))
}
