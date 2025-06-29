package recipe_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

func TestPostOneRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)

	expected := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshall body with err %s", err.Error())
	}

	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d", w.Code)
	}
	actual := &models.Recipe{}
	err = util.GetBody(w.Result().Body, actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	if !expected.Equals(actual) {
		t.Fatalf("Mismatch: got %+v, want %+v", actual, expected)
	}
}

func TestPostDuplicateRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)

	duplicate := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(duplicate)

	body, err := json.Marshal(duplicate)
	if err != nil {
		t.Fatalf("Failed to marshall body with err %s", err.Error())
	}

	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusConflict {
		t.Fatalf("expected %d, got %d", http.StatusConflict, w.Code)
	}
	var actual map[string]string = make(map[string]string)
	err = util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	expected := "Duplicate recipe"
	if actual["error"] != expected {
		t.Fatalf("Mismatch: got %s, want %s", actual, expected)
	}
}

func TestPostInvalidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)

	body := []byte(`{"url":"https://example.com","titles":"Spaghetti"}`)

	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodPost, target, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var actual map[string]string = make(map[string]string)
	err := util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	expected := "following required keys are empty: title,ingredients,instructions"
	if expected != actual["error"] {
		t.Fatalf("Mismatch: got %s, want %s", actual["error"], expected)
	}
}
