package bet

import (
	"time"

	"gorm.io/gorm"
)

// Bet 是下注的 Domain Entity
type Bet struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	UserID    uint           `gorm:"index" json:"userId"`
	MatchID   uint           `gorm:"index" json:"matchId"`
	Choice    string         `gorm:"not null" json:"choice"` // 'A', 'B'
	Amount    int            `gorm:"not null" json:"amount"`
	Status    string         `gorm:"default:'pending'" json:"status"` // 'pending', 'won', 'lost'
	Payout    int            `gorm:"default:0" json:"payout"`
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`

	// Preloaded
	Match BetMatch `gorm:"foreignKey:MatchID" json:"match,omitempty"`
}

// BetMatch 是 Bet 聚合根用於預載的輕量 Match 結構，避免 import cycle
type BetMatch struct {
	ID    uint   `gorm:"primarykey" json:"id"`
	TeamA string `json:"teamA"`
	TeamB string `json:"teamB"`
}

func (BetMatch) TableName() string {
	return "matches"
}
