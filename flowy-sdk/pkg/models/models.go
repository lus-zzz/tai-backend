package models

import (
	"time"
)

// BaseResponse API响应基础结构
type BaseResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
	Success bool        `json:"success"`
	Time    time.Time   `json:"time"`
}

// PaginationMeta 分页元数据
type PaginationMeta struct {
	Page     int `json:"page"`
	PageSize int `json:"page_size"`
	Total    int `json:"total"`
	Pages    int `json:"pages"`
}

// PaginationResponse 分页响应
type PaginationResponse struct {
	BaseResponse
	Meta PaginationMeta `json:"meta"`
}

// Model 模型基础结构
type Model struct {
	ID        string                `json:"id"`
	Name      string                `json:"name"`
	Type      string                `json:"type"`     // chat, embedding
	Provider  string                `json:"provider"` // openai, deepseek, ollama
	Enabled   bool                  `json:"enabled"`
	Config    ModelConnectionConfig `json:"config"`
	CreatedAt time.Time             `json:"created_at"`
	UpdatedAt time.Time             `json:"updated_at"`
}

// ModelConnectionConfig 模型连接配置
type ModelConnectionConfig struct {
	APIKey    string                 `json:"api_key,omitempty"`
	BaseURL   string                 `json:"base_url,omitempty"`
	Model     string                 `json:"model,omitempty"`
	MaxTokens int                    `json:"max_tokens,omitempty"`
	Extra     map[string]interface{} `json:"extra,omitempty"`
}

// SupportedModel 支持的模型
type SupportedModel struct {
	Type        string   `json:"type"`
	Providers   []string `json:"providers"`
	Description string   `json:"description"`
}

// KnowledgeBase 知识库
type KnowledgeBase struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Type        string                 `json:"type"` // qa, product
	Description string                 `json:"description"`
	Config      map[string]interface{} `json:"config"`
	Enabled     bool                   `json:"enabled"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// File 文件信息
type File struct {
	ID           string                 `json:"id"`
	Name         string                 `json:"name"`
	OriginalName string                 `json:"original_name"`
	Type         string                 `json:"type"`
	Size         int64                  `json:"size"`
	OSS          string                 `json:"oss"`
	Status       string                 `json:"status"`
	Enabled      bool                   `json:"enabled"`
	Config       map[string]interface{} `json:"config"`
	ChunkConfig  ChunkConfig            `json:"chunk_config"`
	RecallConfig RecallConfig           `json:"recall_config"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
}

// ChunkConfig 切片配置
type ChunkConfig struct {
	Strategy string `json:"strategy"`
	Size     int    `json:"size"`
	Overlap  int    `json:"overlap"`
}

// RecallConfig 召回配置
type RecallConfig struct {
	Strategy   string  `json:"strategy"`
	Threshold  float64 `json:"threshold"`
	MaxResults int     `json:"max_results"`
	Prompt     string  `json:"prompt"`
}

// QAItem 问答项
type QAItem struct {
	ID        string    `json:"id"`
	Question  string    `json:"question"`
	Answer    string    `json:"answer"`
	FileID    string    `json:"file_id"`
	Enabled   bool      `json:"enabled"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Product 产品
type Product struct {
	ID          string                 `json:"id"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Files       []File                 `json:"files"`
	Properties  []ProductProperty      `json:"properties"`
	Config      map[string]interface{} `json:"config"`
	CreatedAt   time.Time              `json:"created_at"`
	UpdatedAt   time.Time              `json:"updated_at"`
}

// ProductProperty 产品属性
type ProductProperty struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Value string `json:"value"`
	Type  string `json:"type"`
}

// Agent 智能体
type Agent struct {
	ID          string      `json:"id"`
	Name        string      `json:"name"`
	Description string      `json:"description"`
	Type        string      `json:"type"`
	Config      AgentConfig `json:"config"`
	Enabled     bool        `json:"enabled"`
	CreatedAt   time.Time   `json:"created_at"`
	UpdatedAt   time.Time   `json:"updated_at"`
}

// AgentConfig 智能体配置
type AgentConfig struct {
	Model       string                 `json:"model"`
	Prompt      string                 `json:"prompt"`
	Knowledge   []string               `json:"knowledge"`
	MaxTokens   int                    `json:"max_tokens"`
	Temperature float64                `json:"temperature"`
	TopP        float64                `json:"top_p"`
	Extra       map[string]interface{} `json:"extra"`
}

// Session 会话
type Session struct {
	ID        string    `json:"id"`
	AgentID   string    `json:"agent_id"`
	Title     string    `json:"title"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// Message 消息
type Message struct {
	ID        string                 `json:"id"`
	SessionID string                 `json:"session_id"`
	Role      string                 `json:"role"` // user, assistant, system
	Content   string                 `json:"content"`
	Type      string                 `json:"type"` // text, image, file
	Metadata  map[string]interface{} `json:"metadata"`
	Rating    *int                   `json:"rating"` // 1: 好, -1: 差
	CreatedAt time.Time              `json:"created_at"`
}

// FileUploadResponse 文件上传响应
type FileUploadResponse struct {
	FileID string `json:"file_id"`
	OSS    string `json:"oss"`
	Name   string `json:"name"`
	Size   int64  `json:"size"`
}

// AsyncChatRequest 异步对话请求
type AsyncChatRequest struct {
	SessionID string `json:"session_id"`
	Message   string `json:"message"`
	Stream    bool   `json:"stream"`
}

// AsyncChatResponse 异步对话响应
type AsyncChatResponse struct {
	ID       string `json:"id"`
	Content  string `json:"content"`
	Status   string `json:"status"`
	Finished bool   `json:"finished"`
}

// AgentCreateRequest 创建Agent请求
type AgentCreateRequest struct {
	Name             string   `json:"name"`
	Description      string   `json:"description"`
	Type             string   `json:"type,omitempty"`
	Model            string   `json:"model,omitempty"`
	Temperature      float32  `json:"temperature,omitempty"`
	MaxTokens        int      `json:"max_tokens,omitempty"`
	KnowledgeBaseIDs []string `json:"knowledge_base_ids,omitempty"`
	Prompt           string   `json:"prompt,omitempty"`
}

// KnowledgeSearchRequest 知识库搜索请求
type KnowledgeSearchRequest struct {
	Query    string  `json:"query"`
	TopK     int     `json:"top_k"`
	MinScore float64 `json:"min_score"`
}

// KnowledgeSearchResult 知识库搜索结果
type KnowledgeSearchResult struct {
	DocumentID string  `json:"document_id"`
	Title      string  `json:"title"`
	Content    string  `json:"content"`
	Score      float64 `json:"score"`
	ChunkIndex int     `json:"chunk_index"`
}

// KnowledgeSearchResponse 知识库搜索响应
type KnowledgeSearchResponse struct {
	Results []KnowledgeSearchResult `json:"results"`
	Total   int                     `json:"total"`
}
