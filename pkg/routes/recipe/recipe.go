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

func (handler *recipeHandler) hanldeDelete(w http.ResponseWriter, url string) {
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	response := handler.db.Delete(recipe)
	if response.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, response.Error.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusOK, "Deleted")
}

func (handler *recipeHandler) handleUpdate(w http.ResponseWriter, r *http.Request, url string) {
	updateRecipe := &models.Recipe{}
	err := util.GetBody(r, updateRecipe)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}
	updateRecipe.Url = recipe.Url
	response := handler.db.Save(updateRecipe)
	if response.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, response.Error.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, updateRecipe)
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
			handler.handleGet(w, request_url)
		case http.MethodDelete:
			handler.hanldeDelete(w, request_url)
		case http.MethodPut:
			handler.handleUpdate(w, r, request_url)
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
