package recipe_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

func TestDeleteValidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	// Seed the database
	toDelete := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(toDelete)

	//Execute request
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", toDelete.Url)
	req := httptest.NewRequest(http.MethodDelete, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusOK {
		t.Fatalf("expected %d, got %d", http.StatusOK, w.Code)
	}
	actual := ""
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	expected := "Deleted"
	if expected != actual {
		t.Fatalf("Mismatch: got %s, want %s", actual, expected)
	}
}

func TestDeleteNotFounddRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", "unknown")
	req := httptest.NewRequest(http.MethodDelete, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected %d, got %d", http.StatusNotFound, w.Code)
	}
}
