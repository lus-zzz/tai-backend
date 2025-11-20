package interfaces

import (
	"context"

	"chat-backend/models"
)

// ChatServiceInterface 聊天服务接口
type ChatServiceInterface interface {
	// CreateConversation 创建对话
	CreateConversation(ctx context.Context, settings *models.ConversationSettings) (*models.Conversation, error)

	// SendMessage 发送消息并返回SSE流
	SendMessage(ctx context.Context, req *models.ChatRequest, eventChan chan<- models.SSEChatEvent) error

	// ListConversations 获取对话列表
	ListConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error)

	// DeleteConversation 删除对话
	DeleteConversation(ctx context.Context, conversationID string) error

	// GetConversations 获取对话列表 (别名方法)
	GetConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error)

	// UpdateConversationSettings 更新对话设置
	UpdateConversationSettings(ctx context.Context, conversationID string, settings *models.ConversationSettings) error

	// GetConversationSettings 获取对话设置
	GetConversationSettings(ctx context.Context, conversationID string) (*models.ConversationSettings, error)

	// GetConversationHistory 获取对话历史记录
	GetConversationHistory(ctx context.Context, conversationID string) (*models.ConversationHistoryResponse, error)
}
