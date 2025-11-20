# Chat Backend Services 层架构

本文档描述了 chat-backend 服务的抽象层架构，支持 Flowy SDK 和 Langchaingo 两种实现。

## 架构概述

### 设计原则

1. **接口抽象**: 所有服务通过接口定义，屏蔽底层实现差异
2. **可插拔实现**: 支持 Flowy SDK 和 Langchaingo 两种后端
3. **统一配置**: 通过环境变量控制服务类型和配置
4. **依赖注入**: 使用工厂模式管理服务生命周期

### 核心组件

```
chat-backend/services/
├── interfaces/           # 接口定义层
│   ├── chat_service_interface.go
│   ├── knowledge_service_interface.go
│   ├── model_service_interface.go
│   └── default_settings_service_interface.go
├── flowy/              # Flowy SDK 实现
│   ├── flowy_chat_service.go
│   ├── flowy_knowledge_service.go
│   ├── flowy_model_service.go
│   ├── flowy_default_settings_service.go
│   └── README.md
├── langchaingo/         # Langchaingo 实现
│   ├── config.go
│   ├── langchaingo_chat_service.go
│   ├── langchaingo_knowledge_service.go
│   ├── langchaingo_model_service.go
│   ├── langchaingo_default_settings_service.go
│   └── README.md
├── factory.go          # 服务工厂和依赖注入
└── README.md          # 本文档
```

## 接口定义

### ChatServiceInterface

聊天服务接口，提供对话管理功能：

```go
type ChatServiceInterface interface {
    CreateConversation(ctx context.Context, settings *models.ConversationSettings) (*models.Conversation, error)
    SendMessage(ctx context.Context, req *models.ChatRequest, eventChan chan<- models.SSEChatEvent) error
    ListConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error)
    DeleteConversation(ctx context.Context, conversationID string) error
    GetConversations(ctx context.Context, page, pageSize int) (*models.ConversationListResponse, error)
    UpdateConversationSettings(ctx context.Context, conversationID string, settings *models.ConversationSettings) error
    GetConversationSettings(ctx context.Context, conversationID string) (*models.ConversationSettings, error)
    GetConversationHistory(ctx context.Context, conversationID string) (*models.ConversationHistoryResponse, error)
}
```

### KnowledgeServiceInterface

知识库服务接口，提供文档管理功能：

```go
type KnowledgeServiceInterface interface {
    CreateKnowledgeBase(ctx context.Context, req *models.KnowledgeBaseCreateRequest) (*models.KnowledgeBase, error)
    GetKnowledgeBase(ctx context.Context, id int) (*models.KnowledgeBaseDetailResponse, error)
    ListKnowledgeBases(ctx context.Context, page, pageSize int) (*models.KnowledgeBaseListResponse, error)
    DeleteKnowledgeBase(ctx context.Context, id int) error
    UpdateKnowledgeBase(ctx context.Context, id int, req *models.KnowledgeBaseUpdateRequest) error
    UploadFile(ctx context.Context, kbIDStr string, filename string, content []byte) (*models.FileInfo, error)
    ListFiles(ctx context.Context, kbIDStr string, page, pageSize int) (*models.FileListResponse, error)
    DeleteFile(ctx context.Context, id int) error
    GetFileContent(ctx context.Context, id int) (io.ReadCloser, error)
}
```

### ModelServiceInterface

模型服务接口，提供模型管理功能：

```go
type ModelServiceInterface interface {
    ListSupportedChatModels(ctx context.Context) ([]modelSvc.SupportedChatModel, error)
    ListSupportedVectorModels(ctx context.Context) ([]modelSvc.SupportedVectorModel, error)
    ListAvailableAllModels(ctx context.Context) ([]modelSvc.ModelInfo, error)
    SaveModel(ctx context.Context, req *modelSvc.ModelSaveRequest) (int, error)
    DeleteModel(ctx context.Context, id int) error
    ListAvailableChatModels(ctx context.Context) ([]modelSvc.ModelInfo, error)
    ListAvailableVectorModels(ctx context.Context) ([]modelSvc.ModelInfo, error)
    SetModelStatus(ctx context.Context, id int, enable bool) error
}
```

