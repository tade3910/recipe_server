package models

type Recipe struct {
	Title        string   `json:"title"`
	Ingredients  []string `gorm:"type:jsonb" json:"ingredients"` // Stores as JSON array
	Instructions []string `gorm:"type:jsonb" json:"instructions"`
	Url          string   `json:"url" gorm:"primaryKey"`
}
