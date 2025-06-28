package models

import (
	"reflect"
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