### DefaultSettingsServiceInterface

默认设置服务接口，提供配置管理功能：

```go
type DefaultSettingsServiceInterface interface {
    GetDefaultSettings() *models.DefaultSettings
    SaveDefaultSettings(settings *models.DefaultSettings) error
    ResetDefaultSettings() error
}
```

## 实现对比

### Flowy SDK 实现

**特点**：
- 基于 Flowy SDK 的 HTTP API 调用
- 支持完整的 Flowy 生态系统功能
- 适合生产环境使用

**配置**：
```bash
# Flowy SDK 配置
FLOWY_BASE_URL=http://192.168.1.2:8888/api/v1
FLOWY_API_KEY=your_api_key
FLOWY_TOKEN=your_token
```

**服务创建**：
```go
container, err := services.NewServiceContainer(services.ServiceTypeFlowy)
chatService := container.GetChatService()
```

### Langchaingo 实现

**特点**：
- 基于 Langchaingo 框架的本地实现
- 集成 OpenAI、Ollama、Qdrant、Docling
- 适合开发和实验环境使用

**技术栈**：
- **LLM**: OpenAI (gpt-3.5-turbo, gpt-4)
- **Embedding**: Ollama bge-m3:latest
- **Vector Store**: Qdrant
- **Document Processing**: Docling HTTP API
- **Chat Memory**: SQLite

**配置**：
```bash
# Langchaingo 配置
LANGCHAINO_OPENAI_BASE_URL=https://api.openai.com/v1
LANGCHAINO_OPENAI_API_KEY=your_openai_api_key
LANGCHAINO_OPENAI_MODEL=gpt-3.5-turbo

LANGCHAINO_OLLAMA_URL=http://localhost:11434
LANGCHAINO_OLLAMA_MODEL=bge-m3:latest

LANGCHAINO_QDRANT_URL=http://localhost:6333
LANGCHAINO_QDRANT_API_KEY=your_qdrant_api_key

LANGCHAINO_DOCLING_URL=http://localhost:8001
LANGCHAINO_DOCLING_API_KEY=your_docling_api_key

LANGCHAINO_SQLITE_DB_PATH=./chat_history.db
LANGCHAINO_SQLITE_SESSION=default
```

**服务创建**：
```go
container, err := services.NewServiceContainer(services.ServiceTypeLangchaingo)
chatService := container.GetChatService()
```

## 服务工厂

### ServiceContainer

服务容器负责管理所有服务的生命周期：

```go
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
```

### 全局服务管理

```go
// 初始化全局服务
err := services.InitGlobalServices(services.ServiceTypeFlowy)

// 获取全局服务
chatService := services.GetGlobalChatService()
knowledgeService := services.GetGlobalKnowledgeService()
modelService := services.GetGlobalModelService()
defaultSettingsService := services.GetGlobalDefaultSettingsService()

// 获取服务容器
container := services.GetGlobalServiceContainer()

// 健康检查
err := container.HealthCheck(ctx)

// 获取服务信息
info := container.GetServiceInfo()

// 关闭服务
err := services.Shutdown()
```

## 环境变量配置

### 服务类型控制

```bash
# 选择服务实现类型
SERVICE_TYPE=flowy        # 使用 Flowy SDK
# 或
SERVICE_TYPE=langchaingo    # 使用 Langchaingo
```

### 日志配置

```bash
# 日志配置
LOG_DIR=logs
LOG_MAX_SIZE=10485760    # 10MB
LOG_MAX_FILES=10
LOG_ENABLE_STDOUT=true
```

## 使用示例

### 在 Handler 中使用

