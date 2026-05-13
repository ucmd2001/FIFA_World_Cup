package auth

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang_world_cup/internal/common/response"
)

// Handler 處理 Auth 相關的 HTTP 請求
type Handler struct {
	svc *Service
}

func NewHandler(svc *Service) *Handler {
	return &Handler{svc: svc}
}

// Register godoc
// @Summary      User Registration
// @Description  Creates a new user account
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body RegisterInput true "Registration details"
// @Success      201  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      409  {object}  response.ErrorResponse
// @Router       /auth/register [post]
func (h *Handler) Register(c *gin.Context) {
	var input RegisterInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Register(input); err != nil {
		if err.Error() == "Username already taken" {
			response.Error(c, http.StatusConflict, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.JSON(c, http.StatusCreated, "Registration successful", nil)
}

// Login godoc
// @Summary      User Login
// @Description  Authenticates user and returns JWT
// @Tags         Auth
// @Accept       json
// @Produce      json
// @Param        input body LoginInput true "Login credentials"
// @Success      200
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Router       /auth/login [post]
func (h *Handler) Login(c *gin.Context) {
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	if err := h.svc.Login(c, input); err != nil {
		// Service 內部目前仍有處理 Response，我們維持原狀或微調
		return
	}
}
