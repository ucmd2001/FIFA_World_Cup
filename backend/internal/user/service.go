package user

import (
	"errors"
	"time"
)

// Service 定義 User domain 的業務邏輯介面
type Service interface {
	GetMe(userID uint) (*User, error)
	ClaimDaily(userID uint) (*User, error)
	GetAllUsers() ([]User, error)
	GetLeaderboard(top int) ([]User, error)
	UpdateUser(id string, updates map[string]interface{}) error
	DeleteUser(id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetMe(userID uint) (*User, error) {
	return s.repo.FindByID(userID)
}

func (s *service) ClaimDaily(userID uint) (*User, error) {
	u, err := s.repo.FindByID(userID)
	if err != nil {
		return nil, errors.New("user not found")
	}

	now := time.Now()
	if u.DailyClaimedAt != nil {
		y1, m1, d1 := u.DailyClaimedAt.Date()
		y2, m2, d2 := now.Date()
		if y1 == y2 && m1 == m2 && d1 == d2 {
			return nil, errors.New("daily bonus already claimed today")
		}
	}

	u.Points += 1000
	u.DailyClaimedAt = &now
	if err := s.repo.Save(u); err != nil {
		return nil, errors.New("failed to claim daily bonus")
	}
	return u, nil
}

func (s *service) GetAllUsers() ([]User, error) {
	return s.repo.FindAll()
}

func (s *service) GetLeaderboard(top int) ([]User, error) {
	return s.repo.FindTopN(top)
}

func (s *service) UpdateUser(id string, updates map[string]interface{}) error {
	u := &User{}
	if err := s.repo.Update(u, updates); err != nil {
		return errors.New("failed to update user")
	}
	return nil
}

func (s *service) DeleteUser(id string) error {
	return s.repo.Delete(id)
}
