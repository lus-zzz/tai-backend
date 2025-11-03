package models

import (
	"time"

	agentSvc "flowy-sdk/services/agent"
)

// ChatRequest 聊天请求 - 直接使用 SDK 的 AsyncChatRequest
// swagger:model
type ChatRequest struct {
	agentSvc.AsyncChatRequest
}

// ChatMessageRequest 聊天消息请求（用于Swagger文档）
// swagger:model
type ChatMessageRequest struct {
	SessionID int      `json:"sessionId" example:"123"`                                  // 会话ID/对话ID
	Content   string   `json:"content" example:"你好"`                                     // 消息内容
	RequestID string   `json:"requestId" example:"550e8400-e29b-41d4-a716-446655440000"` // 请求ID (UUID)
	Files     []string `json:"files" example:"[]"`                                       // 文件列表（可选）
}

// SSEChatEvent SSE聊天事件 - 直接使用 SDK 的 StreamEvent
// swagger:model
type SSEChatEvent = agentSvc.StreamEvent

// ChatResponse 聊天响应
// swagger:model
type ChatResponse struct {
	ID             string      `json:"id"`                   // 消息ID
	ConversationID string      `json:"conversation_id"`      // 对话ID
	Content        string      `json:"content"`              // 消息内容
	Role           string      `json:"role"`                 // 角色: user/assistant
	Status         string      `json:"status"`               // 状态
	References     []Reference `json:"references,omitempty"` // 引用信息列表
	TokenCount     int         `json:"token_count"`          // Token数量
	CreatedAt      time.Time   `json:"created_at"`           // 创建时间
}

// Reference 引用信息
// swagger:model
type Reference struct {
	DocumentID    string  `json:"document_id"`    // 文档ID
	DocumentTitle string  `json:"document_title"` // 文档标题
	Content       string  `json:"content"`        // 引用内容
	Similarity    float64 `json:"similarity"`     // 相似度
	ChunkIndex    int     `json:"chunk_index"`    // 分块索引
}

