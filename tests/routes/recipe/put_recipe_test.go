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

func SetupPutRouter(db *gorm.DB) http.Handler {
	r := chi.NewRouter()
	handler := routes.NewRecipeHandler(db)
	authHandler := routes.NewAuthHandler(db)
	r.Route("/recipe", func(r chi.Router) {
		r.Use(authHandler.AuthMiddleware)
		r.Put("/{id}", handler.UpdateRecipe)
	})
	chi.Walk(r, func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		log.Printf("Registered Route: [%s] %s\n", method, route)
		return nil
	})
	return r
}

type updateTestCase struct {
	Name           string
	GetTargetID    func(db *gorm.DB) string
	ExpectedStatus int
	ExpectedBody   any
	Updates        any
}

func TestPutRecipe(t *testing.T) {
	db := mocks.TestDb(t)
	router := SetupPutRouter(db)

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

	updatedUrl := "https://example.com/new-recipe"
	updatedTitle := "New Pancakes"
	updatedIngredients := []string{"milk", "eggs"}
	updatedIngredientsSl := models.StringList(updatedIngredients)
	updatedInstructions := []string{"cook", "bake"}
	updatedInstructionsSl := models.StringList(updatedInstructions)
	var emptyString string = ""
	var emptyStringSl models.StringList = []string{}

	getOwnRecipeId := func(db *gorm.DB) string {
		var r models.Recipe
		db.First(&r, "url = ?", mockOwnerRecipe.URL)
		return r.ID
	}

	tests := []updateTestCase{
		{
			Name:           "Update recipe with all fields",
			GetTargetID:    getOwnRecipeId,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: &models.Recipe{
				URL:          updatedUrl,
				Title:        updatedTitle,
				Ingredients:  updatedIngredients,
				Instructions: updatedInstructions,
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			Updates: &dto.UpdateRecipeRequest{
				URL:          &updatedUrl,
				Title:        &updatedTitle,
				Ingredients:  &updatedIngredientsSl,
				Instructions: &updatedInstructionsSl,
			},
		},
		{
			Name:           "Update recipe with new url",
			GetTargetID:    getOwnRecipeId,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: &models.Recipe{
				URL:          updatedUrl,
				Title:        mockOwnerRecipe.Title,
				Ingredients:  mockOwnerRecipe.Ingredients,
				Instructions: mockOwnerRecipe.Instructions,
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			Updates: &dto.UpdateRecipeRequest{
				URL: &updatedUrl,
			},
		},
		{
			Name:           "Update recipe with new title",
			GetTargetID:    getOwnRecipeId,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: &models.Recipe{
				URL:          mockOwnerRecipe.URL,
				Title:        updatedTitle,
				Ingredients:  mockOwnerRecipe.Ingredients,
				Instructions: mockOwnerRecipe.Instructions,
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			Updates: &dto.UpdateRecipeRequest{
				Title: &updatedTitle,
			},
		},
		{
			Name:           "Update recipe with new Ingredients",
			GetTargetID:    getOwnRecipeId,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: &models.Recipe{
				URL:          mockOwnerRecipe.URL,
				Title:        mockOwnerRecipe.Title,
				Ingredients:  updatedIngredients,
				Instructions: mockOwnerRecipe.Instructions,
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			Updates: &dto.UpdateRecipeRequest{
				Ingredients: &updatedIngredientsSl,
			},
		},
		{
			Name:           "Update recipe with new Instructions",
			GetTargetID:    getOwnRecipeId,
			ExpectedStatus: http.StatusOK,
			ExpectedBody: &models.Recipe{
				URL:          mockOwnerRecipe.URL,
				Title:        mockOwnerRecipe.Title,
				Ingredients:  mockOwnerRecipe.Ingredients,
				Instructions: updatedInstructions,
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			Updates: &dto.UpdateRecipeRequest{
				Instructions: &updatedInstructionsSl,
			},
		},

		{
			Name: "Update recipe with missing title",
			Updates: &dto.UpdateRecipeRequest{
				Title: &emptyString,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "Title cannot be empty",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with missing URL",
			Updates: &dto.UpdateRecipeRequest{
				URL: &emptyString,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "URL cannot be empty",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with empty ingredients",
			Updates: &dto.UpdateRecipeRequest{
				Ingredients: &emptyStringSl,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "Ingredients cannot be empty",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with empty instructions",
			Updates: &dto.UpdateRecipeRequest{
				Instructions: &emptyStringSl,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "Instructions cannot be empty",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with owner email",
			Updates: &models.Recipe{
				URL:          "https://example.com/recipe",
				Title:        "Test Title",
				Ingredients:  []string{"ing1"},
				Instructions: []string{"step1"},
				OwnerEmail:   test_util.DefaultTestUser.Email,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "unknown field",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name:           "Update recipe with no body",
			Updates:        nil,
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "request body cannot be empty",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with different owner",
			Updates: &dto.UpdateRecipeRequest{
				Title: &updatedTitle,
			},
			ExpectedStatus: http.StatusForbidden,
			ExpectedBody:   "you are not the owner of this recipe",
			GetTargetID: func(db *gorm.DB) string {
				var r models.Recipe
				db.First(&r, "url = ?", mockOtherUserRecipe.URL)
				return r.ID
			},
		},
		{
			Name:           "Update recipe with empty body",
			Updates:        map[string]string{},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "at least one field (title, ingredients, or instructions) must be provided",
			GetTargetID:    getOwnRecipeId,
		},
		{
			Name: "Update recipe with not found id",
			Updates: &dto.UpdateRecipeRequest{
				Instructions: &updatedInstructionsSl,
			},
			ExpectedStatus: http.StatusNotFound,
			ExpectedBody:   "Recipe not found",
			GetTargetID: func(db *gorm.DB) string {
				return "00000000-0000-0000-0000-000000000000"
			},
		},
		{
			Name: "Update recipe with invalid id",
			Updates: &dto.UpdateRecipeRequest{
				Instructions: &updatedInstructionsSl,
			},
			ExpectedStatus: http.StatusBadRequest,
			ExpectedBody:   "invalid recipe ID format",
			GetTargetID: func(db *gorm.DB) string {
				return "invalid_id"
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			// 1. Clean slate for each sub-test
			test_util.ClearDatabase(db)
			if err := test_util.InsertRecipes(t, db, []*models.Recipe{mockOwnerRecipe, mockOtherUserRecipe}); err != nil {
				t.Fatalf("failed to seed test recipes: %v", err)
			}
			t.Cleanup(func() { test_util.ClearDatabase(db) })

			validToken := test_util.CreateTestToken(t, db, test_util.DefaultTestUser.Email)

			// 2. Prepare the request
			var request io.Reader
			if tt.Updates != nil {
				jsonBody, err := json.Marshal(tt.Updates)
				if err != nil {
					t.Fatalf("failed to marshal request body: %v", err)
				}
				request = bytes.NewBuffer(jsonBody)
			}

			// 2. Make the request
			targetID := tt.GetTargetID(db)
			urlPath := "/recipe/" + targetID
			if targetID == "" {
				urlPath = "/recipe/"
			}
			log.Printf("url path is %s\n", urlPath)

			req := httptest.NewRequest(http.MethodPut, urlPath, request)
			req.Header.Set("Authorization", "Bearer "+validToken)

			w := httptest.NewRecorder()
			router.ServeHTTP(w, req)

			// 3. Make the assertions
			resp := AssertResponse(t, w, tt.Name, tt.ExpectedStatus, tt.ExpectedBody)
			if resp == nil {
				return
			}
			AssertRecipeInDb(
				t, db, resp, tt.Name,
			)
		})
	}

}
