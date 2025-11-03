package agent

import (
	"context"
	"encoding/json"
	"io"
	"strings"

	"flowy-sdk/pkg/client"
	"flowy-sdk/pkg/errors"
	"flowy-sdk/pkg/models"
)

// ==================== 请求和响应类型定义 ====================

// CreateAgentRequest 创建Agent请求
// API: POST /agent/create
type CreateAgentRequest struct {
	Name   string `json:"name"`   // Agent名称
	Desc   string `json:"desc"`   // 描述
	Type   int    `json:"type"`   // Agent类型
	Avatar string `json:"avatar"` // 头像URL
}

// UpdateAgentRequest 更新Agent请求
// API: POST /agent/update
type UpdateAgentRequest struct {
	Name   string `json:"name"`   // Agent名称
	Desc   string `json:"desc"`   // 描述
	Avatar string `json:"avatar"` // 头像URL
	ID     int    `json:"id"`     // 必填: 需要更新的AgentID
}

// AgentInfo Agent信息
type AgentInfo struct {
	ID             int    `json:"id"`             // AgentID
	Name           string `json:"name"`           // Agent名称
	Desc           string `json:"desc"`           // Agent描述
	Type           int    `json:"type"`           // Agent类型
	DefaultSetting int    `json:"defaultSetting"` // Agent默认的配置ID
	TopShow        bool   `json:"topShow"`        // (无意义字段)
	Avatar         string `json:"avatar"`         // 头像URL
}

// AgentListResponse Agent列表响应
type AgentListResponse struct {
	Total   int         `json:"total"`   // 总数
	Records []AgentInfo `json:"records"` // Agent记录列表
}

// AgentDetailResponse Agent详情响应
type AgentDetailResponse struct {
	ID             int    `json:"id"`             // AgentID
	Name           string `json:"name"`           // Agent名称
	Desc           string `json:"desc"`           // Agent描述
	Type           int    `json:"type"`           // Agent类型
	DefaultSetting int    `json:"defaultSetting"` // 默认配置ID
	TopShow        bool   `json:"topShow"`        // (无意义字段)
	Avatar         string `json:"avatar"`         // 头像URL
}

// ModelConfig 模型配置
type ModelConfig struct {
	ID               int     `json:"id"`                  // 模型ID
	TopP             float64 `json:"topP"`                // 采样范围
	Temperature      float64 `json:"temperature"`         // 多样性
	PresencePenalty  float64 `json:"presencePenalty"`     // 词汇控制
	FrequencyPenalty float64 `json:"frequencyPenalty"`    // 重复控制
	ResponseType     string  `json:"responseType"`        // 响应类型
	Model            string  `json:"model,omitempty"`     // 模型名称(部分场景使用)
	MaxTokens        int     `json:"maxTokens,omitempty"` // 最大token数(部分场景使用)
}

// Prompt 提示词配置
type Prompt struct {
	Role int    `json:"role"` // 角色
	Text string `json:"text"` // 提示词文本
}

// PromptConfig 提示词配置
type PromptConfig struct {
	Prompts    []Prompt `json:"prompts"`    // 提示词列表
	PromptVars []string `json:"promptVars"` // 提示词变量
}

// PrologueConfig 开场白配置
type PrologueConfig struct {
	Text            string   `json:"text"`            // 开场白文本
	OverrideForce   bool     `json:"overrideForce"`   // 是否强制覆盖
	SampleQuestions []string `json:"sampleQuestions"` // 示例问题列表
}

// BuiltinTools 内置工具
type BuiltinTools struct {
	Data2Chart bool `json:"data2chart"` // 数据转图表
	Python     bool `json:"python"`     // Python工具
	Datetime   bool `json:"datetime"`   // 日期时间工具
}

// ToolConfig 工具配置
type ToolConfig struct {
	Tools        []string     `json:"tools"`        // 工具列表
	BuiltinTools BuiltinTools `json:"builtinTools"` // 内置工具
}

// IntentConfig 意图配置
type IntentConfig struct {
	LowScoreNum      int      `json:"lowScoreNum"`      // 低分数量
	HighScoreNum     int      `json:"highScoreNum"`     // 高分数量
	LowScoreIntent   int      `json:"lowScoreIntent"`   // 低分意图
	HighScoreIntent  int      `json:"highScoreIntent"`  // 高分意图
	LowScoreMessages []string `json:"lowScoreMessages"` // 低分消息列表
	Enable           bool     `json:"enable"`           // 是否启用
}