// MessageRecord 对话消息记录
// swagger:model
type MessageRecord struct {
	ID        int       `json:"id"`         // 消息ID
	Role      string    `json:"role"`       // 角色: user/assistant
	Content   string    `json:"content"`    // 消息内容
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// ConversationHistoryResponse 对话历史响应
// swagger:model
type ConversationHistoryResponse struct {
	ConversationID string                   `json:"conversation_id"` // 对话ID
	Messages       []agentSvc.SessionRecord `json:"messages"`        // 消息列表
	Total          int                      `json:"total"`           // 消息总数
}

// Conversation 对话信息
// swagger:model
type Conversation struct {
	ID                   int `json:"id"` // 对话ID
	ConversationSettings     // 嵌入对话配置
}

// ConversationListRequest 对话列表请求
// swagger:model
type ConversationListRequest struct {
	Page     int `json:"page" form:"page"`           // 页码
	PageSize int `json:"page_size" form:"page_size"` // 每页数量
}

// ConversationListResponse 对话列表响应
// swagger:model
type ConversationListResponse struct {
	Conversations []Conversation `json:"conversations"` // 对话列表
	Total         int            `json:"total"`         // 总数
	Page          int            `json:"page"`          // 当前页码
	PageSize      int            `json:"page_size"`     // 每页数量
}

// KnowledgeBaseConfig 知识库配置（公用字段）
// swagger:model
type KnowledgeBaseConfig struct {
	Name          string `json:"name"`          // 知识库名称
	Desc          string `json:"desc"`          // 知识库描述
	VectorModel   int    `json:"vectorModel"`   // 向量模型ID
	AgentModel    int    `json:"agentModel"`    // 对话模型ID
	ChunkStrategy string `json:"chunkStrategy"` // 切片策略 固定尺寸:fixed  自然句:period   自然段落:paragraph
	ChunkSize     int    `json:"chunkSize"`     // 切片大小
}

// KnowledgeBase 知识库（列表返回）
// swagger:model
type KnowledgeBase struct {
	ID int `json:"id"`
	KnowledgeBaseConfig
	FileCount int `json:"file_count"`
}

// KnowledgeBaseCreateRequest 创建知识库请求
// swagger:model
type KnowledgeBaseCreateRequest = KnowledgeBaseConfig

// UpdateKnowledgeBaseRequest 更新知识库请求
// swagger:model
type UpdateKnowledgeBaseRequest = KnowledgeBaseConfig

// KnowledgeFile 知识库文件
// swagger:model
type KnowledgeFile struct {
	ID           int       `json:"id"`            // 文件ID
	Name         string    `json:"name"`          // 文件名称
	Size         int       `json:"size"`          // 文件大小（字节）
	Enable       bool      `json:"enable"`        // 启用状态: true=启用, false=禁用
	Status       int       `json:"status"`        // 索引状态: 0=构建中, 1=完成, 2=失败
	UploadedAt   time.Time `json:"uploaded_at"`   // 上传时间
	IndexPercent int       `json:"index_percent"` // 索引进度百分比
	ErrorMessage string    `json:"errorMessage"`  // 错误信息，空字符串表示无错误
}

// APIResponse 通用API响应
// swagger:model
type APIResponse struct {
	Success bool        `json:"success"`           // 请求是否成功
	Message string      `json:"message,omitempty"` // 响应消息
	Data    interface{} `json:"data,omitempty"`    // 响应数据
	Error   string      `json:"error,omitempty"`   // 错误信息
}

// WebSocketMessage WebSocket消息
// swagger:model
type WebSocketMessage struct {
	Type  string      `json:"type"`            // 消息类型
	Data  interface{} `json:"data"`            // 消息数据
	ID    string      `json:"id,omitempty"`    // 消息ID
	Error string      `json:"error,omitempty"` // 错误信息
}

// StreamChatResponse 流式聊天响应
// swagger:model
type StreamChatResponse struct {
	ID      string `json:"id"`      // 消息ID
	Content string `json:"content"` // 完整内容
	Delta   string `json:"delta"`   // 增量内容
	Done    bool   `json:"done"`    // 是否完成
}

// ConversationSettings 对话设置
// @Description 对话设置信息
// swagger:model
type ConversationSettings struct {
	Name             string   `json:"name"`                   // 对话名称
	Desc             string   `json:"desc"`                   // 对话描述
	ModelID          int      `json:"model_id"`               // 模型ID（数字）
	Temperature      float64  `json:"temperature"`            // 多样性
	TopP             float64  `json:"top_p"`                  // 采样范围
	PresencePenalty  float64  `json:"presence_penalty"`       // 词汇控制
	FrequencyPenalty float64  `json:"frequency_penalty"`      // 重复控制
	ResponseType     string   `json:"responseType,omitempty"` // 响应类型  text/json
	Stream           bool     `json:"stream"`                 // 对话输出模式 true=流式输出, false=非流式输出
	KnowledgeBaseIDs []string `json:"knowledge_base_ids"`
	ContextLimit     int      `json:"contextLimit"` // 上下文限制（消息数量）
}

// NewDefaultConversationSettings 创建默认的对话设置
func NewDefaultConversationSettings() *ConversationSettings {
	return &ConversationSettings{
		Name:             "",
		Desc:             "",
		KnowledgeBaseIDs: []string{},
		Temperature:      1,
		TopP:             1,
		PresencePenalty:  0.0,
		FrequencyPenalty: 0.0,
		ResponseType:     "text",
		Stream:           true,
		ContextLimit:     16,
		ModelID:          2,
	}
}

// DefaultModelSettings 默认模型设置
// swagger:model
type DefaultModelSettings struct {
	ChatModelID      int `json:"chat_model_id"`      // 默认对话模型ID
	EmbeddingModelID int `json:"embedding_model_id"` // 默认向量模型ID（嵌入模型）
}

// DefaultConversationConfig 默认对话配置
type DefaultConversationConfig = ConversationSettings

// KnowledgeBaseSettings 知识库设置
// swagger:model
type KnowledgeBaseSettings struct {
	Name          string `json:"name"`          // 知识库名称
	Desc          string `json:"desc"`          // 知识库描述
	VectorModel   int    `json:"vectorModel"`   // 向量模型ID
	AgentModel    int    `json:"agentModel"`    // 对话模型ID
	ChunkStrategy string `json:"chunkStrategy"` // 切片策略
	ChunkSize     int    `json:"chunkSize"`     // 切片大小
}

// DefaultKnowledgeBaseConfig 默认知识库配置
type DefaultKnowledgeBaseConfig = KnowledgeBaseSettings

// DefaultSettings 统一的默认设置（用于持久化）
// swagger:model
type DefaultSettings struct {
	Models        DefaultModelSettings       `json:"models"`         // 默认模型设置
	Conversation  DefaultConversationConfig  `json:"conversation"`   // 默认对话配置
	KnowledgeBase DefaultKnowledgeBaseConfig `json:"knowledge_base"` // 默认知识库配置
	UpdatedAt     time.Time                  `json:"updated_at"`     // 更新时间
}

// DefaultKnowledgeBaseSettings 默认知识库设置（用于持久化）
// swagger:model
type DefaultKnowledgeBaseSettings struct {
	KnowledgeBaseSettings
	UpdatedAt time.Time `json:"updated_at"`
}

// FileToggleEnableRequest 文件启用/禁用请求
// swagger:model
type FileToggleEnableRequest struct {
	Enable *bool `json:"enable" example:"true"` // true=启用, false=禁用
}

// ModelStatusEnableRequest 模型状态启用/禁用请求
// swagger:model ModelStatusEnableRequest
type ModelStatusEnableRequest struct {
	Enable bool `json:"enable" example:"true"` // true=启用, false=禁用
}
