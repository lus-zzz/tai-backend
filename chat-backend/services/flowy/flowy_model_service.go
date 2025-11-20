package flowy

import (
	"context"

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
func (s *FlowyModelService) ListSupportedChatModels(ctx context.Context) ([]modelSvc.SupportedChatModel, error) {
	return s.sdk.Model.ListSupportedChatModels(ctx)
}

// ListSupportedVectorModels 获取支持的向量模型列表
func (s *FlowyModelService) ListSupportedVectorModels(ctx context.Context) ([]modelSvc.SupportedVectorModel, error) {
	return s.sdk.Model.ListSupportedVectorModels(ctx)
}

// ListAvailableAllModels 获取全部可用模型列表
func (s *FlowyModelService) ListAvailableAllModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableAllModels(ctx)
}

// ListAvailableChatModels 获取可用的对话模型列表
func (s *FlowyModelService) ListAvailableChatModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableChatModels(ctx)
}

// ListAvailableVectorModels 获取可用的向量模型列表
func (s *FlowyModelService) ListAvailableVectorModels(ctx context.Context) ([]modelSvc.ModelInfo, error) {
	return s.sdk.Model.ListAvailableVectorModels(ctx)
}

// SaveModel 添加或修改模型
func (s *FlowyModelService) SaveModel(ctx context.Context, req *modelSvc.ModelSaveRequest) (int, error) {
	return s.sdk.Model.SaveModel(ctx, req)
}

// DeleteModel 删除模型
func (s *FlowyModelService) DeleteModel(ctx context.Context, id int) error {
	return s.sdk.Model.DeleteModel(ctx, id)
}

// SetModelStatus 设置模型启用状态
func (s *FlowyModelService) SetModelStatus(ctx context.Context, id int, enable bool) error {
	return s.sdk.Model.SetModelStatus(ctx, id, enable)
}