// KnowledgeConfig 知识库配置
type KnowledgeConfig struct {
	Priority                bool         `json:"priority"`                          // 优先级
	Knowledges              []int        `json:"knowledges"`                        // 知识库ID列表
	MissPrompt              string       `json:"missPrompt"`                        // 缺失提示
	AutoReduceScope         bool         `json:"autoReduceScope"`                   // 自动减少范围
	Intent                  IntentConfig `json:"intent"`                            // 意图配置
	Enable                  bool         `json:"enable"`                            // 是否启用
	EmbeddingMatchThreshold int          `json:"embeddingMatchThreshold,omitempty"` // 嵌入匹配阈值
	EnableReranker          bool         `json:"enableReranker,omitempty"`          // 是否启用重排序
}

// NL2SQLConfig 自然语言转SQL配置
type NL2SQLConfig struct {
	Enable             bool        `json:"enable"`
	EnablePrologue     bool        `json:"enablePrologue"`
	Databases          []string    `json:"databases"`
	ResultDisplayMode  string      `json:"resultDisplayMode"`
	Model              ModelConfig `json:"model,omitempty"`
	AgentModel         ModelConfig `json:"agentModel,omitempty"`
	OverrideQuery      string      `json:"overrideQuery,omitempty"`
	EnableQueryRewrite bool        `json:"enableQueryRewrite,omitempty"`
}

// FileAnalyzerConfig 文件分析配置
type FileAnalyzerConfig struct {
	Enable bool `json:"enable"`
}

// PluginConfig 插件配置
type PluginConfig struct {
	Knowledge    KnowledgeConfig    `json:"knowledge"`
	NL2SQL       NL2SQLConfig       `json:"nl2sql"`
	FileAnalyzer FileAnalyzerConfig `json:"fileAnalyzer"`
}

// MCPServers MCP服务器配置
type MCPServers struct {
	Servers interface{} `json:"servers"`
}

// ChatConfig 多轮对话配置
type ChatConfig struct {
	Stream       bool           `json:"stream"`
	SingleRound  bool           `json:"singleRound"`
	Model        ModelConfig    `json:"model"`
	Prompt       PromptConfig   `json:"prompt"`
	Prologue     PrologueConfig `json:"prologue"`
	Tool         ToolConfig     `json:"tool"`
	MCPServers   MCPServers     `json:"mcpServers"`
	ContextLimit int            `json:"contextLimit"`
	Plugin       PluginConfig   `json:"plugin"`
}

// ExtractItem 抽取字段项
type ExtractItem struct {
	Name      string `json:"name"`
	FieldName string `json:"fieldName"`
	FieldType int    `json:"fieldType"`
	Format    string `json:"format,omitempty"`
	Default   string `json:"default,omitempty"`
	Desc      string `json:"desc,omitempty"`
	Strict    bool   `json:"strict,omitempty"`
	Multi     bool   `json:"multi,omitempty"`
	Required  bool   `json:"required"`
}

// ExtractConfig 实体抽取配置
type ExtractConfig struct {
	Model         ModelConfig   `json:"model"`
	RetryCount    int           `json:"retryCount"`
	FrontScript   string        `json:"frontScript"`
	BackendScript string        `json:"backendScript"`
	Items         []ExtractItem `json:"items"`
	BlockingWords []string      `json:"blockingWords"`
}

// ClassifyConfig 分类配置
type ClassifyConfig struct {
	Model         ModelConfig `json:"model"`
	RetryCount    int         `json:"retryCount"`
	Scope         string      `json:"scope"`
	FrontScript   string      `json:"frontScript,omitempty"`
	BackendScript string      `json:"backendScript,omitempty"`
	Items         []string    `json:"items"`
}

// FormCollectConfig 表单收集配置
type FormCollectConfig struct {
	Model        ModelConfig    `json:"model"`
	Characters   string         `json:"characters"`
	Name         string         `json:"name"`
	RetryCount   int            `json:"retryCount"`
	Items        []string       `json:"items"`
	Prologue     PrologueConfig `json:"prologue"`
	Tool         ToolConfig     `json:"tool"`
	ContextLimit int            `json:"contextLimit"`
	Plugin       PluginConfig   `json:"plugin"`
}

