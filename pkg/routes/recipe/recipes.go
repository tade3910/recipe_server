package recipe

import (
	"fmt"
	"net/http"
	"strings"

	util "github.com/tade3910/recipe_server/pkg"
)

type recipesHandler struct {
}

func NewRecipesHandler() *recipesHandler {
	return &recipesHandler{}
}

func (handler *recipesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	util.RespondWithJSON(w, http.StatusCreated, "Posted")
}

func (handler *recipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	url_split := strings.Split(path, "/")
	fmt.Printf("Url is %s and length is %d and method is %s\n", path, len(url_split), r.Method)
	switch r.Method {
	case http.MethodPost:
		handler.handlePost(w, r)
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
