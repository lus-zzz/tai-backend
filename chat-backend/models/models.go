package models

import (
	"time"
)


// ChatRequest 聊天请求
// swagger:model
type ChatRequest struct {
	// 会话ID
	// required: true
	SessionID int `json:"sessionId" example:"123"`
	// 消息内容
	// required: true
	Content string `json:"content" example:"你好"`
	// 文件列表（可选）
	// required: false
	Files []string `json:"files" example:"[]"`
	// 模型ID
	// required: true
	ModelID int `json:"model_id"`
	// 温度参数
	// required: false
	Temperature float64 `json:"temperature"`
	// 顶部P参数
	// required: false
	TopP float64 `json:"top_p"`
	// 存在惩罚
	// required: false
	PresencePenalty float64 `json:"presence_penalty"`
	// 频率惩罚
	// required: false
	FrequencyPenalty float64 `json:"frequency_penalty"`
	// 最大token数
	// required: false
	MaxTokens int `json:"max_tokens"`
	// 是否流式输出
	// required: false
	Stream bool `json:"stream"`
	// 知识库ID列表
	// required: false
	KnowledgeBaseIDs []int `json:"knowledge_base_ids"`
}

// ChatMessageRequest 聊天消息请求（用于Swagger文档）
// swagger:model
type ChatMessageRequest struct {
	// 会话ID/对话ID
	// required: true
	SessionID int `json:"sessionId" example:"123"`
	// 消息内容
	// required: true
	Content string `json:"content" example:"你好"`
	// 文件列表（可选）
	// required: true
	Files []string `json:"files" example:"[]"`
}

// SSEChatEvent SSE聊天事件
// swagger:model
type SSEChatEvent struct {
	// 事件类型
	// required: true
	Type string `json:"type"`
	// 事件数据
	// required: true
	Data interface{} `json:"data"`
	// 消息ID
	// required: true
	ID string `json:"id"`
	// 错误信息
	// required: false
	Error string `json:"error,omitempty"`
}

// ChatResponse 聊天响应
// swagger:model
type ChatResponse struct {
	// 消息ID
	// required: true
	ID string `json:"id"`
	// 对话ID
	// required: true
	ConversationID string `json:"conversation_id"`
	// 消息内容
	// required: true
	Content string `json:"content"`
	// 角色: user/assistant
	// required: true
	Role string `json:"role"`
	// 状态
	// required: true
	Status string `json:"status"`
	// 引用信息列表
	// required: true
	References []Reference `json:"references"`
	// Token数量
	// required: true
	TokenCount int `json:"token_count"`
	// 创建时间
	// required: true
	CreatedAt time.Time `json:"created_at"`
}

// Reference 引用信息
// swagger:model
type Reference struct {
	// 文档ID
	// required: true
	DocumentID string `json:"document_id"`
	// 文档标题
	// required: true
	DocumentTitle string `json:"document_title"`
	// 引用内容
	// required: true
	Content string `json:"content"`
	// 相似度
	// required: true
	Similarity float64 `json:"similarity"`
	// 分块索引
	// required: true
	ChunkIndex int `json:"chunk_index"`
}

// MessageRecord 对话消息记录
// swagger:model
type MessageRecord struct {
	// 消息ID
	// required: true
	ID int `json:"id"`
	// 角色: user/assistant
	// required: true
	Role string `json:"role"`
	// 消息内容
	// required: true
	Content string `json:"content"`
	// 创建时间
	// required: true
	CreatedAt time.Time `json:"created_at"`
}

// ConversationHistoryResponse 对话历史响应
// swagger:model
type ConversationHistoryResponse struct {
	// 对话ID
	// required: true
	ConversationID string `json:"conversation_id"`
	// 消息列表
	// required: true
	Messages []MessageRecord `json:"messages"`
	// 消息总数
	// required: true
	Total int `json:"total"`
}

// Conversation 对话信息
// swagger:model
type Conversation struct {
	// 对话ID
	// required: true
	ID int `json:"id"`
	// 对话设置
	ConversationSettings `json:"settings"`
}