// IntentionConfig 意图识别配置
type IntentionConfig struct {
	Model        ModelConfig    `json:"model"`
	Characters   string         `json:"characters"`
	RetryCount   int            `json:"retryCount"`
	Items        []string       `json:"items"`
	Prologue     PrologueConfig `json:"prologue"`
	ContextLimit int            `json:"contextLimit"`
}

// IntentRouterConfig 意图路由配置
type IntentRouterConfig struct {
	AgentModel  ModelConfig    `json:"agentModel"`
	EmbedModel  ModelConfig    `json:"embedModel"`
	Prologue    PrologueConfig `json:"prologue"`
	Prompt      PromptConfig   `json:"prompt"`
	Items       []string       `json:"items"`
	SingleRound bool           `json:"singleRound"`
	Stream      bool           `json:"stream,omitempty"`
}

// RouterStep 路由步骤配置
type RouterStep struct {
	// 此结构可根据实际需求扩展
}

// FlowConfig 流程配置
type FlowConfig struct {
	Agents     []int      `json:"agents"`     // Agent ID列表
	Prompt     string     `json:"prompt"`     // 提示词
	Mode       string     `json:"mode"`       // 模式
	RouterStep RouterStep `json:"routerStep"` // 路由步骤
}

// SettingConfig 配置信息
type SettingConfig struct {
	ID           int                 `json:"id"`                     // 配置ID
	Name         string              `json:"name"`                   // 配置名称
	Chat         *ChatConfig         `json:"chat,omitempty"`         // 多轮对话配置
	Flow         *FlowConfig         `json:"flow,omitempty"`         // 流程配置
	Extract      *ExtractConfig      `json:"extract,omitempty"`      // 实体抽取配置
	Classify     *ClassifyConfig     `json:"classify,omitempty"`     // 分类配置
	FormCollect  *FormCollectConfig  `json:"formCollect,omitempty"`  // 表单收集配置
	Intention    *IntentionConfig    `json:"intention,omitempty"`    // 意图识别配置
	IntentRouter *IntentRouterConfig `json:"intentRouter,omitempty"` // 意图路由配置
	AgentID      int                 `json:"agentId"`                // Agent ID
}

// SaveConfigRequest 保存配置请求
// 复用SettingConfig结构
// 注意: 创建时ID可为0,更新时必须提供ID; AgentID和Name为必填字段
type SaveConfigRequest = SettingConfig

