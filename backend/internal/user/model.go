package user

import (
	"time"

	"gorm.io/gorm"
)

// User 是使用者的 Domain Entity
type User struct {
	ID             uint           `gorm:"primarykey" json:"id"`
	Username       string         `gorm:"uniqueIndex;not null" json:"username"`
	Name           string         `json:"name"`
	Email          string         `json:"email"`
	PasswordHash   string         `gorm:"not null" json:"-"`
	Points         int            `gorm:"default:0" json:"points"`
	Role           string         `gorm:"default:'user'" json:"role"` // 'admin' or 'user'
	LastLoginAt    *time.Time     `json:"lastLoginAt"`
	DailyClaimedAt *time.Time     `json:"dailyClaimedAt"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
	DeletedAt      gorm.DeletedAt `gorm:"index" json:"-"`
}
