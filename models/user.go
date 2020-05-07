package models

import "time"

// User represents each Telegram user
type User struct {
	ID              string `gorm:"prmary_key" json:"id"`
	FirstName       string `gorm:"type:varchar(255)" json:"first_name"`
	LastName        string `gorm:"type:varchar(255)" json:"last_name"`
	Username        string `gorm:"type:varchar(255);NOT NULL" json:"username" binding:"required"`
	TelegramID      string `gorm:"type:varchar(255);NOT NULL" json:"telegram_id" binding:"required"`
	EtherminerToken string `gorm:"type:varchar(255)" json:"etherminer_token"`
	CreatedAt       time.Time
	UpdatedAt       time.Time
}

// Users array of User
type Users []User
