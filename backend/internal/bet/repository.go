package bet

import (
	"gorm.io/gorm"
)

// Repository 定義 Bet domain 的資料存取介面
type Repository interface {
	Create(b *Bet) error
	FindByUserID(userID uint) ([]Bet, error)
	FindMatchByID(matchID uint) (*MatchForBet, error)
	DeductUserPoints(userID uint, amount int) (rowsAffected int64, err error)
	GetUserByID(userID uint) (*UserForBet, error)
}

// MatchForBet 用於 Bet 業務邏輯查詢的賽事資訊（避免 import cycle）
type MatchForBet struct {
	ID        uint
	Status    string
	StartTime interface{}
}

// UserForBet 用於 Bet 業務邏輯查詢的使用者資訊（避免 import cycle）
type UserForBet struct {
	ID     uint
	Points int
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) Create(b *Bet) error {
	return r.db.Create(b).Error
}

func (r *GormRepository) FindByUserID(userID uint) ([]Bet, error) {
	var bets []Bet
	if err := r.db.Preload("Match").Where("user_id = ?", userID).Order("created_at desc").Find(&bets).Error; err != nil {
		return nil, err
	}
	return bets, nil
}

func (r *GormRepository) FindMatchByID(matchID uint) (*MatchForBet, error) {
	var m struct {
		ID        uint
		Status    string
		StartTime interface{}
	}
	if err := r.db.Table("matches").Select("id, status, start_time").Where("id = ? AND deleted_at IS NULL", matchID).First(&m).Error; err != nil {
		return nil, err
	}
	return &MatchForBet{ID: m.ID, Status: m.Status, StartTime: m.StartTime}, nil
}

func (r *GormRepository) DeductUserPoints(userID uint, amount int) (int64, error) {
	result := r.db.Table("users").
		Where("id = ? AND points >= ?", userID, amount).
		UpdateColumn("points", gorm.Expr("points - ?", amount))
	return result.RowsAffected, result.Error
}

func (r *GormRepository) GetUserByID(userID uint) (*UserForBet, error) {
	var u struct {
		ID     uint
		Points int
	}
	if err := r.db.Table("users").Select("id, points").Where("id = ?", userID).First(&u).Error; err != nil {
		return nil, err
	}
	return &UserForBet{ID: u.ID, Points: u.Points}, nil
}
