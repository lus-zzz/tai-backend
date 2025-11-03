package services

import (
	"context"

	"flowy-sdk"
	modelSvc "flowy-sdk/services/model"
)

// ModelService 模型服务
type ModelService struct {
	sdk *flowy.SDK
}

// NewModelService 创建模型服务
func NewModelService(sdk *flowy.SDK) *ModelService {
	return &ModelService{
		sdk: sdk,
	}
}

// ListSupportedChatModels 获取支持的聊天模型列表
func (s *ModelService) ListSupportedChatModels(ctx context.Context) ([]modelSvc.SupportedChatModel, error) {
	return s.sdk.Model.ListSupportedChatModels(ctx)
}

// ListSupportedVectorModels 获取支持的向量模型列表
func (s *ModelService) ListSupportedVectorModels(ctx context.Context) ([]modelSvc.SupportedVectorModel, error) {
	return s.sdk.Model.ListSupportedVectorModels(ctx)
}

// ListAvailableAllModels 获取全部可用模型列表
func (s *ModelService) ListAvailableAllModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableAllModels(ctx)
}

// ListAvailableChatModels 获取可用的对话模型列表
func (s *ModelService) ListAvailableChatModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableChatModels(ctx)
}

// ListAvailableVectorModels 获取可用的向量模型列表
func (s *ModelService) ListAvailableVectorModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableVectorModels(ctx)
}

// SaveModel 添加或修改模型
func (s *ModelService) SaveModel(ctx context.Context, req *modelSvc.ModelSaveRequest) (int, error) {
	return s.sdk.Model.SaveModel(ctx, req)
}

// DeleteModel 删除模型
func (s *ModelService) DeleteModel(ctx context.Context, id int) error {
	return s.sdk.Model.DeleteModel(ctx, id)
}

// SetModelStatus 设置模型启用状态
func (s *ModelService) SetModelStatus(ctx context.Context, id int, enable bool) error {
	return s.sdk.Model.SetModelStatus(ctx, id, enable)
}
