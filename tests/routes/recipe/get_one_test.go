package recipe_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

func validateSameRecipe(expected *models.Recipe, actual *models.Recipe, t *testing.T) {
	if actual.Url != expected.Url {
		t.Errorf("url mismatch: got %q, want %q", actual.Url, expected.Url)
	}
	if actual.Title != expected.Title {
		t.Errorf("title mismatch: got %q, want %q", actual.Title, expected.Title)
	}
	if !reflect.DeepEqual(actual.Ingredients, expected.Ingredients) {
		t.Errorf("ingredients mismatch: got %+v, want %+v", actual.Ingredients, expected.Ingredients)
	}
	if !reflect.DeepEqual(actual.Instructions, expected.Instructions) {
		t.Errorf("Instructions mismatch: got %+v, want %+v", actual.Instructions, expected.Instructions)
	}
}

func TestGetValidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	// Seed the database
	expected := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(expected)

	//Execute request
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", expected.Url)
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	actual := &models.Recipe{}
	err := util.GetBody(w.Result().Body, actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	validateSameRecipe(expected, actual, t)
}

func TestGetInValidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe?url=https://example.com"
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	actual := &models.Recipe{}
	err := util.GetBody(w.Result().Body, actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	validateSameRecipe(expected, actual, t)
}