// NewDefaultSettingConfig 创建默认配置
func NewDefaultSettingConfig(agentID int, name string) *SettingConfig {
	return &SettingConfig{
		ID:      0,
		Name:    name,
		AgentID: agentID,
		Chat: &ChatConfig{
			Stream:      true,
			SingleRound: false,
			Model: ModelConfig{
				ID:               2,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				ResponseType:     "text",
			},
			Prompt: PromptConfig{
				Prompts: []Prompt{
					{Role: 0, Text: ""},
					{Role: 1, Text: ""},
				},
				PromptVars: []string{},
			},
			Prologue: PrologueConfig{
				Text:            "",
				OverrideForce:   false,
				SampleQuestions: []string{},
			},
			Tool: ToolConfig{
				Tools: []string{},
				BuiltinTools: BuiltinTools{
					Data2Chart: false,
					Python:     false,
					Datetime:   false,
				},
			},
			ContextLimit: 16,
			Plugin: PluginConfig{
				Knowledge: KnowledgeConfig{
					Priority:        false,
					Knowledges:      []int{},
					MissPrompt:      "",
					AutoReduceScope: false,
					Intent: IntentConfig{
						LowScoreNum:      50,
						HighScoreNum:     100,
						LowScoreIntent:   0,
						HighScoreIntent:  0,
						LowScoreMessages: []string{},
					},
					Enable: false,
				},
				NL2SQL: NL2SQLConfig{
					Enable:            false,
					EnablePrologue:    false,
					Databases:         []string{},
					ResultDisplayMode: "queryPanel",
				},
				FileAnalyzer: FileAnalyzerConfig{
					Enable: false,
				},
			},
		},
		Flow: &FlowConfig{
			Agents:     []int{},
			Prompt:     "",
			Mode:       "auto",
			RouterStep: RouterStep{},
		},
		Extract: &ExtractConfig{
			Model: ModelConfig{
				ID:               0,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Model:            "",
				ResponseType:     "text",
			},
			RetryCount:    3,
			FrontScript:   "",
			BackendScript: "",
			Items:         []ExtractItem{},
			BlockingWords: []string{},
		},
		Classify: &ClassifyConfig{
			Model: ModelConfig{
				ID:               0,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Model:            "",
				ResponseType:     "text",
			},
			RetryCount:    3,
			Scope:         "",
			FrontScript:   "",
			BackendScript: "",
			Items:         []string{},
		},
		FormCollect: &FormCollectConfig{
			Model: ModelConfig{
				ID:               0,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Model:            "",
				ResponseType:     "text",
			},
			Characters: "",
			Name:       "",
			RetryCount: 3,
			Items:      []string{},
			Prologue: PrologueConfig{
				Text:            "",
				OverrideForce:   false,
				SampleQuestions: []string{},
			},
			Tool: ToolConfig{
				Tools: []string{},
				BuiltinTools: BuiltinTools{
					Data2Chart: false,
					Python:     false,
					Datetime:   false,
				},
			},
			ContextLimit: 16,
			Plugin: PluginConfig{
				Knowledge: KnowledgeConfig{
					Priority:        false,
					Knowledges:      []int{},
					MissPrompt:      "",
					AutoReduceScope: false,
					Intent: IntentConfig{
						LowScoreNum:      50,
						HighScoreNum:     100,
						LowScoreIntent:   0,
						HighScoreIntent:  0,
						LowScoreMessages: []string{},
					},
				},
			},
		},
		Intention: &IntentionConfig{
			Model: ModelConfig{
				ID:               0,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Model:            "",
				ResponseType:     "text",
			},
			Characters: "",
			RetryCount: 3,
			Items:      []string{},
			Prologue: PrologueConfig{
				Text:            "",
				OverrideForce:   false,
				SampleQuestions: []string{},
			},
			ContextLimit: 16,
		},
		IntentRouter: &IntentRouterConfig{
			AgentModel: ModelConfig{
				ID:               0,
				TopP:             1,
				Temperature:      1,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				Model:            "",
				ResponseType:     "text",
			},
			EmbedModel: ModelConfig{
				ID:               0,
				TopP:             0,
				Temperature:      0,
				FrequencyPenalty: 0,
				PresencePenalty:  0,
				MaxTokens:        0,
				Model:            "",
				ResponseType:     "",
			},
			Prologue: PrologueConfig{
				Text:            "",
				OverrideForce:   false,
				SampleQuestions: []string{},
			},
			Prompt: PromptConfig{
				Prompts: []Prompt{
					{Role: 0, Text: ""},
					{Role: 1, Text: ""},
				},
				PromptVars: []string{},
			},
			Items:       []string{},
			SingleRound: true,
		},
	}
}

// PromptVar 提示词变量
type PromptVar struct {
	ID         string      `json:"id"`         // 变量ID
	Type       string      `json:"type"`       // 变量类型
	Name       string      `json:"name"`       // 变量名称
	Prompt     string      `json:"prompt"`     // 提示词
	EnumValues interface{} `json:"enumValues"` // 枚举值
	Value      string      `json:"value"`      // 变量值
}

// CreateSessionRequest 创建会话请求
type CreateSessionRequest struct {
	SettingID  int         `json:"settingId"`  // 必填: 配置ID
	PromptVars []PromptVar `json:"promptVars"` // 必填: 提示词变量列表 (可为空数组)
}

// NL2SQLPrologue NL2SQL开场白
type NL2SQLPrologue struct {
	Info      string      `json:"info"`
	Prologues interface{} `json:"prologues"`
}

