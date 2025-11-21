package flowy

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	"chat-backend/models"
	"chat-backend/services/interfaces"
	"chat-backend/utils"
	"flowy-sdk"
	knowledgeSvc "flowy-sdk/services/knowledge"
)

// FlowyKnowledgeService 基于 flowy-sdk 的知识库服务实现
type FlowyKnowledgeService struct {
	sdk *flowy.SDK
}

// NewFlowyKnowledgeService 创建 Flowy 知识库服务
func NewFlowyKnowledgeService(sdk *flowy.SDK) interfaces.KnowledgeServiceInterface {
	return &FlowyKnowledgeService{
		sdk: sdk,
	}
}

// ListKnowledgeBases 获取知识库列表
func (s *FlowyKnowledgeService) ListKnowledgeBases(ctx context.Context) ([]models.KnowledgeBase, error) {
	utils.InfoWith("获取知识库列表")

	// 调用Flowy SDK获取知识库列表
	response, err := s.sdk.Knowledge.ListKnowledgeBases(ctx)
	if err != nil {
		utils.ErrorWith("获取知识库列表失败", "error", err)
		return nil, fmt.Errorf("获取知识库列表失败: %w", err)
	}

	// 转换数据 - 保持与 KnowledgeBaseCreateRequest 一致的字段结构
	var knowledgeBases []models.KnowledgeBase
	for _, kb := range response {
		knowledgeBase := models.KnowledgeBase{
			ID: kb.ID,
			KnowledgeBaseConfig: models.KnowledgeBaseConfig{
				Name:          kb.Name,
				Desc:          kb.Desc,
				VectorModel:   kb.VectorModel,
				AgentModel:    kb.AgentModel,
				ChunkStrategy: kb.ChunkStrategy,
				ChunkSize:     kb.ChunkSize,
			},
			FileCount: kb.FileCount,
		}
		knowledgeBases = append(knowledgeBases, knowledgeBase)
	}

	utils.InfoWith("获取知识库列表成功", "count", len(knowledgeBases))
	return knowledgeBases, nil
}

// CreateKnowledgeBase 创建知识库
func (s *FlowyKnowledgeService) CreateKnowledgeBase(ctx context.Context, req *models.KnowledgeBaseCreateRequest) (*models.KnowledgeBase, error) {
	utils.InfoWith("创建知识库",
		"name", req.Name,
		"description", req.Desc,
		"chunkSize", req.ChunkSize,
		"vectorModel", req.VectorModel,
		"agentModel", req.AgentModel)

	// 根据模型ID获取模型名称（用于日志记录）
	embeddingModelName := ""
	if req.VectorModel > 0 {
		name, err := s.getModelNameByID(ctx, req.VectorModel)
		if err != nil {
			utils.WarnWith("获取嵌入模型名称失败", "model_id", req.VectorModel, "error", err)
		} else {
			embeddingModelName = name
			utils.InfoWith("使用嵌入模型", "model_id", req.VectorModel, "model_name", embeddingModelName)
		}
	}

	agentModelName := ""
	if req.AgentModel > 0 {
		name, err := s.getModelNameByID(ctx, req.AgentModel)
		if err != nil {
			utils.WarnWith("获取对话模型名称失败", "model_id", req.AgentModel, "error", err)
		} else {
			agentModelName = name
			utils.InfoWith("使用对话模型", "model_id", req.AgentModel, "model_name", agentModelName)
		}
	}

	// 调用Flowy SDK创建知识库，使用工厂函数创建默认配置
	createReq := knowledgeSvc.NewDefaultKnowledgeBaseCreateRequest(req.Name, req.Desc)

	// 使用前端传入的配置覆盖默认值
	createReq.VectorModel = req.VectorModel     // 使用前端传入的嵌入模型ID（向量模型）
	createReq.AgentModel = req.AgentModel       // 使用前端传入的对话模型ID
	createReq.ChunkStrategy = req.ChunkStrategy // 使用前端传入的分块策略
	createReq.ChunkSize = req.ChunkSize         // 使用前端传入的分块大小

	kbID, err := s.sdk.Knowledge.CreateKnowledgeBase(ctx, createReq)
	if err != nil {
		utils.ErrorWith("创建知识库失败", "name", req.Name, "error", err)
		return nil, fmt.Errorf("创建知识库失败: %w", err)
	}

	result := &models.KnowledgeBase{
		ID:                  kbID,
		KnowledgeBaseConfig: *req, // 直接使用请求的配置
		FileCount:           0,
	}

	utils.InfoWith("创建知识库成功",
		"id", result.ID,
		"name", result.Name,
		"chunk_size", req.ChunkSize,
		"embedding_model", embeddingModelName,
		"agent_model", agentModelName)
	return result, nil
}

