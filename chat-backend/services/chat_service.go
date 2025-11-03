package services

import (
	"context"
	"fmt"
	"strconv"

	"chat-backend/models"
	"chat-backend/utils"
	"flowy-sdk"
	agentSvc "flowy-sdk/services/agent"
)

// ChatService 聊天服务
type ChatService struct {
	sdk                    *flowy.SDK
	defaultSettingsService *DefaultSettingsService
}

// NewChatService 创建聊天服务
func NewChatService(sdk *flowy.SDK, defaultSettingsService *DefaultSettingsService) *ChatService {
	return &ChatService{
		sdk:                    sdk,
		defaultSettingsService: defaultSettingsService,
	}
}

// CreateConversation 创建对话
func (s *ChatService) CreateConversation(ctx context.Context, settings *models.ConversationSettings) (*models.Conversation, error) {
	// 如果没有提供设置，使用持久化的默认配置
	if settings == nil {
		defaultSettings := s.defaultSettingsService.GetDefaultSettings()
		settings = &models.ConversationSettings{
			Name:             "新对话",
			ModelID:          defaultSettings.Models.ChatModelID,
			Temperature:      defaultSettings.Conversation.Temperature,
			KnowledgeBaseIDs: []int{},
			TopP:             defaultSettings.Conversation.TopP,
			FrequencyPenalty: defaultSettings.Conversation.FrequencyPenalty,
			PresencePenalty:  defaultSettings.Conversation.PresencePenalty,
			ResponseType:     defaultSettings.Conversation.ResponseType,
			Stream:           defaultSettings.Conversation.Stream,
			ContextLimit:     defaultSettings.Conversation.ContextLimit,
		}
		utils.LogInfo("使用默认配置创建对话: 新对话")
	} else {
		utils.LogInfo("创建对话: %s", settings.Name)
	}

	// 步骤1: 创建一个新的Agent
	createAgentReq := &agentSvc.CreateAgentRequest{
		Name:   settings.Name,
		Desc:   settings.Desc,
		Type:   0,  // 默认类型为0 (多轮对话)
		Avatar: "", // 默认无头像
	}

	agentID, err := s.sdk.Agent.CreateAgent(ctx, createAgentReq)
	if err != nil {
		return nil, fmt.Errorf("创建Agent失败: %w", err)
	}

	utils.InfoWith("Agent创建成功", "agent_id", agentID, "title", settings.Name)

	// 步骤2: 为这个Agent创建一个配置，使用传入的设置参数
	// 根据模型ID获取模型名称
	modelName, err := s.getModelNameByID(ctx, settings.ModelID)
	if err != nil {
		// 如果获取模型名称失败,清理已创建的Agent
		_ = s.sdk.Agent.DeleteAgent(ctx, agentID)
		return nil, fmt.Errorf("获取模型名称失败: %w", err)
	}

	utils.InfoWith("使用模型", "model_id", settings.ModelID, "model_name", modelName)

	// 处理知识库配置
	knowledgeIDs := []int{}
	enableKnowledge := false
	if len(settings.KnowledgeBaseIDs) > 0 {
		enableKnowledge = true
		for _, kbID := range settings.KnowledgeBaseIDs {
			knowledgeIDs = append(knowledgeIDs, kbID)
		}
	}

	// 使用工厂函数创建默认配置
	saveConfigReq := agentSvc.NewDefaultSettingConfig(agentID, settings.Name)

	// 覆盖需要自定义的配置
	saveConfigReq.Chat.Stream = settings.Stream
	saveConfigReq.Chat.Model.ID = settings.ModelID
	// saveConfigReq.Chat.Model.Model = modelName
	saveConfigReq.Chat.Model.Temperature = settings.Temperature
	saveConfigReq.Chat.Model.TopP = settings.TopP
	saveConfigReq.Chat.Model.PresencePenalty = settings.PresencePenalty
	saveConfigReq.Chat.Model.FrequencyPenalty = settings.FrequencyPenalty
	saveConfigReq.Chat.Model.ResponseType = settings.ResponseType
	saveConfigReq.Chat.ContextLimit = settings.ContextLimit

	// 设置系统提示词
	saveConfigReq.Chat.Prompt.Prompts = []agentSvc.Prompt{
		{
			Role: 0,
			Text: "你是一个有帮助的AI助手。",
		},
	}

	// 设置知识库配置
	saveConfigReq.Chat.Plugin.Knowledge.Enable = enableKnowledge
	saveConfigReq.Chat.Plugin.Knowledge.Knowledges = knowledgeIDs

	// 输出配置 saveConfigReq 的json
	// jsonData, _ := json.MarshalIndent(saveConfigReq, "", "  ")
	// fmt.Printf("保存配置: %s\n", jsonData)

	settingID, err := s.sdk.Agent.SaveConfig(ctx, saveConfigReq)
	if err != nil {
		// 如果配置创建失败,需要清理已创建的Agent
		_ = s.sdk.Agent.DeleteAgent(ctx, agentID)
		return nil, fmt.Errorf("创建配置失败: %w", err)
	}

	utils.InfoWith("配置创建成功", "setting_id", settingID, "agent_id", agentID)

	// 步骤3: 创建会话
	createSessionReq := &agentSvc.CreateSessionRequest{
		SettingID:  settingID,
		PromptVars: []agentSvc.PromptVar{},
	}

	session, err := s.sdk.Agent.CreateSession(ctx, createSessionReq)
	if err != nil {
		// 如果会话创建失败,需要清理已创建的配置和Agent
		_ = s.sdk.Agent.DeleteConfig(ctx, settingID)
		_ = s.sdk.Agent.DeleteAgent(ctx, agentID)
		return nil, fmt.Errorf("创建会话失败: %w", err)
	}

	utils.InfoWith("会话创建成功", "session_id", session.ID, "setting_id", settingID, "agent_id", agentID)

	// 转换为Conversation模型
	return &models.Conversation{
		ID:                   session.ID,
		ConversationSettings: *settings,
	}, nil
}

