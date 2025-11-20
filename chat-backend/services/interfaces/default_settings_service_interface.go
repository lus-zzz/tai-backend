package interfaces

import "chat-backend/models"

// DefaultSettingsServiceInterface 默认配置服务接口
type DefaultSettingsServiceInterface interface {
	// GetDefaultSettings 获取默认配置
	GetDefaultSettings() *models.DefaultSettings

	// UpdateDefaultSettings 更新默认配置
	UpdateDefaultSettings(settings *models.DefaultSettings) error

	// ResetToDefaults 重置为内置默认配置
	ResetToDefaults() error
}
