package utils

import (
	"fmt"
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

// ErrorCodeMapping 错误码映射信息
type ErrorCodeMapping struct {
	CodeNum int    // 数字错误代码
	Message string // 默认中文提示消息
}

// errorMapping 错误码映射表 - 集中管理
var errorMapping = map[ErrorCode]ErrorCodeMapping{
	// 通用错误
	ErrInvalidRequest: {CodeNum: 400, Message: "请求参数无效"},
	ErrInternalServer: {CodeNum: 500, Message: "服务器内部错误"},
	ErrUnauthorized:   {CodeNum: 401, Message: "未授权访问"},
	ErrForbidden:      {CodeNum: 403, Message: "禁止访问"},
	ErrNotFound:       {CodeNum: 404, Message: "资源不存在"},
	ErrTimeout:        {CodeNum: 408, Message: "请求超时"},
	ErrRateLimited:    {CodeNum: 429, Message: "请求过于频繁"},

	// 聊天相关错误
	ErrConversationNotFound: {CodeNum: 404, Message: "对话不存在"},
	ErrConversationCreate:   {CodeNum: 500, Message: "创建对话失败"},
	ErrMessageSend:          {CodeNum: 500, Message: "发送消息失败"},
	ErrMessageInvalid:       {CodeNum: 400, Message: "消息内容无效"},

	// 知识库相关错误
	ErrKnowledgeBaseNotFound: {CodeNum: 404, Message: "知识库不存在"},
	ErrKnowledgeBaseCreate:   {CodeNum: 500, Message: "创建知识库失败"},
	ErrKnowledgeBaseUpdate:   {CodeNum: 500, Message: "更新知识库失败"},
	ErrKnowledgeBaseDelete:   {CodeNum: 500, Message: "删除知识库失败"},
	ErrFileUpload:            {CodeNum: 500, Message: "文件上传失败"},
	ErrFileNotFound:          {CodeNum: 404, Message: "文件不存在"},
	ErrFileSize:              {CodeNum: 400, Message: "文件大小超出限制"},
	ErrFileType:              {CodeNum: 400, Message: "不支持的文件类型"},

	// 设置相关错误
	ErrModelNotFound:   {CodeNum: 404, Message: "模型不存在"},
	ErrSettingsUpdate:  {CodeNum: 500, Message: "更新设置失败"},
	ErrSettingsInvalid: {CodeNum: 400, Message: "设置参数无效"},

	// Flowy SDK相关错误
	ErrFlowyConnection: {CodeNum: 503, Message: "无法连接到AI服务"},
	ErrFlowyAuth:       {CodeNum: 401, Message: "AI服务认证失败"},
	ErrFlowyAPI:        {CodeNum: 500, Message: "AI服务调用失败"},
}

// APIError 自定义错误类型 - 简化版本
type APIError struct {
	ErrorCode ErrorCode // 错误代码
	Err       error     // 原始错误
}

// Error 实现error接口
func (e *APIError) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("%s: %v", e.ErrorCode, e.Err)
	}
	return string(e.ErrorCode)
}

// NewAPIError 创建新的API错误 - 简化版本
func NewAPIError(errorCode ErrorCode, err error) *APIError {
	return &APIError{
		ErrorCode: errorCode,
		Err:       err,
	}
}

// GetErrorMapping 获取错误码映射信息
func GetErrorMapping(errorCode ErrorCode) ErrorCodeMapping {
	if mapping, ok := errorMapping[errorCode]; ok {
		return mapping
	}
	// 默认返回内部服务器错误
	return errorMapping[ErrInternalServer]
}

// WrapError 包装标准错误为APIError
func WrapError(err error, code ErrorCode) *APIError {
	if err == nil {
		return nil
	}

	// 如果已经是APIError，直接返回
	if apiErr, ok := err.(*APIError); ok {
		return apiErr
	}

	// 创建新的APIError
	return NewAPIError(code, err)
}
