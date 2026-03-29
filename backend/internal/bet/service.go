package bet

import (
	"errors"
	"time"
)

// PlaceBetInput 下注所需的輸入
type PlaceBetInput struct {
	MatchID uint   `json:"matchId" binding:"required"`
	Choice  string `json:"choice" binding:"required"`
	Amount  int    `json:"amount" binding:"required,gt=0"`
}

// PlaceBetResult 下注成功的回傳資料
type PlaceBetResult struct {
	Bet       *Bet `json:"bet"`
	NewPoints int  `json:"newPoints"`
}

// Service 定義 Bet domain 的業務邏輯介面
type Service interface {
	PlaceBet(userID uint, input PlaceBetInput) (*PlaceBetResult, error)
	GetUserBets(userID uint) ([]Bet, error)
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) PlaceBet(userID uint, input PlaceBetInput) (*PlaceBetResult, error) {
	// 1. 確認賽事 & 狀態
	m, err := s.repo.FindMatchByID(input.MatchID)
	if err != nil {
		return nil, errors.New("match not found")
	}
	if m.Status == "ended" {
		return nil, errors.New("match has already ended")
	}

	// 2. 確認賽事未開始（時間校驗）
	if startTime, ok := m.StartTime.(string); ok {
		t, _ := time.Parse("2006-01-02T15:04:05Z07:00", startTime)
		if t.Before(time.Now()) {
			return nil, errors.New("match has already started, betting is closed")
		}
	}

	// 3. 原子扣點
	rowsAffected, err := s.repo.DeductUserPoints(userID, input.Amount)
	if err != nil {
		return nil, errors.New("failed to update points")
	}
	if rowsAffected == 0 {
		return nil, errors.New("insufficient points or user not found")
	}

	// 4. 取得最新餘額
	u, err := s.repo.GetUserByID(userID)
	if err != nil {
		return nil, errors.New("failed to get user")
	}

	// 5. 建立下注紀錄
	b := &Bet{
		UserID:  userID,
		MatchID: input.MatchID,
		Choice:  input.Choice,
		Amount:  input.Amount,
		Status:  "pending",
	}
	if err := s.repo.Create(b); err != nil {
		return nil, errors.New("failed to create bet")
	}

	return &PlaceBetResult{Bet: b, NewPoints: u.Points}, nil
}

func (s *service) GetUserBets(userID uint) ([]Bet, error) {
	return s.repo.FindByUserID(userID)
}
