package langchaingo

import (
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"chat-backend/models"
	"chat-backend/services/interfaces"
	"chat-backend/utils"
)

// LangchaingoDefaultSettingsService 基于 langchaingo 的默认配置服务实现
type LangchaingoDefaultSettingsService struct {
	filePath string
	mu       sync.RWMutex
	settings *models.DefaultSettings
}

// NewLangchaingoDefaultSettingsService 创建 Langchaingo 默认配置服务
func NewLangchaingoDefaultSettingsService() interfaces.DefaultSettingsServiceInterface {
	configDir := "config"

	// 尝试创建外部配置目录（用于保存用户修改）
	_ = os.MkdirAll(configDir, 0755)

	service := &LangchaingoDefaultSettingsService{
		filePath: filepath.Join(configDir, "langchaingo_default_settings.json"),
	}

	// 尝试从外部文件加载
	if err := service.load(); err != nil {
		// 如果外部文件不存在，使用内置默认配置
		utils.LogInfo("使用内置默认配置")
		service.settings = service.getBuiltinDefaults()
	}

	return service
}

// getBuiltinDefaults 获取内置默认配置
func (s *LangchaingoDefaultSettingsService) getBuiltinDefaults() *models.DefaultSettings {
	return &models.DefaultSettings{
		Models: models.DefaultModelSettings{
			ChatModelID:      1, // 默认对话模型ID (OpenAI)
			EmbeddingModelID: 1, // 默认嵌入模型ID (bge-m3)
		},
		Conversation: models.DefaultConversationConfig{
			Name:             "新对话",   // 默认对话名称
			Desc:             "",      // 默认对话描述
			ModelID:          1,       // 默认模型ID (OpenAI)
			Temperature:      0.7,     // 默认温度（多样性）
			TopP:             1.0,     // 默认采样范围
			FrequencyPenalty: 0.0,     // 默认重复控制
			PresencePenalty:  0.0,     // 默认词汇控制
			ResponseType:     "text",  // 默认响应类型
			Stream:           true,    // 默认启用流式输出
			KnowledgeBaseIDs: []int{}, // 默认不关联知识库
			ContextLimit:     16,      // 默认上下文限制16条消息
		},
		KnowledgeBase: models.DefaultKnowledgeBaseConfig{
			Name:          "默认知识库", // 默认知识库名称
			Desc:          "",      // 默认知识库描述
			VectorModel:   1,       // 默认向量模型ID (bge-m3)
			AgentModel:    1,       // 默认对话模型ID (OpenAI)
			ChunkStrategy: "fixed", // 默认切片策略：固定大小
			ChunkSize:     512,     // 默认分块大小512
		},
		UpdatedAt: time.Now(),
	}
}

// GetDefaultSettings 获取默认配置
func (s *LangchaingoDefaultSettingsService) GetDefaultSettings() *models.DefaultSettings {
	s.mu.RLock()
	defer s.mu.RUnlock()

	if s.settings == nil {
		s.settings = s.getBuiltinDefaults()
	}

	// 返回副本
	settingsCopy := *s.settings
	return &settingsCopy
}

// UpdateDefaultSettings 更新默认配置
func (s *LangchaingoDefaultSettingsService) UpdateDefaultSettings(settings *models.DefaultSettings) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if settings == nil {
		return nil
	}

	// 直接替换整个配置
	settings.UpdatedAt = time.Now()
	s.settings = settings

	// 持久化到文件
	if err := s.save(); err != nil {
		utils.ErrorWith("保存默认配置失败", "error", err)
		return err
	}

	utils.LogInfo("默认配置已更新并保存")
	return nil
}

// ResetToDefaults 重置为内置默认配置
func (s *LangchaingoDefaultSettingsService) ResetToDefaults() error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.settings = s.getBuiltinDefaults()

	// 持久化到文件
	if err := s.save(); err != nil {
		utils.ErrorWith("重置默认配置失败", "error", err)
		return err
	}

	utils.LogInfo("默认配置已重置")
	return nil
}

// load 从文件加载配置
func (s *LangchaingoDefaultSettingsService) load() error {
	data, err := os.ReadFile(s.filePath)
	if err != nil {
		return err
	}

	var settings models.DefaultSettings
	if err := json.Unmarshal(data, &settings); err != nil {
		return err
	}

	s.settings = &settings
	utils.LogInfo("默认配置已从文件加载: %s", s.filePath)
	return nil
}

// save 保存配置到文件
func (s *LangchaingoDefaultSettingsService) save() error {
	if s.settings == nil {
		return nil
	}

	data, err := json.MarshalIndent(s.settings, "", "  ")
	if err != nil {
		return err
	}

	if err := os.WriteFile(s.filePath, data, 0644); err != nil {
		return err
	}

	return nil
}