// ConversationListRequest 对话列表请求
// swagger:model
type ConversationListRequest struct {
	// 页码
	// required: true
	Page int `json:"page" form:"page"`
	// 每页数量
	// required: true
	PageSize int `json:"page_size" form:"page_size"`
}

// ConversationListResponse 对话列表响应
// swagger:model
type ConversationListResponse struct {
	// 对话列表
	// required: true
	Conversations []Conversation `json:"conversations"`
	// 总数
	// required: true
	Total int `json:"total"`
	// 当前页码
	// required: true
	Page int `json:"page"`
	// 每页数量
	// required: true
	PageSize int `json:"page_size"`
}

// KnowledgeBaseConfig 知识库配置（公用字段）
// swagger:model
type KnowledgeBaseConfig struct {
	// 知识库名称
	// required: true
	Name string `json:"name"`
	// 知识库描述
	// required: true
	Desc string `json:"desc"`
	// 向量模型ID
	// required: true
	VectorModel int `json:"vectorModel"`
	// 对话模型ID
	// required: true
	AgentModel int `json:"agentModel"`
	// 切片策略 固定尺寸:fixed  自然句:period   自然段落:paragraph
	// required: true
	ChunkStrategy string `json:"chunkStrategy"`
	// 切片大小
	// required: true
	ChunkSize int `json:"chunkSize"`
}

// KnowledgeBase 知识库（列表返回）
// swagger:model
type KnowledgeBase struct {
	// 知识库ID
	// required: true
	ID int `json:"id"`
	// 知识库配置
	KnowledgeBaseConfig `json:"config"`
	// 文件数量
	// required: true
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
	// 文件ID
	// required: true
	ID int `json:"id"`
	// 文件名称
	// required: true
	Name string `json:"name"`
	// 文件大小（字节）
	// required: true
	Size int `json:"size"`
	// 启用状态: true=启用, false=禁用
	// required: true
	Enable bool `json:"enable"`
	// 索引状态: 0=构建中, 1=完成, 2=失败
	// required: true
	Status int `json:"status"`
	// 上传时间
	// required: true
	UploadedAt time.Time `json:"uploaded_at"`
	// 索引进度百分比
	// required: true
	IndexPercent int `json:"index_percent"`
	// 错误信息，空字符串表示无错误
	// required: true
	ErrorMessage string `json:"errorMessage"`
}

// BatchUploadFilesRequest 批量上传文件请求
// swagger:model
type BatchUploadFilesRequest struct {
	// 上传文件路径列表（本地文件系统路径）
	// required: true
	FilePaths []string `json:"file_paths"`
}

// BatchUploadResult 单个文件上传结果
// swagger:model
type BatchUploadResult struct {
	// 文件路径
	// required: true
	FilePath string `json:"file_path"`
	// 上传是否成功
	// required: true
	Success bool `json:"success"`
	// 上传结果信息
	// required: true
	Message string `json:"message"`
	// 文件信息（成功时返回）
	// required: false
	File *KnowledgeFile `json:"file,omitempty"`
	// 错误信息（失败时返回）
	// required: false
	Error string `json:"error,omitempty"`
}

// BatchUploadResponse 批量上传响应
// swagger:model
type BatchUploadResponse struct {
	// 上传总数
	// required: true
	Total int `json:"total"`
	// 成功数
	// required: true
	SuccessCount int `json:"success_count"`
	// 失败数
	// required: true
	FailureCount int `json:"failure_count"`
	// 每个文件的上传结果
	// required: true
	Results []BatchUploadResult `json:"results"`
}

// APIResponse 通用API响应（成功和错误都使用这个结构）
// swagger:model
type APIResponse struct {
	// 请求是否成功
	// required: true
	Success bool `json:"success"`
	// 响应消息
	// required: true
	Message string `json:"message"`
	// 响应数据（成功时有值，失败时为null）
	// required: false
	Data interface{} `json:"data,omitempty"`
	// 错误代码（失败时有值）
	// required: false
	ErrorCode string `json:"error_code,omitempty"`
	// 错误详情（失败时有值）
	// required: false
	Details string `json:"details,omitempty"`
	// 时间戳
	// required: true
	Timestamp string `json:"timestamp"`
}

