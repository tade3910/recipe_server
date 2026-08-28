package tests

import (
	"fmt"
	"math/rand/v2"
	"strings"
	"testing"
	"time"

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

func RandomString(length int) string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	var sb strings.Builder
	sb.Grow(length)
	for range length {
		// rand.IntN provides uniform pseudo-random distributions
		sb.WriteByte(charset[rand.IntN(len(charset))])
	}
	return sb.String()
}

func InsertUser(t *testing.T, db *gorm.DB, email string) error {
	t.Helper()
	user := &models.User{
		Email:    email,
		Name:     RandomString(10) + " " + RandomString(10),
		Password: RandomString(20), // In a real scenario, this should be a hashed password
	}
	return db.FirstOrCreate(user).Error
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
		if r.OwnerEmail != DefaultTestUser.Email {
			err := InsertUser(t, db, r.OwnerEmail)
			if err != nil {
				return fmt.Errorf("failed to insert owner user %s: %w", r.OwnerEmail, err)
			}
		}
	}

	// 3. Insert recipes
	return db.Create(recipes).Error
}

func CreateTestToken(t *testing.T, db *gorm.DB, email string) string {
	t.Helper()

	// 1. Ensure user exists first so foreign key doesn't fail
	if err := InsertUser(t, db, email); err != nil {
		t.Fatalf("failed to ensure user exists for token creation: %v", err)
	}

	// 2. Generate unique token per call to avoid unique constraint collisions
	tokenString := fmt.Sprintf("test_access_token_%s_%s", email, RandomString(6))

	token := &models.AuthToken{
		UserEmail: email,
		Token:     tokenString,
		TokenType: "access",
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(1 * time.Hour),
	}

	if err := db.Create(token).Error; err != nil {
		t.Fatalf("failed to insert test auth token: %v", err)
	}

	return tokenString
}
