package models

import (
	"reflect"
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type Recipe struct {
	// 1. UUID Primary Key with GORM Auto-Generation
	ID string `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`

	// 2. Composite Unique Index on (OwnerEmail, Url)
	URL string `json:"url" gorm:"not null;index:idx_owner_url,unique;check:url <> ''"`

	Title string `json:"title" gorm:"not null;check:title <> '';comment:Recipe title must not be empty"`

	// 3. Array Checks on JSONB fields
	Ingredients  StringList `json:"ingredients" gorm:"type:jsonb;not null;check:jsonb_array_length(ingredients) > 0;comment:Must have at least one ingredient"`
	Instructions StringList `json:"instructions" gorm:"type:jsonb;not null;check:jsonb_array_length(instructions) > 0;comment:Must have at least one instruction"`

	// 4. Foreign Key & Composite Index
	OwnerEmail string `json:"owner_email" gorm:"not null;index:idx_owner_url,unique;comment:Owner of the recipe"`
	Owner      User   `json:"owner,omitempty" gorm:"foreignKey:OwnerEmail;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`

	// 5. Standard Timestamps (Best Practice)
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// BeforeCreate hook ensures Go generates a UUID if PostgreSQL's gen_random_uuid() isn't invoked during mock testing
func (r *Recipe) BeforeCreate(tx *gorm.DB) error {
	if r.ID == "" {
		r.ID = uuid.NewString()
	}
	return nil
}

func (actual *Recipe) Equals(other *Recipe) bool {
	if actual.URL != other.URL {
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
	if actual.OwnerEmail != other.OwnerEmail {
		return false
	}
	return true
}
