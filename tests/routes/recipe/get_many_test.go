package recipe_test

import (
	"net/http"
	"net/http/httptest"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

func TestGetNoRecipes(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	actual := &[]models.Recipe{}
	err := util.GetBody(w.Result().Body, actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	if len(*actual) != 0 {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
}

func validateSameRecipes(expected []models.Recipe, actual []models.Recipe, t *testing.T) {
	for _, actual_recipe := range actual {
		found := false
		for _, expected_recipe := range expected {
			if actual_recipe.Equals(&expected_recipe) {
				found = true
				break
			}
		}
		if !found {
			t.Fatalf("Mismatch: got %+v, want %+v", actual, expected)
		}
	}
}

func TestGetOneRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	expected := &[]models.Recipe{
		{
			Url:          "https://example.com",
			Title:        "Test Pancakes",
			Ingredients:  []string{"flour", "milk", "egg"},
			Instructions: []string{"mix ingredients", "cook on pan"},
		},
	}
	db.Create(expected)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	actual := []models.Recipe{}
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	validateSameRecipes(*expected, actual, t)
}

func TestGetManyRecipes(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodGet, target, nil)
	w := httptest.NewRecorder()

	expected := &[]models.Recipe{
		{
			Url:          "https://example.com",
			Title:        "Test Pancake",
			Ingredients:  []string{"flour", "milk", "egg"},
			Instructions: []string{"mix ingredients", "cook on pan"},
		},
		{
			Url:          "https://many_burgers.com",
			Title:        "Test Burgers",
			Ingredients:  []string{"patty", "bun", "sauce"},
			Instructions: []string{"cook patty", "cook bun", "apply sauce"},
		},
		{
			Url:          "https://ghana.com",
			Title:        "Test boiled eggs",
			Ingredients:  []string{"egg", "water"},
			Instructions: []string{"heat up water", "put egg in hot water and let cook"},
		},
	}
	db.Create(expected)
	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected 202, got %d", w.Code)
	}
	actual := []models.Recipe{}
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	validateSameRecipes(*expected, actual, t)
}
