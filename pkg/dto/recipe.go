package dto

import (
	"fmt"
	"strings"

	"github.com/tade3910/recipe_server/pkg/models"
)

type CreateRecipeRequest struct {
	URL          string            `json:"url"`
	Title        string            `json:"title"`
	Ingredients  models.StringList `json:"ingredients"`
	Instructions models.StringList `json:"instructions"`
}

func (r *CreateRecipeRequest) Validate() error {
	var reqErrors []string
	if strings.TrimSpace(r.URL) == "" {
		reqErrors = append(reqErrors, "url")
	}
	if strings.TrimSpace(r.Title) == "" {
		reqErrors = append(reqErrors, "title")
	}
	if len(r.Ingredients) == 0 {
		reqErrors = append(reqErrors, "ingredients")
	}
	if len(r.Instructions) == 0 {
		reqErrors = append(reqErrors, "instructions")
	}

	if len(reqErrors) > 0 {
		return fmt.Errorf("the following required fields are missing or empty: %s", strings.Join(reqErrors, ", "))
	}
	return nil
}

func (r *CreateRecipeRequest) ToRecipe(userEmail string) *models.Recipe {
	return &models.Recipe{
		URL:          strings.TrimSpace(r.URL),
		Title:        strings.TrimSpace(r.Title),
		Ingredients:  r.Ingredients,
		Instructions: r.Instructions,
		OwnerEmail:   userEmail,
	}
}

type UpdateRecipeRequest struct {
	URL          *string            `json:"url"`
	Title        *string            `json:"title"`
	Ingredients  *models.StringList `json:"ingredients"`
	Instructions *models.StringList `json:"instructions"`
}