```go
// 在 handler 中使用服务
func (h *ChatHandler) HandleChat(c *gin.Context) {
    // 获取全局聊天服务
    chatService := services.GetGlobalChatService()
    
    // 创建对话
    conversation, err := chatService.CreateConversation(c.Request.Context(), settings)
    if err != nil {
        c.JSON(500, gin.H{"error": err.Error()})
        return
    }
    
    // 发送消息
    eventChan := make(chan models.SSEChatEvent)
    go func() {
        defer close(eventChan)
        err := chatService.SendMessage(c.Request.Context(), req, eventChan)
        if err != nil {
            // 处理错误
        }
    }()
    
    // SSE 流式响应
    c.Stream(func(w http.ResponseWriter) error {
        // 处理事件流
        return nil
    })
}
```

### 切换服务实现

```go
// 在 main.go 中初始化
func main() {
    // 从环境变量确定服务类型
    serviceType := services.ServiceTypeFlowy
    if os.Getenv("SERVICE_TYPE") == "langchaingo" {
        serviceType = services.ServiceTypeLangchaingo
    }
    
    // 初始化全局服务
    if err := services.InitGlobalServices(serviceType); err != nil {
        log.Fatalf("Failed to initialize services: %v", err)
    }
    
    // 启动 HTTP 服务器
    // ...
    
    // 清理资源
    defer services.Shutdown()
}
```

## 开发指南

### 添加新的服务实现

1. **创建接口**: 在 `interfaces/` 目录下定义新的接口
2. **实现接口**: 在对应的实现包中实现接口
3. **注册服务**: 在工厂中注册新服务
4. **更新配置**: 如需要，添加配置结构

### 测试

```go
// 单元测试示例
func TestChatService(t *testing.T) {
    // 创建测试容器
    container, err := services.NewServiceContainer(services.ServiceTypeFlowy)
    require.NoError(t, err)
    
    // 获取服务
    chatService := container.GetChatService()
    require.NotNil(t, chatService)
    
    // 测试功能
    conversation, err := chatService.CreateConversation(context.Background(), nil)
    assert.NoError(t, err)
    assert.NotNil(t, conversation)
}
```

## 性能考虑

### Flowy SDK
- **网络延迟**: HTTP 调用可能存在网络延迟
- **连接池**: 使用 HTTP 连接池优化性能
- **缓存**: 适当缓存不经常变化的数据

### Langchaingo
- **本地计算**: 所有计算在本地进行，无网络延迟
- **内存使用**: 嵌入模型和大语言模型需要较多内存
- **并发处理**: 利用 Go 协程实现高并发

## 故障排除

### 常见问题

1. **服务初始化失败**
   - 检查环境变量配置
   - 确认依赖服务是否运行
   - 查看日志文件

2. **配置验证失败**
   - 检查配置文件格式
   - 验证必需的配置项
   - 确认配置值有效性

3. **服务调用失败**
   - 检查网络连接
   - 验证 API 密钥
   - 查看服务健康状态

### 调试技巧

```go
// 启用详细日志
utils.InitLogger(&utils.LogConfig{
    LogDir:       "logs",
    EnableStdout: true,
})

// 获取服务信息
container := services.GetGlobalServiceContainer()
info := container.GetServiceInfo()
fmt.Printf("Service Info: %+v\n", info)

// 健康检查
err := container.HealthCheck(context.Background())
if err != nil {
    fmt.Printf("Health check failed: %v\n", err)
}
```

## 未来扩展

### 计划功能

1. **多租户支持**: 支持多个租户隔离
2. **服务发现**: 动态服务发现和注册
3. **监控集成**: 集成 Prometheus 和 Grafana
4. **链路追踪**: 添加 OpenTelemetry 支持
5. **热重载**: 支持配置热重载

### 扩展点

- **新服务实现**: 可以添加新的后端实现
- **插件系统**: 支持第三方插件
- **中间件**: 添加请求/响应中间件
- **事件系统**: 实现事件驱动架构

## 总结

chat-backend 服务层通过接口抽象实现了：

- **灵活性**: 支持多种后端实现
- **可维护性**: 清晰的分层架构
- **可扩展性**: 易于添加新功能
- **可测试性**: 接口便于单元测试

这种架构设计使得系统可以根据不同需求选择合适的实现，同时保持 API 的兼容性和一致性。
