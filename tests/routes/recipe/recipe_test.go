package recipe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/tade3910/recipe_server/pkg/models"
	"github.com/tade3910/recipe_server/pkg/routes/recipe"
	util "github.com/tade3910/recipe_server/tests"
)

func TestGetRecipeWithRealDB(t *testing.T) {
	db := util.TestInit(t)
	defer util.DeleteRecipes(db)
	// Seed the database
	expected := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(expected)

	// Act
	handler := recipe.NewRecipesHandler(db)
	req := httptest.NewRequest(http.MethodGet, "/recipe?url="+expected.Url, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}

}
