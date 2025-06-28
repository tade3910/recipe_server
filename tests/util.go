package tests

import (
	"testing"

	util "github.com/tade3910/recipe_server/pkg"

	"github.com/tade3910/recipe_server/pkg/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var test_db *gorm.DB

func TestInit(t *testing.T) *gorm.DB {
	if test_db != nil {
		return test_db
	}

	dsn := util.LoadEnvs().TestUrl
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to connect to test DB: %v", err)
	}

	if err := db.AutoMigrate(&models.Recipe{}); err != nil {
		t.Fatalf("failed to migrate schema: %v", err)
	}

	return db
}

func DeleteRecipes(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE recipes RESTART IDENTITY CASCADE")
}
