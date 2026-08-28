package recipes_test

import (
	"encoding/json"
	"fmt"
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

func setupGetRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	handler := routes.NewRecipesHandler(db)
	authHandler := routes.NewAuthHandler(db)
	r.Route("/recipes", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Get("/", handler.GetRecipes)
	})
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Registered Route: [%s] %s\n", method, route)
		return nil
	})
	return r
}

type GetRecipesTestEnv struct {
	ValidToken   string
	ownerRecipes map[string]*models.Recipe
}

// seed test returning token for auth
func setupTestEnviroment(t *testing.T, db *gorm.DB) *GetRecipesTestEnv {
	t.Helper()
	test_util.ClearDatabase(db)
	ownerRecipes := make(map[string]*models.Recipe)

	for i := range 15 {
		r := models.Recipe{
			Title:        fmt.Sprintf("Owner Recipe %d", i),
			URL:          fmt.Sprintf("https://example.com/recipe/%d", i),
			Ingredients:  models.StringList{"salt", "water"},
			Instructions: models.StringList{"mix", "cook"},
			OwnerEmail:   test_util.DefaultTestUser.Email,
		}
		err := test_util.InsertRecipes(
			t, db,
			[]*models.Recipe{&r},
		)
		if err != nil {
			t.Fatalf("failed to seed recipe with %s", err.Error())
		}
		db.Create(&r)
		var recipe models.Recipe
		err = db.Where("URL = ?", r.URL).First(&recipe).Error
		if err != nil {
			t.Fatalf("Failed to get recipe from db with :%s", err.Error())
		}
		ownerRecipes[recipe.ID] = &recipe
	}

	// 3. Seed 2 Recipes for Other User (Isolation Check)
	for i := range 2 {
		r := models.Recipe{
			Title:        fmt.Sprintf("Other Recipe %d", i),
			URL:          fmt.Sprintf("https://example.com/other/%d", i),
			Ingredients:  models.StringList{"sugar"},
			Instructions: models.StringList{"bake"},
			OwnerEmail:   "other_user@example.com",
		}
		db.Create(&r)
	}
	return &GetRecipesTestEnv{
		ValidToken:   test_util.CreateTestToken(t, db, test_util.DefaultTestUser.Email),
		ownerRecipes: ownerRecipes,
	}
}

func TestGetRecipes(t *testing.T) {
	db := mocks.TestDb(t)
	router := setupGetRouter(db)

	t.Cleanup(func() { test_util.ClearDatabase(db) })
	testEnv := setupTestEnviroment(t, db)
	tests := []struct {
		Name           string
		QueryString    string
		ExpectedStatus int
		ExpectedCount  int // Number of items expected in response slice
	}{
		{
			Name:           "Get with no query",
			QueryString:    "",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  10,
		},
		{
			Name:           "Get with limit lt num recipes",
			QueryString:    "?limit=3",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  3,
		},
		{
			Name:           "3. Get with limit gt num recipes",
			QueryString:    "?limit=20",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  15,
		},
		{
			Name:           "3. Get with limit lt 1",
			QueryString:    "?limit=-1",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  10,
		},
		{
			Name:           "Get with page lt 1",
			QueryString:    "?page=-1",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  10,
		},
		{
			Name:           "5. Get with page > total pages -> returns empty list []",
			QueryString:    "?limit=2&page=99",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  0,
		},
		{
			Name:           "6. Get with unknown query string -> ignores unknown params and uses defaults",
			QueryString:    "?foo=bar&unknown_key=123",
			ExpectedStatus: http.StatusOK,
			ExpectedCount:  10,
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			urlPath := "/recipes" + tt.QueryString
			req := httptest.NewRequest(http.MethodGet, urlPath, nil)
			req.Header.Set("Authorization", "Bearer "+testEnv.ValidToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Basic status check
			if w.Code != tt.ExpectedStatus {
				t.Fatalf("%s: expected status %d, got %d. Body: %s", tt.Name, tt.ExpectedStatus, w.Code, w.Body.String())
			}

			// Validate array item count
			var fetchedRecipes []models.Recipe
			if err := json.Unmarshal(w.Body.Bytes(), &fetchedRecipes); err != nil {
				t.Fatalf("%s: failed to unmarshal response: %v", tt.Name, err)
			}

			if len(fetchedRecipes) != tt.ExpectedCount {
				t.Errorf("%s: expected %d recipes, got %d", tt.Name, tt.ExpectedCount, len(fetchedRecipes))
			}
		})
	}
}

func TestGetRecipesPagination(t *testing.T) {
	db := mocks.TestDb(t)
	router := setupGetRouter(db)

	t.Cleanup(func() { test_util.ClearDatabase(db) })
	testEnv := setupTestEnviroment(t, db)
	testName := "Test pagination"
	t.Run(testName, func(t *testing.T) {
		limit := 2
		pages := len(testEnv.ownerRecipes)/limit + 2
		count := len(testEnv.ownerRecipes)
		receivedRecipes := make(map[string]*models.Recipe)
		for page := range pages {
			urlPath := fmt.Sprintf("/recipes?limit=%d&page=%d", limit, page+1)
			expectedCount := limit
			if count < limit {
				expectedCount = max(count, 0)
			}
			count -= limit
			req := httptest.NewRequest(http.MethodGet, urlPath, nil)
			req.Header.Set("Authorization", "Bearer "+testEnv.ValidToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// Basic status check
			if w.Code != http.StatusOK {
				t.Fatalf("%s: expected status %d, got %d. Body: %s", testName, http.StatusOK, w.Code, w.Body.String())
			}

			// Validate array item count
			var fetchedRecipes []models.Recipe
			if err := json.Unmarshal(w.Body.Bytes(), &fetchedRecipes); err != nil {
				t.Fatalf("%s: failed to unmarshal response: %v", testName, err)
			}

			if len(fetchedRecipes) != expectedCount {
				t.Errorf("%s: expected %d recipes, got %d", testName, expectedCount, len(fetchedRecipes))
			}
			// Validate individual recipe fields
			for _, recipe := range fetchedRecipes {
				expectedRecipe, ok := testEnv.ownerRecipes[recipe.ID]
				if !ok {
					t.Fatalf("Can not find recipe matching %v in expected\n%v", recipe, testEnv.ownerRecipes)
				}
				if !recipe.Equals(expectedRecipe) {
					t.Fatalf("%s: expected recipe %+v, got %+v", testName, expectedRecipe, recipe)
				}
				_, ok = receivedRecipes[recipe.ID]
				if ok {
					t.Fatalf("Received duplicate recipe: %v, all received\n%v", recipe, receivedRecipes)
				}
				receivedRecipes[recipe.ID] = &recipe
			}
		}
		if len(receivedRecipes) != len(testEnv.ownerRecipes) {
			t.Fatalf("received:%d, expected:%d", len(receivedRecipes), len(testEnv.ownerRecipes))
		}
	})
}
