package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"time"

	"chat-backend/models"
	"chat-backend/services"
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// errorResponseWrapper represents an error response.
//
// swagger:response ErrorResponse
type errorResponseWrapper struct {
	// in: body
	Body models.APIResponse
}

// conversationResponseWrapper represents a conversation response.
//
// swagger:response ConversationResponse
type ConversationResponse struct {
	// in: body
	Body struct {
		Success bool                `json:"success"`
		Message string              `json:"message,omitempty"`
		Data    models.Conversation `json:"data,omitempty"`
	}
}

// conversationListResponseWrapper represents a conversation list response.
//
// swagger:response ConversationListResponse
type ConversationListResponse struct {
	// in: body
	Body struct {
		Success bool                            `json:"success"`
		Message string                          `json:"message,omitempty"`
		Data    models.ConversationListResponse `json:"data,omitempty"`
	}
}

// conversationSettingsResponseWrapper represents a conversation settings response.
//
// swagger:response ConversationSettingsResponse
type ConversationSettingsResponse struct {
	// in: body
	Body struct {
		Success bool                        `json:"success"`
		Message string                      `json:"message,omitempty"`
		Data    models.ConversationSettings `json:"data,omitempty"`
	}
}

// conversationHistoryResponseWrapper represents a conversation history response.
//
// swagger:response ConversationHistoryResponse
type ConversationHistoryResponse struct {
	// in: body
	Body struct {
		Success bool                               `json:"success"`
		Message string                             `json:"message,omitempty"`
		Data    models.ConversationHistoryResponse `json:"data,omitempty"`
	}
}

// ChatHandler 处理聊天相关的HTTP请求。
type ChatHandler struct {
	chatService            *services.ChatService
	defaultSettingsService *services.DefaultSettingsService
}

// NewChatHandler 创建并返回一个新的聊天处理器实例。
func NewChatHandler(chatService *services.ChatService, defaultSettingsService *services.DefaultSettingsService) *ChatHandler {
	return &ChatHandler{
		chatService:            chatService,
		defaultSettingsService: defaultSettingsService,
	}
}

// CreateConversation 处理创建新对话的HTTP请求。
//
// 该方法从请求体中解析 ConversationSettings,调用聊天服务创建对话,
// 并返回创建的对话信息(包含ID和配置)。
//
// swagger:route POST /api/v1/chat/conversations Chat createConversation
//
// Creates a new conversation with the given settings.
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
//     description: Conversation settings
//     required: true
//     type: ConversationSettings
//
// Responses:
//
//	200: ConversationResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ChatHandler) CreateConversation(c *gin.Context) {

	var settings models.ConversationSettings
	// 使用 BindJSON 避免 validate 标签验证
	if err := c.BindJSON(&settings); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conversation, err := h.chatService.CreateConversation(ctx, &settings)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrConversationCreate, "创建对话失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, conversation, "对话创建成功")
}

// GetConversations 返回分页的对话列表。
//
// swagger:route GET /api/v1/chat/conversations Chat getConversations
//
// Returns a paginated list of conversations.
//
// ---
// produces:
// - application/json
// parameters:
//   - +name: page
//     in: query
//     description: Page number
//     required: false
//     type: integer
//     default: 1
//   - +name: page_size
//     in: query
//     description: Number of items per page
//     required: false
//     type: integer
//     default: 20
//
// responses:
//
//	200: ConversationListResponse
//	500: ErrorResponse
func (h *ChatHandler) GetConversations(c *gin.Context) {
	page := 1
	pageSize := 20

	if p := c.Query("page"); p != "" {
		if pInt, err := strconv.Atoi(p); err == nil && pInt > 0 {
			page = pInt
		}
	}

	if ps := c.Query("page_size"); ps != "" {
		if psInt, err := strconv.Atoi(ps); err == nil && psInt > 0 && psInt <= 100 {
			pageSize = psInt
		}
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	conversations, err := h.chatService.GetConversations(ctx, page, pageSize)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取对话列表失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, conversations)
}

// DeleteConversation 根据ID删除指定的对话。
//
// swagger:route DELETE /api/v1/chat/conversations/{id} Chat deleteConversation
//
// Deletes a conversation by its ID.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Conversation ID
//     required: true
//     type: string
//
// Responses:
//
//	200: APIResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ChatHandler) DeleteConversation(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少对话ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.chatService.DeleteConversation(ctx, conversationID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrConversationNotFound, "删除对话失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, nil, "对话删除成功")
}