// SendMessage 发送消息并返回SSE流
func (s *ChatService) SendMessage(ctx context.Context, req *models.ChatRequest, eventChan chan<- models.SSEChatEvent) error {
	utils.InfoWith("开始流式发送消息", "session_id", req.SessionID, "content", req.Content)

	defer close(eventChan)

	// 调用 SDK 的流式对话接口，直接传递事件通道（类型别名，无需转换）
	err := s.sdk.Agent.ChatAsync(ctx, &req.AsyncChatRequest, eventChan)
	if err != nil {
		utils.ErrorWith("流式发送消息失败", "session_id", req.SessionID, "error", err)
		return err
	}

	utils.InfoWith("流式消息发送完成", "session_id", req.SessionID)
	return nil
}

// ListConversations 获取对话列表
func (s *ChatService) ListConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error) {
	utils.LogInfo("获取对话列表")

	// 获取所有Agent的会话
	agentList, err := s.sdk.Agent.ListAgentsByPage(ctx, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取Agent列表失败: %w", err)
	}

	var allConversations []models.Conversation

	// 遍历agents获取会话
	for _, agent := range agentList.Records {

		// 获取Agent的配置
		configs, err := s.sdk.Agent.ListConfigs(ctx, agent.ID)
		if err != nil {
			continue // 跳过错误的agent
		}

		// 获取每个配置的会话列表
		for _, config := range configs {
			sessions, err := s.sdk.Agent.ListSessions(ctx, config.ID, nil)
			if err != nil {
				continue
			}

			for _, session := range sessions {
				// 使用工厂函数创建默认设置，然后从session和config中提取真实的值
				settings := models.NewDefaultConversationSettings()
				settings.Name = agent.Name // 使用会话名称
				settings.Desc = agent.Desc // 使用Agent描述

				// 从配置中提取真实的值
				if config.Chat != nil {
					settings.Temperature = config.Chat.Model.Temperature
					settings.TopP = config.Chat.Model.TopP
					settings.PresencePenalty = config.Chat.Model.PresencePenalty
					settings.FrequencyPenalty = config.Chat.Model.FrequencyPenalty
					settings.ResponseType = config.Chat.Model.ResponseType
					settings.Stream = config.Chat.Stream
					settings.ContextLimit = config.Chat.ContextLimit
					settings.ModelID = config.Chat.Model.ID

					// 提取知识库配置
					if config.Chat.Plugin.Knowledge.Enable && len(config.Chat.Plugin.Knowledge.Knowledges) > 0 {
						settings.KnowledgeBaseIDs = []int{}
						for _, kbID := range config.Chat.Plugin.Knowledge.Knowledges {
							settings.KnowledgeBaseIDs = append(settings.KnowledgeBaseIDs, kbID)
						}
					}
				}

				conv := models.Conversation{
					ID:                   session.ID,
					ConversationSettings: *settings,
				}
				allConversations = append(allConversations, conv)
			}
		}
	}

	// 简单分页
	start := (page - 1) * pageSize
	end := start + pageSize
	if start > len(allConversations) {
		start = len(allConversations)
	}
	if end > len(allConversations) {
		end = len(allConversations)
	}

	pagedConversations := allConversations[start:end]

	return &models.ConversationListResponse{
		Conversations: pagedConversations,
		Total:         len(allConversations),
		Page:          page,
		PageSize:      pageSize,
	}, nil
}

