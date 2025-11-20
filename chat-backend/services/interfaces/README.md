# Services Interfaces

这个包定义了 chat-backend 服务层的所有接口，用于实现服务抽象和多种实现的支持。

## 接口列表

### ChatServiceInterface
聊天服务接口，定义了对话管理相关的所有方法：
- 创建对话
- 发送消息
- 获取对话列表和历史
- 更新对话设置
- 删除对话

### KnowledgeServiceInterface  
知识库服务接口，定义了知识库和文档管理相关的所有方法：
- 知识库 CRUD 操作
- 文件上传和管理
- 批量操作支持

### ModelServiceInterface
模型服务接口，定义了模型管理相关的所有方法：
- 获取可用模型列表
- 模型的增删改查
- 模型状态管理

## 实现要求

所有实现都必须：
1. 严格遵循接口定义的方法签名
2. 使用 `chat-backend/models` 包中的数据结构
3. 正确处理 `context.Context` 参数
4. 返回符合 Go 错误处理约定的错误信息

## 当前实现

- **Flowy 实现**: `../flowy/` - 基于 flowy-sdk 的实现
- **Langchaingo 实现**: `../langchaingo/` - 基于 langchaingo 的实现（开发中）
