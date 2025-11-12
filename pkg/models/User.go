package models

import "time"

type User struct {
	Email    string   `gorm:"primaryKey;not null"`
	Name     string   `json:"name"`
	Password string   `gorm:"not null"` // hashed
	Recipes  []Recipe `gorm:"foreignKey:Owner;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
}

type AuthToken struct {
	ID        string `gorm:"primaryKey"`
	UserEmail string `gorm:"not null"` // foreign key column
	User      User   `gorm:"foreignKey:UserEmail;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string `gorm:"unique;not null"`
	CreatedAt time.Time
}
