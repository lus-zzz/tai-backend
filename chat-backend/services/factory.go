package services

import (
	"context"
	"fmt"
	"os"

	"chat-backend/pkg/database"
	"chat-backend/services/flowy"
	"chat-backend/services/interfaces"
	"chat-backend/services/langchaingo"
	"chat-backend/utils"
	flowySDK "flowy-sdk"
	"flowy-sdk/pkg/config"
)

// ServiceType 服务类型枚举
type ServiceType string

const (
	ServiceTypeFlowy      ServiceType = "flowy"
	ServiceTypeLangchaingo ServiceType = "langchaingo"
)

// ServiceContainer 服务容器
type ServiceContainer struct {
	serviceType ServiceType
	
	// 接口实例
	chatService             interfaces.ChatServiceInterface
	knowledgeService         interfaces.KnowledgeServiceInterface
	modelService            interfaces.ModelServiceInterface
	defaultSettingsService  interfaces.DefaultSettingsServiceInterface
	
	// 配置
	flowyConfig      *config.Config
	langchaingoConfig *langchaingo.LangchaingoConfig
}

// NewServiceContainer 创建服务容器
func NewServiceContainer(serviceType ServiceType) (*ServiceContainer, error) {
	container := &ServiceContainer{
		serviceType: serviceType,
	}

	// 根据服务类型初始化
	switch serviceType {
	case ServiceTypeFlowy:
		return container.initFlowyServices()
	case ServiceTypeLangchaingo:
		return container.initLangchaingoServices()
	default:
		return nil, fmt.Errorf("unsupported service type: %s", serviceType)
	}
}

// initFlowyServices 初始化 Flowy 服务
func (sc *ServiceContainer) initFlowyServices() (*ServiceContainer, error) {
	utils.LogInfo("初始化 Flowy 服务")

	// 加载 Flowy 配置
	flowyCfg := config.DefaultConfig().LoadFromEnv()

	// 创建 Flowy SDK 实例
	sdk := flowySDK.New(flowyCfg)

	// 创建默认设置服务
	sc.defaultSettingsService = flowy.NewFlowyDefaultSettingsService()

	// 创建其他服务
	sc.chatService = flowy.NewFlowyChatService(sdk, sc.defaultSettingsService)
	sc.knowledgeService = flowy.NewFlowyKnowledgeService(sdk)
	sc.modelService = flowy.NewFlowyModelService(sdk)

	utils.InfoWith("Flowy 服务初始化完成", "chat_service", "flowy", "knowledge_service", "flowy", "model_service", "flowy")
	return sc, nil
}

// initLangchaingoServices 初始化 Langchaingo 服务
func (sc *ServiceContainer) initLangchaingoServices() (*ServiceContainer, error) {
	utils.LogInfo("初始化 Langchaingo 服务")

	// 加载 Langchaingo 配置
	langchaingoCfg := langchaingo.GetLangchaingoConfig()
	
	// 验证配置
	if err := langchaingoCfg.ValidateConfig(); err != nil {
		utils.WarnWith("Langchaingo 配置验证失败，使用默认值", "error", err.Error())
	}
	sc.langchaingoConfig = langchaingoCfg

	// 创建数据库连接
	db, err := database.NewDatabase(langchaingoCfg.SQLite.DBPath)
	if err != nil {
		return nil, fmt.Errorf("创建数据库连接失败: %w", err)
	}

	// 创建默认设置服务
	sc.defaultSettingsService = langchaingo.NewLangchaingoDefaultSettingsService()

	// 创建其他服务
	sc.chatService = langchaingo.NewLangchaingoChatService(langchaingoCfg, sc.defaultSettingsService)
	sc.knowledgeService = langchaingo.NewLangchaingoKnowledgeService(langchaingoCfg)
	sc.modelService = langchaingo.NewLangchaingoModelService(db, langchaingoCfg)

	utils.InfoWith("Langchaingo 服务初始化完成", "chat_service", "langchaingo", "knowledge_service", "langchaingo", "model_service", "langchaingo")
	return sc, nil
}

// GetChatService 获取聊天服务
func (sc *ServiceContainer) GetChatService() interfaces.ChatServiceInterface {
	return sc.chatService
}

// GetKnowledgeService 获取知识库服务
func (sc *ServiceContainer) GetKnowledgeService() interfaces.KnowledgeServiceInterface {
	return sc.knowledgeService
}

// GetModelService 获取模型服务
func (sc *ServiceContainer) GetModelService() interfaces.ModelServiceInterface {
	return sc.modelService
}

// GetDefaultSettingsService 获取默认设置服务
func (sc *ServiceContainer) GetDefaultSettingsService() interfaces.DefaultSettingsServiceInterface {
	return sc.defaultSettingsService
}

// GetServiceType 获取服务类型
func (sc *ServiceContainer) GetServiceType() ServiceType {
	return sc.serviceType
}

// HealthCheck 健康检查
func (sc *ServiceContainer) HealthCheck(ctx context.Context) error {
	utils.InfoWith("执行服务健康检查", "service_type", sc.serviceType)

	// 检查默认设置服务
	if sc.defaultSettingsService == nil {
		return fmt.Errorf("default settings service is nil")
	}

	// 尝试获取默认设置
	_ = sc.defaultSettingsService.GetDefaultSettings()
	// 接口方法不返回错误，所以直接检查服务是否可用
	if sc.defaultSettingsService == nil {
		return fmt.Errorf("default settings service is nil")
	}

	utils.InfoWith("服务健康检查通过", "service_type", sc.serviceType)
	return nil
}

