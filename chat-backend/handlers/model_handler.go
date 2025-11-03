package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	modelSvc "flowy-sdk/services/model"

	"github.com/gin-gonic/gin"
)

// SupportedModelListResponse represents a model list response.
//
// swagger:response SupportedModelListResponse
type SupportedModelListResponse struct {
	// in: body
	Body struct {
		Success bool                          `json:"success"`
		Message string                        `json:"message,omitempty"`
		Data    []modelSvc.SupportedChatModel `json:"data,omitempty"`
	}
}

// SupportedVectorModelListResponse represents a vector model list response.
//
// swagger:response SupportedVectorModelListResponse
type SupportedVectorModelListResponse struct {
	// in: body
	Body struct {
		Success bool                            `json:"success"`
		Message string                          `json:"message,omitempty"`
		Data    []modelSvc.SupportedVectorModel `json:"data,omitempty"`
	}
}

// AvailableModelListResponse represents an available model list response.
//
// swagger:response AvailableModelListResponse
type AvailableModelListResponse struct {
	// in: body
	Body struct {
		Success bool                 `json:"success"`
		Message string               `json:"message,omitempty"`
		Data    []modelSvc.ModelInfo `json:"data,omitempty"`
	}
}

// modelSaveResponseWrapper represents a model save response.
//
// swagger:response ModelSaveResponse
type modelSaveResponseWrapper struct {
	// in: body
	Body struct {
		Success bool   `json:"success"`
		Message string `json:"message,omitempty"`
		Data    struct {
			ID int `json:"id"`
		} `json:"data,omitempty"`
	}
}

// ModelHandler 处理模型相关的HTTP请求。
type ModelHandler struct {
	modelService *services.ModelService
}

// NewModelHandler 创建并返回一个新的模型处理器实例。
func NewModelHandler(modelService *services.ModelService) *ModelHandler {
	return &ModelHandler{
		modelService: modelService,
	}
}

// ListSupportedChatModels 返回支持的聊天模型列表。
//
// swagger:route GET /api/v1/models/supported/chat Models listSupportedChatModels
// 聊天模型列表
//
// Produces:
// - application/json
//
// Responses:
//
//	200: SupportedModelListResponse
//	500: ErrorResponse
func (h *ModelHandler) ListSupportedChatModels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	models, err := h.modelService.ListSupportedChatModels(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取支持的聊天模型列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, models)
}

// ListSupportedVectorModels 返回支持的向量模型列表。
//
// swagger:route GET /api/v1/models/supported/vector Models listSupportedVectorModels
// 向量模型列表
//
// Produces:
// - application/json
//
// Responses:
//
//	200: SupportedVectorModelListResponse
//	500: ErrorResponse
func (h *ModelHandler) ListSupportedVectorModels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	models, err := h.modelService.ListSupportedVectorModels(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取支持的向量模型列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, models)
}

// ListAvailableAllModels 返回所有可用模型的列表。
//
// swagger:route GET /api/v1/models/available/all Models listAvailableAllModels
// 全部可用模型
//
// Produces:
// - application/json
//
// Responses:
//
//	200: AvailableModelListResponse
//	500: ErrorResponse
func (h *ModelHandler) ListAvailableAllModels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	models, err := h.modelService.ListAvailableAllModels(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取全部可用模型列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, models)
}

// ListAvailableChatModels 返回可用的聊天模型列表。
//
// swagger:route GET /api/v1/models/available/chat Models listAvailableChatModels
// 可用聊天模型
//
// Produces:
// - application/json
//
// Responses:
//
//	200: AvailableModelListResponse
//	500: ErrorResponse
func (h *ModelHandler) ListAvailableChatModels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	models, err := h.modelService.ListAvailableChatModels(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取可用的对话模型列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, models)
}

// ListAvailableVectorModels 返回可用的向量模型列表。
//
// swagger:route GET /api/v1/models/available/vector Models listAvailableVectorModels
// 可用向量模型
//
// Produces:
// - application/json
//
// Responses:
//
//	200: AvailableModelListResponse
//	500: ErrorResponse
func (h *ModelHandler) ListAvailableVectorModels(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	models, err := h.modelService.ListAvailableVectorModels(ctx)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取可用的向量模型列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, models)
}

// SaveModel 保存或更新模型配置。
//
// swagger:route POST /api/v1/models Models saveModel
// 保存模型
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
//     description: 模型信息
//     required: true
//     type: ModelSaveRequest
//
// Responses:
//
//	200: ModelSaveResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ModelHandler) SaveModel(c *gin.Context) {
	var req modelSvc.ModelSaveRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	modelID, err := h.modelService.SaveModel(ctx, &req)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "保存模型失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, gin.H{"id": modelID}, "模型保存成功")
}

// DeleteModel 根据ID删除指定的模型。
//
// swagger:route DELETE /api/v1/models/{id} Models deleteModel
// 删除模型
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: 模型ID
//     required: true
//     type: integer
//
// Responses:
//
//	200: APIResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ModelHandler) DeleteModel(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "无效的模型ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.modelService.DeleteModel(ctx, id)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "删除模型失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, nil, "模型删除成功")
}

// SetModelStatus 设置模型的启用状态。
//
// swagger:route PUT /api/v1/models/{id}/status Models setModelStatus
// 设置模型状态
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
//     description: 模型ID
//     required: true
//     type: integer
//   - +name: body
//     in: body
//     description: 状态设置
//     required: true
//     type: ModelStatusEnableRequest
//
// Responses:
//
//	200: APIResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ModelHandler) SetModelStatus(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "无效的模型ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	var req models.ModelStatusEnableRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.RespondWithValidationError(c, err)
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	err = h.modelService.SetModelStatus(ctx, id, req.Enable)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "设置模型状态失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, nil, "模型状态更新成功")
}