// WebSocketMessage WebSocket消息
// swagger:model
type WebSocketMessage struct {
	// 消息类型
	// required: true
	Type string `json:"type"`
	// 消息数据
	// required: true
	Data interface{} `json:"data"`
	// 消息ID
	// required: true
	ID string `json:"id"`
	// 错误信息
	// required: true
	Error string `json:"error"`
}

// StreamChatResponse 流式聊天响应
// swagger:model
type StreamChatResponse struct {
	// 消息ID
	// required: true
	ID string `json:"id"`
	// 完整内容
	// required: true
	Content string `json:"content"`
	// 增量内容
	// required: true
	Delta string `json:"delta"`
	// 是否完成
	// required: true
	Done bool `json:"done"`
}

// ConversationSettings 对话设置
// swagger:model
type ConversationSettings struct {
	// 对话名称
	// required: true
	Name string `json:"name"`
	// 对话描述
	// required: true
	Desc string `json:"desc"`
	// 模型ID（数字）
	// required: true
	ModelID int `json:"model_id"`
	// 多样性
	// required: true
	Temperature float64 `json:"temperature"`
	// 采样范围
	// required: true
	TopP float64 `json:"top_p"`
	// 词汇控制
	// required: true
	PresencePenalty float64 `json:"presence_penalty"`
	// 重复控制
	// required: true
	FrequencyPenalty float64 `json:"frequency_penalty"`
	// 响应类型  text/json
	// required: true
	ResponseType string `json:"responseType"`
	// 对话输出模式 true=流式输出, false=非流式输出
	// required: true
	Stream bool `json:"stream"`
	// 知识库ID列表
	// required: true
	KnowledgeBaseIDs []int `json:"knowledge_base_ids"`
	// 上下文限制（消息数量）
	// required: true
	ContextLimit int `json:"contextLimit"`
}

