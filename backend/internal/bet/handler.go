package bet

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang_world_cup/internal/common/response"
	"golang_world_cup/internal/common/utils"
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
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	userID, ok := utils.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	result, err := h.svc.PlaceBet(userID, input)
	if err != nil {
		switch err.Error() {
		case "match not found":
			response.Error(c, http.StatusNotFound, err.Error())
		case "match has already ended", "match has already started, betting is closed", "insufficient points or user not found":
			response.Error(c, http.StatusBadRequest, err.Error())
		default:
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(c, http.StatusCreated, "Bet placed successfully", result)
}

// GetUserBets 取得使用者所有下注紀錄
func (h *Handler) GetUserBets(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	bets, err := h.svc.GetUserBets(userID)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch bets")
		return
	}
	if bets == nil {
		bets = []Bet{}
	}
	response.JSON(c, http.StatusOK, "Success", bets)
}
