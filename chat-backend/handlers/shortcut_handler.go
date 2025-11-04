package handlers

import (
	"net/http"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// ShortcutHandler 处理快捷方式相关的HTTP请求。
type ShortcutHandler struct {
	shortcutService *services.ShortcutService
}

// NewShortcutHandler 创建并返回一个新的快捷方式处理器实例。
func NewShortcutHandler(shortcutService *services.ShortcutService) *ShortcutHandler {
	return &ShortcutHandler{
		shortcutService: shortcutService,
	}
}

// RecommendSettings 根据用户输入推荐设置。
//
// swagger:route POST /shortcut/recommend Shortcut recommendSettings
//
// 快捷方式推荐
//
// 根据用户输入和推荐数目，得到推荐的设置列表
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
//     description: 用户输入和推荐数量
//     required: true
//     type: ChatInput
//
// Responses:
//
//	200: RecommendDataSuccessResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ShortcutHandler) RecommendSettings(c *gin.Context) {
	var input models.ChatInput

	// 绑定JSON数据
	if err := c.ShouldBindJSON(&input); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 验证输入参数
	if input.UserInput == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "用户输入不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	if input.RecommendNum <= 0 {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "推荐数量必须大于0", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	// 调用快捷方式服务进行API转发
	recommendData, err := h.shortcutService.RecommendSettings(&input)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取推荐设置失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, recommendData, "推荐设置获取成功")
}

// GetSupportedSettings 获取所有支持的设置。
//
// swagger:route POST /shortcut/supportedSetting Shortcut getSupportedSettings
//
// 获取支持的设置列表
//
// 获取当前支持的设置列表
//
// Produces:
// - application/json
//
// Responses:
//
//	200: SettingNameSuccessResponse
//	500: ErrorResponse
func (h *ShortcutHandler) GetSupportedSettings(c *gin.Context) {
	// 调用快捷方式服务进行API转发
	supportedSettings, err := h.shortcutService.GetSupportedSettings()
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取支持的设置列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, supportedSettings, "支持的设置列表获取成功")
}
