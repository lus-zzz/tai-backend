package langchaingo

import (
	"context"

	"chat-backend/models"
	"chat-backend/pkg/database"
	"chat-backend/services/interfaces"
	"chat-backend/utils"
)

// LangchaingoModelService 基于 langchaingo 的模型服务实现
type LangchaingoModelService struct {
	db      *database.Database
	config  *LangchaingoConfig
}

// NewLangchaingoModelService 创建 Langchaingo 模型服务
func NewLangchaingoModelService(db *database.Database, config *LangchaingoConfig) interfaces.ModelServiceInterface {
	return &LangchaingoModelService{
		db:     db,
		config:  config,
	}
}

// ListSupportedChatModels 获取支持的聊天模型列表
func (s *LangchaingoModelService) ListSupportedChatModels(ctx context.Context) ([]models.SupportedChatModel, error) {
	return s.db.GetSupportedChatModels()
}

// ListSupportedVectorModels 获取支持的向量模型列表
func (s *LangchaingoModelService) ListSupportedVectorModels(ctx context.Context) ([]models.SupportedVectorModel, error) {
	return s.db.GetSupportedVectorModels()
}

// ListAvailableAllModels 获取全部可用模型列表
func (s *LangchaingoModelService) ListAvailableAllModels(ctx context.Context) ([]models.ModelInfo, error) {
	return s.db.GetAllModels()
}

// ListAvailableChatModels 获取可用的对话模型列表
func (s *LangchaingoModelService) ListAvailableChatModels(ctx context.Context) ([]models.ModelInfo, error) {
	return s.db.GetAvailableModelsByType(0) // 0 = chat
}

// ListAvailableVectorModels 获取可用的向量模型列表
func (s *LangchaingoModelService) ListAvailableVectorModels(ctx context.Context) ([]models.ModelInfo, error) {
	return s.db.GetAvailableModelsByType(1) // 1 = embedding
}

// SaveModel 添加或修改模型
func (s *LangchaingoModelService) SaveModel(ctx context.Context, req *models.ModelSaveRequest) (int, error) {
	modelID, err := s.db.SaveModel(req)
	if err != nil {
		utils.ErrorWith("保存模型失败", "error", err.Error())
		return 0, err
	}
	
	utils.InfoWith("保存模型成功", "model_id", modelID, "model_name", req.Name)
	return modelID, nil
}

// DeleteModel 删除模型
func (s *LangchaingoModelService) DeleteModel(ctx context.Context, id int) error {
	err := s.db.DeleteModel(id)
	if err != nil {
		utils.ErrorWith("删除模型失败", "model_id", id, "error", err.Error())
		return err
	}
	
	utils.InfoWith("删除模型成功", "model_id", id)
	return nil
}

// SetModelStatus 设置模型启用状态
func (s *LangchaingoModelService) SetModelStatus(ctx context.Context, id int, enable bool) error {
	err := s.db.UpdateModelStatus(id, enable)
	if err != nil {
		status := "禁用"
		if enable {
			status = "启用"
		}
		utils.ErrorWith("设置模型状态失败", "model_id", id, "status", status, "error", err.Error())
		return err
	}
	
	status := "禁用"
	if enable {
		status = "启用"
	}
	utils.InfoWith("设置模型状态成功", "model_id", id, "status", status)
	return nil
}
