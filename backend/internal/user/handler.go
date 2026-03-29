package user

import (
	"net/http"

	"github.com/gin-gonic/gin"
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
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	u, err := h.svc.GetMe(userID)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "User not found"})
		return
	}
	c.JSON(http.StatusOK, u)
}

// ClaimDaily 領取每日 1000 點
func (h *Handler) ClaimDaily(c *gin.Context) {
	userIDFloat, _ := c.Get("userID")
	userID := uint(userIDFloat.(float64))

	u, err := h.svc.ClaimDaily(userID)
	if err != nil {
		if err.Error() == "daily bonus already claimed today" {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "Daily bonus claimed successfully",
		"points":  u.Points,
	})
}

// GetAllUsers 取得所有使用者列表 (Admin Only)
func (h *Handler) GetAllUsers(c *gin.Context) {
	users, err := h.svc.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}

	var response []gin.H
	for _, u := range users {
		response = append(response, gin.H{
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
	c.JSON(http.StatusOK, response)
}

// GetLeaderboard 取得排行榜前 10 名
func (h *Handler) GetLeaderboard(c *gin.Context) {
	users, err := h.svc.GetLeaderboard(10)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch leaderboard"})
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
	c.JSON(http.StatusOK, leaderboard)
}

// UpdateUserInput 定義可更新的欄位
type UpdateUserInput struct {
	Name   string `json:"name"`
	Email  string `json:"email"`
	Points *int   `json:"points"`
	Role   string `json:"role"`
}

// UpdateUser 管理員更新使用者資訊
func (h *Handler) UpdateUser(c *gin.Context) {
	userID := c.Param("id")
	var input UpdateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

	// 用 id 當作條件更新
	updates["id"] = userID
	if err := h.svc.UpdateUser(userID, updates); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User updated successfully"})
}

// DeleteUser 管理員刪除使用者
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.svc.DeleteUser(userID); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User deleted successfully"})
}
