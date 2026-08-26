package dto

import "github.com/tade3910/recipe_server/pkg/models"

type UpdateRecipeRequest struct {
	Title        *string            `json:"title"`
	Ingredients  *models.StringList `json:"ingredients"`
	Instructions *models.StringList `json:"instructions"`
}

type CreateRecipeRequest struct {
	URL          string            `json:"url"`
	Title        string            `json:"title"`
	Ingredients  models.StringList `json:"ingredients"`
	Instructions models.StringList `json:"instructions"`
}
