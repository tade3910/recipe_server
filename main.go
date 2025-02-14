package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"
	"github.com/tade3910/recipe_server/pkg/routes/recipe"
)

func main() {
	godotenv.Load()
	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("Could not read port from .env file")
	}
	router := http.NewServeMux()
	router.Handle("/recipe", recipe.NewRecipesHandler())
	router.Handle("/recipe/", recipe.NewRecipeHandler())
	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}
	fmt.Printf("Server listening on http://localhost:%s/\n", port)
	server.ListenAndServe()
}
