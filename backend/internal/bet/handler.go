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

// PlaceBet godoc
// @Summary      Place a bet
// @Description  User places a bet on a match (A/B)
// @Tags         Bet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body PlaceBetInput true "Bet details"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /bet [post]
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

// GetUserBets godoc
// @Summary      Get user bets
// @Description  Returns all bets placed by the logged-in user
// @Tags         Bet
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /bet [get]
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
