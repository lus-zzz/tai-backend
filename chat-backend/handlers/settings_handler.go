package handlers

import (
	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/services/interfaces"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// SettingsHandler 处理系统设置相关的HTTP请求。
type SettingsHandler struct {
	defaultSettingsService interfaces.DefaultSettingsServiceInterface
}

// NewSettingsHandler 创建并返回一个新的设置处理器实例。
func NewSettingsHandler(defaultSettingsService interfaces.DefaultSettingsServiceInterface) *SettingsHandler {
	return &SettingsHandler{
		defaultSettingsService: defaultSettingsService,
	}
}

// NewSettingsHandlerFromGlobal 使用全局服务创建设置处理器实例
func NewSettingsHandlerFromGlobal() *SettingsHandler {
	return &SettingsHandler{
		defaultSettingsService: services.GetGlobalDefaultSettingsService(),
	}
}

// GetDefaultSettings 获取系统的默认设置。
//
// swagger:route GET /settings/defaults Settings getDefaultSettings
//
// 获取默认设置
//
// 获取系统当前的默认配置信息，包括模型设置、聊天参数等
//
// Produces:
// - application/json
//
// Responses:
//
//	200: DefaultSettings
func (h *SettingsHandler) GetDefaultSettings(c *gin.Context) (interface{}, error) {
	settings := h.defaultSettingsService.GetDefaultSettings()
	return settings, nil
}

// UpdateDefaultSettings 更新系统的默认设置。
//
// swagger:route PUT /settings/defaults Settings updateDefaultSettings
//
// 更新默认设置
//
// 更新系统的默认配置信息，新的设置将应用于后续创建的对话
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
//     description: 默认设置
//     required: true
//     type: DefaultSettings
//
// Responses:
//
//	200: DefaultSettings
//	400: ResponseBody
func (h *SettingsHandler) UpdateDefaultSettings(c *gin.Context) (interface{}, error) {
	var settings models.DefaultSettings

	// 直接使用 BindJSON,只做 JSON 解析,不做 validate 验证
	if err := c.BindJSON(&settings); err != nil {
		return nil, utils.NewAPIError(utils.ErrInvalidRequest, err)
	}

	if err := h.defaultSettingsService.UpdateDefaultSettings(&settings); err != nil {
		return nil, utils.NewAPIError(utils.ErrInternalServer, err)
	}

	// 返回更新后的完整配置
	updatedSettings := h.defaultSettingsService.GetDefaultSettings()
	return updatedSettings, nil
}

// ResetDefaultSettings 将系统设置重置为默认值。
//
// swagger:route POST /settings/defaults/reset Settings resetDefaultSettings
//
// 重置默认设置
//
// 将系统配置重置为出厂默认值，恢复所有设置到初始状态
//
// Produces:
// - application/json
//
// Responses:
//
//	200: DefaultSettings
//	400: ResponseBody
func (h *SettingsHandler) ResetDefaultSettings(c *gin.Context) (interface{}, error) {
	err := h.defaultSettingsService.ResetToDefaults()
	if err != nil {
		return nil, utils.NewAPIError(utils.ErrInternalServer, err)
	}

	settings := h.defaultSettingsService.GetDefaultSettings()
	return settings, nil
}
