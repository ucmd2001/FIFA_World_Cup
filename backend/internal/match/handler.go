package match

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Handler 負責處理 Match 相關的 HTTP 請求
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GetMatches 取得所有賽事（含注池統計）
func (h *Handler) GetMatches(c *gin.Context) {
	matches, err := h.svc.GetAllMatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch matches"})
		return
	}
	if matches == nil {
		matches = []MatchResponse{}
	}
	c.JSON(http.StatusOK, matches)
}

type CreateMatchInput struct {
	TeamA     string    `json:"teamA" binding:"required"`
	TeamB     string    `json:"teamB" binding:"required"`
	StartTime time.Time `json:"startTime" binding:"required"`
}

// CreateMatch 管理員新增賽事
func (h *Handler) CreateMatch(c *gin.Context) {
	var input CreateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	m, err := h.svc.CreateMatch(input.TeamA, input.TeamB, input.StartTime)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, m)
}

type ResolveMatchInput struct {
	Result string `json:"result" binding:"required"` // 'A', 'B'
}

// ResolveMatch 管理員結算賽果
func (h *Handler) ResolveMatch(c *gin.Context) {
	matchID := c.Param("id")
	var input ResolveMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.svc.ResolveMatch(matchID, input.Result); err != nil {
		if err.Error() == "match not found" {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		} else if err.Error() == "match already ended" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Match resolved successfully"})
}

// SyncMatches 從 football-data.org 同步賽程
func (h *Handler) SyncMatches(c *gin.Context) {
	synced, err := h.svc.SyncMatches()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Matches synced successfully",
		"count":   len(synced),
	})
}
