package model

import (
	"context"
	"encoding/json"

	"flowy-sdk/pkg/client"
	"flowy-sdk/pkg/errors"
	"flowy-sdk/pkg/models"
)

// ============================================================================
// 数据类型定义
// ============================================================================

// CredentialSchemaItem 模型凭证配置项
type CredentialSchemaItem struct {
	Name  string `json:"name"`  // 参数名称
	Value string `json:"value"` // 参数值
	Desc  string `json:"desc"`  // 参数描述
}

// LLMProperty 聊天模型属性
type LLMProperty struct {
	MaxToken        int    `json:"maxToken"`        // 最大token数
	ContextLength   int    `json:"contextLength"`   // 上下文长度
	TokenProportion int    `json:"tokenProportion"` // token比例
	Stream          bool   `json:"stream"`          // 是否支持流式输出
	Model           string `json:"model"`           // 模型名称
	SupportTool     bool   `json:"supportTool"`     // 是否支持工具调用
	FixQwqThink     bool   `json:"fixQwqThink"`     // 是否修复Qwq思考模式
}

// EmbeddingProperty 向量模型属性
type EmbeddingProperty struct {
	EmbeddingMaxLength int    `json:"embeddingMaxLength"` // 最大嵌入长度
	EmbeddingDimension int    `json:"embeddingDimension"` // 嵌入维度
	BatchLimit         int    `json:"batchLimit"`         // 批处理限制
	Model              string `json:"model"`              // 模型名称
}

// SupportedChatModel 支持的聊天模型
type SupportedChatModel struct {
	Name             string                 `json:"name"`             // 模型名称：OpenAI, DeepSeek, Ollama
	Identify         string                 `json:"identify"`         // 模型标识：openai-chat, deepseek-chat, ollama-chat
	CredentialSchema []CredentialSchemaItem `json:"credentialSchema"` // 凭证配置
	LLMProperty      LLMProperty            `json:"llmProperty"`      // LLM属性
	OncurrencyLimit  int                    `json:"OncurrencyLimit"`  // 并发限制
}

// SupportedVectorModel 支持的向量模型
type SupportedVectorModel struct {
	Name             string                 `json:"name"`             // 模型名称：Ollama
	Identify         string                 `json:"identify"`         // 模型标识：ollama-embedding
	CredentialSchema []CredentialSchemaItem `json:"credentialSchema"` // 凭证配置
	LLMProperty      EmbeddingProperty      `json:"llmProperty"`      // 嵌入属性
}

// ModelInfo 模型信息
type ModelInfo struct {
	ID          int                    `json:"id"`          // 模型ID
	Name        string                 `json:"name"`        // 模型名称
	Symbol      string                 `json:"symbol"`      // 模型符号标识
	Endpoint    string                 `json:"endpoint"`    // 模型服务端点
	Enable      bool                   `json:"enable"`      // 是否启用
	Credentials []CredentialSchemaItem `json:"credentials"` // 模型配置参数
	Type        int                    `json:"type"`        // 模型类型: 0=聊天模型, 1=嵌入模型
	Role        string                 `json:"role"`        // 角色信息
}

// ModelSaveRequest 添加/修改模型请求
// API: POST /model/save
// swagger:model
type ModelSaveRequest struct {
	// 模型ID: 0=新建, >0=修改
	// required: true
	ID int `json:"id"`
	// 模型名称
	// required: true
	Name string `json:"name"`
	// 模型类型: 0=聊天模型, 1=嵌入模型
	// required: true
	Type int `json:"type"`
	// 模型编码
	// required: true
	Symbol string `json:"symbol"`
	// 模型接入点
	// required: true
	Endpoint string `json:"endpoint"`
	// 是否启用
	// required: true
	Enable bool `json:"enable"`
	// 模型鉴权信息
	// required: true
	Credentials string `json:"credentials"`
}

// ModelDeleteRequest 删除模型请求
// API: POST /model/delete
type ModelDeleteRequest struct {
	ID int `json:"id"` // 模型ID
}

// ModelStatusRequest 切换模型状态请求
// API: POST /model/setModelStatus
type ModelStatusRequest struct {
	ID     int  `json:"id"`     // 模型ID
	Enable bool `json:"enable"` // 是否启用
}

// ============================================================================
// Service接口定义
// ============================================================================

