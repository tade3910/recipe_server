package main

import (
	"fmt"
	"net/http"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/databse"
	"github.com/tade3910/recipe_server/pkg/routes/recipe"
)

func main() {
	loadedEnvs := util.LoadEnvs()
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
