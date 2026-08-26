package tests

import (
	"fmt"
	"testing"

	"github.com/tade3910/recipe_server/pkg/models"
	"gorm.io/gorm"
)

var DefaultTestUser = models.User{
	Email:    "test@example.com",
	Name:     "Test User",
	Password: "hashed_password_here",
}

func DeleteRecipes(db *gorm.DB) {
	db.Exec("DELETE FROM recipes")
}

func DeleteUsers(db *gorm.DB) {
	db.Exec("DELETE FROM users")
}

func ClearDatabase(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE recipes, users, auth_tokens RESTART IDENTITY CASCADE")
}

func InsertRecipes(t *testing.T, db *gorm.DB, recipes []*models.Recipe) error {
	t.Helper()
	if len(recipes) == 0 {
		return nil
	}

	// 1. Ensure foreign key owner exists in DB first
	if err := db.FirstOrCreate(&DefaultTestUser, models.User{Email: DefaultTestUser.Email}).Error; err != nil {
		return fmt.Errorf("failed to seed test user: %w", err)
	}

	// 2. Attach default owner email if not explicitly provided
	for _, r := range recipes {
		if r.OwnerEmail == "" {
			r.OwnerEmail = DefaultTestUser.Email
		}
	}

	// 3. Insert recipes
	return db.Create(recipes).Error
}

type TestDetails[T any] struct {
	Name             string
	Target           string
	ExpectedStatus   int
	MockRecipes      []*models.Recipe
	ExpectedResponse T
}