// DeleteConversation 删除对话
func (s *ChatService) DeleteConversation(ctx context.Context, conversationID string) error {
	utils.LogInfo("删除对话: %s", conversationID)

	sessionID, err := strconv.Atoi(conversationID)
	if err != nil {
		return fmt.Errorf("无效的会话ID: %w", err)
	}

	// 步骤1: 查找会话对应的Agent和Setting
	sessionInfo, err := s.findSessionInfo(ctx, sessionID)
	if err != nil {
		return err
	}

	utils.InfoWith("找到会话关联", "session_id", sessionID, "setting_id", sessionInfo.SettingID, "agent_id", sessionInfo.AgentID)

	// 步骤2: 删除会话
	err = s.sdk.Agent.DeleteSession(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("删除会话失败: %w", err)
	}

	utils.InfoWith("会话已删除", "session_id", sessionID)

	// 步骤3: 删除配置
	err = s.sdk.Agent.DeleteConfig(ctx, sessionInfo.SettingID)
	if err != nil {
		utils.ErrorWith("删除配置失败", "setting_id", sessionInfo.SettingID, "error", err)
		// 继续删除Agent，即使配置删除失败
	} else {
		utils.InfoWith("配置已删除", "setting_id", sessionInfo.SettingID)
	}

	// 步骤4: 删除Agent
	err = s.sdk.Agent.DeleteAgent(ctx, sessionInfo.AgentID)
	if err != nil {
		utils.ErrorWith("删除Agent失败", "agent_id", sessionInfo.AgentID, "error", err)
		return fmt.Errorf("删除Agent失败: %w", err)
	}

	utils.InfoWith("Agent已删除", "agent_id", sessionInfo.AgentID)
	utils.LogInfo("对话及其关联的配置和Agent已全部删除: %s", conversationID)

	return nil
}

// GetConversations 获取对话列表 (别名方法)
func (s *ChatService) GetConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error) {
	return s.ListConversations(ctx, page, pageSize)
}

