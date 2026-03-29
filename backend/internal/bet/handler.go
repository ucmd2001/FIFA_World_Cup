package bet

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

// Handler 負責處理 Bet 相關的 HTTP 請求
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// PlaceBet 使用者進行下注
func (h *Handler) PlaceBet(c *gin.Context) {
	var input PlaceBetInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userIDFloat, exists := c.Get("userID")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}
	userID := uint(userIDFloat.(float64))

	result, err := h.svc.PlaceBet(userID, input)
	if err != nil {
		switch err.Error() {
		case "match not found":
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		case "match has already ended", "match has already started, betting is closed", "insufficient points or user not found":
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		default:
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":   "Bet placed successfully",
		"bet":       result.Bet,
		"newPoints": result.NewPoints,
	})
}

// GetUserBets 取得使用者所有下注紀錄
func (h *Handler) GetUserBets(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	bets, err := h.svc.GetUserBets(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch bets"})
		return
	}
	if bets == nil {
		bets = []Bet{}
	}
	c.JSON(http.StatusOK, bets)
}
