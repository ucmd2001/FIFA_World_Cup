package match

import (
	"gorm.io/gorm"
)

// Repository 定義 Match domain 的資料存取介面
type Repository interface {
	FindAll() ([]Match, error)
	FindByID(id string) (*Match, error)
	FindByIDWithBets(id string) (*Match, error)
	Create(m *Match) error
	UpdateStatus(m *Match, status, result string) error
	UpdateBetStatus(betID uint, status string, payout int) error
	AddPointsToUser(userID uint, amount int) error
}

type GormRepository struct {
	db *gorm.DB
}

func NewRepository(db *gorm.DB) Repository {
	return &GormRepository{db: db}
}

func (r *GormRepository) FindAll() ([]Match, error) {
	var matches []Match
	if err := r.db.Preload("Bets").Order("start_time asc").Find(&matches).Error; err != nil {
		return nil, err
	}
	return matches, nil
}

func (r *GormRepository) FindByID(id string) (*Match, error) {
	var m Match
	if err := r.db.First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *GormRepository) FindByIDWithBets(id string) (*Match, error) {
	var m Match
	if err := r.db.Preload("Bets").First(&m, id).Error; err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *GormRepository) Create(m *Match) error {
	return r.db.Create(m).Error
}

func (r *GormRepository) UpdateStatus(m *Match, status, result string) error {
	return r.db.Model(m).Updates(map[string]interface{}{
		"status": status,
		"result": result,
	}).Error
}

func (r *GormRepository) UpdateBetStatus(betID uint, status string, payout int) error {
	return r.db.Table("bets").Where("id = ?", betID).Updates(map[string]interface{}{
		"status": status,
		"payout": payout,
	}).Error
}

func (r *GormRepository) AddPointsToUser(userID uint, amount int) error {
	return r.db.Table("users").Where("id = ?", userID).
		UpdateColumn("points", gorm.Expr("points + ?", amount)).Error
}