// UpdateKnowledgeBase 更新知识库
func (s *FlowyKnowledgeService) UpdateKnowledgeBase(ctx context.Context, id int, req *models.UpdateKnowledgeBaseRequest) (*models.KnowledgeBase, error) {
	utils.InfoWith("更新知识库",
		"id", id,
		"name", req.Name,
		"desc", req.Desc,
		"agentModel", req.AgentModel,
		"chunkStrategy", req.ChunkStrategy,
		"chunkSize", req.ChunkSize)

	// 调用Flowy SDK更新知识库
	updateReq := knowledgeSvc.NewDefaultKnowledgeBaseUpdateRequest(id, req.Name, req.Desc)

	// 使用前端传入的配置覆盖默认值（注意：VectorModel 不能更新）
	updateReq.AgentModel = req.AgentModel
	updateReq.ChunkStrategy = req.ChunkStrategy
	updateReq.ChunkSize = req.ChunkSize

	updatedID, err := s.sdk.Knowledge.UpdateKnowledgeBase(ctx, updateReq)
	if err != nil {
		utils.ErrorWith("更新知识库失败", "id", id, "error", err)
		return nil, fmt.Errorf("更新知识库失败: %w", err)
	}

	result := &models.KnowledgeBase{
		ID:                  updatedID,
		KnowledgeBaseConfig: *req, // 直接使用请求的配置
	}

	utils.InfoWith("更新知识库成功", "id", result.ID)
	return result, nil
}

// DeleteKnowledgeBase 删除知识库
func (s *FlowyKnowledgeService) DeleteKnowledgeBase(ctx context.Context, id int) error {
	utils.InfoWith("删除知识库", "id", id)

	// 调用Flowy SDK删除知识库
	err := s.sdk.Knowledge.DeleteKnowledgeBase(ctx, id)
	if err != nil {
		utils.ErrorWith("删除知识库失败", "id", id, "error", err)
		return fmt.Errorf("删除知识库失败: %w", err)
	}

	utils.InfoWith("删除知识库成功", "id", id)
	return nil
}

// GetKnowledgeBaseFiles 获取知识库文件列表
func (s *FlowyKnowledgeService) GetKnowledgeBaseFiles(ctx context.Context, id int) ([]models.KnowledgeFile, error) {
	utils.InfoWith("获取知识库文件列表", "id", id)

	// 调用Flowy SDK获取文件列表
	files, err := s.sdk.Knowledge.ListFiles(ctx, id, "zh")
	if err != nil {
		utils.ErrorWith("获取知识库文件列表失败", "id", id, "error", err)
		return nil, fmt.Errorf("获取知识库文件列表失败: %w", err)
	}

	// 转换数据
	var knowledgeFiles []models.KnowledgeFile
	for _, file := range files {
		kf := models.KnowledgeFile{
			ID:           file.ID,
			Name:         file.Name,
			Size:         file.FileSize,
			Enable:       file.Enable,
			Status:       file.Status,
			UploadedAt:   parseTime(file.CreateTime),
			IndexPercent: file.IndexPercent,
			ErrorMessage: file.ErrorMessage,
		}
		knowledgeFiles = append(knowledgeFiles, kf)
	}

	utils.InfoWith("获取知识库文件列表成功", "id", id, "count", len(knowledgeFiles))
	return knowledgeFiles, nil
}

