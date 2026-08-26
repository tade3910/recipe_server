package routes

import (
	"errors"
	"fmt"
	"net/http"
	"strings"

	"github.com/go-chi/chi"
	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/dto"
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

func isUniqueConstraintError(err error) bool {
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return true
	}
	// Catches postgres unique constraint violation error code 23505
	return strings.Contains(err.Error(), "23505") || strings.Contains(err.Error(), "idx_owner_url")
}

func (handler *recipeHandler) getByIdAndOwner(id string, ownerEmail string) (*models.Recipe, error) {
	var recipe models.Recipe
	err := handler.db.Where("id = ? AND owner_email = ?", id, ownerEmail).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (handler *recipeHandler) getById(id string) (*models.Recipe, error) {
	var recipe models.Recipe
	err := handler.db.Where("id = ?", id).First(&recipe).Error
	if err != nil {
		return nil, err
	}
	return &recipe, nil
}

func (handler *recipeHandler) requireId(w http.ResponseWriter, r *http.Request) (string, bool) {
	id := chi.URLParam(r, "id")
	if id == "" {
		util.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("missing recipe id"))
		return "", false
	}
	return id, true
}

func (handler *recipeHandler) CreateRecipe(w http.ResponseWriter, r *http.Request) {
	userEmail, ok := RequireUserEmail(w, r)
	if !ok {
		return
	}
	var req dto.CreateRecipeRequest
	if err := util.GetBody(r.Body, &req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	var reqErrors []string
	if strings.TrimSpace(req.URL) == "" {
		reqErrors = append(reqErrors, "url")
	}
	if strings.TrimSpace(req.Title) == "" {
		reqErrors = append(reqErrors, "title")
	}
	if len(req.Ingredients) == 0 {
		reqErrors = append(reqErrors, "ingredients")
	}
	if len(req.Instructions) == 0 {
		reqErrors = append(reqErrors, "instructions")
	}

	if len(reqErrors) > 0 {
		err := fmt.Errorf("the following required fields are missing or empty: %s", strings.Join(reqErrors, ", "))
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	recipe := &models.Recipe{
		URL:          strings.TrimSpace(req.URL),
		Title:        strings.TrimSpace(req.Title),
		Ingredients:  req.Ingredients,
		Instructions: req.Instructions,
		OwnerEmail:   userEmail,
	}

	result := handler.db.Create(recipe)
	if result.Error != nil {
		if isUniqueConstraintError(result.Error) {
			util.RespondWithError(w, http.StatusConflict, fmt.Errorf("a recipe with this URL already exists for your account"))
			return
		}
		util.RespondWithError(w, http.StatusInternalServerError, result.Error)
		return
	}

	util.RespondWithJSON(w, http.StatusCreated, recipe)
}

func (handler *recipeHandler) GetRecipe(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.requireId(w, r)
	if !ok {
		return
	}

	userEmail, ok := RequireUserEmail(w, r)
	if !ok {
		return
	}
	recipe, err := handler.getByIdAndOwner(id, userEmail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			if _, err := handler.getById(id); err != nil {
				util.RespondWithError(w, http.StatusNotFound, fmt.Errorf("Recipe not found"))
			} else {
				util.RespondWithError(w, http.StatusForbidden, fmt.Errorf("Unauthorized: you are not the owner of this recipe"))
			}
		}
		util.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}
	util.RespondWithJSON(w, http.StatusOK, recipe)
}

func (handler *recipeHandler) DeleteRecipe(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.requireId(w, r)
	if !ok {
		return
	}
	userEmail, ok := RequireUserEmail(w, r)
	if !ok {
		return
	}
	result := handler.db.Where("id = ? AND owner_email = ?", id, userEmail).Delete(&models.Recipe{})
	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error)
		return
	}
	if result.RowsAffected == 0 {
		if _, err := handler.getById(id); err != nil {
			util.RespondWithError(w, http.StatusNotFound, fmt.Errorf("Recipe not found"))
		} else {
			util.RespondWithError(w, http.StatusForbidden, fmt.Errorf("Unauthorized: you are not the owner of this recipe"))
		}
		return
	}
	util.RespondWithJSON(w, http.StatusOK, "Deleted")
}

func (handler *recipeHandler) UpdateRecipe(w http.ResponseWriter, r *http.Request) {
	id, ok := handler.requireId(w, r)
	if !ok {
		return
	}

	userEmail, ok := RequireUserEmail(w, r)
	if !ok {
		return
	}

	var req dto.UpdateRecipeRequest
	if err := util.GetBody(r.Body, &req); err != nil {
		util.RespondWithError(w, http.StatusBadRequest, err)
		return
	}

	updates := make(map[string]any)

	if req.URL != nil {
		updates["url"] = *req.URL
	}

	if req.Title != nil {
		if *req.Title == "" {
			util.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("title cannot be empty"))
			return
		}
		updates["title"] = *req.Title
	}

	if req.Title != nil {
		if *req.Title == "" {
			util.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("title cannot be empty"))
			return
		}
		updates["title"] = *req.Title
	}

	if req.Ingredients != nil {
		updates["ingredients"] = *req.Ingredients
	}

	if req.Instructions != nil {
		updates["instructions"] = *req.Instructions
	}

	if len(updates) == 0 {
		util.RespondWithError(w, http.StatusBadRequest, fmt.Errorf("at least one field (title, ingredients, or instructions) must be provided"))
		return
	}

	result := handler.db.Model(&models.Recipe{}).
		Where("id = ? AND owner_email = ?", id, userEmail).
		Updates(updates)

	if result.Error != nil {
		util.RespondWithError(w, http.StatusInternalServerError, result.Error)
		return
	}

	if result.RowsAffected == 0 {
		if _, err := handler.getById(id); err != nil {
			util.RespondWithError(w, http.StatusNotFound, fmt.Errorf("Recipe not found"))
		} else {
			util.RespondWithError(w, http.StatusForbidden, fmt.Errorf("Unauthorized: you are not the owner of this recipe"))
		}
		return
	}

	updatedRecipe, err := handler.getByIdAndOwner(id, userEmail)
	if err != nil {
		util.RespondWithError(w, http.StatusInternalServerError, err)
		return
	}

	util.RespondWithJSON(w, http.StatusOK, updatedRecipe)
}

func (h *recipeHandler) Routes() http.Handler {
	r := chi.NewRouter()
	r.Post("/", h.CreateRecipe)
	r.Get("/{id}", h.GetRecipe)
	r.Put("/{id}", h.UpdateRecipe)
	r.Delete("/{id}", h.DeleteRecipe)
	return r
}
