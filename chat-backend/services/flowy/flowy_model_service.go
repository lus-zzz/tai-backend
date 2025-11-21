package flowy

import (
	"context"

	"chat-backend/models"
	"chat-backend/services/interfaces"
	"flowy-sdk"
	modelSvc "flowy-sdk/services/model"
)

// FlowyModelService 基于 flowy-sdk 的模型服务实现
type FlowyModelService struct {
	sdk *flowy.SDK
}

// NewFlowyModelService 创建 Flowy 模型服务
func NewFlowyModelService(sdk *flowy.SDK) interfaces.ModelServiceInterface {
	return &FlowyModelService{
		sdk: sdk,
	}
}

// ListSupportedChatModels 获取支持的聊天模型列表
func (s *FlowyModelService) ListSupportedChatModels(ctx context.Context) ([]models.SupportedChatModel, error) {
	flowyModels, err := s.sdk.Model.ListSupportedChatModels(ctx)
	if err != nil {
		return nil, err
	}

	// 转换 flowy 模型到本地模型
	var localModels []models.SupportedChatModel
	for _, model := range flowyModels {
		localModels = append(localModels, models.SupportedChatModel{
			Name:     model.Name,
			Identify:  model.Identify,
			LLMProperty: models.LLMProperty{
				Model:           model.LLMProperty.Model,
				MaxToken:        model.LLMProperty.MaxToken,
				ContextLength:   model.LLMProperty.ContextLength,
				TokenProportion: model.LLMProperty.TokenProportion,
				Stream:          model.LLMProperty.Stream,
				SupportTool:     model.LLMProperty.SupportTool,
				FixQwqThink:    model.LLMProperty.FixQwqThink,
			},
		})
	}

	return localModels, nil
}

// ListSupportedVectorModels 获取支持的向量模型列表
func (s *FlowyModelService) ListSupportedVectorModels(ctx context.Context) ([]models.SupportedVectorModel, error) {
	flowyModels, err := s.sdk.Model.ListSupportedVectorModels(ctx)
	if err != nil {
		return nil, err
	}

	// 转换 flowy 模型到本地模型
	var localModels []models.SupportedVectorModel
	for _, model := range flowyModels {
		localModels = append(localModels, models.SupportedVectorModel{
			Name:     model.Name,
			Identify:  model.Identify,
			LLMProperty: models.EmbeddingProperty{
				Model:              model.LLMProperty.Model,
				EmbeddingMaxLength:  model.LLMProperty.EmbeddingMaxLength,
				EmbeddingDimension: model.LLMProperty.EmbeddingDimension,
				BatchLimit:         model.LLMProperty.BatchLimit,
			},
		})
	}

	return localModels, nil
}

// ListAvailableAllModels 获取全部可用模型列表
func (s *FlowyModelService) ListAvailableAllModels(ctx context.Context) ([]models.ModelInfo, error) {
	flowyModels, err := s.sdk.Model.ListAvailableAllModels(ctx)
	if err != nil {
		return nil, err
	}

	// 转换 flowy 模型到本地模型
	var localModels []models.ModelInfo
	for _, model := range flowyModels {
		localModels = append(localModels, models.ModelInfo{
			ID:       model.ID,
			Name:     model.Name,
			Symbol:   model.Symbol,
			Endpoint: model.Endpoint,
			Enable:   model.Enable,
			Type:     model.Type,
			Role:     model.Role,
		})
	}

	return localModels, nil
}

// ListAvailableChatModels 获取可用的对话模型列表
func (s *FlowyModelService) ListAvailableChatModels(ctx context.Context) ([]models.ModelInfo, error) {
	flowyModels, err := s.sdk.Model.ListAvailableChatModels(ctx)
	if err != nil {
		return nil, err
	}

	// 转换 flowy 模型到本地模型
	var localModels []models.ModelInfo
	for _, model := range flowyModels {
		localModels = append(localModels, models.ModelInfo{
			ID:       model.ID,
			Name:     model.Name,
			Symbol:   model.Symbol,
			Endpoint: model.Endpoint,
			Enable:   model.Enable,
			Type:     model.Type,
			Role:     model.Role,
		})
	}

	return localModels, nil
}

// ListAvailableVectorModels 获取可用的向量模型列表
func (s *FlowyModelService) ListAvailableVectorModels(ctx context.Context) ([]models.ModelInfo, error) {
	flowyModels, err := s.sdk.Model.ListAvailableVectorModels(ctx)
	if err != nil {
		return nil, err
	}

	// 转换 flowy 模型到本地模型
	var localModels []models.ModelInfo
	for _, model := range flowyModels {
		localModels = append(localModels, models.ModelInfo{
			ID:       model.ID,
			Name:     model.Name,
			Symbol:   model.Symbol,
			Endpoint: model.Endpoint,
			Enable:   model.Enable,
			Type:     model.Type,
			Role:     model.Role,
		})
	}

	return localModels, nil
}

// SaveModel 添加或修改模型
func (s *FlowyModelService) SaveModel(ctx context.Context, req *models.ModelSaveRequest) (int, error) {
	// 转换本地请求到 flowy 请求
	flowyReq := &modelSvc.ModelSaveRequest{
		ID:       req.ID,
		Name:     req.Name,
		Type:     req.Type,
		Symbol:   req.Symbol,
		Endpoint: req.Endpoint,
		Enable:   req.Enable,
		Credentials: req.Credentials,
	}

	return s.sdk.Model.SaveModel(ctx, flowyReq)
}

// DeleteModel 删除模型
func (s *FlowyModelService) DeleteModel(ctx context.Context, id int) error {
	return s.sdk.Model.DeleteModel(ctx, id)
}

// SetModelStatus 设置模型启用状态
func (s *FlowyModelService) SetModelStatus(ctx context.Context, id int, enable bool) error {
	return s.sdk.Model.SetModelStatus(ctx, id, enable)
}