// NewDefaultConversationSettings 创建默认的对话设置
func NewDefaultConversationSettings() *ConversationSettings {
	return &ConversationSettings{
		Name:             "",
		Desc:             "",
		KnowledgeBaseIDs: []int{},
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
	// 默认对话模型ID
	// required: true
	ChatModelID int `json:"chat_model_id"`
	// 默认向量模型ID（嵌入模型）
	// required: true
	EmbeddingModelID int `json:"embedding_model_id"`
}

// DefaultConversationConfig 默认对话配置
type DefaultConversationConfig = ConversationSettings

// KnowledgeBaseSettings 知识库设置
// swagger:model
type KnowledgeBaseSettings struct {
	// 知识库名称
	// required: true
	Name string `json:"name"`
	// 知识库描述
	// required: true
	Desc string `json:"desc"`
	// 向量模型ID
	// required: true
	VectorModel int `json:"vectorModel"`
	// 对话模型ID
	// required: true
	AgentModel int `json:"agentModel"`
	// 切片策略
	// required: true
	ChunkStrategy string `json:"chunkStrategy"`
	// 切片大小
	// required: true
	ChunkSize int `json:"chunkSize"`
}

// DefaultKnowledgeBaseConfig 默认知识库配置
type DefaultKnowledgeBaseConfig = KnowledgeBaseSettings

// DefaultSettings 统一的默认设置（用于持久化）
// swagger:model
type DefaultSettings struct {
	// 默认模型设置
	// required: true
	Models DefaultModelSettings `json:"models"`
	// 默认对话配置
	// required: true
	Conversation DefaultConversationConfig `json:"conversation"`
	// 默认知识库配置
	// required: true
	KnowledgeBase DefaultKnowledgeBaseConfig `json:"knowledge_base"`
	// 更新时间
	// required: true
	UpdatedAt time.Time `json:"updated_at"`
}

// DefaultKnowledgeBaseSettings 默认知识库设置（用于持久化）
// swagger:model
type DefaultKnowledgeBaseSettings struct {
	KnowledgeBaseSettings
	// 更新时间
	// required: true
	UpdatedAt time.Time `json:"updated_at"`
}

// FileToggleEnableRequest 文件启用/禁用请求
// swagger:model
type FileToggleEnableRequest struct {
	// true=启用, false=禁用
	// required: true
	Enable *bool `json:"enable" example:"true"`
}

// ModelInfo 模型信息
// swagger:model
type ModelInfo struct {
	// 模型ID
	// required: true
	ID int `json:"id"`
	// 模型名称
	// required: true
	Name string `json:"name"`
	// 模型类型
	// required: true
	Type string `json:"type"`
	// 模型描述
	// required: true
	Description string `json:"description"`
	// 是否启用
	// required: true
	Enabled bool `json:"enabled"`
}

// SupportedChatModel 支持的聊天模型
// swagger:model
type SupportedChatModel struct {
	// 模型ID
	// required: true
	ID int `json:"id"`
	// 模型名称
	// required: true
	Name string `json:"name"`
	// 模型描述
	// required: true
	Description string `json:"description"`
}

// SupportedVectorModel 支持的向量模型
// swagger:model
type SupportedVectorModel struct {
	// 模型ID
	// required: true
	ID int `json:"id"`
	// 模型名称
	// required: true
	Name string `json:"name"`
	// 模型描述
	// required: true
	Description string `json:"description"`
	// 向量维度
	// required: true
	Dimensions int `json:"dimensions"`
}

// ModelID 模型ID响应数据
// swagger:model
type ModelID struct {
	// 模型ID
	// required: true
	ID int `json:"id" example:"1"`
}

// ====== Swagger 响应定义 ======

// 通用响应类型在 utils 包中定义
// - SuccessResponse：统一成功响应 (utils/response.go)
// - ErrorResponse：统一错误响应 (utils/errors.go)

// ConversationSuccessResponse 对话创建成功响应
// swagger:response ConversationSuccessResponse
type ConversationSuccessResponse struct {
	// 请求是否成功
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 对话数据
		// required: true
		Data Conversation `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ConversationListSuccessResponse 对话列表获取成功响应
// swagger:response ConversationListSuccessResponse
type ConversationListSuccessResponse struct {
	// 对话列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 对话列表数据
		// required: true
		Data ConversationListResponse `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// EmptySuccessResponse 空数据成功响应
// swagger:response EmptySuccessResponse
type EmptySuccessResponse struct {
	// 空数据响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 空数据
		// required: false
		Data interface{} `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ConversationSettingsSuccessResponse 对话设置成功响应
// swagger:response ConversationSettingsSuccessResponse
type ConversationSettingsSuccessResponse struct {
	// 对话设置响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 对话设置数据
		// required: true
		Data ConversationSettings `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ConversationHistorySuccessResponse 对话历史成功响应
// swagger:response ConversationHistorySuccessResponse
type ConversationHistorySuccessResponse struct {
	// 对话历史响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 对话历史数据
		// required: true
		Data ConversationHistoryResponse `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// KnowledgeBaseListSuccessResponse 知识库列表成功响应
// swagger:response KnowledgeBaseListSuccessResponse
type KnowledgeBaseListSuccessResponse struct {
	// 知识库列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 知识库列表数据
		// required: true
		Data []KnowledgeBase `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// KnowledgeBaseSuccessResponse 知识库成功响应
// swagger:response KnowledgeBaseSuccessResponse
type KnowledgeBaseSuccessResponse struct {
	// 知识库响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 知识库数据
		// required: true
		Data KnowledgeBase `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// KnowledgeFileListSuccessResponse 知识库文件列表成功响应
// swagger:response KnowledgeFileListSuccessResponse
type KnowledgeFileListSuccessResponse struct {
	// 知识库文件列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 知识库文件列表数据
		// required: true
		Data []KnowledgeFile `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// KnowledgeFileSuccessResponse 知识库文件成功响应
// swagger:response KnowledgeFileSuccessResponse
type KnowledgeFileSuccessResponse struct {
	// 知识库文件响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 知识库文件数据
		// required: true
		Data KnowledgeFile `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// DefaultSettingsSuccessResponse 默认设置成功响应
// swagger:response DefaultSettingsSuccessResponse
type DefaultSettingsSuccessResponse struct {
	// 默认设置响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 默认设置数据
		// required: true
		Data DefaultSettings `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ModelSaveSuccessResponse 模型保存成功响应
// swagger:response ModelSaveSuccessResponse
type ModelSaveSuccessResponse struct {
	// 模型保存响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 模型ID数据
		// required: true
		Data ModelID `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ModelListSuccessResponse 模型列表成功响应 (通用)
// swagger:response ModelListSuccessResponse
type ModelListSuccessResponse struct {
	// 模型列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 模型列表数据
		// required: true
		Data []ModelInfo `json:"data"` // ModelInfo类型的列表 (适用于可用模型接口)
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// SupportedChatModelListSuccessResponse 支持的聊天模型列表成功响应
// swagger:response SupportedChatModelListSuccessResponse
type SupportedChatModelListSuccessResponse struct {
	// 支持的聊天模型列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 支持的聊天模型列表数据
		// required: true
		Data []SupportedChatModel `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// SupportedVectorModelListSuccessResponse 支持的向量模型列表成功响应
// swagger:response SupportedVectorModelListSuccessResponse
type SupportedVectorModelListSuccessResponse struct {
	// 支持的向量模型列表响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 支持的向量模型列表数据
		// required: true
		Data []SupportedVectorModel `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// ====== 通用响应 ======

// 所有API现在使用统一的响应结构：
// - SuccessResponse：定义在 utils/response.go
// - ErrorResponse：定义在 utils/errors.go

// ModelStatusEnableRequest 模型状态启用/禁用请求
// swagger:model
type ModelStatusEnableRequest struct {
	// true=启用, false=禁用
	// required: true
	Enable bool `json:"enable" example:"true"`
}

// ChatInput 快捷方式推荐输入
// swagger:model
type ChatInput struct {
	// 用户输入的自然语言描述，会有截断（模型输入tokens限制为512，包括prompt），如果用户输入太长
	// required: true
	UserInput string `json:"user_input" example:"我想要调整模型的温度参数"`
	// 根据用户的输入，返回推荐设置的数量
	// required: true
	RecommendNum int32 `json:"recommend_num" example:"3"`
}

// RecommendData 推荐数据响应
// swagger:model
type RecommendData struct {
	// 推荐的设置名列表，与recommend_score分数按顺序对应
	// required: true
	SettingName []string `json:"setting_name" example:"[\"temperature\", \"max_tokens\", \"top_p\"]"`
	// 推荐的设置名的分数，与setting_name分数按顺序对应
	// required: true
	RecommendScore []float64 `json:"recommend_score" example:"[0.95, 0.87, 0.76]"`
	// 是否为用户执行推荐名列表第一个设置，True为执行，False为不执行
	// required: true
	ProcessForUser bool `json:"process_for_user" example:"true"`
}

// SettingName 支持的设置名响应
// swagger:model
type SettingName struct {
	// 支持的全部设置列表
	// required: true
	SupportedSettingName []string `json:"supported_setting_name" example:"[\"temperature\", \"max_tokens\", \"top_p\", \"frequency_penalty\", \"presence_penalty\"]"`
}

// RecommendDataSuccessResponse 推荐数据成功响应
// swagger:response RecommendDataSuccessResponse
type RecommendDataSuccessResponse struct {
	// 推荐数据响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 推荐数据
		// required: true
		Data RecommendData `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// SettingNameSuccessResponse 支持设置名成功响应
// swagger:response SettingNameSuccessResponse
type SettingNameSuccessResponse struct {
	// 支持设置名响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 设置名数据
		// required: true
		Data SettingName `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}
