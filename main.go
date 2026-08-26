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

	auth_handler := routes.NewAuthHandler(db)
	r.Route("/auth", func(r chi.Router) {
		r.Post("/signup", auth_handler.SignUp)
		r.Post("/signin", auth_handler.SignIn)
		r.Post("/refresh", auth_handler.SignIn)
		r.Group(func(r chi.Router) {
			r.Use(auth_handler.AuthMiddleware)
			r.Post("/logout", auth_handler.Logout)
		})
	})

	recipeHandler := routes.NewRecipeHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Use(auth_handler.AuthMiddleware)
		r.Mount("/", recipeHandler.Routes())
	})

	recipesHandler := routes.NewRecipesHandler(db)
	r.Route("/recipes", func(r chi.Router) {
		r.Use(auth_handler.AuthMiddleware)
		r.Mount("/", recipesHandler.Routes())
	})

	scraperHandler := routes.NewscraperHandler(db, &services.RecipeScraper{})
	r.Route("/scrape", func(r chi.Router) {
		r.Use(auth_handler.AuthMiddleware)
		r.Mount("/", scraperHandler.Routes())
	})

	addr := ":" + loadedEnvs.Port
	fmt.Printf("Server listening on http://localhost%s\n", addr)
	log.Fatal(http.ListenAndServe(addr, r))
}
