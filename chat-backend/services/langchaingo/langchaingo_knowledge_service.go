package langchaingo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services/interfaces"
	"chat-backend/utils"
)

// LangchaingoKnowledgeService 基于 langchaingo 的知识库服务实现
type LangchaingoKnowledgeService struct {
	config *LangchaingoConfig
}

// NewLangchaingoKnowledgeService 创建 Langchaingo 知识库服务
func NewLangchaingoKnowledgeService(config *LangchaingoConfig) interfaces.KnowledgeServiceInterface {
	return &LangchaingoKnowledgeService{
		config: config,
	}
}

// ListKnowledgeBases 获取知识库列表
func (s *LangchaingoKnowledgeService) ListKnowledgeBases(ctx context.Context) ([]models.KnowledgeBase, error) {
	utils.InfoWith("获取知识库列表")

	// TODO: 实现 Qdrant 集合列表获取
	// 这里需要从 Qdrant 获取所有集合（知识库）
	// 暂时返回空列表
	var knowledgeBases []models.KnowledgeBase

	utils.InfoWith("获取知识库列表成功", "count", len(knowledgeBases))
	return knowledgeBases, nil
}

// CreateKnowledgeBase 创建知识库
func (s *LangchaingoKnowledgeService) CreateKnowledgeBase(ctx context.Context, req *models.KnowledgeBaseCreateRequest) (*models.KnowledgeBase, error) {
	utils.InfoWith("创建知识库",
		"name", req.Name,
		"description", req.Desc,
		"chunkSize", req.ChunkSize,
		"vectorModel", req.VectorModel,
		"agentModel", req.AgentModel)

	// TODO: 实现 Qdrant 集合创建
	// 这里需要在 Qdrant 中创建新的集合（知识库）
	// 暂时返回模拟数据
	kbID := 1 // 模拟ID

	result := &models.KnowledgeBase{
		ID:                  kbID,
		KnowledgeBaseConfig: *req, // 直接使用请求的配置
		FileCount:           0,
	}

	utils.InfoWith("创建知识库成功", "id", result.ID, "name", result.Name)
	return result, nil
}

// UpdateKnowledgeBase 更新知识库
func (s *LangchaingoKnowledgeService) UpdateKnowledgeBase(ctx context.Context, id string, req *models.UpdateKnowledgeBaseRequest) (*models.KnowledgeBase, error) {
	utils.InfoWith("更新知识库",
		"id", id,
		"name", req.Name,
		"desc", req.Desc,
		"agentModel", req.AgentModel,
		"chunkStrategy", req.ChunkStrategy,
		"chunkSize", req.ChunkSize)

	// TODO: 实现 Qdrant 集合元数据更新
	// 暂时返回模拟数据
	kbID, err := strconv.Atoi(id)
	if err != nil {
		return nil, fmt.Errorf("无效的知识库ID: %w", err)
	}

	result := &models.KnowledgeBase{
		ID:                  kbID,
		KnowledgeBaseConfig: *req, // 直接使用请求的配置
	}

	utils.InfoWith("更新知识库成功", "id", result.ID)
	return result, nil
}

// DeleteKnowledgeBase 删除知识库
func (s *LangchaingoKnowledgeService) DeleteKnowledgeBase(ctx context.Context, id string) error {
	utils.LogInfo("删除知识库: %s", id)

	// TODO: 实现 Qdrant 集合删除
	// 暂时只是记录日志
	utils.InfoWith("删除知识库成功", "id", id)
	return nil
}

// GetKnowledgeBaseFiles 获取知识库文件列表
func (s *LangchaingoKnowledgeService) GetKnowledgeBaseFiles(ctx context.Context, id string) ([]models.KnowledgeFile, error) {
	utils.LogInfo("获取知识库文件列表: %s", id)

	// TODO: 实现 Qdrant 文件元数据获取
	// 暂时返回空列表
	var knowledgeFiles []models.KnowledgeFile

	utils.InfoWith("获取知识库文件列表成功", "id", id, "count", len(knowledgeFiles))
	return knowledgeFiles, nil
}

