package utils

import (
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// ErrorCode 错误代码类型
type ErrorCode string

const (
	// 通用错误
	ErrInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrInternalServer ErrorCode = "INTERNAL_SERVER_ERROR"
	ErrUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrForbidden      ErrorCode = "FORBIDDEN"
	ErrNotFound       ErrorCode = "NOT_FOUND"
	ErrTimeout        ErrorCode = "TIMEOUT"
	ErrRateLimited    ErrorCode = "RATE_LIMITED"

	// 聊天相关错误
	ErrConversationNotFound ErrorCode = "CONVERSATION_NOT_FOUND"
	ErrConversationCreate   ErrorCode = "CONVERSATION_CREATE_FAILED"
	ErrMessageSend          ErrorCode = "MESSAGE_SEND_FAILED"
	ErrMessageInvalid       ErrorCode = "MESSAGE_INVALID"

	// 知识库相关错误
	ErrKnowledgeBaseNotFound ErrorCode = "KNOWLEDGE_BASE_NOT_FOUND"
	ErrKnowledgeBaseCreate   ErrorCode = "KNOWLEDGE_BASE_CREATE_FAILED"
	ErrKnowledgeBaseUpdate   ErrorCode = "KNOWLEDGE_BASE_UPDATE_FAILED"
	ErrKnowledgeBaseDelete   ErrorCode = "KNOWLEDGE_BASE_DELETE_FAILED"
	ErrFileUpload            ErrorCode = "FILE_UPLOAD_FAILED"
	ErrFileNotFound          ErrorCode = "FILE_NOT_FOUND"
	ErrFileSize              ErrorCode = "FILE_SIZE_EXCEEDED"
	ErrFileType              ErrorCode = "FILE_TYPE_NOT_SUPPORTED"

	// 设置相关错误
	ErrModelNotFound   ErrorCode = "MODEL_NOT_FOUND"
	ErrSettingsUpdate  ErrorCode = "SETTINGS_UPDATE_FAILED"
	ErrSettingsInvalid ErrorCode = "SETTINGS_INVALID"

	// Flowy SDK相关错误
	ErrFlowyConnection ErrorCode = "FLOWY_CONNECTION_FAILED"
	ErrFlowyAuth       ErrorCode = "FLOWY_AUTH_FAILED"
	ErrFlowyAPI        ErrorCode = "FLOWY_API_ERROR"
)

// ErrorInfo 错误信息配置
type ErrorInfo struct {
	Code           ErrorCode
	Message        string // 中文提示消息
	HTTPStatusCode int    // HTTP状态码
}

// errorInfoMap 错误信息映射表
var errorInfoMap = map[ErrorCode]ErrorInfo{
	// 通用错误
	ErrInvalidRequest: {
		Code:           ErrInvalidRequest,
		Message:        "请求参数无效",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrInternalServer: {
		Code:           ErrInternalServer,
		Message:        "服务器内部错误",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrUnauthorized: {
		Code:           ErrUnauthorized,
		Message:        "未授权访问",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	ErrForbidden: {
		Code:           ErrForbidden,
		Message:        "禁止访问",
		HTTPStatusCode: http.StatusForbidden,
	},
	ErrNotFound: {
		Code:           ErrNotFound,
		Message:        "资源不存在",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrTimeout: {
		Code:           ErrTimeout,
		Message:        "请求超时",
		HTTPStatusCode: http.StatusRequestTimeout,
	},
	ErrRateLimited: {
		Code:           ErrRateLimited,
		Message:        "请求过于频繁",
		HTTPStatusCode: http.StatusTooManyRequests,
	},

	// 聊天相关错误
	ErrConversationNotFound: {
		Code:           ErrConversationNotFound,
		Message:        "对话不存在",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrConversationCreate: {
		Code:           ErrConversationCreate,
		Message:        "创建对话失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrMessageSend: {
		Code:           ErrMessageSend,
		Message:        "发送消息失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrMessageInvalid: {
		Code:           ErrMessageInvalid,
		Message:        "消息内容无效",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// 知识库相关错误
	ErrKnowledgeBaseNotFound: {
		Code:           ErrKnowledgeBaseNotFound,
		Message:        "知识库不存在",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrKnowledgeBaseCreate: {
		Code:           ErrKnowledgeBaseCreate,
		Message:        "创建知识库失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrKnowledgeBaseUpdate: {
		Code:           ErrKnowledgeBaseUpdate,
		Message:        "更新知识库失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrKnowledgeBaseDelete: {
		Code:           ErrKnowledgeBaseDelete,
		Message:        "删除知识库失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrFileUpload: {
		Code:           ErrFileUpload,
		Message:        "文件上传失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrFileNotFound: {
		Code:           ErrFileNotFound,
		Message:        "文件不存在",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrFileSize: {
		Code:           ErrFileSize,
		Message:        "文件大小超出限制",
		HTTPStatusCode: http.StatusBadRequest,
	},
	ErrFileType: {
		Code:           ErrFileType,
		Message:        "不支持的文件类型",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// 设置相关错误
	ErrModelNotFound: {
		Code:           ErrModelNotFound,
		Message:        "模型不存在",
		HTTPStatusCode: http.StatusNotFound,
	},
	ErrSettingsUpdate: {
		Code:           ErrSettingsUpdate,
		Message:        "更新设置失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
	ErrSettingsInvalid: {
		Code:           ErrSettingsInvalid,
		Message:        "设置参数无效",
		HTTPStatusCode: http.StatusBadRequest,
	},

	// Flowy SDK相关错误
	ErrFlowyConnection: {
		Code:           ErrFlowyConnection,
		Message:        "无法连接到AI服务",
		HTTPStatusCode: http.StatusServiceUnavailable,
	},
	ErrFlowyAuth: {
		Code:           ErrFlowyAuth,
		Message:        "AI服务认证失败",
		HTTPStatusCode: http.StatusUnauthorized,
	},
	ErrFlowyAPI: {
		Code:           ErrFlowyAPI,
		Message:        "AI服务调用失败",
		HTTPStatusCode: http.StatusInternalServerError,
	},
}

// GetErrorInfo 获取错误信息
func GetErrorInfo(code ErrorCode) ErrorInfo {
	if info, ok := errorInfoMap[code]; ok {
		return info
	}
	// 默认返回内部服务器错误
	return errorInfoMap[ErrInternalServer]
}

// GetHTTPStatus 获取错误对应的HTTP状态码
func GetHTTPStatus(code ErrorCode) int {
	return GetErrorInfo(code).HTTPStatusCode
}

// GetMessage 获取错误提示消息
func GetMessage(code ErrorCode) string {
	return GetErrorInfo(code).Message
}

// ErrorResponse 错误响应结构
type ErrorResponse struct {
	Success   bool      `json:"success"`              // 请求是否成功
	ErrorCode ErrorCode `json:"error_code"`           // 错误代码
	Message   string    `json:"message"`              // 错误消息
	Details   string    `json:"details,omitempty"`    // 错误详情
	RequestID string    `json:"request_id,omitempty"` // 请求ID
	Timestamp time.Time `json:"timestamp"`            // 时间戳
}

// APIError 自定义错误类型
type APIError struct {
	Code       ErrorCode
	Message    string
	Details    string
	HTTPStatus int
	Cause      error
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Cause != nil {
		return fmt.Sprintf("%s: %s (caused by: %v)", e.Code, e.Message, e.Cause)
	}
	return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

// WithCause 添加原因错误
func (e *APIError) WithCause(cause error) *APIError {
	e.Cause = cause
	if cause != nil {
		e.Details = cause.Error()
	}
	return e
}

// WithDetails 添加详细信息
func (e *APIError) WithDetails(details string) *APIError {
	e.Details = details
	return e
}

// NewAPIError 创建新的API错误
func NewAPIError(code ErrorCode, message string, httpStatus int) *APIError {
	// 如果没有提供自定义消息，使用默认消息
	if message == "" {
		message = GetMessage(code)
	}
	// 如果没有提供HTTP状态码，使用默认状态码
	if httpStatus == 0 {
		httpStatus = GetHTTPStatus(code)
	}
	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
	}
}

// RespondWithError 返回错误响应
func RespondWithError(c *gin.Context, err *APIError) {
	// 如果错误消息为空，使用默认消息
	message := err.Message
	if message == "" {
		message = GetMessage(err.Code)
	}

	response := ErrorResponse{
		Success:   false,
		ErrorCode: err.Code,
		Message:   message,
		Details:   err.Details,
		RequestID: c.GetString("request_id"),
		Timestamp: time.Now(),
	}

	LogError("API错误: %v", err)

	c.JSON(err.HTTPStatus, response)
}

// RespondWithValidationError 返回验证错误响应
func RespondWithValidationError(c *gin.Context, err error) {
	apiErr := NewAPIError(ErrInvalidRequest, "请求参数验证失败", http.StatusBadRequest).
		WithCause(err)
	RespondWithError(c, apiErr)
}

// ParseFlowyError 解析Flowy SDK错误并转换为APIError
func ParseFlowyError(err error) *APIError {
	if err == nil {
		return nil
	}

	errorMsg := err.Error()
	errorLower := strings.ToLower(errorMsg)

	var code ErrorCode

	// API相关错误
	if strings.Contains(errorLower, "connection refused") || strings.Contains(errorLower, "connection") {
		code = ErrFlowyConnection
	} else if strings.Contains(errorLower, "unauthorized") || strings.Contains(errorLower, "auth") {
		// 认证相关错误
		code = ErrFlowyAuth
	} else if strings.Contains(errorLower, "timeout") {
		// 超时错误
		code = ErrTimeout
	} else if strings.Contains(errorLower, "not found") {
		// 资源不存在错误
		code = ErrNotFound
	} else if strings.Contains(errorLower, "file size") {
		// 文件相关错误
		code = ErrFileSize
	} else if strings.Contains(errorLower, "file type") {
		code = ErrFileType
	} else if strings.Contains(errorLower, "knowledge") {
		// 知识库相关错误
		code = ErrFlowyAPI
	} else if strings.Contains(errorLower, "conversation") || strings.Contains(errorLower, "session") {
		// 对话相关错误
		code = ErrFlowyAPI
	} else {
		// 默认为内部服务器错误
		code = ErrInternalServer
	}

	return NewAPIError(code, "", 0).WithCause(err)
}

// WrapError 包装标准错误为APIError
func WrapError(err error, code ErrorCode, message string) *APIError {
	if err == nil {
		return nil
	}

	// 如果已经是APIError，直接返回
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	// 尝试从Flowy错误解析
	if strings.Contains(err.Error(), "flowy") || strings.Contains(err.Error(), "sdk") {
		return ParseFlowyError(err)
	}

	// 创建新的APIError
	httpStatus := GetHTTPStatus(code)
	if message == "" {
		message = GetMessage(code)
	}

	return &APIError{
		Code:       code,
		Message:    message,
		HTTPStatus: httpStatus,
		Cause:      err,
		Details:    err.Error(),
	}
}
