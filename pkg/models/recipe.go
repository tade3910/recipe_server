package models

import (
	"fmt"
	"reflect"
	"strings"
)

type Recipe struct {
	Url          string     `json:"url" gorm:"primaryKey"`
	Title        string     `json:"title"`
	Ingredients  StringList `gorm:"type:jsonb" json:"ingredients"`
	Instructions StringList `gorm:"type:jsonb" json:"instructions"`
}

func (actual *Recipe) Equals(other *Recipe) bool {
	if actual.Url != other.Url {
		return false
	}
	if actual.Title != other.Title {
		return false
	}
	if !reflect.DeepEqual(actual.Ingredients, other.Ingredients) {
		return false
	}
	if !reflect.DeepEqual(actual.Instructions, other.Instructions) {
		return false
	}
	return true
}

func (actual *Recipe) HasRecipeError() error {
	errors := []string{}
	if actual.Url == "" {
		errors = append(errors, "url")
	}
	if actual.Title == "" {
		errors = append(errors, "title")
	}
	if len(actual.Ingredients) == 0 {
		errors = append(errors, "ingredients")
	}
	if len(actual.Instructions) == 0 {
		errors = append(errors, "instructions")
	}
	if len(errors) != 0 {
		return fmt.Errorf("following required keys are empty: %s", strings.Join(errors, ","))

	}
	return nil
}
