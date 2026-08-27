package routes

import (
	"net/http"
	"strconv"

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

func (handler *recipesHandler) GetRecipes(w http.ResponseWriter, r *http.Request) {
	pageStr := r.URL.Query().Get("page")
	limitStr := r.URL.Query().Get("limit")
	userEmail, ok := RequireUserEmail(w, r)
	if !ok {
		return
	}

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		limit = 10
	}

	offset := (page - 1) * limit

	var recipes []models.Recipe
	result := handler.db.Model(&models.Recipe{}).
		Where("owner_email = ?", userEmail).
		Limit(limit).
		Offset(offset).
		Find(&recipes)

	if result.Error != nil {
		util.RespondWithError(w, http.StatusNotFound, result.Error)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, recipes)
}
