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

// GetMe godoc
// @Summary      Get current user profile
// @Description  Returns profile data of the logged-in user
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      401  {object}  response.ErrorResponse
// @Failure      404  {object}  response.ErrorResponse
// @Router       /user/me [get]
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

// ClaimDaily godoc
// @Summary      Claim daily points
// @Description  Allows user to claim their 1000 daily points
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      401  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /user/daily [post]
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

// GetAllUsers godoc
// @Summary      Get all users
// @Description  Returns list of all users, Admin only
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /user/all [get]
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

// GetLeaderboard godoc
// @Summary      Get Leaderboard
// @Description  Returns top 10 users by points
// @Tags         User
// @Accept       json
// @Produce      json
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /user/leaderboard [get]
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

// UpdateUser godoc
// @Summary      Update user
// @Description  Update user details (Admin only)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Param        input body UpdateUserInput true "Update details"
// @Success      200  {object}  response.Response
// @Failure      400  {object}  response.ErrorResponse
// @Failure      500  {object}  response.ErrorResponse
// @Router       /user/{id} [put]
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

// DeleteUser godoc
// @Summary      Delete user
// @Description  Delete user by ID (Admin only)
// @Tags         User
// @Accept       json
// @Produce      json
// @Security     BearerAuth
// @Param        id   path      int  true  "User ID"
// @Success      200  {object}  response.Response
// @Failure      500  {object}  response.ErrorResponse
// @Router       /user/{id} [delete]
func (h *Handler) DeleteUser(c *gin.Context) {
	userID := c.Param("id")
	if err := h.svc.DeleteUser(userID); err != nil {
		response.Error(c, http.StatusInternalServerError, "Failed to delete user")
		return
	}
	response.JSON(c, http.StatusOK, "User deleted successfully", nil)
}
