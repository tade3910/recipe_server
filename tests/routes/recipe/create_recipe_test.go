package recipe_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-chi/chi"
	"github.com/tade3910/recipe_server/pkg/dto"
	"github.com/tade3910/recipe_server/pkg/models"
	"github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
	"github.com/tade3910/recipe_server/tests/mocks"
	"gorm.io/gorm"
)

type createTestCase struct {
	Name           string
	Recipe         any
	ExpectedStatus int
	ExpectedBody   any
}

func SetupCreateRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	handler := routes.NewRecipeHandler(db)
	authHandler := routes.NewAuthHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Post("/", handler.CreateRecipe)
	})
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Registered Route: [%s] %s\n", method, route)
		return nil
	})
	return r
}

func TestCreateRecipe(t *testing.T) {
	db := mocks.TestDb(t)
	router := SetupCreateRouter(db)

	validCreateRequest := &dto.CreateRecipeRequest{
		URL:          "https://example.com/recipe",
		Title:        "Test Title",
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
	}

	validRecipe := &models.Recipe{
		URL:          "https://example.com/recipe",
		Title:        "Test Title",
		Ingredients:  []string{"ing1"},
		Instructions: []string{"step1"},
		OwnerEmail:   test_util.DefaultTestUser.Email,
	}

	tests := []createTestCase{
		{
			Name:           "Create recipe with all required fields",
			Recipe:         validCreateRequest,
			ExpectedStatus: http.StatusCreated,
			ExpectedBody:   validRecipe,
		},
		{
			Name: "Create recipe with missing URL",
			Recipe: &dto.CreateRecipeRequest{
				Title:        "Test Title",
				Ingredients:  []string{"ing1"},
				Instructions: []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: url",
		},
		{
			Name: "Create recipe with missing title",
			Recipe: &dto.CreateRecipeRequest{
				URL:          "https://example.com/recipe",
				Ingredients:  []string{"ing1"},
				Instructions: []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: title",
		},
		{
			Name: "Create recipe with missing ingredients",
			Recipe: &dto.CreateRecipeRequest{
				URL:          "https://example.com/recipe",
				Title:        "Test Title",
				Instructions: []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: ingredients",
		},
		{
			Name: "Create recipe with missing instructions",
			Recipe: &dto.CreateRecipeRequest{
				URL:         "https://example.com/recipe",
				Title:       "Test Title",
				Ingredients: []string{"ing1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: instructions",
		},
		{
			Name: "Create recipe with owner email",
			Recipe: &models.Recipe{
				URL:          "https://example.com/recipe",
				Title:        "Test Title",
				Ingredients:  []string{"ing1"},
				Instructions: []string{"step1"},
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "unknown field",
		},
		{
			Name:           "Create recipe with no body",
			Recipe:         nil,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "request body cannot be empty",
		},
		{
			Name: "Create recipe with no URL",
			Recipe: map[string]any{
				"Title":        "Test Title",
				"Ingredients":  []string{"ing1"},
				"Instructions": []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: url",
		},
		{
			Name: "Create recipe with no title",
			Recipe: map[string]any{
				"URL":          "https://example.com/recipe",
				"Ingredients":  []string{"ing1"},
				"Instructions": []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: title",
		},
		{
			Name: "Create recipe with no ingredients",
			Recipe: map[string]any{
				"URL":          "https://example.com/recipe",
				"Title":        "Test Title",
				"Instructions": []string{"step1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: ingredients",
		},
		{
			Name: "Create recipe with no instructions",
			Recipe: map[string]any{
				"URL":         "https://example.com/recipe",
				"Title":       "Test Title",
				"Ingredients": []string{"ing1"},
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "required fields are missing or empty: instructions",
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// 1. Clean slate for each sub-test
			test_util.ClearDatabase(db)
			t.Cleanup(func() { test_util.ClearDatabase(db) })

			err := test_util.InsertUser(t, db, test_util.DefaultTestUser.Email)
			if err != nil {
				t.Fatalf("failed to insert new user: %v", err)
			}
			validToken := test_util.CreateTestToken(t, db, test_util.DefaultTestUser.Email)

			// 2. Build target URL path dynamically
			var body io.Reader
			if tt.Recipe != nil {
				jsonBody, err := json.Marshal(tt.Recipe)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				body = bytes.NewBuffer(jsonBody)
			}

			// 2. Pass bytes.NewBuffer(jsonBody) as the io.Reader
			req := httptest.NewRequest(http.MethodPost, "/recipe/", body)
			req.Header.Set("Authorization", "Bearer "+validToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 3. Make HTTP request
			if w.Code != tt.ExpectedStatus {
				t.Fatalf("%s: expected status %d, got %d. Body: %s", tt.Name, tt.ExpectedStatus, w.Code, w.Body.String())
			}

			// 4. Test ExpectedBody
			resp := AssertResponse(t, w, tt.Name, tt.ExpectedStatus, tt.ExpectedBody)
			if resp == nil {
				return
			}
			// 5. Test actually in db
			AssertRecipeInDb(
				t, db, resp, tt.Name,
			)
		})
	}
}
