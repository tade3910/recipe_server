package recipe_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	test_util "github.com/tade3910/recipe_server/tests"
	"github.com/tade3910/recipe_server/tests/mocks"
	"gorm.io/gorm"
)

func validateResponse(testDetail *test_util.TestDetails[string], response *http.Response, db *gorm.DB, t *testing.T) {
	t.Helper()
	if testDetail.ExpectedResponse == "" {
		return
	}
	actual := ""
	err := util.GetBody(response.Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	if actual != testDetail.ExpectedResponse {
		t.Fatalf("Mismatch: got %+s, want %+s", actual, testDetail.ExpectedResponse)
	}

	// Check if the mock recipes are still in the database
	for _, mockRecipe := range testDetail.MockRecipes {
		var recipe models.Recipe
		// Query directly by the string Primary Key ID
		err := db.First(&recipe, "id = ?", mockRecipe.ID).Error
		if err == nil {
			t.Errorf("Recipe with ID %s still exists in the database", mockRecipe.Id)
		} else if err != gorm.ErrRecordNotFound {
			t.Errorf("Error checking if recipe with ID %s exists: %s", mockRecipe.Id, err.Error())
		}
	}
}

func TestDeleteRecipe(t *testing.T) {
	db := mocks.TestDb(t)
	tests := []test_util.TestDetails[string]{
		{
			Name:           "Missing target URL",
			Target:         "/recipe",
			ExpectedStatus: http.StatusBadRequest,
		},
		{
			Name:           "Delete valid recipe",
			Target:         fmt.Sprintf("/recipe?url=%s", "https://example.com"),
			ExpectedStatus: http.StatusOK,
			MockRecipes: []*models.Recipe{
				{
					URL:          "https://example.com",
					Title:        "Test Pancakes",
					Ingredients:  []string{"flour", "milk", "egg"},
					Instructions: []string{"mix ingredients", "cook on pan"},
				},
			},
			ExpectedResponse: "Deleted",
		},
		{
			Name:             "Delete not found recipe",
			Target:           fmt.Sprintf("/recipe?url=%s", "unknown"),
			ExpectedStatus:   http.StatusNotFound,
			ExpectedResponse: "Not Found",
		},
		{
			Name:           "Delete recipe from differnt user",
			Target:         fmt.Sprintf("/recipe?url=%s", "https://example.com"),
			ExpectedStatus: http.StatusNotFound,
			MockRecipes: []*models.Recipe{
				{
					URL:          "https://example.com",
					Title:        "Test Pancakes",
					Ingredients:  []string{"flour", "milk", "egg"},
					Instructions: []string{"mix ingredients", "cook on pan"},
					OwnerEmail:   "different_user@example.com",
				},
			},
		},
	}

	router := setupRouter(db)

	test_util.DeleteRecipes(db)
	defer test_util.DeleteRecipes(db)
	if router == nil {
		t.Fatal("handler is nil")
	}

	// Run the test function
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			if err := test_util.InsertRecipes(t, db, tt.MockRecipes); err != nil {
				t.Fatalf("failed to insert mock recipes: %v", err)
			}
			req := httptest.NewRequest(http.MethodDelete, tt.Target, nil)
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != tt.ExpectedStatus {
				t.Errorf("%s: expected %d, got %d", tt.Name, tt.ExpectedStatus, w.Code)
			}
			validateResponse(&tt, w.Result(), db, t)
		})
	}
}
