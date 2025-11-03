package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// defaultSettingsResponseWrapper represents a default settings response.
//
// swagger:response DefaultSettingsResponse
type defaultSettingsResponseWrapper struct {
	// in: body
	Body struct {
		Success bool                   `json:"success"`
		Message string                 `json:"message,omitempty"`
		Data    models.DefaultSettings `json:"data,omitempty"`
	}
}

// SettingsHandler 处理系统设置相关的HTTP请求。
type SettingsHandler struct {
	defaultSettingsService *services.DefaultSettingsService
}

// NewSettingsHandler 创建并返回一个新的设置处理器实例。
func NewSettingsHandler(defaultSettingsService *services.DefaultSettingsService) *SettingsHandler {
	return &SettingsHandler{
		defaultSettingsService: defaultSettingsService,
	}
}

// GetDefaultSettings 获取系统的默认设置。
//
// swagger:route GET /api/v1/settings/defaults Settings getDefaultSettings
//
// Gets the default settings.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: DefaultSettingsResponse
func (h *SettingsHandler) GetDefaultSettings(c *gin.Context) {
	settings := h.defaultSettingsService.GetDefaultSettings()
	utils.RespondWithSuccess(c, settings, "获取默认配置成功")
}

// UpdateDefaultSettings 更新系统的默认设置。
//
// swagger:route PUT /api/v1/settings/defaults Settings updateDefaultSettings
//
// Updates the default settings.
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
//     description: Default settings
//     required: true
//     type: DefaultSettings
//
// Responses:
//
//	200: DefaultSettingsResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *SettingsHandler) UpdateDefaultSettings(c *gin.Context) {
	var settings models.DefaultSettings

	// 直接使用 BindJSON,只做 JSON 解析,不做 validate 验证
	if err := c.BindJSON(&settings); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	if err := h.defaultSettingsService.UpdateDefaultSettings(&settings); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "更新默认配置失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 返回更新后的完整配置
	updatedSettings := h.defaultSettingsService.GetDefaultSettings()
	utils.RespondWithSuccess(c, updatedSettings, "默认配置已更新")
}

// ResetDefaultSettings 将系统设置重置为默认值。
//
// swagger:route POST /api/v1/settings/defaults/reset Settings resetDefaultSettings
//
// Resets the default settings.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: DefaultSettingsResponse
//	500: ErrorResponse
func (h *SettingsHandler) ResetDefaultSettings(c *gin.Context) {
	err := h.defaultSettingsService.ResetToDefaults()
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "重置默认配置失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	settings := h.defaultSettingsService.GetDefaultSettings()
	utils.RespondWithSuccess(c, settings, "默认配置已重置")
}
