package user

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"golang_world_cup/internal/common/response"
	"golang_world_cup/internal/common/utils"
)

// Handler 負責處理 User 相關的 HTTP 請求
type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

// GetMe 取得當前使用者資訊
func (h *Handler) GetMe(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	u, err := h.svc.GetMe(userID)
	if err != nil {
		response.Error(c, http.StatusNotFound, "User not found")
		return
	}
	response.JSON(c, http.StatusOK, "Success", u)
}

// ClaimDaily 領取每日 1000 點
func (h *Handler) ClaimDaily(c *gin.Context) {
	userID, ok := utils.GetUserID(c)
	if !ok {
		response.Error(c, http.StatusUnauthorized, "Unauthorized")
		return
	}

	u, err := h.svc.ClaimDaily(userID)
	if err != nil {
		if err.Error() == "daily bonus already claimed today" {
			response.Error(c, http.StatusBadRequest, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}
	response.JSON(c, http.StatusOK, "Daily bonus claimed successfully", gin.H{
		"points": u.Points,
	})
}

// GetAllUsers 取得所有使用者列表 (Admin Only)
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch users")
		return
	}

	var data []gin.H
	for _, u := range users {
		data = append(data, gin.H{
			"id":             u.ID,
			"username":       u.Username,
			"name":           u.Name,
			"email":          u.Email,
			"role":           u.Role,
			"points":         u.Points,
			"lastLoginAt":    u.LastLoginAt,
			"dailyClaimedAt": u.DailyClaimedAt,
			"createdAt":      u.CreatedAt,
		})
	}
	response.JSON(c, http.StatusOK, "Success", data)
}

// GetLeaderboard 取得排行榜前 10 名
func (h *Handler) GetLeaderboard(c *gin.Context) {
	users, err := h.svc.GetLeaderboard(10)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to fetch leaderboard")
		return
	}

	var leaderboard []gin.H
	for _, u := range users {
		leaderboard = append(leaderboard, gin.H{
			"id":       u.ID,
			"username": u.Username,
			"points":   u.Points,
		})
	}
	response.JSON(c, http.StatusOK, "Success", leaderboard)
}

// UpdateUser 管理員更新使用者資訊
func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updates := map[string]interface{}{}
	if input.Name != "" {
		updates["name"] = input.Name
	}
	if input.Email != "" {
		updates["email"] = input.Email
	}
	if input.Points != nil {
		updates["points"] = *input.Points
	}
	if input.Role != "" {
		updates["role"] = input.Role
	}

	updates["id"] = userID
	if err := h.svc.UpdateUser(userID, updates); err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}
	response.JSON(c, http.StatusOK, "User updated successfully", nil)
}

// DeleteUser 管理員刪除使用者
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.svc.DeleteUser(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	response.JSON(c, http.StatusOK, "User deleted successfully", nil)
}
