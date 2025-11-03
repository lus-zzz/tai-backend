package errors

import (
	"fmt"
	"net/http"
)

// ErrorCode 定义错误代码类型
type ErrorCode string

const (
	// 通用错误代码
	ErrCodeUnknown        ErrorCode = "UNKNOWN"
	ErrCodeInvalidRequest ErrorCode = "INVALID_REQUEST"
	ErrCodeUnauthorized   ErrorCode = "UNAUTHORIZED"
	ErrCodeForbidden      ErrorCode = "FORBIDDEN"
	ErrCodeNotFound       ErrorCode = "NOT_FOUND"
	ErrCodeInternalError  ErrorCode = "INTERNAL_ERROR"
	ErrCodeTimeout        ErrorCode = "TIMEOUT"
	ErrCodeNetworkError   ErrorCode = "NETWORK_ERROR"

	// 业务错误代码
	ErrCodeModelNotFound ErrorCode = "MODEL_NOT_FOUND"
	ErrCodeFileNotFound  ErrorCode = "FILE_NOT_FOUND"
	ErrCodeAgentNotFound ErrorCode = "AGENT_NOT_FOUND"
	ErrCodeInvalidFile   ErrorCode = "INVALID_FILE"
)

// FlowyError Flowy SDK 自定义错误类型
type FlowyError struct {
	Code       ErrorCode `json:"code"`
	Message    string    `json:"message"`
	Details    string    `json:"details,omitempty"`
	StatusCode int       `json:"status_code,omitempty"`
}

// Error 实现 error 接口
func (e *FlowyError) Error() string {
	if e.Details != "" {
		return fmt.Sprintf("[%s] %s: %s", e.Code, e.Message, e.Details)
	}
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

// New 创建新的 FlowyError
func New(code ErrorCode, message string) *FlowyError {
	return &FlowyError{
		Code:    code,
		Message: message,
	}
}

// WithDetails 添加详细信息
func (e *FlowyError) WithDetails(details string) *FlowyError {
	e.Details = details
	return e
}

// FromHTTPStatus 根据HTTP状态码创建错误
func FromHTTPStatus(statusCode int, message string) *FlowyError {
	var code ErrorCode
	switch statusCode {
	case http.StatusBadRequest:
		code = ErrCodeInvalidRequest
	case http.StatusUnauthorized:
		code = ErrCodeUnauthorized
	case http.StatusForbidden:
		code = ErrCodeForbidden
	case http.StatusNotFound:
		code = ErrCodeNotFound
	case http.StatusInternalServerError:
		code = ErrCodeInternalError
	default:
		code = ErrCodeUnknown
	}

	return &FlowyError{
		Code:       code,
		Message:    message,
		StatusCode: statusCode,
	}
}
