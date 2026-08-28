package recipe_test

import (
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/tade3910/recipe_server/pkg/models"
	"github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
	"github.com/tade3910/recipe_server/tests/mocks"
	"gorm.io/gorm"
)

type deleteTestCase struct {
	Name           string
	GetTargetID    func(db *gorm.DB) string
	ExpectedStatus int
	ExpectedBody   string
}

func setupDeleteRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	handler := routes.NewRecipeHandler(db)
	authHandler := routes.NewAuthHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Delete("/{id}", handler.DeleteRecipe)
	})
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Registered Route: [%s] %s\n", method, route)
		return nil
	})
	return r
}

func TestDeleteRecipe(t *testing.T) {
	db := mocks.TestDb(t)
	router := setupDeleteRouter(db)

	mockOwnerRecipe := &models.Recipe{
		URL:          "https://example.com/my-recipe",
		Title:        "Pancakes",
		Ingredients:  []string{"flour"},
		Instructions: []string{"cook"},
		OwnerEmail:   test_util.DefaultTestUser.Email,
	}

	mockOtherUserRecipe := &models.Recipe{
		URL:          "https://example.com/other-recipe",
		Title:        "Waffles",
		Ingredients:  []string{"flour"},
		Instructions: []string{"cook"},
		OwnerEmail:   "other_user@example.com",
	}

	tests := []deleteTestCase{
		{
			Name: "Delete valid owned recipe succeeds",
			GetTargetID: func(db *gorm.DB) string {
				var r models.Recipe
				db.First(&r, "url = ?", mockOwnerRecipe.URL)
				return r.ID
			},
			ExpectedStatus: http.StatusOK,
			ExpectedBody:   `"Deleted"`,
		},
		{
			Name: "Delete recipe owned by another user returns 404",
			GetTargetID: func(db *gorm.DB) string {
				var r models.Recipe
				db.First(&r, "url = ?", mockOtherUserRecipe.URL)
				return r.ID
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   `you are not the owner of this recipe`,
		},
		{
			Name: "Delete non-existent recipe ID returns 404",
			GetTargetID: func(db *gorm.DB) string {
				return "00000000-0000-0000-0000-000000000000"
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   `"Recipe not found"`,
		},
		{
			Name: "Delete recipe with missing ID in path",
			GetTargetID: func(db *gorm.DB) string {
				return ""
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   "404 page not found\n",
		},
		{
			Name:           "Delete recipe with invalid id",
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "invalid recipe ID format",
			GetTargetID: func(db *gorm.DB) string {
				return "invalid_id"
			},
		},
	}
	t.Cleanup(func() { test_util.ClearDatabase(db) })
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// 1. Clean slate for each sub-test
			test_util.ClearDatabase(db)
			if err := test_util.InsertRecipes(t, db, []*models.Recipe{mockOwnerRecipe, mockOtherUserRecipe}); err != nil {
				t.Fatalf("failed to seed test recipes: %v", err)
			}

			var expectedCount int64
			db.Model(&models.Recipe{}).Count(&expectedCount)

			validToken := test_util.CreateTestToken(t, db, test_util.DefaultTestUser.Email)

			// 2. Build target URL path dynamically
			targetID := tt.GetTargetID(db)
			urlPath := "/recipe/" + targetID
			if targetID == "" {
				urlPath = "/recipe/"
			}

			// 3. Make HTTP request
			req := httptest.NewRequest(http.MethodDelete, urlPath, nil)
			req.Header.Set("Authorization", "Bearer "+validToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			AssertResponse(t, w, tt.Name, tt.ExpectedStatus, tt.ExpectedBody)

			// 6. DB sanity check for successful deletes
			if tt.ExpectedStatus == http.StatusOK {
				var count int64
				db.Model(&models.Recipe{}).Where("id = ?", targetID).Count(&count)
				if count != 0 {
					t.Errorf("%s: recipe was expected to be deleted from DB, but still exists", tt.Name)
				}
				expectedCount--
			}

			var finalCount int64
			db.Model(&models.Recipe{}).Count(&finalCount)

			if finalCount != expectedCount {
				t.Errorf("expected DB count to decrease by 1 (want %d, got %d)", expectedCount, finalCount)
			}
		})
	}
}
