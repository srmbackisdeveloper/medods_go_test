package models

import (
	"github.com/google/uuid"
	"gorm.io/gorm"
	"time"
)

type Account struct {
	ID        string      `gorm:"type:uuid;primaryKey" json:"id"`
	Email     string      `gorm:"type:varchar(100);uniqueIndex" json:"email"`
	Number    string      `gorm:"type:varchar(20)" json:"number"`
	CreatedAt time.Time   `gorm:"autoCreateTime" json:"created_at"`
	Tokens    []UserToken `gorm:"foreignKey:UserID" json:"tokens"`
}

func (account *Account) BeforeCreate(tx *gorm.DB) (err error) {
	account.ID = uuid.NewString()
	return
}

type UserToken struct {
	UserID             string    `gorm:"primaryKey;type:uuid" json:"user_id"`
	TokenID            string    `gorm:"type:text;not null" json:"token_id"`
	HashedRefreshToken string    `gorm:"type:text;not null" json:"hashed_refresh_token"`
	CreatedAt          time.Time `gorm:"autoCreateTime" json:"created_at"`
}
