package recipe_test

import (
	"encoding/json"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/tade3910/recipe_server/pkg/models"
	"gorm.io/gorm"
)

// Validate response status and body against expected JSON structure
func AssertResponse(t *testing.T, w *httptest.ResponseRecorder, testName string, expectedStatus int, expectedBody any) *models.Recipe {
	t.Helper()
	if w.Code != expectedStatus {
		t.Fatalf("%s: expected status %d, got %d. Body: %s", testName, expectedStatus, w.Code, w.Body.String())
	}

	switch expected := expectedBody.(type) {
	case string:
		if !strings.Contains(w.Body.String(), expected) {
			t.Fatalf("%s: expected body to contain %q, got %q", testName, expected, w.Body.String())
		}
	case *models.Recipe:
		var actual models.Recipe
		err := json.Unmarshal(w.Body.Bytes(), &actual)
		if err != nil {
			t.Fatalf("%s: failed to unmarshal response body: %v", testName, err)
		}
		if !actual.Equals(expected) {
			t.Fatalf("%s: expected recipe %+v, got %+v", testName, expected, actual)
		}
		return &actual
	default:
		t.Fatalf("%s: unexpected type for ExpectedBody: %T", testName, expectedBody)
	}
	return nil
}

func AssertRecipeInDb(
	t *testing.T,
	db *gorm.DB,
	resp *models.Recipe,
	testName string,
) {
	t.Helper()
	var recipe models.Recipe
	err := db.Where("id = ?", resp.ID).First(&recipe).Error
	if err != nil {
		t.Fatalf("Failed to get recipe from db with :%s", err.Error())
	}
	if !recipe.Equals(resp) {
		t.Fatalf("%s: expected recipe %+v, got %+v", testName, recipe, resp)
	}
}