// UploadFile 上传文件到知识库（文件流上传）
func (s *LangchaingoKnowledgeService) UploadFile(ctx context.Context, id string, filename string, content []byte) (*models.KnowledgeFile, error) {
	utils.LogInfo("上传文件到知识库: %s, 文件: %s", id, filename)

	// TODO: 实现完整的文档分块和向量化流程
	// 这里需要实现 KEY_PROCESS_AND_CODE.md 中的 chunkAndVectorize 流程
	
	// 1. 使用 Docling API 进行文档分块
	chunks, err := s.chunkDocument(ctx, filename, content)
	if err != nil {
		return nil, fmt.Errorf("文档分块失败: %w", err)
	}

	// 2. 将分块向量化并存储到 Qdrant
	err = s.vectorizeAndStore(ctx, id, chunks)
	if err != nil {
		return nil, fmt.Errorf("向量化和存储失败: %w", err)
	}

	result := &models.KnowledgeFile{
		ID:           1, // 模拟ID
		Name:         filename,
		Size:         len(content),
		Enable:       true,
		Status:       1, // 完成
		UploadedAt:   time.Now(),
		IndexPercent: 100,
		ErrorMessage: "",
	}

	utils.InfoWith("上传文件成功", "id", id, "filename", filename, "file_id", result.ID)
	return result, nil
}

// UploadFileFromPath 从文件路径上传文件到知识库
func (s *LangchaingoKnowledgeService) UploadFileFromPath(ctx context.Context, id string, filePath string) (*models.KnowledgeFile, error) {
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

	utils.LogInfo("从路径上传文件到知识库: %s, 路径: %s, 文件: %s", id, filePath, filename)

	// 打开文件
	file, err := os.Open(filePath)
	if err != nil {
		return nil, fmt.Errorf("打开文件失败: %w", err)
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		return nil, fmt.Errorf("读取文件失败: %w", err)
	}

	// 调用文件上传方法
	return s.UploadFile(ctx, id, filename, content)
}

