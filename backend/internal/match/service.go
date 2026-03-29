package match

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"time"
)

// Service 定義 Match domain 的業務邏輯介面
type Service interface {
	GetAllMatches() ([]MatchResponse, error)
	CreateMatch(teamA, teamB string, startTime time.Time) (*Match, error)
	ResolveMatch(matchID string, result string) error
	SyncMatches() ([]Match, error)
}

// MatchResponse 是送給前端的搜尋結果（含 pool 統計）
type MatchResponse struct {
	ID        uint      `json:"id"`
	TeamA     string    `json:"teamA"`
	TeamB     string    `json:"teamB"`
	StartTime time.Time `json:"startTime"`
	Status    string    `json:"status"`
	Result    string    `json:"result"`
	Pool      PoolStats `json:"pool"`
}

type PoolStats struct {
	A int `json:"A"`
	B int `json:"B"`
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) GetAllMatches() ([]MatchResponse, error) {
	matches, err := s.repo.FindAll()
	if err != nil {
		return nil, err
	}

	var result []MatchResponse
	for _, m := range matches {
		totalA, totalB := 0, 0
		for _, b := range m.Bets {
			switch b.Choice {
			case "A":
				totalA += b.Amount
			case "B":
				totalB += b.Amount
			}
		}
		result = append(result, MatchResponse{
			ID:        m.ID,
			TeamA:     m.TeamA,
			TeamB:     m.TeamB,
			StartTime: m.StartTime,
			Status:    m.Status,
			Result:    m.Result,
			Pool:      PoolStats{A: totalA, B: totalB},
		})
	}
	return result, nil
}

func (s *service) CreateMatch(teamA, teamB string, startTime time.Time) (*Match, error) {
	m := &Match{
		TeamA:     teamA,
		TeamB:     teamB,
		StartTime: startTime,
		Status:    "active",
	}
	if err := s.repo.Create(m); err != nil {
		return nil, errors.New("failed to create match")
	}
	return m, nil
}

func (s *service) ResolveMatch(matchID string, result string) error {
	m, err := s.repo.FindByIDWithBets(matchID)
	if err != nil {
		return errors.New("match not found")
	}
	if m.Status == "ended" {
		return errors.New("match already ended")
	}

	// 計算各選項總注金
	totalWinningPool := 0
	totalLosingPool := 0
	for _, b := range m.Bets {
		if b.Choice == result {
			totalWinningPool += b.Amount
		} else {
			totalLosingPool += b.Amount
		}
	}

	// 更新每張 Bet 的狀態與 Payout，並派發積分給贏家
	for _, b := range m.Bets {
		var status string
		var payout int

		if b.Choice == result {
			status = "won"
			if totalWinningPool > 0 {
				profitFloat := (float64(b.Amount) * float64(totalLosingPool)) / float64(totalWinningPool)
				payout = b.Amount + int(profitFloat+0.5)
			} else {
				payout = b.Amount
			}
			if err := s.repo.AddPointsToUser(b.UserID, payout); err != nil {
				return errors.New("failed to distribute points")
			}
		} else {
			status = "lost"
			payout = 0
		}

		if err := s.repo.UpdateBetStatus(b.ID, status, payout); err != nil {
			return errors.New("failed to update bet status")
		}
	}

	return s.repo.UpdateStatus(m, "ended", result)
}

func (s *service) SyncMatches() ([]Match, error) {
	apiKey := os.Getenv("FOOTBALL_DATA_API_KEY")
	var matchesToAdd []Match

	if apiKey == "" {
		// Mock data when no API key
		now := time.Now()
		matchesToAdd = []Match{
			{TeamA: "阿根廷", TeamB: "法國", StartTime: now.Add(24 * time.Hour), Status: "active"},
			{TeamA: "巴西", TeamB: "克羅埃西亞", StartTime: now.Add(48 * time.Hour), Status: "active"},
			{TeamA: "英格蘭", TeamB: "葡萄牙", StartTime: now.Add(72 * time.Hour), Status: "pending"},
			{TeamA: "摩洛哥", TeamB: "西班牙", StartTime: now.Add(96 * time.Hour), Status: "pending"},
		}
	} else {
		req, err := http.NewRequest("GET", "https://api.football-data.org/v4/competitions/2000/matches", nil)
		if err != nil {
			return nil, err
		}
		req.Header.Set("X-Auth-Token", apiKey)

		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			return nil, err
		}
		defer resp.Body.Close()

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("football-data API response error")
		}

		body, _ := io.ReadAll(resp.Body)
		var apiResponse struct {
			Matches []struct {
				HomeTeam struct {
					Name string `json:"name"`
				} `json:"homeTeam"`
				AwayTeam struct {
					Name string `json:"name"`
				} `json:"awayTeam"`
				UtcDate string `json:"utcDate"`
				Status  string `json:"status"`
			} `json:"matches"`
		}
		if err := json.Unmarshal(body, &apiResponse); err != nil {
			return nil, err
		}

		for _, m := range apiResponse.Matches {
			if m.HomeTeam.Name == "" || m.AwayTeam.Name == "" {
				continue
			}
			t, _ := time.Parse(time.RFC3339, m.UtcDate)
			status := "pending"
			if m.Status == "IN_PLAY" || m.Status == "TIMED" {
				status = "active"
			} else if m.Status == "FINISHED" {
				continue
			}
			matchesToAdd = append(matchesToAdd, Match{TeamA: m.HomeTeam.Name, TeamB: m.AwayTeam.Name, StartTime: t, Status: status})
		}
	}

	for i := range matchesToAdd {
		if err := s.repo.Create(&matchesToAdd[i]); err != nil {
			return nil, errors.New("failed to save synced matches")
		}
	}
	return matchesToAdd, nil
}