// UpdateConversationSettings 更新对话设置
func (s *ChatService) UpdateConversationSettings(ctx context.Context, conversationID string, settings *models.ConversationSettings) error {
	utils.LogInfo("更新对话设置: %s", conversationID)

	// 将conversationID转换为int
	sessionID, err := strconv.Atoi(conversationID)
	if err != nil {
		return fmt.Errorf("无效的会话ID: %w", err)
	}

	// 查找会话对应的配置信息
	sessionInfo, err := s.findSessionInfo(ctx, sessionID)
	if err != nil {
		return fmt.Errorf("查找会话配置失败: %w", err)
	}

	// 如果 Name 或 Desc 有变化，则更新 Agent
	if settings.Name != "" && settings.Name != sessionInfo.AgentName || settings.Desc != sessionInfo.AgentDesc {
		updateAgentReq := &agentSvc.UpdateAgentRequest{
			ID:     sessionInfo.AgentID,
			Name:   settings.Name,
			Desc:   settings.Desc,
			Avatar: "", // 保持现有头像
		}

		// 如果没有提供名称，使用原有名称
		if settings.Name == "" {
			updateAgentReq.Name = sessionInfo.AgentName
		}

		_, err := s.sdk.Agent.UpdateAgent(ctx, updateAgentReq)
		if err != nil {
			utils.ErrorWith("更新Agent失败", "agent_id", sessionInfo.AgentID, "error", err)
			// 不中断流程，继续更新配置
		} else {
			utils.InfoWith("Agent已更新", "agent_id", sessionInfo.AgentID, "name", updateAgentReq.Name, "desc", updateAgentReq.Desc)
		}
	}

	// 构建更新的配置请求，基于现有配置
	updateReq := &agentSvc.SaveConfigRequest{
		AgentID: sessionInfo.AgentID,
		ID:      sessionInfo.SettingID,
		Name:    sessionInfo.Config.Name,
	}

	// 如果有Chat配置，则更新它
	if sessionInfo.Config.Chat != nil {
		chatConfig := sessionInfo.Config.Chat

		// 更新模型配置
		chatConfig.Model.Temperature = settings.Temperature
		chatConfig.Model.TopP = settings.TopP
		chatConfig.Model.FrequencyPenalty = settings.FrequencyPenalty
		chatConfig.Model.PresencePenalty = settings.PresencePenalty
		chatConfig.Model.ResponseType = settings.ResponseType
		chatConfig.Model.ID = settings.ModelID

		// 更新流式输出设置
		chatConfig.Stream = settings.Stream

		// 更新上下文限制
		if settings.ContextLimit > 0 {
			chatConfig.ContextLimit = settings.ContextLimit
		}

		// 更新知识库配置
		if len(settings.KnowledgeBaseIDs) > 0 {
			chatConfig.Plugin.Knowledge.Enable = true
			chatConfig.Plugin.Knowledge.Knowledges = []int{}
			for _, kbID := range settings.KnowledgeBaseIDs {
				chatConfig.Plugin.Knowledge.Knowledges = append(chatConfig.Plugin.Knowledge.Knowledges, kbID)
			}
		} else {
			chatConfig.Plugin.Knowledge.Enable = false
			chatConfig.Plugin.Knowledge.Knowledges = []int{}
		}

		updateReq.Chat = chatConfig
	} else {
		// 如果没有Chat配置，使用工厂函数创建默认配置
		// 根据模型ID获取模型名称
		// modelName := ""
		// if settings.ModelID > 0 {
		// 	name, err := s.getModelNameByID(ctx, settings.ModelID)
		// 	if err != nil {
		// 		utils.ErrorWith("获取模型名称失败", "model_id", settings.ModelID, "error", err)
		// 	} else {
		// 		modelName = name
		// 	}
		// }

		// 使用工厂函数创建临时配置以获取默认的Chat配置
		tempConfig := agentSvc.NewDefaultSettingConfig(sessionInfo.AgentID, sessionInfo.Config.Name)
		updateReq.Chat = tempConfig.Chat

		// 覆盖需要自定义的配置
		updateReq.Chat.Stream = settings.Stream
		updateReq.Chat.Model.ID = settings.ModelID
		// updateReq.Chat.Model.Model = modelName
		updateReq.Chat.Model.Temperature = settings.Temperature
		updateReq.Chat.Model.TopP = settings.TopP
		updateReq.Chat.Model.PresencePenalty = settings.PresencePenalty
		updateReq.Chat.Model.FrequencyPenalty = settings.FrequencyPenalty
		updateReq.Chat.Model.ResponseType = settings.ResponseType
		updateReq.Chat.ContextLimit = settings.ContextLimit

		// 设置知识库配置
		if len(settings.KnowledgeBaseIDs) > 0 {
			for _, kbID := range settings.KnowledgeBaseIDs {
				updateReq.Chat.Plugin.Knowledge.Knowledges = append(updateReq.Chat.Plugin.Knowledge.Knowledges, kbID)
			}
		}
	}

	// 保留其他类型的配置
	if sessionInfo.Config.Extract != nil {
		updateReq.Extract = sessionInfo.Config.Extract
	}
	if sessionInfo.Config.Classify != nil {
		updateReq.Classify = sessionInfo.Config.Classify
	}
	if sessionInfo.Config.FormCollect != nil {
		updateReq.FormCollect = sessionInfo.Config.FormCollect
	}
	if sessionInfo.Config.Intention != nil {
		updateReq.Intention = sessionInfo.Config.Intention
	}
	if sessionInfo.Config.IntentRouter != nil {
		updateReq.IntentRouter = sessionInfo.Config.IntentRouter
	}

	// 保存配置
	_, err = s.sdk.Agent.SaveConfig(ctx, updateReq)
	if err != nil {
		return fmt.Errorf("保存配置失败: %w", err)
	}

	utils.InfoWith("对话设置已更新", "conversation_id", conversationID, "setting_id", sessionInfo.SettingID)
	return nil
}

