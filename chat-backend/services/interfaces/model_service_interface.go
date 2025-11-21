package interfaces

import (
	"context"

	"chat-backend/models"
)

// ModelServiceInterface 模型服务接口
type ModelServiceInterface interface {
	// ListSupportedChatModels 获取支持的聊天模型列表
	ListSupportedChatModels(ctx context.Context) ([]models.SupportedChatModel, error)

	// ListSupportedVectorModels 获取支持的向量模型列表
	ListSupportedVectorModels(ctx context.Context) ([]models.SupportedVectorModel, error)

	// ListAvailableAllModels 获取全部可用模型列表
	ListAvailableAllModels(ctx context.Context) ([]models.ModelInfo, error)

	// ListAvailableChatModels 获取可用的对话模型列表
	ListAvailableChatModels(ctx context.Context) ([]models.ModelInfo, error)

	// ListAvailableVectorModels 获取可用的向量模型列表
	ListAvailableVectorModels(ctx context.Context) ([]models.ModelInfo, error)

	// SaveModel 添加或修改模型
	SaveModel(ctx context.Context, req *models.ModelSaveRequest) (int, error)

	// DeleteModel 删除模型
	DeleteModel(ctx context.Context, id int) error

	// SetModelStatus 设置模型启用状态
	SetModelStatus(ctx context.Context, id int, enable bool) error
}
