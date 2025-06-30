package recipe_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	util "github.com/tade3910/recipe_server/pkg"
	"github.com/tade3910/recipe_server/pkg/models"
	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

func TestUpdateValidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	// Seed the database
	toUpdate := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(toUpdate)

	expected := &models.Recipe{
		Url:          toUpdate.Url,
		Title:        "Updated Pancakes",
		Ingredients:  toUpdate.Ingredients,
		Instructions: toUpdate.Instructions,
	}
	body, err := json.Marshal(expected)
	if err != nil {
		t.Fatalf("Failed to marshall body with err %s", err.Error())
	}
	//Execute request
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", toUpdate.Url)
	req := httptest.NewRequest(http.MethodPut, target, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusAccepted {
		t.Fatalf("expected %d, got %d", http.StatusAccepted, w.Code)
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

func TestUpdateValidRecipeWithInvalidRecipe(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	// Seed the database
	toUpdate := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(toUpdate)

	invalid := &models.Recipe{
		Url:          toUpdate.Url,
		Ingredients:  toUpdate.Ingredients,
		Instructions: toUpdate.Instructions,
	}
	body, err := json.Marshal(invalid)
	if err != nil {
		t.Fatalf("Failed to marshall body with err %s", err.Error())
	}
	//Execute request
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", toUpdate.Url)
	req := httptest.NewRequest(http.MethodPut, target, bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected %d, got %d", http.StatusBadRequest, w.Code)
	}
	var actual map[string]string = make(map[string]string)
	err = util.GetBody(w.Result().Body, &actual)
	if err != nil {
		t.Fatalf("Unexpected error thrown when getting result %s", err.Error())
	}
	expected := "following required keys are empty: title"
	if expected != actual["error"] {
		t.Fatalf("Mismatch: got %s, want %s", actual["error"], expected)
	}
}

func TestUpdateNotFoundRecipe(t *testing.T) {
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

func TestUpdateEmptyBody(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	// Seed the database
	toUpdate := &models.Recipe{
		Url:          "https://example.com",
		Title:        "Test Pancakes",
		Ingredients:  []string{"flour", "milk", "egg"},
		Instructions: []string{"mix ingredients", "cook on pan"},
	}
	db.Create(toUpdate)

	//Execute request
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", toUpdate.Url)
	req := httptest.NewRequest(http.MethodPut, target, nil)
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
	expected := "invalid body"
	if expected != actual["error"] {
		t.Fatalf("Mismatch: got %s, want %s", actual["error"], expected)
	}
}
