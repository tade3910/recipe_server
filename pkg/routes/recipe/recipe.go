package recipe

import (
	"fmt"
	"net/http"
	"strings"

	util "github.com/tade3910/recipe_server/pkg"
)

type recipeHandler struct {
}

func NewRecipeHandler() *recipeHandler {
	return &recipeHandler{}
}

func (handler *recipeHandler) handleGet(w http.ResponseWriter, r *http.Request, url string) {
	response_string := fmt.Sprintf("Gotten url %s", url)
	util.RespondWithJSON(w, http.StatusAccepted, response_string)
}

func (handler *recipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	url_split := strings.Split(path, "/")
	fmt.Printf("Url is %s and length is %d and method is %s\n", path, len(url_split), r.Method)
	switch len(url_split) {
	case 3:
		switch r.Method {
		case http.MethodGet:
			handler.handleGet(w, r, url_split[2])
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
