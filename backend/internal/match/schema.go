package match

import "time"

// CreateMatchInput 管理員新增賽事輸入
type CreateMatchInput struct {
	TeamA     string    `json:"teamA" binding:"required"`
	TeamB     string    `json:"teamB" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
}

// ResolveMatchInput 管理員結算賽果輸入
type ResolveMatchInput struct {
	Result string `json:"result" binding:"required"` // 'A', 'B'
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

// PoolStats 賽事金額池統計
type PoolStats struct {
	A int `json:"A"`
	B int `json:"B"`
}
