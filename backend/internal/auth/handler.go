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

// Register 處理使用者註冊
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

// Login 處理使用者登入
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
