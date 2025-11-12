package routes

import (
	"net/http"

	"github.com/go-chi/chi"
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

func (handler *recipeHandler) getById(id string) (*models.Recipe, error) {
	recipe := &models.Recipe{
		Id: id,
	}
	result := handler.db.First(recipe)
	if result.Error != nil {
		return nil, result.Error
	}
	return recipe, nil
}

// Returns recipe if owner has recipe with matching url
func (handler *recipeHandler) createdRecipe(owner_email string, url string) *models.Recipe {
	recipe := &models.Recipe{
		Url:        url,
		OwnerEmail: owner_email,
	}
	result := handler.db.First(recipe)
	if result.Error != nil {
		return recipe
	}
	return nil
}

func (handler *recipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	recipe := &models.Recipe{}
	err := util.GetBody(r.Body, recipe)
	if err == nil {
		err = recipe.HasRecipeError()
	}
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	duplicate_recipe := handler.createdRecipe(recipe.OwnerEmail, recipe.Url)
	if duplicate_recipe != nil {
		util.RespondWithError(w, http.StatusConflict, duplicate_recipe)
		return
	}
	result := handler.db.Create(recipe)
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error)
	} else {
		util.RespondWithJSON(w, http.StatusCreated, recipe)
	}
}

func (handler *recipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	recipe, err := handler.getById(id)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err)
		return
	}
	util.RespondWithJSON(w, http.StatusAccepted, recipe)
}

func (handler *recipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	recipe, err := handler.getById(id)
	if err != nil {
		util.RespondWithError(w, http.StatusNotFound, err)
		return
	}
	response := handler.db.Delete(recipe)
	if response.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, response.Error)
		return
	}
	util.RespondWithJSON(w, http.StatusOK, "Deleted")
}

func (handler *recipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	updateRecipe := &models.Recipe{}
	err := util.GetBody(r.Body, updateRecipe)
	if err == nil {
		err = updateRecipe.HasRecipeError()
	}
	if err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}
	recipe, err := handler.getById(id)
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
