package models

type Recipe struct {
	Url          string     `json:"url" gorm:"primaryKey"`
	Title        string     `json:"title"`
	Ingredients  StringList `gorm:"type:jsonb" json:"ingredients"`
	Instructions StringList `gorm:"type:jsonb" json:"instructions"`
}
