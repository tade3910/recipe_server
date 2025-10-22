package recipe_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
	"github.com/tade3910/recipe_server/tests/mocks"
)

func TestScrapeValidUrl(t *testing.T) {
	db := mocks.TestDb(t)
	scraper := &mocks.MockRecipeScraper{}
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db, scraper)
	target := "/recipe?url=https://example.com"
	req := httptest.NewRequest(http.MethodPost, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusCreated {
		t.Fatalf("expected %d, got %d", http.StatusCreated, w.Code)
	}
	actual := map[string][][]string{}
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	actual_ingredients, ok := actual["ingredients"]
	if !ok {
		t.Fatalf("Ingredients not in response")
	}
	if !reflect.DeepEqual(actual_ingredients, mocks.MockIngredients) {
		t.Fatalf("actual ingredients: %s are not the same as expected", actual_ingredients)
	}
	actual_instructions, ok := actual["instructions"]
	if !ok {
		t.Fatalf("instructions not in response")
	}
	if !reflect.DeepEqual(actual_instructions, mocks.MockInstructions) {
		t.Fatalf("actual instructions: %s are not the same as expected", actual_instructions)
	}
}

func TestScrapeInvalidUrl(t *testing.T) {
	db := mocks.TestDb(t)
	scraper := &mocks.MockRecipeScraper{}
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db, scraper)
	target := fmt.Sprintf("/recipe?url=%s", mocks.ErrorUrl)
	req := httptest.NewRequest(http.MethodPost, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)
	if w.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", w.Code)
	}
	actual := map[string]string{}
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	actual_error, ok := actual["error"]
	if !ok {
		t.Fatalf("error not in response")
	}
	if actual_error != mocks.ScrapingError {
		t.Fatalf("Unexpected error received: %s", actual_error)
	}
}
