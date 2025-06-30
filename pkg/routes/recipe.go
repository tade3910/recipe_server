package recipe

import (
	"net/http"

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

func (handler *recipesHandler) getByUrl(url string) (*models.Recipe, error) {
	recipe := &models.Recipe{
		Url: url,
	}
	result := handler.db.First(recipe)
	if result.Error != nil {
		return nil, result.Error
	}
	return recipe, nil
}

func (handler *recipesHandler) PostRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}
	err := util.GetBody(r.Body, recipe)
	if err == nil {
		err = recipe.HasRecipeError()
	}
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	duplicate_recipe, _ := handler.getByUrl(recipe.Url)
	if duplicate_recipe != nil {
		util.RespondWithError(w, http.StatusConflict, "Duplicate recipe")
		return
	}
	result := handler.db.Create(recipe)
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
	} else {
		util.RespondWithJSON(w, http.StatusCreated, recipe)
	}
}

func (handler *recipesHandler) getAllRecipes(w http.ResponseWriter) {
	// TODO:Tokenization to not post all
	var recipes []models.Recipe
	result := handler.db.Find(&recipes)
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, recipes)
}

func (handler *recipesHandler) getRecipe(w http.ResponseWriter, url string) {
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, recipe)
}

func (handler *recipesHandler) deleteRecipe(w http.ResponseWriter, url string) {
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err.Error())
		return
	}
	response := handler.db.Delete(recipe)
	if response.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, response.Error.Error())
		return
	}
	util.RespondWithJSON(w, http.StatusOK, "Deleted")
}

func (handler *recipesHandler) updateRecipe(w http.ResponseWriter, r *http.Request, url string) {
	updateRecipe := &models.Recipe{}
	err := util.GetBody(r.Body, updateRecipe)
	if err == nil {
		err = updateRecipe.HasRecipeError()
	}
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err.Error())
		return
	}
	recipe, err := handler.getByUrl(url)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err.Error())
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

func (handler *recipesHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	url := r.URL.Query().Get("url")
	switch url {
	case "":
		switch r.Method {
		case http.MethodPost:
			handler.PostRecipe(w, r)
		case http.MethodGet:
			handler.getAllRecipes(w)
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	default:
		switch r.Method {
		case http.MethodGet:
			handler.getRecipe(w, url)
		case http.MethodDelete:
			handler.deleteRecipe(w, url)
		case http.MethodPut:
			handler.updateRecipe(w, r, url)
		default:
			util.RespondWithError(w, http.StatusMethodNotAllowed, "Invalid method")
		}
	}
}
