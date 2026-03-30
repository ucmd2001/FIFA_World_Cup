package utils

import (
	"github.com/gin-gonic/gin"
)

// GetUserID 從 Gin Context 中提取 userID (從 AuthMiddleware 注入)
func GetUserID(c *gin.Context) (uint, bool) {
	val, exists := c.Get("userID")
	if !exists {
		return 0, false
	}

	// JWT 宣告中儲存的是 float64 (因為 JSON 反序列化)
	if id, ok := val.(float64); ok {
		return uint(id), true
	}

	if id, ok := val.(uint); ok {
		return id, true
	}

	return 0, false
}
