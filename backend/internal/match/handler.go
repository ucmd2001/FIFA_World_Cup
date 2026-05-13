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

// GetMatches godoc
// @Summary      Get all matches
// @Description  Returns list of matches including betting pool stats
// @Tags         Match
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /match [get]
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

// CreateMatch godoc
// @Summary      Create match
// @Description  Admin creates a new match
// @Tags         Match
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        input body CreateMatchInput true "Match details"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /match [post]
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

// ResolveMatch godoc
// @Summary      Resolve match
// @Description  Admin resolves a match result (A/B)
// @Tags         Match
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "Match ID"
// @Param        input body ResolveMatchInput true "Result A/B"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /match/{id}/resolve [put]
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

// SyncMatches godoc
// @Summary      Sync matches
// @Description  Admin manually updates matches from external API
// @Tags         Match
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /match/sync [post]
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
