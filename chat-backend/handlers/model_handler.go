package handlers

import (
	"context"
	"net/http"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/services/interfaces"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// ModelHandler 处理模型相关的HTTP请求。
type ModelHandler struct {
	modelService interfaces.ModelServiceInterface
}

// NewModelHandler 创建并返回一个新的模型处理器实例。
func NewModelHandler(modelService interfaces.ModelServiceInterface) *ModelHandler {
	return &ModelHandler{
		modelService: modelService,
	}
}

// NewModelHandlerFromGlobal 使用全局服务创建模型处理器实例
func NewModelHandlerFromGlobal() *ModelHandler {
	return &ModelHandler{
		modelService: services.GetGlobalModelService(),
	}
}

// ListSupportedChatModels 返回支持的聊天模型列表。
//
// swagger:route GET /models/supported/chat Models listSupportedChatModels
//
// 获取支持的聊天模型
//
// 获取系统支持的所有聊天模型列表，包括模型名称、类型等信息
//
// Produces:
// - application/json
//
// Responses:
//
//	200: SupportedChatModelListSuccessResponse
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
// swagger:route GET /models/supported/vector Models listSupportedVectorModels
//
// 获取支持的向量模型
//
// 获取系统支持的所有向量模型列表，用于知识库向量化处理
//
// Produces:
// - application/json
//
// Responses:
//
//	200: SupportedVectorModelListSuccessResponse
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
// swagger:route GET /models/available/all Models listAvailableAllModels
//
// 获取所有可用模型
//
// 获取当前系统中所有可用的模型列表，包括聊天模型和向量模型
//
// Produces:
// - application/json
//
// Responses:
//
//	200: ModelListSuccessResponse
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
// swagger:route GET /models/available/chat Models listAvailableChatModels
//
// 获取可用聊天模型
//
// 获取当前系统中可用的聊天模型列表，仅包含已配置且可用的聊天模型
//
// Produces:
// - application/json
//
// Responses:
//
//	200: ModelListSuccessResponse
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
// swagger:route GET /models/available/vector Models listAvailableVectorModels
//
// 获取可用向量模型
//
// 获取当前系统中可用的向量模型列表，用于知识库文档向量化
//
// Produces:
// - application/json
//
// Responses:
//
//	200: ModelListSuccessResponse
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
// swagger:route POST /models Models saveModel
//
// 保存模型配置
//
// 保存或更新模型的配置信息，包括API密钥、端点等设置
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
//	200: ModelSaveSuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ModelHandler) SaveModel(c *gin.Context) {
	var req models.ModelSaveRequest
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

	utils.RespondWithSuccess(c, models.ModelID{ID: modelID}, "模型保存成功")
}

// DeleteModel 根据ID删除指定的模型。
//
// swagger:route DELETE /models/{id} Models deleteModel
//
// 删除模型配置
//
// 根据模型ID删除指定的模型配置信息
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
//	200: EmptySuccessResponse
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
// swagger:route PUT /models/{id}/status Models setModelStatus
//
// 设置模型状态
//
// 启用或禁用指定的模型，控制模型是否可用于对话
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
//	200: EmptySuccessResponse
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
