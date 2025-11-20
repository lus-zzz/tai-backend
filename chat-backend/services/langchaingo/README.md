# Langchaingo 服务实现

基于 Langchaingo 框架的服务实现，提供与 Flowy SDK 兼容的接口。

## 架构概述

Langchaingo 实现包提供了以下服务：

### 核心组件

1. **ChatService** (`langchaingo_chat_service.go`)
   - 基于 OpenAI LLM 的对话服务
   - 使用 SQLite 进行对话记忆管理
   - 支持流式输出
   - 集成 Qdrant 向量检索

2. **KnowledgeService** (`langchaingo_knowledge_service.go`)
   - 基于 Docling API 的文档分块
   - 使用 Ollama bge-m3 进行向量化
   - 集成 Qdrant 向量存储
   - 支持批量文件上传

3. **ModelService** (`langchaingo_model_service.go`)
   - 管理支持的聊天和向量模型
   - 提供 OpenAI 和 Ollama 模型信息
   - 支持模型状态管理

4. **DefaultSettingsService** (`langchaingo_default_settings_service.go`)
   - 管理默认配置
   - 持久化到本地文件
   - 支持配置重置

5. **Config** (`config.go`)
   - 统一的配置管理
   - 环境变量支持
   - 配置验证

## 技术栈

### Langchaingo 组件

- **LLM**: OpenAI (gpt-3.5-turbo, gpt-4)
- **Embedding**: Ollama bge-m3:latest
- **Vector Store**: Qdrant
- **Document Processing**: Docling HTTP API
- **Chat Memory**: SQLite

### 配置环境变量

```bash
# OpenAI 配置
LANGCHAINO_OPENAI_BASE_URL=https://api.openai.com/v1
LANGCHAINO_OPENAI_API_KEY=your_openai_api_key
LANGCHAINO_OPENAI_MODEL=gpt-3.5-turbo

# Ollama 配置
LANGCHAINO_OLLAMA_URL=http://localhost:11434
LANGCHAINO_OLLAMA_MODEL=bge-m3:latest

# Qdrant 配置
LANGCHAINO_QDRANT_URL=http://localhost:6333
LANGCHAINO_QDRANT_API_KEY=your_qdrant_api_key

# Docling 配置
LANGCHAINO_DOCLING_URL=http://localhost:8001
LANGCHAINO_DOCLING_API_KEY=your_docling_api_key

# SQLite 配置
LANGCHAINO_SQLITE_DB_PATH=./chat_history.db
LANGCHAINO_SQLITE_SESSION=default
```

## 核心流程实现

### 1. 文档分块和向量化流程 (chunkAndVectorize)

基于 `KEY_PROCESS_AND_CODE.md` 中的流程：

```go
func (s *LangchaingoKnowledgeService) chunkAndVectorize(filename string, content []byte) error {
    // 1. 使用 Docling API 进行文档分块
    chunks, err := s.chunkDocument(ctx, filename, content)
    
    // 2. 使用 Ollama bge-m3 生成嵌入向量
    // 3. 存储到 Qdrant 向量数据库
    err = s.vectorizeAndStore(ctx, collectionName, chunks)
}
```

### 2. 对话流程 (chat)

```go
func (s *LangchaingoChatService) chat(collectionName string) error {
    // 1. 初始化 OpenAI LLM
    llm, err := s.initializeLLM(ctx)
    
    // 2. 初始化 Ollama 嵌入模型
    embedder, err := s.initializeEmbedder(ctx)
    
    // 3. 连接到 Qdrant 向量数据库
    store, err := s.connectToQdrant(ctx, collectionName)
    
    // 4. 创建 SQLite 对话记忆
    history, err := s.createChatHistory(ctx, sessionID)
    
    // 5. 使用 ConversationalRetrievalQA 进行问答
    executor := chains.NewConversationalRetrievalQAFromLLM(llm, retriever, conversation)
}
```

## 使用示例

### 创建服务实例

```go
// 获取配置
config := GetLangchaingoConfig()
config.ValidateConfig()

// 创建服务
chatService := NewLangchaingoChatService(config, defaultSettingsService)
knowledgeService := NewLangchaingoKnowledgeService(config)
modelService := NewLangchaingoModelService(config)
defaultSettingsService := NewLangchaingoDefaultSettingsService()
```

### 对话示例

```go
// 创建对话
conversation, err := chatService.CreateConversation(ctx, &models.ConversationSettings{
    Name:    "测试对话",
    ModelID: 1,
    Stream:   true,
})

// 发送消息
eventChan := make(chan models.SSEChatEvent)
err = chatService.SendMessage(ctx, &models.ChatRequest{
    SessionID: conversation.ID,
    Content:   "你好，请介绍一下你自己",
}, eventChan)

// 处理流式响应
for event := range eventChan {
    switch event.EventType {
    case "resp_splash":
        fmt.Println("开始响应")
    case "resp_increment":
        fmt.Print(event.Message)
    case "resp_finish":
        fmt.Println("\n响应完成")
    }
}
```

### 知识库示例

```go
// 创建知识库
kb, err := knowledgeService.CreateKnowledgeBase(ctx, &models.KnowledgeBaseCreateRequest{
    Name:        "测试知识库",
    Desc:        "用于测试的知识库",
    VectorModel: 1, // bge-m3
    AgentModel:  1, // gpt-3.5-turbo
    ChunkSize:   512,
})

// 上传文件
file, err := knowledgeService.UploadFile(ctx, strconv.Itoa(kb.ID), "document.pdf", content)
```

## 特性

### 已实现

- ✅ 接口抽象层
- ✅ 配置管理
- ✅ 基础服务结构
- ✅ 模拟流式对话
- ✅ 文档分块接口
- ✅ 模型管理接口

### 待实现

- 🔄 实际的 Langchaingo OpenAI 集成
- 🔄 实际的 Langchaingo Ollama 集成
- 🔄 实际的 Langchaingo Qdrant 集成
- 🔄 实际的 Langchaingo SQLite 集成
- 🔄 完整的 RAG 流程实现
- 🔄 错误处理和重试机制
- 🔄 性能优化和监控

## 与 Flowy SDK 的兼容性

Langchaingo 实现完全兼容 Flowy SDK 的接口规范：

- 使用相同的请求/响应结构
- 保持相同的 API 签名
- 支持相同的错误处理模式
- 兼容现有的 Handler 层

## 部署要求

### 依赖服务

1. **OpenAI API**: 可访问的 OpenAI API 端点
2. **Ollama**: 运行 bge-m3:latest 模型
3. **Qdrant**: 向量数据库服务
4. **Docling**: 文档处理服务

### 系统要求

- Go 1.19+
- 足够的内存用于嵌入模型
- 网络连接到外部服务

## 开发指南

### 添加新的 Langchaingo 功能

1. 在相应的服务文件中实现 TODO 项
2. 更新配置结构（如果需要）
3. 添加相应的测试
4. 更新文档

### 调试

启用详细日志：

```go
utils.SetLogLevel("debug")
```

## 性能考虑

- **批处理**: 文档分块支持批处理
- **缓存**: 模型初始化结果缓存
- **连接池**: 数据库连接复用
- **异步**: 非阻塞的流式处理

## 安全注意事项

- API 密钥通过环境变量管理
- 输入验证和清理
- 错误信息不泄露敏感数据
- 访问控制和审计日志
