package recipe

import (
	"fmt"
	"net/http"
	"strings"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	"gorm.io/gorm"
)

type recipeHandler struct {
	db *gorm.DB
}

func NewRecipeHandler(db *gorm.DB) *recipeHandler {
	return &recipeHandler{
		db: db,
	}
}

func (handler *recipeHandler) getByUrl(url string) (*models.Recipe, error) {
	recipe := &models.Recipe{
		Url: url,
	}
	result := handler.db.First(recipe)
	if result.Error != nil {
		return nil, result.Error
	}
	return recipe, nil
}

func (handler *recipeHandler) handleGet(w http.ResponseWriter, url string) {
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, recipe)
}

func (handler *recipeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := strings.TrimPrefix(r.URL.Path, "/recipe/")
	url_split := strings.Split(url, "/")
	fmt.Printf("Url is %s and length is %d and method is %s and split is %s\n", url, len(url_split), r.Method, url_split)
	switch len(url_split) {
	case 1:
		request_url := url_split[0]
		if !util.IsUrl(request_url) {
			util.RespondWithError(w, http.StatusBadRequest, fmt.Sprintf("Url: %s is not a valid url", request_url))
			return
		}
		switch r.Method {
		case http.MethodGet:
			handler.handleGet(w, url_split[0])
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
