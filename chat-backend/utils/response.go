package utils

import (
	"net/http"
	"os"
	"time"

	"chat-backend/models"

	"github.com/gin-gonic/gin"
)

// RespondWithSuccess 返回成功响应
func RespondWithSuccess(c *gin.Context, data interface{}, message ...string) {
	msg := "操作成功"
	if len(message) > 0 {
		msg = message[0]
	}

	response := models.APIResponse{
		Success: true,
		Message: msg,
		Data:    data,
	}

	c.JSON(http.StatusOK, response)
}

// GetEnvOrDefault 获取环境变量或默认值
func GetEnvOrDefault(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}