// UploadFile 上传文件到知识库（文件流上传）
func (s *FlowyKnowledgeService) UploadFile(ctx context.Context, id int, filename string, reader io.Reader) (*models.KnowledgeFile, error) {
	utils.InfoWith("上传文件到知识库", "id", id, "filename", filename)

	// 调用Flowy SDK上传文件
	uploadData, err := s.sdk.Knowledge.UploadFile(ctx, reader, filename, id, 0, "zh")
	if err != nil {
		utils.ErrorWith("上传文件失败", "id", id, "filename", filename, "error", err)
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	result := &models.KnowledgeFile{
		ID:           uploadData.ID,
		Name:         uploadData.Name,
		Size:         uploadData.FileSize,
		Enable:       uploadData.Enable,
		Status:       uploadData.Status,
		UploadedAt:   time.Now(),
		IndexPercent: uploadData.IndexPercent,
		ErrorMessage: uploadData.ErrorMessage,
	}

	utils.InfoWith("上传文件成功", "id", id, "filename", filename, "file_id", uploadData.ID)
	return result, nil
}

// UploadFileFromPath 从文件路径上传文件到知识库
func (s *FlowyKnowledgeService) UploadFileFromPath(ctx context.Context, id int, filePath string) (*models.KnowledgeFile, error) {
	// 验证文件路径
	if filePath == "" {
		return nil, fmt.Errorf("文件路径不能为空")
	}

	// 验证文件是否存在
	fileInfo, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("文件不存在: %s", filePath)
		}
		return nil, fmt.Errorf("无法访问文件: %w", err)
	}

	// 验证是否是文件（不是目录）
	if fileInfo.IsDir() {
		return nil, fmt.Errorf("路径是一个目录，不是文件: %s", filePath)
	}

	// 获取文件名
	filename := filepath.Base(filePath)

	utils.LogInfo("从路径上传文件到知识库: %d, 路径: %s, 文件: %s", id, filePath, filename)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 调用Flowy SDK上传文件，直接传递文件 reader
	uploadData, err := s.sdk.Knowledge.UploadFile(ctx, file, filename, id, 0, "zh")
	if err != nil {
		utils.ErrorWith("从路径上传文件失败", "id", id, "path", filePath, "filename", filename, "error", err)
		return nil, fmt.Errorf("上传文件失败: %w", err)
	}

	result := &models.KnowledgeFile{
		ID:           uploadData.ID,
		Name:         uploadData.Name,
		Size:         uploadData.FileSize,
		Enable:       uploadData.Enable,
		Status:       uploadData.Status,
		UploadedAt:   time.Now(),
		IndexPercent: uploadData.IndexPercent,
		ErrorMessage: uploadData.ErrorMessage,
	}

	utils.InfoWith("从路径上传文件成功", "id", id, "path", filePath, "filename", filename, "file_id", uploadData.ID)
	return result, nil
}