// GetServiceInfo 获取服务信息
func (sc *ServiceContainer) GetServiceInfo() map[string]interface{} {
	info := map[string]interface{}{
		"service_type": sc.serviceType,
		"services": map[string]string{
			"chat_service":             getServiceTypeName(sc.chatService),
			"knowledge_service":         getServiceTypeName(sc.knowledgeService),
			"model_service":            getServiceTypeName(sc.modelService),
			"default_settings_service": getServiceTypeName(sc.defaultSettingsService),
		},
	}

	// 添加配置信息
	switch sc.serviceType {
	case ServiceTypeFlowy:
		if sc.flowyConfig != nil {
			info["config"] = map[string]interface{}{
				"server_url": sc.flowyConfig.BaseURL,
			}
		}
	case ServiceTypeLangchaingo:
		if sc.langchaingoConfig != nil {
			info["config"] = map[string]interface{}{
				"llm_url":       sc.langchaingoConfig.LLM.BaseURL,
				"embedding_url":  sc.langchaingoConfig.Embedding.BaseURL,
				"qdrant_url":    sc.langchaingoConfig.Qdrant.URL,
				"docling_url":   sc.langchaingoConfig.Docling.BaseURL,
				"sqlite_db":     sc.langchaingoConfig.SQLite.DBPath,
			}
		}
	}

	return info
}

// getServiceTypeName 获取服务类型名称
func getServiceTypeName(service interface{}) string {
	if service == nil {
		return "nil"
	}
	
	switch service.(type) {
	case *flowy.FlowyChatService:
		return "flowy"
	case *flowy.FlowyKnowledgeService:
		return "flowy"
	case *flowy.FlowyModelService:
		return "flowy"
	case *flowy.FlowyDefaultSettingsService:
		return "flowy"
	case *langchaingo.LangchaingoChatService:
		return "langchaingo"
	case *langchaingo.LangchaingoKnowledgeService:
		return "langchaingo"
	case *langchaingo.LangchaingoModelService:
		return "langchaingo"
	case *langchaingo.LangchaingoDefaultSettingsService:
		return "langchaingo"
	default:
		return "unknown"
	}
}

// 全局服务容器实例
var globalServiceContainer *ServiceContainer

// InitGlobalServices 初始化全局服务
func InitGlobalServices(serviceType ServiceType) error {
	utils.InfoWith("初始化全局服务", "service_type", serviceType)

	container, err := NewServiceContainer(serviceType)
	if err != nil {
		return fmt.Errorf("failed to create service container: %w", err)
	}

	// 执行健康检查
	ctx := context.Background()
	if err := container.HealthCheck(ctx); err != nil {
		utils.WarnWith("服务健康检查失败", "error", err.Error())
		// 不返回错误，允许服务继续运行
	}

	globalServiceContainer = container
	utils.InfoWith("全局服务初始化完成", "service_type", serviceType)
	return nil
}

// GetGlobalServiceContainer 获取全局服务容器
func GetGlobalServiceContainer() *ServiceContainer {
	if globalServiceContainer == nil {
		// 如果没有初始化，尝试从环境变量确定服务类型
		serviceType := getDefaultServiceType()
		utils.WarnWith("全局服务容器未初始化，使用默认服务类型", "service_type", serviceType)
		
		if err := InitGlobalServices(serviceType); err != nil {
			utils.ErrorWith("初始化默认服务失败", "error", err.Error())
			return nil
		}
	}
	return globalServiceContainer
}

// getDefaultServiceType 从环境变量获取默认服务类型
func getDefaultServiceType() ServiceType {
	serviceType := os.Getenv("SERVICE_TYPE")
	if serviceType == "" {
		serviceType = "flowy" // 默认使用 flowy
	}
	
	switch serviceType {
	case "flowy":
		return ServiceTypeFlowy
	case "langchaingo":
		return ServiceTypeLangchaingo
	default:
		utils.WarnWith("未知的 SERVICE_TYPE，使用默认值", "service_type", serviceType)
		return ServiceTypeFlowy
	}
}

// GetGlobalChatService 获取全局聊天服务
func GetGlobalChatService() interfaces.ChatServiceInterface {
	container := GetGlobalServiceContainer()
	if container == nil {
		return nil
	}
	return container.GetChatService()
}

// GetGlobalKnowledgeService 获取全局知识库服务
func GetGlobalKnowledgeService() interfaces.KnowledgeServiceInterface {
	container := GetGlobalServiceContainer()
	if container == nil {
		return nil
	}
	return container.GetKnowledgeService()
}

// GetGlobalModelService 获取全局模型服务
func GetGlobalModelService() interfaces.ModelServiceInterface {
	container := GetGlobalServiceContainer()
	if container == nil {
		return nil
	}
	return container.GetModelService()
}

// GetGlobalDefaultSettingsService 获取全局默认设置服务
func GetGlobalDefaultSettingsService() interfaces.DefaultSettingsServiceInterface {
	container := GetGlobalServiceContainer()
	if container == nil {
		return nil
	}
	return container.GetDefaultSettingsService()
}

// Shutdown 关闭服务
func Shutdown() error {
	if globalServiceContainer != nil {
		utils.InfoWith("关闭服务", "service_type", globalServiceContainer.serviceType)
		
		// 这里可以添加具体的清理逻辑
		// 比如关闭数据库连接、释放资源等
		
		globalServiceContainer = nil
		utils.LogInfo("服务已关闭")
	}
	return nil
}
