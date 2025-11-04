package utils

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// SuccessResponse 统一成功响应结构
// swagger:model
type SuccessResponse struct {
	// 请求是否成功
	// required: true
	Success bool `json:"success"`
	// 响应消息
	// required: true
	Message string `json:"message"`
	// 响应数据
	// required: true
	Data interface{} `json:"data"`
	// 时间戳
	// required: true
	Timestamp time.Time `json:"timestamp"`
}

// RespondWithSuccess 返回成功响应
func RespondWithSuccess(c *gin.Context, data interface{}, message ...string) {
	msg := "操作成功"
	if len(message) > 0 {
		msg = message[0]
	}

	response := SuccessResponse{
		Success:   true,
		Message:   msg,
		Data:      data,
		Timestamp: time.Now(),
	}

	c.JSON(http.StatusOK, response)
}

// RespondWithBadRequest 返回400错误响应
func RespondWithBadRequest(c *gin.Context, message string) {
	apiErr := NewAPIError(ErrInvalidRequest, message, http.StatusBadRequest)
	RespondWithError(c, apiErr)
}

// RespondWithInternalError 返回500错误响应
func RespondWithInternalError(c *gin.Context, message string) {
	apiErr := NewAPIError(ErrInternalServer, message, http.StatusInternalServerError)
	RespondWithError(c, apiErr)
}

// GetCurrentTime 获取当前时间
func GetCurrentTime() time.Time {
	return time.Now()
}