// BatchUploadFilesFromPath 批量从文件路径上传文件到知识库
func (s *FlowyKnowledgeService) BatchUploadFilesFromPath(ctx context.Context, id int, filePaths []string) (*models.BatchUploadResponse, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("文件路径列表不能为空")
	}

	utils.LogInfo("开始批量上传文件到知识库: %d, 文件数: %d", id, len(filePaths))

	response := &models.BatchUploadResponse{
		Total:   len(filePaths),
		Results: make([]models.BatchUploadResult, 0, len(filePaths)),
	}

	// 创建超时上下文用于整个批量操作
	batchCtx, cancel := context.WithTimeout(context.Background(), time.Duration(len(filePaths)*60)*time.Second)
	defer cancel()

	// 上传每个文件
	for _, filePath := range filePaths {
		result := models.BatchUploadResult{
			FilePath: filePath,
		}

		// 为每个文件单独设置超时
		fileCtx, fileCancel := context.WithTimeout(batchCtx, 60*time.Second)

		uploadedFile, err := s.UploadFileFromPath(fileCtx, id, filePath)
		fileCancel()

		if err != nil {
			result.Success = false
			result.Message = "上传失败"
			result.Error = err.Error()
			response.FailureCount++
			utils.ErrorWith("批量上传文件失败", "path", filePath, "error", err)
		} else {
			result.Success = true
			result.Message = "上传成功"
			result.File = uploadedFile
			response.SuccessCount++
			utils.InfoWith("批量上传文件成功", "path", filePath, "file_id", uploadedFile.ID)
		}

		response.Results = append(response.Results, result)
	}

	utils.InfoWith("批量上传完成", "知识库ID", id, "总数", response.Total, "成功", response.SuccessCount, "失败", response.FailureCount)
	return response, nil
}

// DeleteFile 删除知识库文件
func (s *FlowyKnowledgeService) DeleteFile(ctx context.Context, fileID int) error {
	utils.InfoWith("删除知识库文件", "file_id", fileID)

	if fileID == 0 {
		utils.ErrorWith("无效的文件ID", "file_id", fileID)
		return fmt.Errorf("无效的文件ID: %d", fileID)
	}

	// 调用Flowy SDK删除文件
	err := s.sdk.Knowledge.DeleteFile(ctx, fileID)
	if err != nil {
		utils.ErrorWith("删除文件失败", "file_id", fileID, "error", err)
		return fmt.Errorf("删除文件失败: %w", err)
	}

	utils.InfoWith("删除文件成功", "file_id", fileID)
	return nil
}

// ToggleFileEnable 切换文件启用状态
func (s *FlowyKnowledgeService) ToggleFileEnable(ctx context.Context, fileID int, enable bool) error {
	utils.InfoWith("切换文件启用状态", "file_id", fileID, "enable", enable)

	if fileID == 0 {
		utils.ErrorWith("无效的文件ID", "file_id", fileID)
		return fmt.Errorf("无效的文件ID: %d", fileID)
	}

	// 调用Flowy SDK切换文件启用状态
	err := s.sdk.Knowledge.ToggleFileEnable(ctx, fileID, enable)
	if err != nil {
		utils.ErrorWith("切换文件启用状态失败", "file_id", fileID, "enable", enable, "error", err)
		return fmt.Errorf("切换文件启用状态失败: %w", err)
	}

	status := "禁用"
	if enable {
		status = "启用"
	}
	utils.InfoWith("切换文件启用状态成功", "file_id", fileID, "status", status)
	return nil
}

func getFileType(filename string) string {
	// 简单的文件类型判断
	if len(filename) > 4 {
		ext := filename[len(filename)-4:]
		return ext
	}
	return "unknown"
}

func parseTime(timeStr string) time.Time {
	// 尝试解析时间字符串
	t, err := time.Parse("2006-01-02 15:04:05", timeStr)
	if err != nil {
		return time.Now()
	}
	return t
}

// getModelNameByID 根据模型ID获取模型名称
func (s *FlowyKnowledgeService) getModelNameByID(ctx context.Context, modelID int) (string, error) {
	if modelID <= 0 {
		return "", fmt.Errorf("无效的模型ID: %d", modelID)
	}

	// 获取所有可用模型
	models, err := s.sdk.Model.ListAvailableAllModels(ctx)
	if err != nil {
		return "", fmt.Errorf("获取模型列表失败: %w", err)
	}

	// 查找匹配的模型
	for _, model := range models {
		if model.ID == modelID {
			return model.Symbol, nil
		}
	}

	return "", fmt.Errorf("未找到ID为 %d 的模型", modelID)
}
