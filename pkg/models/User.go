package models

import "time"

type User struct {
	Email    string `gorm:"primaryKey;type:varchar(255);not null"`
	Name     string `json:"name"`
	Password string `gorm:"not null"` // hashed password

	//foreignKey points to OwnerEmail, OnDelete changed to CASCADE
	Recipes []Recipe    `gorm:"foreignKey:OwnerEmail;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Tokens  []AuthToken `gorm:"foreignKey:UserEmail;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
}

// AuthToken represents access or refresh tokens
type AuthToken struct {
	ID        string `gorm:"primaryKey;type:uuid;default:gen_random_uuid()"`
	UserEmail string `gorm:"not null"`
	User      User   `gorm:"foreignKey:UserEmail;references:Email;constraint:OnUpdate:CASCADE,OnDelete:CASCADE"`
	Token     string `gorm:"unique;not null"` // the token string
	TokenType string `gorm:"not null"`        // "access" or "refresh"
	CreatedAt time.Time
	ExpiresAt time.Time
}
