package match

import (
	"time"

	"gorm.io/gorm"
)

// Match 是賽事的 Domain Entity
type Match struct {
	ID        uint           `gorm:"primarykey" json:"id"`
	TeamA     string         `gorm:"not null" json:"teamA"`
	TeamB     string         `gorm:"not null" json:"teamB"`
	StartTime time.Time      `gorm:"not null" json:"startTime"`
	Status    string         `gorm:"default:'pending'" json:"status"` // 'pending', 'active', 'ended'
	Result    string         `json:"result"`                           // 'A', 'B'
	Options   string         `json:"options"`                          // Reserved for future
	CreatedAt time.Time      `json:"createdAt"`
	UpdatedAt time.Time      `json:"updatedAt"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
	Bets      []MatchBet     `gorm:"foreignKey:MatchID" json:"bets,omitempty"`
}

// MatchBet 是 Match 聚合根用於計算下注的輕量結構
// 避免 import cycle（不引入 bet package）
type MatchBet struct {
	ID      uint   `gorm:"primarykey"`
	MatchID uint   `gorm:"index"`
	UserID  uint   `gorm:"index"`
	Choice  string `gorm:"not null"`
	Amount  int    `gorm:"not null"`
	Status  string `gorm:"default:'pending'"`
	Payout  int    `gorm:"default:0"`
}

// TableName 指向 bets 資料表，避免 GORM 依據結構體名稱產生錯誤的表名
func (MatchBet) TableName() string {
	return "bets"
}