// SessionResponse 会话响应
type SessionResponse struct {
	SettingID      int            `json:"settingId"`
	AgentID        int            `json:"agentId"`
	ID             int            `json:"id"`
	PromptVars     []string       `json:"promptVars"`
	Title          string         `json:"title"`
	LastMessage    string         `json:"lastMessage"`
	Prologue       string         `json:"prologue"`
	NL2SQLPrologue NL2SQLPrologue `json:"nl2sqlPrologue"`
	ContextFiles   interface{}    `json:"contextFiles"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	SettingID      int            `json:"settingId"`
	AgentID        int            `json:"agentId"`
	ID             int            `json:"id"`
	PromptVars     []string       `json:"promptVars"`
	Title          string         `json:"title"`
	LastMessage    string         `json:"lastMessage"`
	Prologue       string         `json:"prologue"`
	NL2SQLPrologue NL2SQLPrologue `json:"nl2sqlPrologue"`
	ContextFiles   []string       `json:"contextFiles"`
}

// AsyncChatRequest 异步对话请求
type AsyncChatRequest struct {
	SessionID int      `json:"sessionId"` // 会话ID
	Content   string   `json:"content"`   // 消息内容
	RequestID string   `json:"requestId"` // 请求ID (UUID)
	Files     []string `json:"files"`     // 文件列表
}

// SessionRecord 对话记录
type SessionRecord struct {
	ID                 int                    `json:"id"`                 // 记录ID
	AgentID            int                    `json:"agentId"`            // AgentID
	Sender             int                    `json:"sender"`             // 发送者: 1=用户, 2=助手
	Index              int                    `json:"index"`              // 索引
	Message            string                 `json:"message"`            // 消息内容
	Question           string                 `json:"question"`           // 问题
	Role               string                 `json:"role"`               // 角色
	Content            string                 `json:"content"`            // 内容
	MultiContent       []MultiContentItem     `json:"multi_content"`      // 多内容
	Error              bool                   `json:"error"`              // 是否错误
	ErrorMessage       string                 `json:"errorMessage"`       // 错误消息
	SessionID          int                    `json:"sessionId"`          // 会话ID
	Usage              UsageInfo              `json:"usage"`              // 使用信息
	RequestID          string                 `json:"requestId"`          // 请求ID
	UserData           map[string]interface{} `json:"userData"`           // 用户数据
	Pending            bool                   `json:"pending"`            // 是否等待中
	Plugins            PluginsInfo            `json:"plugins"`            // 插件信息
	Tools              interface{}            `json:"tools"`              // 工具
	Extra              *ExtraInfo             `json:"extra"`              // 额外信息
	SuggestedQuestions []string               `json:"suggestedQuestions"` // 建议问题
	FinishReason       string                 `json:"finishReason"`       // 完成原因
}

// UsageInfo 使用信息
type UsageInfo struct {
	PromptTokens     int `json:"promptTokens"`     // 提示词token数
	CompletionTokens int `json:"completionTokens"` // 完成token数
	TotalTokens      int `json:"totalTokens"`      // 总token数
	Duration         int `json:"duration"`         // 持续时间(秒)
}

// PluginsInfo 插件信息
type PluginsInfo struct {
	FileAnalysis []interface{} `json:"file_analysis"` // 文件分析
	Knowledge    []interface{} `json:"knowledge"`     // 知识库
	NL2SQL       []interface{} `json:"nl2sql"`        // NL2SQL
	OnlineSearch []interface{} `json:"onlineSearch"`  // 在线搜索
}

// ExtraInfo 额外信息
type ExtraInfo struct {
	ChartData interface{}   `json:"chartData"` // 图表数据
	Tools     []interface{} `json:"tools"`     // 工具列表
}

// MultiContentItem 多内容项
type MultiContentItem struct {
	Text string `json:"text"` // 文本内容
	Type string `json:"type"` // 类型 (如: "question")
}

// StructData 结构化数据
type StructData struct {
	JSON    interface{} `json:"json"`    // JSON 数据
	Code    interface{} `json:"code"`    // 代码数据
	MarkMap interface{} `json:"markMap"` // 思维导图数据
}

// StreamEvent SSE流事件
// swagger:model
type StreamEvent struct {
	EventType          string                 `json:"eventType,omitempty"` // 事件类型 (resp_splash, resp_increment, resp_finish)
	Message            string                 `json:"message"`             // 消息内容
	Question           string                 `json:"question"`            // 问题
	SuggestedQuestions []string               `json:"suggestedQuestions"`  // 建议问题
	ShadowMessage      string                 `json:"shadowMessage"`       // 影子消息
	RequestID          string                 `json:"requestId"`           // 请求ID
	Index              int                    `json:"index"`               // 索引
	Error              bool                   `json:"error"`               // 是否错误
	Tools              []interface{}          `json:"tools"`               // 工具列表
	Plugins            map[string]interface{} `json:"plugins"`             // 插件信息
	Pending            bool                   `json:"pending"`             // 是否等待中
	AgentID            int                    `json:"agentId"`             // Agent ID
	SessionID          int                    `json:"sessionId"`           // 会话ID
	SettingID          int                    `json:"settingId"`           // 配置ID
	Usage              UsageInfo              `json:"usage"`               // 使用信息
	ChartData          interface{}            `json:"chartData"`           // 图表数据
	StructData         StructData             `json:"structData"`          // 结构化数据
	FinishReason       string                 `json:"finishReason"`        // 完成原因
}

// SessionRecordsResponse 会话记录响应
type SessionRecordsResponse struct {
	Type    int             `json:"type"`    // 类型
	Records []SessionRecord `json:"records"` // 记录列表
}

// Service Agent服务接口
type Service interface {
	// ==================== Agent管理 ====================
	// Agent列表(分页)
	// API: POST /agent/listByPage
	ListAgentsByPage(ctx context.Context, current, size int) (*AgentListResponse, error)

	// Agent详情
	// API: POST /agent/detail
	GetAgentDetail(ctx context.Context, agentID string) (*AgentDetailResponse, error)

	// 创建Agent
	// API: POST /agent/create
	CreateAgent(ctx context.Context, req *CreateAgentRequest) (int, error)

	// 更新Agent
	// API: POST /agent/update
	UpdateAgent(ctx context.Context, req *UpdateAgentRequest) (int, error)

	// 删除Agent
	// API: POST /agent/delete
	DeleteAgent(ctx context.Context, agentID int) error

	// ==================== 配置管理 ====================
	// 配置列表查询
	// API: POST /agent/setting/list
	ListConfigs(ctx context.Context, agentID int) ([]SettingConfig, error)

	// 新建或保存配置
	// API: POST /agent/setting/save
	SaveConfig(ctx context.Context, req *SaveConfigRequest) (int, error)

	// 配置删除
	// API: POST /agent/setting/delete
	DeleteConfig(ctx context.Context, configID int) error

	// ==================== 会话管理 ====================
	// 会话列表
	// API: POST /agent/session/list
	ListSessions(ctx context.Context, settingID int, idLimits []int) ([]SessionInfo, error)

	// 创建会话
	// API: POST /agent/session/create
	CreateSession(ctx context.Context, req *CreateSessionRequest) (*SessionResponse, error)

	// 删除会话
	// API: POST /agent/session/delete
	DeleteSession(ctx context.Context, sessionID int) error

	// ==================== 对话管理 ====================
	// 对话记录
	// API: POST /agent/session/record/list
	GetSessionRecords(ctx context.Context, sessionID int) (*SessionRecordsResponse, error)

	// SSE流式对话
	// API: POST /agent/chatAsync (SSE)
	ChatAsync(ctx context.Context, req *AsyncChatRequest, eventChan chan<- StreamEvent) error

	// 赞
	// API: POST /blade-flowy/agent/session/record/pros
	LikeMessage(ctx context.Context, messageID string) error

	// 踩
	// API: POST /blade-flowy/agent/session/record/cons
	DislikeMessage(ctx context.Context, messageID string) error
}

// ServiceImpl Agent服务实现
type ServiceImpl struct {
	client client.HTTPClient
}

// NewService 创建Agent服务
func NewService(client client.HTTPClient) Service {
	return &ServiceImpl{
		client: client,
	}
}

// ==================== Agent管理 ====================

// ListAgentsByPage 获取Agent列表(分页)
// API: POST /agent/listByPage
func (s *ServiceImpl) ListAgentsByPage(ctx context.Context, current, size int) (*AgentListResponse, error) {
	req := map[string]interface{}{
		"current": current,
		"size":    size,
	}

	resp, err := s.client.Post(ctx, "/agent/listByPage", req)
	if err != nil {
		return nil, err
	}

	var result AgentListResponse
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetAgentDetail 获取Agent详情
// API: POST /agent/detail
func (s *ServiceImpl) GetAgentDetail(ctx context.Context, agentID string) (*AgentDetailResponse, error) {
	if agentID == "" {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "agent ID is required")
	}

	req := map[string]interface{}{
		"id": agentID,
	}

	resp, err := s.client.Post(ctx, "/agent/detail", req)
	if err != nil {
		return nil, err
	}

	var result AgentDetailResponse
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateAgent 创建Agent
// API: POST /agent/create
func (s *ServiceImpl) CreateAgent(ctx context.Context, req *CreateAgentRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "agent name is required")
	}

	resp, err := s.client.Post(ctx, "/agent/create", req)
	if err != nil {
		return 0, err
	}

	var agentID int
	if err := s.parseResponseData(resp, &agentID); err != nil {
		return 0, err
	}

	return agentID, nil
}

// UpdateAgent 更新Agent
// API: POST /agent/update
func (s *ServiceImpl) UpdateAgent(ctx context.Context, req *UpdateAgentRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.ID == 0 {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "agent ID is required")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "agent name is required")
	}

	resp, err := s.client.Post(ctx, "/agent/update", req)
	if err != nil {
		return 0, err
	}

	var agentID int
	if err := s.parseResponseData(resp, &agentID); err != nil {
		return 0, err
	}

	return agentID, nil
}

// DeleteAgent 删除Agent
// API: POST /agent/delete
func (s *ServiceImpl) DeleteAgent(ctx context.Context, agentID int) error {
	if agentID == 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "agent ID is required")
	}

	req := map[string]interface{}{
		"id": agentID,
	}

	_, err := s.client.Post(ctx, "/agent/delete", req)
	return err
}

// ==================== 配置管理 ====================

// ListConfigs 获取配置列表
func (s *ServiceImpl) ListConfigs(ctx context.Context, agentID int) ([]SettingConfig, error) {
	if agentID == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "agent ID is required")
	}

	req := map[string]interface{}{
		"agentId": agentID,
	}

	resp, err := s.client.Post(ctx, "/agent/setting/list", req)
	if err != nil {
		return nil, err
	}

	var configs []SettingConfig
	if err := s.parseResponseData(resp, &configs); err != nil {
		return nil, err
	}

	return configs, nil
}

// SaveConfig 保存配置
func (s *ServiceImpl) SaveConfig(ctx context.Context, req *SaveConfigRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.AgentID == 0 {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "agent ID is required")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "config name is required")
	}

	resp, err := s.client.Post(ctx, "/agent/setting/save", req)
	if err != nil {
		return 0, err
	}

	var configID int
	if err := s.parseResponseData(resp, &configID); err != nil {
		return 0, err
	}

	return configID, nil
}

// DeleteConfig 删除配置
func (s *ServiceImpl) DeleteConfig(ctx context.Context, configID int) error {
	if configID == 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "config ID is required")
	}

	req := map[string]interface{}{
		"id": configID,
	}

	_, err := s.client.Post(ctx, "/agent/setting/delete", req)
	return err
}

// ==================== 会话管理 ====================

// ListSessions 获取会话列表
func (s *ServiceImpl) ListSessions(ctx context.Context, settingID int, idLimits []int) ([]SessionInfo, error) {
	if settingID == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "setting ID is required")
	}

	req := map[string]interface{}{
		"settingId": settingID,
		"idLimits":  idLimits,
	}

	resp, err := s.client.Post(ctx, "/agent/session/list", req)
	if err != nil {
		return nil, err
	}

	var sessions []SessionInfo
	if err := s.parseResponseData(resp, &sessions); err != nil {
		return nil, err
	}

	return sessions, nil
}

// CreateSession 创建会话
func (s *ServiceImpl) CreateSession(ctx context.Context, req *CreateSessionRequest) (*SessionResponse, error) {
	if req == nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.SettingID == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "setting ID is required")
	}

	resp, err := s.client.Post(ctx, "/agent/session/create", req)
	if err != nil {
		return nil, err
	}

	var session SessionResponse
	if err := s.parseResponseData(resp, &session); err != nil {
		return nil, err
	}

	return &session, nil
}

// DeleteSession 删除会话
func (s *ServiceImpl) DeleteSession(ctx context.Context, sessionID int) error {
	if sessionID == 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "session ID is required")
	}

	req := map[string]interface{}{
		"id": sessionID,
	}

	_, err := s.client.Post(ctx, "/agent/session/delete", req)
	return err
}

// ==================== 对话管理 ====================

// GetSessionRecords 获取对话记录
func (s *ServiceImpl) GetSessionRecords(ctx context.Context, sessionID int) (*SessionRecordsResponse, error) {
	if sessionID == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "session ID is required")
	}

	req := map[string]interface{}{
		"sessionId": sessionID,
	}

	resp, err := s.client.Post(ctx, "/agent/session/record/list", req)
	if err != nil {
		return nil, err
	}

	var result SessionRecordsResponse
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ChatAsync SSE流式对话
func (s *ServiceImpl) ChatAsync(ctx context.Context, req *AsyncChatRequest, eventChan chan<- StreamEvent) error {
	if req == nil {
		return errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.SessionID == 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "session ID is required")
	}

	if req.Content == "" {
		return errors.New(errors.ErrCodeInvalidRequest, "content is required")
	}

	if req.Files == nil {
		req.Files = []string{}
	}

	// 获取 SSE 流
	stream, err := s.client.PostSSE(ctx, "/agent/chatAsync", req)
	if err != nil {
		return err
	}
	defer stream.Close()

	// 读取并解析 SSE 事件
	return s.parseSSEStream(ctx, stream, eventChan)
}

// parseSSEStream 解析 SSE 流
func (s *ServiceImpl) parseSSEStream(ctx context.Context, stream io.ReadCloser, eventChan chan<- StreamEvent) error {
	buf := make([]byte, 4096)
	var line string
	var currentEventType string

	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}

		n, err := stream.Read(buf)
		if n > 0 {
			data := string(buf[:n])
			lines := strings.Split(line+data, "\n")

			for i := 0; i < len(lines)-1; i++ {
				trimmedLine := strings.TrimSpace(lines[i])

				// 检查是否是 event: 行
				if strings.HasPrefix(trimmedLine, "event:") {
					currentEventType = strings.TrimSpace(strings.TrimPrefix(trimmedLine, "event:"))
					continue
				}

				// 处理 data: 行
				if err := s.processSSELine(lines[i], currentEventType, eventChan); err != nil {
					return err
				}

				// 空行表示事件结束，重置事件类型
				if trimmedLine == "" {
					currentEventType = ""
				}
			}

			line = lines[len(lines)-1]
		}

		if err != nil {
			if err == io.EOF {
				break
			}
			return err
		}
	}

	return nil
}

// processSSELine 处理单行 SSE 数据
func (s *ServiceImpl) processSSELine(line string, eventType string, eventChan chan<- StreamEvent) error {
	line = strings.TrimSpace(line)

	if line == "" || !strings.HasPrefix(line, "data:") {
		return nil
	}

	// 提取 data: 后面的内容
	data := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if data == "" || data == "[DONE]" {
		// 发送完成事件
		eventChan <- StreamEvent{
			EventType:    eventType,
			FinishReason: "stop",
			Pending:      false,
		}
		return nil
	}

	// 解析 JSON 数据
	var event StreamEvent
	if err := json.Unmarshal([]byte(data), &event); err != nil {
		// 如果解析失败，将原始数据作为消息发送
		eventChan <- StreamEvent{
			EventType: eventType,
			Message:   data,
			Error:     true,
		}
		return nil
	}

	// 设置事件类型
	event.EventType = eventType
	eventChan <- event
	return nil
}

// LikeMessage 点赞对话记录
func (s *ServiceImpl) LikeMessage(ctx context.Context, messageID string) error {
	if messageID == "" {
		return errors.New(errors.ErrCodeInvalidRequest, "message ID is required")
	}

	req := map[string]interface{}{
		"id": messageID,
	}

	_, err := s.client.Post(ctx, "/blade-flowy/agent/session/record/pros", req)
	return err
}

// DislikeMessage 踩对话记录
func (s *ServiceImpl) DislikeMessage(ctx context.Context, messageID string) error {
	if messageID == "" {
		return errors.New(errors.ErrCodeInvalidRequest, "message ID is required")
	}

	req := map[string]interface{}{
		"id": messageID,
	}

	_, err := s.client.Post(ctx, "/blade-flowy/agent/session/record/cons", req)
	return err
}

// parseResponseData 解析响应数据
func (s *ServiceImpl) parseResponseData(resp *models.BaseResponse, target interface{}) error {
	if resp == nil {
		return errors.New(errors.ErrCodeInternalError, "response is nil")
	}

	if !resp.Success {
		return errors.New(errors.ErrCodeInternalError, resp.Message)
	}

	if resp.Data == nil {
		return errors.New(errors.ErrCodeInternalError, "response data is nil")
	}

	// 将数据转换为JSON然后反序列化到目标结构
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return errors.New(errors.ErrCodeInternalError, "failed to marshal response data").WithDetails(err.Error())
	}

	if err := json.Unmarshal(data, target); err != nil {
		return errors.New(errors.ErrCodeInternalError, "failed to unmarshal response data").WithDetails(err.Error())
	}

	return nil
}
