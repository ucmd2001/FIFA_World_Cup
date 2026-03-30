package bet

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
