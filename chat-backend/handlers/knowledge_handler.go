package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// knowledgeBaseResponseWrapper represents a knowledge base response.
//
// swagger:response KnowledgeBaseResponse
type knowledgeBaseResponseWrapper struct {
	// in: body
	Body struct {
		Success bool                 `json:"success"`
		Message string               `json:"message,omitempty"`
		Data    models.KnowledgeBase `json:"data,omitempty"`
	}
}

// knowledgeBaseListResponseWrapper represents a knowledge base list response.
//
// swagger:response KnowledgeBaseListResponse
type knowledgeBaseListResponseWrapper struct {
	// in: body
	Body struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message,omitempty"`
		Data    []models.KnowledgeBase `json:"data,omitempty"`
	}
}

// knowledgeFileResponseWrapper represents a knowledge file response.
//
// swagger:response KnowledgeFileResponse
type knowledgeFileResponseWrapper struct {
	// in: body
	Body struct {
		Success bool                 `json:"success"`
		Message string               `json:"message,omitempty"`
		Data    models.KnowledgeFile `json:"data,omitempty"`
	}
}

// knowledgeFileListResponseWrapper represents a knowledge file list response.
//
// swagger:response KnowledgeFileListResponse
type knowledgeFileListResponseWrapper struct {
	// in: body
	Body struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message,omitempty"`
		Data    []models.KnowledgeFile `json:"data,omitempty"`
	}
}

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
// swagger:route GET /api/v1/knowledge/bases Knowledge listKnowledgeBases
//
// Returns a list of all knowledge bases.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: KnowledgeBaseListResponse
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
// swagger:route POST /api/v1/knowledge/bases Knowledge createKnowledgeBase
//
// Creates a new knowledge base.
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
//     description: Knowledge base information
//     required: true
//     type: KnowledgeBaseCreateRequest
//
// Responses:
//
//	200: KnowledgeBaseResponse
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
// swagger:route PUT /api/v1/knowledge/bases/{id} Knowledge updateKnowledgeBase
//
// Updates a knowledge base.
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
//     description: Knowledge base ID
//     required: true
//     type: string
//   - +name: body
//     in: body
//     description: Knowledge base update information
//     required: true
//     type: UpdateKnowledgeBaseRequest
//
// Responses:
//
//	200: KnowledgeBaseResponse
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
// swagger:route DELETE /api/v1/knowledge/bases/{id} Knowledge deleteKnowledgeBase
//
// Deletes a knowledge base.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Knowledge base ID
//     required: true
//     type: string
//
// Responses:
//
//	200: APIResponse
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
// swagger:route GET /api/v1/knowledge/bases/{id}/files Knowledge getKnowledgeBaseFiles
//
// Returns a list of files in a knowledge base.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Knowledge base ID
//     required: true
//     type: string
//
// Responses:
//
//	200: KnowledgeFileListResponse
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
// swagger:route POST /api/v1/knowledge/bases/{id}/files Knowledge uploadFile
//
// Uploads a file to a knowledge base.
//
// Consumes:
// - multipart/form-data
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Knowledge base ID
//     required: true
//     type: string
//   - +name: file
//     in: formData
//     description: File to upload
//     required: true
//     type: file
//
// Responses:
//
//	200: KnowledgeFileResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *KnowledgeHandler) UploadFile(c *gin.Context) {
	kbID := c.Param("id")
	if kbID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少知识库ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	file, header, err := c.Request.FormFile("file")
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "文件上传失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}
	defer file.Close()

	// 读取文件内容
	content := make([]byte, header.Size)
	_, err = file.Read(content)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "读取文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	uploadedFile, err := h.knowledgeService.UploadFile(ctx, kbID, header.Filename, content)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrFileUpload, "上传文件失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, uploadedFile, "文件上传成功")
}

// DeleteFile 从知识库中删除指定的文件。
//
// swagger:route DELETE /api/v1/knowledge/files/{file_id} Knowledge deleteFile
//
// Deletes a file from a knowledge base.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: file_id
//     in: path
//     description: File ID
//     required: true
//     type: integer
//
// Responses:
//
//	200: APIResponse
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
// swagger:route PUT /api/v1/knowledge/files/{file_id}/toggle Knowledge toggleFileEnable
//
// Toggles the enabled state of a file.
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
//     description: File ID
//     required: true
//     type: string
//   - +name: body
//     in: body
//     description: Enable/disable request
//     required: true
//     type: FileToggleEnableRequest
//
// Responses:
//
//	200: APIResponse
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
