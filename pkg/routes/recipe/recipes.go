package recipe

import (
	"fmt"
	"net/http"
	"strings"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	"gorm.io/gorm"
)

type recipesHandler struct {
	db *gorm.DB
}

func NewRecipesHandler(db *gorm.DB) *recipesHandler {
	return &recipesHandler{
		db: db,
	}
}

func (handler *recipesHandler) handlePost(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}
	err := util.GetBody(r, recipe)
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
	}
	result := handler.db.Create(recipe)
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
	} else {
		util.RespondWithJSON(w, http.StatusCreated, recipe)
	}
}

func (handler *recipesHandler) handleGet(w http.ResponseWriter) {
	var recipes []models.Recipe
	result := handler.db.Find(&recipes)
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, recipes)
}

func (handler *recipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	path := strings.TrimSuffix(r.URL.Path, "/")
	url_split := strings.Split(path, "/")
	fmt.Printf("Url is %s and length is %d and method is %s\n", path, len(url_split), r.Method)
	switch r.Method {
	case http.MethodPost:
		handler.handlePost(w, r)
	case http.MethodGet:
		handler.handleGet(w)
	default:
		util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
	}
}