// Service 模型管理服务接口
type Service interface {
	// 可添加的模型列表 (查看可以支持的模型格式列表：deepseek, openAI, ollama)
	// API: POST /model/supportedChatModels
	ListSupportedChatModels(ctx context.Context) ([]SupportedChatModel, error)

	// 可添加的向量模型列表
	// API: POST /model/supportedVectorModels
	ListSupportedVectorModels(ctx context.Context) ([]SupportedVectorModel, error)

	// 全部可用模型列表 (在模型管理已经注册的模型，当前全部可用模型列表)
	// API: POST /model/availableAllModels
	ListAvailableAllModels(ctx context.Context) ([]ModelInfo, error)

	// 添加/修改模型
	// API: POST /model/save
	SaveModel(ctx context.Context, req *ModelSaveRequest) (int, error)

	// 删除模型
	// API: POST /model/delete
	DeleteModel(ctx context.Context, id int) error

	// 可用的对话模型的列表 (已经在模型管理注册的用于对话的模型)
	// API: POST /model/availableChatModels
	ListAvailableChatModels(ctx context.Context) ([]ModelInfo, error)

	// 可用的向量模型列表 (已经在模型管理注册的向量模型列表)
	// API: POST /model/availableVectorModels
	ListAvailableVectorModels(ctx context.Context) ([]ModelInfo, error)

	// 切换模型启用/停用
	// API: POST /model/setModelStatus
	SetModelStatus(ctx context.Context, id int, enable bool) error
}

// ============================================================================
// Service实现
// ============================================================================

// ServiceImpl 模型管理服务实现
type ServiceImpl struct {
	client client.HTTPClient
}

// NewService 创建模型管理服务
func NewService(client client.HTTPClient) Service {
	return &ServiceImpl{
		client: client,
	}
}

// ListSupportedChatModels 可添加的模型列表
// API: POST /model/supportedChatModels
func (s *ServiceImpl) ListSupportedChatModels(ctx context.Context) ([]SupportedChatModel, error) {
	resp, err := s.client.Post(ctx, "/model/supportedChatModels", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var models []SupportedChatModel
	if err := s.parseResponseData(resp, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// ListSupportedVectorModels 可添加的向量模型列表
// API: POST /model/supportedVectorModels
func (s *ServiceImpl) ListSupportedVectorModels(ctx context.Context) ([]SupportedVectorModel, error) {
	resp, err := s.client.Post(ctx, "/model/supportedVectorModels", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var models []SupportedVectorModel
	if err := s.parseResponseData(resp, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// ListAvailableAllModels 全部可用模型列表
// API: POST /model/availableAllModels
func (s *ServiceImpl) ListAvailableAllModels(ctx context.Context) ([]ModelInfo, error) {
	resp, err := s.client.Post(ctx, "/model/availableAllModels", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	if err := s.parseResponseData(resp, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// SaveModel 添加/修改模型
// API: POST /model/save
func (s *ServiceImpl) SaveModel(ctx context.Context, req *ModelSaveRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "name is required")
	}

	if req.Symbol == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "symbol is required")
	}

	if req.Endpoint == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "endpoint is required")
	}

	if req.Credentials == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "credentials is required")
	}

	resp, err := s.client.Post(ctx, "/model/save", req)
	if err != nil {
		return 0, err
	}

	var modelID int
	if err := s.parseResponseData(resp, &modelID); err != nil {
		return 0, err
	}

	return modelID, nil
}

// DeleteModel 删除模型
// API: POST /model/delete
func (s *ServiceImpl) DeleteModel(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "invalid model id")
	}

	req := &ModelDeleteRequest{
		ID: id,
	}

	resp, err := s.client.Post(ctx, "/model/delete", req)
	if err != nil {
		return err
	}

	// 验证响应
	if !resp.Success {
		return errors.New(errors.ErrCodeInternalError, resp.Message)
	}

	return nil
}

// ListAvailableChatModels 可用的对话模型列表
// API: POST /model/availableChatModels
func (s *ServiceImpl) ListAvailableChatModels(ctx context.Context) ([]ModelInfo, error) {
	resp, err := s.client.Post(ctx, "/model/availableChatModels", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	if err := s.parseResponseData(resp, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// ListAvailableVectorModels 可用的向量模型列表
// API: POST /model/availableVectorModels
func (s *ServiceImpl) ListAvailableVectorModels(ctx context.Context) ([]ModelInfo, error) {
	resp, err := s.client.Post(ctx, "/model/availableVectorModels", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var models []ModelInfo
	if err := s.parseResponseData(resp, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// SetModelStatus 切换模型启用/停用
// API: POST /model/setModelStatus
func (s *ServiceImpl) SetModelStatus(ctx context.Context, id int, enable bool) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "invalid model id")
	}

	req := &ModelStatusRequest{
		ID:     id,
		Enable: enable,
	}

	resp, err := s.client.Post(ctx, "/model/setModelStatus", req)
	if err != nil {
		return err
	}

	// 验证响应
	if !resp.Success {
		return errors.New(errors.ErrCodeInternalError, resp.Message)
	}

	return nil
}

// ============================================================================
// 辅助方法
// ============================================================================

// parseResponseData 解析响应数据
func (s *ServiceImpl) parseResponseData(resp *models.BaseResponse, target interface{}) error {
	if resp == nil {
		return errors.New(errors.ErrCodeInternalError, "response is nil")
	}

	if !resp.Success {
		return errors.New(errors.ErrCodeInternalError, resp.Message)
	}

	if resp.Data == nil {
		return nil // 某些API返回null data是正常的
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
