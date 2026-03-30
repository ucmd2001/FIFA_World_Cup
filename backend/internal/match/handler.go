package match

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang_world_cup/internal/common/response"
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
		response.Error(c, http.StatusInternalServerError, "Failed to fetch matches")
		return
	}
	if matches == nil {
		matches = []MatchResponse{}
	}
	response.JSON(c, http.StatusOK, "Success", matches)
}

// CreateMatch 管理員新增賽事
func (h *Handler) CreateMatch(c *gin.Context) {
	var input CreateMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	m, err := h.svc.CreateMatch(input.TeamA, input.TeamB, input.StartTime)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusCreated, "Match created successfully", m)
}

// ResolveMatch 管理員結算賽果
func (h *Handler) ResolveMatch(c *gin.Context) {
	matchID := c.Param("id")
	var input ResolveMatchInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}
	if err := h.svc.ResolveMatch(matchID, input.Result); err != nil {
		if err.Error() == "match not found" {
			response.Error(c, http.StatusNotFound, err.Error())
		} else if err.Error() == "match already ended" {
			response.Error(c, http.StatusBadRequest, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.JSON(c, http.StatusOK, "Match resolved successfully", nil)
}

// SyncMatches 從 football-data.org 同步賽程
func (h *Handler) SyncMatches(c *gin.Context) {
	synced, err := h.svc.SyncMatches()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "Matches synced successfully", gin.H{
		"count": len(synced),
	})
}
