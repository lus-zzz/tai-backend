package langchaingo

import (
	"context"

	"chat-backend/services/interfaces"
	"chat-backend/utils"
	modelSvc "flowy-sdk/services/model"
)

// LangchaingoModelService 基于 langchaingo 的模型服务实现
type LangchaingoModelService struct {
	config *LangchaingoConfig
}

// NewLangchaingoModelService 创建 Langchaingo 模型服务
func NewLangchaingoModelService(config *LangchaingoConfig) interfaces.ModelServiceInterface {
	return &LangchaingoModelService{
		config: config,
	}
}

// ListSupportedChatModels 获取支持的聊天模型列表
func (s *LangchaingoModelService) ListSupportedChatModels(ctx context.Context) ([]modelSvc.SupportedChatModel, error) {
	// TODO: 实现 langchaingo 支持的聊天模型列表
	// 这里需要返回 langchaingo 框架支持的聊天模型
	// 暂时返回空列表
	var models []modelSvc.SupportedChatModel
	
	// 示例：OpenAI 模型
	models = append(models, modelSvc.SupportedChatModel{
		Name:     "gpt-3.5-turbo",
		Identify: "openai-chat",
		LLMProperty: modelSvc.LLMProperty{
			Model:           "gpt-3.5-turbo",
			MaxToken:        4096,
			ContextLength:   4096,
			TokenProportion: 1,
			Stream:          true,
			SupportTool:     true,
			FixQwqThink:    false,
		},
	})
	
	models = append(models, modelSvc.SupportedChatModel{
		Name:     "gpt-4",
		Identify: "openai-chat",
		LLMProperty: modelSvc.LLMProperty{
			Model:           "gpt-4",
			MaxToken:        8192,
			ContextLength:   8192,
			TokenProportion: 1,
			Stream:          true,
			SupportTool:     true,
			FixQwqThink:    false,
		},
	})

	utils.InfoWith("获取支持的聊天模型列表", "count", len(models))
	return models, nil
}

// ListSupportedVectorModels 获取支持的向量模型列表
func (s *LangchaingoModelService) ListSupportedVectorModels(ctx context.Context) ([]modelSvc.SupportedVectorModel, error) {
	// TODO: 实现 langchaingo 支持的向量模型列表
	// 这里需要返回 langchaingo 框架支持的向量模型
	// 暂时返回空列表
	var models []modelSvc.SupportedVectorModel
	
	// 示例：bge-m3 模型
	models = append(models, modelSvc.SupportedVectorModel{
		Name:     "bge-m3:latest",
		Identify: "ollama-embedding",
		LLMProperty: modelSvc.EmbeddingProperty{
			Model:              "bge-m3:latest",
			EmbeddingMaxLength:  8192,
			EmbeddingDimension: 1024,
			BatchLimit:         100,
		},
	})

	utils.InfoWith("获取支持的向量模型列表", "count", len(models))
	return models, nil
}

// ListAvailableAllModels 获取全部可用模型列表
func (s *LangchaingoModelService) ListAvailableAllModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	// TODO: 实现 langchaingo 可用模型列表
	// 这里需要检查配置中模型的可用性
	// 暂时返回空列表
	var models []modelSvc.ModelInfo

	// 聊天模型
	chatModels, err := s.ListSupportedChatModels(ctx)
	if err != nil {
		return nil, err
	}

	for i, model := range chatModels {
		models = append(models, modelSvc.ModelInfo{
			ID:       i + 1, // 生成ID
			Name:     model.Name,
			Symbol:   model.Identify,
			Endpoint: s.config.OpenAI.BaseURL,
			Enable:   true,
			Type:     0, // 聊天模型类型为0
			Role:     "chat",
		})
	}

	// 向量模型
	vectorModels, err := s.ListSupportedVectorModels(ctx)
	if err != nil {
		return nil, err
	}

	for i, model := range vectorModels {
		models = append(models, modelSvc.ModelInfo{
			ID:       i + 100, // 生成ID，避免与聊天模型冲突
			Name:     model.Name,
			Symbol:   model.Identify,
			Endpoint: s.config.Ollama.BaseURL,
			Enable:   true,
			Type:     1, // 向量模型类型为1
			Role:     "embedding",
		})
	}

	utils.InfoWith("获取全部可用模型列表", "count", len(models))
	return models, nil
}

// ListAvailableChatModels 获取可用的对话模型列表
func (s *LangchaingoModelService) ListAvailableChatModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	// TODO: 实现 langchaingo 可用聊天模型列表
	// 暂时返回支持的聊天模型
	chatModels, err := s.ListSupportedChatModels(ctx)
	if err != nil {
		return nil, err
	}

	var models []modelSvc.ModelInfo
	for i, model := range chatModels {
		models = append(models, modelSvc.ModelInfo{
			ID:       i + 1,
			Name:     model.Name,
			Symbol:   model.Identify,
			Endpoint: s.config.OpenAI.BaseURL,
			Enable:   true,
			Type:     0, // 聊天模型类型为0
			Role:     "chat",
		})
	}

	utils.InfoWith("获取可用聊天模型列表", "count", len(models))
	return models, nil
}

// ListAvailableVectorModels 获取可用的向量模型列表
func (s *LangchaingoModelService) ListAvailableVectorModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	// TODO: 实现 langchaingo 可用向量模型列表
	// 暂时返回支持的向量模型
	vectorModels, err := s.ListSupportedVectorModels(ctx)
	if err != nil {
		return nil, err
	}

	var models []modelSvc.ModelInfo
	for i, model := range vectorModels {
		models = append(models, modelSvc.ModelInfo{
			ID:       i + 100,
			Name:     model.Name,
			Symbol:   model.Identify,
			Endpoint: s.config.Ollama.BaseURL,
			Enable:   true,
			Type:     1, // 向量模型类型为1
			Role:     "embedding",
		})
	}

	utils.InfoWith("获取可用向量模型列表", "count", len(models))
	return models, nil
}

// SaveModel 添加或修改模型
func (s *LangchaingoModelService) SaveModel(ctx context.Context, req *modelSvc.ModelSaveRequest) (int, error) {
	// TODO: 实现 langchaingo 模型保存
	// langchaingo 通常不需要动态添加模型，因为模型是预配置的
	// 暂时返回模拟ID
	modelID := 1
	
	utils.InfoWith("保存模型", "model_id", modelID, "model_name", req.Name)
	return modelID, nil
}

// DeleteModel 删除模型
func (s *LangchaingoModelService) DeleteModel(ctx context.Context, id int) error {
	// TODO: 实现 langchaingo 模型删除
	// langchaingo 通常不需要动态删除模型
	utils.InfoWith("删除模型", "model_id", id)
	return nil
}

// SetModelStatus 设置模型启用状态
func (s *LangchaingoModelService) SetModelStatus(ctx context.Context, id int, enable bool) error {
	// TODO: 实现 langchaingo 模型状态设置
	// langchaingo 通常不需要动态设置模型状态
	status := "禁用"
	if enable {
		status = "启用"
	}
	utils.InfoWith("设置模型状态", "model_id", id, "status", status)
	return nil
}
