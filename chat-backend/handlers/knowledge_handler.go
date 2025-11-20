package handlers

import (
	"context"
	"io"
	"mime/multipart"
	"net/http"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// KnowledgeHandler 处理知识库相关的HTTP请求。
type KnowledgeHandler struct {
	knowledgeService *services.KnowledgeService
}

// NewKnowledgeHandler 创建并返回一个新的知识库处理器实例。
func NewKnowledgeHandler(knowledgeService *services.KnowledgeService) *KnowledgeHandler {
	return &KnowledgeHandler{
		knowledgeService: knowledgeService,
	}
}

// ListKnowledgeBases 返回所有知识库的列表。
//
// swagger:route GET /knowledge/bases Knowledge listKnowledgeBases
//
// 获取知识库列表
//
// 获取系统中所有已创建的知识库列表，包括知识库基本信息
//
// Produces:
// - application/json
//
// Responses:
//
//	200: KnowledgeBaseListSuccessResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) ListKnowledgeBases(c *gin.Context) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	knowledgeBases, err := h.knowledgeService.ListKnowledgeBases(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFlowyAPI, "获取知识库列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, knowledgeBases)
}

// CreateKnowledgeBase 创建一个新的知识库。
//
// swagger:route POST /knowledge/bases Knowledge createKnowledgeBase
//
// 创建知识库
//
// 创建一个新的知识库，用于存储和管理相关文档
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: body
//     in: body
//     description: 知识库信息
//     required: true
//     type: KnowledgeBaseCreateRequest
//
// Responses:
//
//	200: KnowledgeBaseSuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) CreateKnowledgeBase(c *gin.Context) {
	var req models.KnowledgeBaseCreateRequest
	// 使用 BindJSON 避免 validate 标签验证
	if err := c.BindJSON(&req); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	if req.Name == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "知识库名称不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	knowledgeBase, err := h.knowledgeService.CreateKnowledgeBase(ctx, &req)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrKnowledgeBaseCreate, "创建知识库失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, knowledgeBase, "知识库创建成功")
}

// UpdateKnowledgeBase 更新指定知识库的信息。
//
// swagger:route PUT /knowledge/bases/{id} Knowledge updateKnowledgeBase
//
// 更新知识库
//
// 更新指定知识库的名称、描述等基本信息
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: 知识库ID
//     required: true
//     type: string
//   - +name: body
//     in: body
//     description: 更新信息
//     required: true
//     type: UpdateKnowledgeBaseRequest
//
// Responses:
//
//	200: KnowledgeBaseSuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) UpdateKnowledgeBase(c *gin.Context) {
	kbID := c.Param("id")
	if kbID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少知识库ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	var req models.UpdateKnowledgeBaseRequest
	// 使用 BindJSON 避免 validate 标签验证
	if err := c.BindJSON(&req); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	if req.Name == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "知识库名称不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	knowledgeBase, err := h.knowledgeService.UpdateKnowledgeBase(ctx, kbID, &req)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrKnowledgeBaseUpdate, "更新知识库失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, knowledgeBase, "知识库更新成功")
}

// DeleteKnowledgeBase 根据ID删除指定的知识库。
//
// swagger:route DELETE /knowledge/bases/{id} Knowledge deleteKnowledgeBase
//
// 删除知识库
//
// 根据知识库ID删除指定的知识库及其所有相关文件
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: 知识库ID
//     required: true
//     type: string
//
// Responses:
//
//	200: EmptySuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) DeleteKnowledgeBase(c *gin.Context) {
	kbID := c.Param("id")
	if kbID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少知识库ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.knowledgeService.DeleteKnowledgeBase(ctx, kbID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrKnowledgeBaseDelete, "删除知识库失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, nil, "知识库删除成功")
}

// GetKnowledgeBaseFiles 返回指定知识库中的文件列表。
//
// swagger:route GET /knowledge/bases/{id}/files Knowledge getKnowledgeBaseFiles
//
// 获取知识库文件
//
// 获取指定知识库中所有已上传文件的列表和状态信息
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: 知识库ID
//     required: true
//     type: string
//
// Responses:
//
//	200: KnowledgeFileListSuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) GetKnowledgeBaseFiles(c *gin.Context) {
	kbID := c.Param("id")
	if kbID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少知识库ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	files, err := h.knowledgeService.GetKnowledgeBaseFiles(ctx, kbID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFlowyAPI, "获取文件列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, files)
}

// UploadFile 上传文件到指定的知识库。
//
// swagger:route POST /knowledge/bases/{id}/files Knowledge uploadFile
//
// 上传文件到知识库
//
// 支持两种上传方式：
// 文件流上传：使用 multipart/form-data 格式，字段名为 'file',多个文件使用多个 'file' 字段
// 文件路径上传：使用 application/json 格式，body 中包含 file_paths 数组
//
// Consumes:
// - multipart/form-data
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: 知识库ID
//     required: true
//     type: string
//   - +name: file
//     in: formData
//     description: 上传文件（支持单个或多个文件）
//     required: false
//     type: file
//   - +name: body
//     in: body
//     description: 单文件路径或多文件路径列表（用于路径上传）
//     required: false
//     type: BatchUploadFilesRequest
//
// Responses:
//
//	200: KnowledgeFileSuccessResponse or BatchUploadResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) UploadFile(c *gin.Context) {
	kbID := c.Param("id")
	if kbID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少知识库ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 判断上传方式：先检查是否是表单数据（文件流上传），否则检查JSON（路径上传）
	contentType := c.ContentType()

	if contentType == "application/json" {
		// JSON上传方式，支持单文件或多文件路径
		h.uploadFilesFromPath(c, kbID)
	} else {
		// 文件流上传方式，支持单文件或多文件
		h.uploadFilesFromStream(c, kbID)
	}
}

