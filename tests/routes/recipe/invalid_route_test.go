package recipe_test

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	recipe "github.com/tade3910/recipe_server/pkg/routes"
	test_util "github.com/tade3910/recipe_server/tests"
)

const EXPECTED_STATUS = http.StatusMethodNotAllowed

func TestUpdateNoUrl(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodPut, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != EXPECTED_STATUS {
		t.Fatalf("expected %d, got %d", EXPECTED_STATUS, w.Code)
	}
}

func TestDeleteNoUrl(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := "/recipe"
	req := httptest.NewRequest(http.MethodDelete, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != EXPECTED_STATUS {
		t.Fatalf("expected %d, got %d", EXPECTED_STATUS, w.Code)
	}
}

func TestPostUrl(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?url=%s", "invalid")
	req := httptest.NewRequest(http.MethodPost, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != EXPECTED_STATUS {
		t.Fatalf("expected %d, got %d", EXPECTED_STATUS, w.Code)
	}
}

func TestInvalidQueryUpdate(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?query=%s", "invalid")
	req := httptest.NewRequest(http.MethodPut, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != EXPECTED_STATUS {
		t.Fatalf("expected %d, got %d", EXPECTED_STATUS, w.Code)
	}
}

func TestInvalidQueryDelete(t *testing.T) {
	db := test_util.TestInit(t)
	defer test_util.DeleteRecipes(db)
	handler := recipe.NewRecipesHandler(db)
	target := fmt.Sprintf("/recipe?query=%s", "invalid")
	req := httptest.NewRequest(http.MethodDelete, target, nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	if w.Code != EXPECTED_STATUS {
		t.Fatalf("expected %d, got %d", EXPECTED_STATUS, w.Code)
	}
}