// SendMessage 向对话发送消息并返回SSE流式响应。
//
// swagger:route POST /api/v1/chat/messages Chat sendMessage
//
// Sends a message to a conversation and returns an SSE stream.
//
// Consumes:
// - application/json
//
// Produces:
// - text/event-stream
//
// Parameters:
//   - +name: body
//     in: body
//     description: Chat message request
//     required: true
//     type: ChatRequest
//
// Responses:
//
//	200:
//	  description: SSE stream
//	400:
//	  description: Invalid request
//	500:
//	  description: Internal server error
func (h *ChatHandler) SendMessage(c *gin.Context) {
	var req models.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "请求数据解析失败"})
		return
	}

	if req.Content == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "消息内容不能为空"})
		return
	}

	if req.SessionID == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "会话ID不能为空"})
		return
	}

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 创建事件通道
	eventChan := make(chan models.SSEChatEvent, 10)

	// 使用请求的上下文，当客户端断开连接时自动取消
	ctx := c.Request.Context()

	// 在goroutine中发送消息
	go func() {
		defer func() {
			if r := recover(); r != nil {
				utils.ErrorWith("SendMessage panic", "error", r)
			}
		}()

		if err := h.chatService.SendMessage(ctx, &req, eventChan); err != nil {
			utils.ErrorWith("流式发送消息失败", "error", err)
		}
	}()

	// 发送SSE事件
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "不支持流式响应"})
		return
	}

	for {
		select {
		case event, ok := <-eventChan:
			if !ok {
				return
			}

			eventData, err := json.Marshal(event)
			if err != nil {
				utils.ErrorWith("序列化事件失败", "error", err)
				continue
			}

			// 使用SDK中的事件类型，如果没有则根据状态推断
			eventType := event.EventType
			if eventType == "" {
				// 兼容逻辑：如果SDK没有提供事件类型，则根据状态推断
				eventType = "resp_increment" // 默认为增量响应
				if event.Pending {
					eventType = "resp_splash" // 等待状态
				}
				if event.FinishReason != "" {
					eventType = "resp_finish" // 完成状态
				}
			}

			// 按照SSE格式输出: event字段 + data字段
			fmt.Fprintf(c.Writer, "event:%s\ndata: %s\n\n", eventType, eventData)
			flusher.Flush()

			// 根据 FinishReason 或 Error 判断是否结束
			if event.FinishReason != "" || event.Error {
				return
			}

		case <-ctx.Done():
			return
		}
	}
}

// GetConversationSettings 获取指定对话的设置信息。
//
// swagger:route GET /api/v1/chat/conversations/{id}/settings Chat getConversationSettings
//
// Gets the settings for a specific conversation.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Conversation ID
//     required: true
//     type: string
//
// Responses:
//
//	200: ConversationSettingsResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ChatHandler) GetConversationSettings(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少对话ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	settings, err := h.chatService.GetConversationSettings(ctx, conversationID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取对话设置失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, settings)
}

// UpdateConversationSettings 更新指定对话的设置信息。
//
// swagger:route PUT /api/v1/chat/conversations/{id}/settings Chat updateConversationSettings
//
// Updates the settings for a specific conversation.
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
//     description: Conversation ID
//     required: true
//     type: string
//   - +name: body
//     in: body
//     description: Conversation settings
//     required: true
//     type: ConversationSettings
//
// Responses:
//
//	200: ConversationSettingsResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ChatHandler) UpdateConversationSettings(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "缺少对话ID", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	var settings models.ConversationSettings
	// 直接使用 BindJSON,只做 JSON 解析,不做 validate 验证
	if err := c.BindJSON(&settings); err != nil {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "JSON 解析失败", http.StatusBadRequest).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err := h.chatService.UpdateConversationSettings(ctx, conversationID, &settings)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrSettingsUpdate, "更新对话设置失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, settings, "对话设置更新成功")
}

// GetConversationHistory 获取指定对话的消息历史记录。
//
// swagger:route GET /api/v1/chat/conversations/{id}/history Chat getConversationHistory
//
// Gets the message history for a specific conversation.
//
// Produces:
// - application/json
//
// Parameters:
//   - +name: id
//     in: path
//     description: Conversation ID
//     required: true
//     type: string
//
// Responses:
//
//	200: ConversationHistoryResponse
//	400: ErrorResponse
//	500: ErrorResponse
func (h *ChatHandler) GetConversationHistory(c *gin.Context) {
	conversationID := c.Param("id")
	if conversationID == "" {
		apiErr := utils.NewAPIError(utils.ErrInvalidRequest, "对话ID不能为空", http.StatusBadRequest)
		utils.RespondWithError(c, apiErr)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	history, err := h.chatService.GetConversationHistory(ctx, conversationID)
	if err != nil {
		apiErr := utils.NewAPIError(utils.ErrInternalServer, "获取对话历史失败", http.StatusInternalServerError).WithCause(err)
		utils.RespondWithError(c, apiErr)
		return
	}

	utils.RespondWithSuccess(c, history, "获取对话历史成功")
}