// uploadFilesFromStream 从文件流上传文件（支持单文件和多文件）
func (h *KnowledgeHandler) uploadFilesFromStream(c *gin.Context, kbID string) {
	// 获取所有上传的文件
	form, err := c.MultipartForm()
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "解析表单失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	files := form.File["file"]
	if len(files) == 0 {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "没有找到上传的文件", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 多个文件，返回批量上传结果
	h.uploadMultipleFilesFromStream(c, kbID, files)
}

// uploadSingleFileFromStream 上传单个文件流
func (h *KnowledgeHandler) uploadSingleFileFromStream(c *gin.Context, kbID string, fileHeader *multipart.FileHeader) {
	file, err := fileHeader.Open()
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "打开文件失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}
	defer file.Close()

	// 读取文件内容
	content, err := io.ReadAll(file)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "读取文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	uploadedFile, err := h.knowledgeService.UploadFile(ctx, kbID, fileHeader.Filename, content)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "上传文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, uploadedFile, "文件上传成功")
}

// uploadMultipleFilesFromStream 上传多个文件流
func (h *KnowledgeHandler) uploadMultipleFilesFromStream(c *gin.Context, kbID string, fileHeaders []*multipart.FileHeader) {
	response := &models.BatchUploadResponse{
		Total:   len(fileHeaders),
		Results: make([]models.BatchUploadResult, 0, len(fileHeaders)),
	}

	// 创建超时上下文，根据文件数量调整超时时间
	timeoutSeconds := len(fileHeaders) * 60
	batchCtx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	// 上传每个文件
	for _, fileHeader := range fileHeaders {
		result := models.BatchUploadResult{
			FilePath: fileHeader.Filename,
		}

		// 为每个文件单独设置超时
		fileCtx, fileCancel := context.WithTimeout(batchCtx, 60*time.Second)

		// 打开文件
		file, err := fileHeader.Open()
		if err != nil {
			result.Success = false
			result.Message = "打开文件失败"
			result.Error = err.Error()
			response.FailureCount++
			fileCancel()
			response.Results = append(response.Results, result)
			continue
		}

		// 读取文件内容
		content, err := io.ReadAll(file)
		file.Close()

		if err != nil {
			result.Success = false
			result.Message = "读取文件失败"
			result.Error = err.Error()
			response.FailureCount++
			fileCancel()
			response.Results = append(response.Results, result)
			continue
		}

		// 上传文件
		uploadedFile, err := h.knowledgeService.UploadFile(fileCtx, kbID, fileHeader.Filename, content)
		fileCancel()

		if err != nil {
			result.Success = false
			result.Message = "上传失败"
			result.Error = err.Error()
			response.FailureCount++
			utils.ErrorWith("批量上传文件失败", "filename", fileHeader.Filename, "error", err)
		} else {
			result.Success = true
			result.Message = "上传成功"
			result.File = uploadedFile
			response.SuccessCount++
			utils.InfoWith("批量上传文件成功", "filename", fileHeader.Filename, "file_id", uploadedFile.ID)
		}

		response.Results = append(response.Results, result)
	}

	utils.RespondWithSuccess(c, response, "批量上传完成")
}

// uploadFilesFromPath 从文件路径上传文件（支持单文件和多文件，统一使用 file_paths）
func (h *KnowledgeHandler) uploadFilesFromPath(c *gin.Context, kbID string) {
	var req models.BatchUploadFilesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	if len(req.FilePaths) == 0 {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "文件路径列表不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 统一处理：直接调用批量上传，遍历处理所有文件
	timeoutSeconds := len(req.FilePaths) * 60
	ctx, cancel := context.WithTimeout(context.Background(), time.Duration(timeoutSeconds)*time.Second)
	defer cancel()

	response, err := h.knowledgeService.BatchUploadFilesFromPath(ctx, kbID, req.FilePaths)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "上传文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, response, "文件上传完成")
}

// DeleteFile 从知识库中删除指定的文件。
//
// swagger:route DELETE /knowledge/files/{file_id} Knowledge deleteFile
//
// 删除知识库文件
//
// 从知识库中删除指定的文件及其向量数据
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: file_id
//     in: path
//     description: 文件ID
//     required: true
//     type: integer
//
// Responses:
//
//	200: EmptySuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) DeleteFile(c *gin.Context) {
	fileID := c.Param("file_id")

	id, err := strconv.Atoi(fileID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "无效的文件ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.knowledgeService.DeleteFile(ctx, id)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileNotFound, "删除文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, nil, "文件删除成功")
}

// ToggleFileEnable 切换文件的启用状态。
//
// swagger:route PUT /knowledge/files/{file_id}/toggle Knowledge toggleFileEnable
//
// 切换文件启用状态
//
// 启用或禁用知识库中的指定文件，控制其是否参与问答检索
//
// Consumes:
// - application/json
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: file_id
//     in: path
//     description: 文件ID
//     required: true
//     type: string
//   - +name: body
//     in: body
//     description: 状态切换
//     required: true
//     type: FileToggleEnableRequest
//
// Responses:
//
//	200: EmptySuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) ToggleFileEnable(c *gin.Context) {
	fileID := c.Param("file_id")

	id, err := strconv.Atoi(fileID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "无效的文件ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	var req models.FileToggleEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithValidationError(c, err)
		return
	}

	if req.Enable == nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "enable参数不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = h.knowledgeService.ToggleFileEnable(ctx, id, *req.Enable)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFlowyAPI, "切换文件状态失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	message := "文件已禁用"
	if *req.Enable {
		message = "文件已启用"
	}
	utils.RespondWithSuccess(c, nil, message)
}
