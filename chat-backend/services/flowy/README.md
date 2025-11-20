# Flowy SDK Services

这个包包含了基于 flowy-sdk 的服务实现，遵循 `../interfaces` 中定义的接口。

## 服务列表

### FlowyChatService
- **实现接口**: `interfaces.ChatServiceInterface`
- **功能**: 基于 flowy-sdk 的对话管理功能
- **特点**: 
  - 支持 Agent/Session/Config 三层架构
  - 流式消息传输
  - 完整的对话生命周期管理

### FlowyKnowledgeService
- **实现接口**: `interfaces.KnowledgeServiceInterface`
- **功能**: 基于 flowy-sdk 的知识库和文档管理
- **特点**:
  - 支持多种文档格式上传
  - 批量文件处理
  - 文件状态管理

### FlowyModelService
- **实现接口**: `interfaces.ModelServiceInterface`
- **功能**: 基于 flowy-sdk 的模型管理
- **特点**:
  - 支持聊天模型和向量模型
  - 模型状态控制
  - 动态模型添加

### FlowyDefaultSettingsService
- **实现接口**: `interfaces.DefaultSettingsServiceInterface`
- **功能**: 基于 flowy-sdk 的默认配置管理
- **特点**:
  - 持久化配置存储
  - 内置默认配置
  - 配置重置功能

## 使用方法

```go
import (
    "flowy-sdk"
    "chat-backend/services/flowy"
)

// 创建 SDK 实例
sdk := flowy.New(config)

// 创建服务实例
chatService := flowy.NewFlowyChatService(sdk, defaultSettingsService)
knowledgeService := flowy.NewFlowyKnowledgeService(sdk)
modelService := flowy.NewFlowyModelService(sdk)
defaultSettingsService := flowy.NewFlowyDefaultSettingsService()

// 使用接口方法
conversation, err := chatService.CreateConversation(ctx, settings)
knowledgeBases, err := knowledgeService.ListKnowledgeBases(ctx)
models, err := modelService.ListAvailableAllModels(ctx)
```

## 依赖关系

- **flowy-sdk**: 核心 SDK 依赖
- **chat-backend/models**: 统一的数据模型
- **chat-backend/services/interfaces**: 接口定义
- **chat-backend/utils**: 工具函数和日志

## 配置要求

Flowy 实现需要以下环境配置：
- `FLOWY_BASE_URL`: Flowy API 基础URL
- `FLOWY_API_KEY`: API 密钥（可选）
- `FLOWY_TOKEN`: 认证令牌

## 错误处理

所有服务都遵循统一的错误处理模式：
- 使用 `fmt.Errorf` 包装底层错误
- 提供详细的错误上下文信息
- 保持与 flowy-sdk 的错误兼容性
