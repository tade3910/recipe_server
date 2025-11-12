package models

import (
	"fmt"
	"reflect"
	"strings"
)

type Recipe struct {
	Id           string     `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Url          string     `json:"url" gorm:"not null;index:idx_owner_url,unique"`
	Title        string     `json:"title" gorm:"not null;check:title <> '';comment:Recipe title must not be empty"`
	Ingredients  StringList `json:"ingredients" gorm:"type:jsonb;not null;check:json_array_length(ingredients) > 0;comment:Must have at least one ingredient"`
	Instructions StringList `json:"instructions" gorm:"type:jsonb;not null;check:json_array_length(instructions) > 0;comment:Must have at least one instruction"`
	Owner        string     `json:"owner" gorm:"not null;index:idx_owner_url,unique;comment:Owner of the recipe (used for uniqueness constraint with URL)"`
}

func (actual *Recipe) Equals(other *Recipe) bool {
	if actual.Id != other.Id {
		return false
	}
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
	if actual.Owner != other.Owner {
		return false
	}
	return true
}

func (actual *Recipe) HasRecipeError() error {
	errors := []string{}
	if strings.TrimSpace(actual.Url) == "" {
		errors = append(errors, "url")
	}
	if strings.TrimSpace(actual.Title) == "" {
		errors = append(errors, "title")
	}
	if len(actual.Ingredients) == 0 {
		errors = append(errors, "ingredients")
	}
	if len(actual.Instructions) == 0 {
		errors = append(errors, "instructions")
	}
	if strings.TrimSpace(actual.Owner) == "" {
		errors = append(errors, "owner")
	}
	if len(errors) != 0 {
		return fmt.Errorf("following required keys are empty: %s", strings.Join(errors, ","))
	}
	return nil
}