// GetConversationSettings 获取对话设置
func (s *ChatService) GetConversationSettings(ctx context.Context, conversationID string) (*models.ConversationSettings, error) {
	utils.LogInfo("获取对话设置: %s", conversationID)

	// 将conversationID转换为int
	sessionID, err := strconv.Atoi(conversationID)
	if err != nil {
		return nil, fmt.Errorf("无效的会话ID: %w", err)
	}

	// 查找会话对应的配置信息
	sessionInfo, err := s.findSessionInfo(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("查找会话配置失败: %w", err)
	}

	// 将SDK的配置转换为ConversationSettings
	settings := models.NewDefaultConversationSettings()

	// 从配置中提取设置
	if sessionInfo.Config.Chat != nil {
		// 提取模型配置 - 根据模型名称获取模型ID
		// if sessionInfo.Config.Chat.Model.Model != "" {
		// 	modelID, err := s.getModelIDByName(ctx, sessionInfo.Config.Chat.Model.Model)
		// 	if err != nil {
		// 		utils.ErrorWith("根据模型名称获取模型ID失败", "model_name", sessionInfo.Config.Chat.Model.Model, "error", err)
		// 		settings.ModelID = 1 // 使用默认模型ID
		// 	} else {
		// 		settings.ModelID = modelID
		// 	}
		// } else {
		// 	settings.ModelID = 1 // 默认模型ID
		// }
		settings.ModelID = sessionInfo.ModelID
		settings.Name = sessionInfo.AgentName
		settings.Desc = sessionInfo.AgentDesc
		settings.Temperature = sessionInfo.Config.Chat.Model.Temperature
		settings.TopP = sessionInfo.Config.Chat.Model.TopP
		settings.FrequencyPenalty = sessionInfo.Config.Chat.Model.FrequencyPenalty
		settings.PresencePenalty = sessionInfo.Config.Chat.Model.PresencePenalty
		settings.ResponseType = sessionInfo.Config.Chat.Model.ResponseType
		settings.Stream = sessionInfo.Config.Chat.Stream

		// 提取上下文限制
		if sessionInfo.Config.Chat.ContextLimit > 0 {
			settings.ContextLimit = sessionInfo.Config.Chat.ContextLimit
		}

		// 提取知识库配置
		if sessionInfo.Config.Chat.Plugin.Knowledge.Enable && len(sessionInfo.Config.Chat.Plugin.Knowledge.Knowledges) > 0 {
			settings.KnowledgeBaseIDs = []int{}
			for _, kbID := range sessionInfo.Config.Chat.Plugin.Knowledge.Knowledges {
				settings.KnowledgeBaseIDs = append(settings.KnowledgeBaseIDs, kbID)
			}
		}
	}

	utils.InfoWith("成功获取对话设置", "conversation_id", conversationID)
	return settings, nil
}