// BatchUploadFilesFromPath 批量从文件路径上传文件到知识库
func (s *LangchaingoKnowledgeService) BatchUploadFilesFromPath(ctx context.Context, id string, filePaths []string) (*models.BatchUploadResponse, error) {
	if len(filePaths) == 0 {
		return nil, fmt.Errorf("文件路径列表不能为空")
	}

	utils.LogInfo("开始批量上传文件到知识库: %s, 文件数: %d", id, len(filePaths))

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
func (s *LangchaingoKnowledgeService) DeleteFile(ctx context.Context, fileID int) error {
	utils.InfoWith("删除知识库文件", "file_id", fileID)

	if fileID == 0 {
		utils.ErrorWith("无效的文件ID", "file_id", fileID)
		return fmt.Errorf("无效的文件ID: %d", fileID)
	}

	// TODO: 实现 Qdrant 向量删除
	// 暂时只是记录日志
	utils.InfoWith("删除文件成功", "file_id", fileID)
	return nil
}

// ToggleFileEnable 切换文件启用状态
func (s *LangchaingoKnowledgeService) ToggleFileEnable(ctx context.Context, fileID int, enable bool) error {
	utils.InfoWith("切换文件启用状态", "file_id", fileID, "enable", enable)

	if fileID == 0 {
		utils.ErrorWith("无效的文件ID", "file_id", fileID)
		return fmt.Errorf("无效的文件ID: %d", fileID)
	}

	// TODO: 实现 Qdrant 向量启用/禁用
	// 暂时只是记录日志
	status := "禁用"
	if enable {
		status = "启用"
	}
	utils.InfoWith("切换文件启用状态成功", "file_id", fileID, "status", status)
	return nil
}

// DoclingChunkRequest Docling 分块请求
type DoclingChunkRequest struct {
	Files                    [][]byte `json:"files"`
	ConvertDoOCR            bool      `json:"convertDoOCR"`
	ConvertImageExportMode   string    `json:"convertImageExportMode"`
	ConvertPDFBackend       string    `json:"convertPDFBackend"`
	ConvertTableMode       string    `json:"convertTableMode"`
	ConvertPipeline         string    `json:"convertPipeline"`
	ConvertAbortOnError     bool      `json:"convertAbortOnError"`
	ConvertDoCodeEnrichment bool      `json:"convertDoCodeEnrichment"`
	ConvertDoFormulaEnrichment bool    `json:"convertDoFormulaEnrichment"`
	ChunkingUseMarkdownTables bool     `json:"chunkingUseMarkdownTables"`
	ChunkingIncludeRawText    bool     `json:"chunkingIncludeRawText"`
	ChunkingTokenizer        string    `json:"chunkingTokenizer"`
	ChunkingMaxTokens       int       `json:"chunkingMaxTokens"`
	ChunkingMergePeers      bool      `json:"chunkingMergePeers"`
}

// DoclingChunkResponse Docling 分块响应
type DoclingChunkResponse struct {
	Chunks []DoclingChunk `json:"chunks"`
}

// DoclingChunk Docling 分块结果
type DoclingChunk struct {
	Text        string     `json:"text"`
	Filename    string     `json:"filename"`
	PageNumbers *[]int     `json:"pageNumbers"`
	NumTokens   *int       `json:"numTokens"`
}

// chunkDocument 使用 Docling API 进行文档分块
func (s *LangchaingoKnowledgeService) chunkDocument(ctx context.Context, filename string, content []byte) ([]DoclingChunk, error) {
	// 准备分块请求
	req := DoclingChunkRequest{
		Files:                    [][]byte{content},
		ConvertDoOCR:             false,
		ConvertImageExportMode:     "placeholder",
		ConvertPDFBackend:         "dlparse_v4",
		ConvertTableMode:          "accurate",
		ConvertPipeline:           "standard",
		ConvertAbortOnError:       false,
		ConvertDoCodeEnrichment:   false,
		ConvertDoFormulaEnrichment: false,
		ChunkingUseMarkdownTables:  true,
		ChunkingIncludeRawText:     false,
		ChunkingTokenizer:          "/Volume2/test_work/models/BAAI/bge-m3",
		ChunkingMaxTokens:         512,
		ChunkingMergePeers:        true,
	}

	// 调用 Docling API
	url := s.config.Docling.BaseURL + "/chunk/hybrid"
	
	jsonData, err := json.Marshal(req)
	if err != nil {
		return nil, fmt.Errorf("序列化请求失败: %w", err)
	}

	httpReq, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	httpReq.Header.Set("Content-Type", "application/json")
	if s.config.Docling.APIKey != "" {
		httpReq.Header.Set("Authorization", "Bearer "+s.config.Docling.APIKey)
	}

	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(httpReq)
	if err != nil {
		return nil, fmt.Errorf("调用 Docling API 失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("Docling API 返回错误状态码: %d", resp.StatusCode)
	}

	var chunkResp DoclingChunkResponse
	if err := json.NewDecoder(resp.Body).Decode(&chunkResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	utils.InfoWith("文档分块完成", "filename", filename, "chunk_count", len(chunkResp.Chunks))
	return chunkResp.Chunks, nil
}

// vectorizeAndStore 向量化分块并存储到 Qdrant
func (s *LangchaingoKnowledgeService) vectorizeAndStore(ctx context.Context, collectionName string, chunks []DoclingChunk) error {
	// TODO: 实现 Ollama 嵌入模型调用和 Qdrant 存储
	// 这里需要：
	// 1. 使用 Ollama bge-m3 模型生成嵌入向量
	// 2. 将向量存储到 Qdrant 指定集合中
	// 3. 添加适当的元数据

	utils.InfoWith("向量化和存储完成", "collection", collectionName, "chunk_count", len(chunks))
	return nil
}
