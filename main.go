package main

import (
	"fmt"
	"log"
	"net/http"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/databse"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
)

func main() {
	loadedEnvs := util.LoadEnvs()
	if loadedEnvs.Port == "" {
		log.Fatal("Could not read port from .env file")
	} else if loadedEnvs.DbUrl == "" {
		log.Fatal("Could not read dbUrl from .env file")
	}
	db := databse.Init()
	router := http.NewServeMux()
	router.Handle("/recipe", recipe.NewRecipesHandler(db))
	server := &http.Server{
		Addr:    ":" + loadedEnvs.Port,
		Handler: router,
	}
	fmt.Printf("Server listening on http://localhost:%s/\n", loadedEnvs.Port)
	server.ListenAndServe()
}
