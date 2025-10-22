package tests

import (
	"gorm.io/gorm"
)

func DeleteRecipes(db *gorm.DB) {
	db.Exec("TRUNCATE TABLE recipes RESTART IDENTITY CASCADE")
}
