package response

import (
	"github.com/gin-gonic/gin"
)

// JSON 統一回傳格式
func JSON(c *gin.Context, code int, message string, data interface{}) {
	c.JSON(code, gin.H{
		"code":    code,
		"message": message,
		"data":    data,
	})
}

// Error 統一錯誤回傳格式
func Error(c *gin.Context, code int, message string) {
	c.JSON(code, gin.H{
		"code":  code,
		"error": message,
	})
}
