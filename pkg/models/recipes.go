package models

type User struct {
	ID      string   `json:"id" gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	Name    string   `json:"name"`
	Email   string   `json:"email"`
	Recipes []Recipe `gorm:"foreignKey:Owner;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}