// SessionInfo 会话信息
type SessionInfo struct {
	AgentID   int
	SettingID int
	Config    *agentSvc.SettingConfig
	AgentName string
	AgentDesc string
	ModelID   int
}

// findSessionInfo 查找会话对应的AgentID、SettingID、配置信息以及Agent的名称和描述
func (s *ChatService) findSessionInfo(ctx context.Context, sessionID int) (*SessionInfo, error) {
	agentList, err := s.sdk.Agent.ListAgentsByPage(ctx, 1, 100)
	if err != nil {
		return nil, fmt.Errorf("获取Agent列表失败: %w", err)
	}

	// 遍历所有Agent的配置和会话来找到对应的session
	for _, agent := range agentList.Records {

		configs, err := s.sdk.Agent.ListConfigs(ctx, agent.ID)
		if err != nil {
			continue
		}

		for _, cfg := range configs {

			sessions, err := s.sdk.Agent.ListSessions(ctx, cfg.ID, nil)
			if err != nil {
				continue
			}

			for _, session := range sessions {
				if session.ID == sessionID {
					modelID := 0
					if cfg.Chat != nil && cfg.Chat.Model.ID > 0 {
						modelID = cfg.Chat.Model.ID
					}
					return &SessionInfo{
						AgentID:   agent.ID,
						SettingID: cfg.ID,
						Config:    &cfg,
						AgentName: agent.Name,
						AgentDesc: agent.Desc,
						ModelID:   modelID,
					}, nil
				}
			}
		}
	}

	return nil, fmt.Errorf("未找到会话 %d 对应的配置和Agent", sessionID)
}

// getModelNameByID 根据模型ID获取模型名称
func (s *ChatService) getModelNameByID(ctx context.Context, modelID int) (string, error) {
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
			return model.Symbol, nil // 返回模型的Symbol字段作为模型名称
		}
	}

	return "", fmt.Errorf("未找到ID为 %d 的模型", modelID)
}

// getModelIDByName 根据模型名称获取模型ID
func (s *ChatService) getModelIDByName(ctx context.Context, modelName string) (int, error) {
	if modelName == "" {
		return 0, fmt.Errorf("模型名称不能为空")
	}

	// 获取所有可用模型
	models, err := s.sdk.Model.ListAvailableAllModels(ctx)
	if err != nil {
		return 0, fmt.Errorf("获取模型列表失败: %w", err)
	}

	// 查找匹配的模型
	for _, model := range models {
		if model.Symbol == modelName || model.Name == modelName {
			return model.ID, nil
		}
	}

	// 如果找不到，返回默认模型ID
	utils.WarnWith("未找到匹配的模型，使用默认模型ID", "model_name", modelName)
	return 1, nil // 默认返回ID为1的模型
}

// GetConversationHistory 获取对话历史记录
func (s *ChatService) GetConversationHistory(ctx context.Context, conversationID string) (*models.ConversationHistoryResponse, error) {
	utils.LogInfo("获取对话历史: %s", conversationID)

	// 将conversationID转换为int
	sessionID, err := strconv.Atoi(conversationID)
	if err != nil {
		return nil, fmt.Errorf("无效的会话ID: %w", err)
	}

	// 调用SDK获取会话记录
	recordsResp, err := s.sdk.Agent.GetSessionRecords(ctx, sessionID)
	if err != nil {
		return nil, fmt.Errorf("获取会话记录失败: %w", err)
	}

	// 直接使用SDK的输出
	response := &models.ConversationHistoryResponse{
		ConversationID: conversationID,
		Messages:       recordsResp.Records,
		Total:          len(recordsResp.Records),
	}

	utils.InfoWith("成功获取对话历史", "conversation_id", conversationID, "message_count", len(recordsResp.Records))
	return response, nil
}
