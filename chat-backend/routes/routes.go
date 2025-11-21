package routes

import (
	"chat-backend/handlers"
	"chat-backend/utils"
	"embed"
	"encoding/json"
	"io/fs"
	"net/http"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// Router 路由器
type Router struct {
	engine           *gin.Engine
	chatHandler      *handlers.ChatHandler
	settingsHandler  *handlers.SettingsHandler
	knowledgeHandler *handlers.KnowledgeHandler
	modelHandler     *handlers.ModelHandler
	versionHandler   *handlers.VersionHandler
	staticFS         embed.FS
	docsFS           embed.FS
}

// NewRouter 创建新路由
func NewRouter(
	chatHandler *handlers.ChatHandler,
	settingsHandler *handlers.SettingsHandler,
	knowledgeHandler *handlers.KnowledgeHandler,
	modelHandler *handlers.ModelHandler,
	versionHandler *handlers.VersionHandler,
	staticFS embed.FS,
	docsFS embed.FS,
) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())
	engine.Use(utils.ResponseHandlerMiddleware()) // 添加响应处理中间件

	// CORS中间件
	engine.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "*")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	})

	router := &Router{
		engine:           engine,
		chatHandler:      chatHandler,
		settingsHandler:  settingsHandler,
		knowledgeHandler: knowledgeHandler,
		modelHandler:     modelHandler,
		versionHandler:   versionHandler,
		staticFS:         staticFS,
		docsFS:           docsFS,
	}

	router.setupRoutes()
	return router
}

// setupRoutes 设置路由
func (r *Router) setupRoutes() {
	api := r.engine.Group("/api/v1")

	// 聊天相关路由
	chat := api.Group("/chat")
	{
		chat.POST("/conversations", utils.WrapHandler(r.chatHandler.CreateConversation))
		chat.GET("/conversations", utils.WrapHandler(r.chatHandler.GetConversations))
		chat.DELETE("/conversations/:id", utils.WrapHandler(r.chatHandler.DeleteConversation))
		chat.GET("/conversations/:id/history", utils.WrapHandler(r.chatHandler.GetConversationHistory))
		chat.POST("/messages", r.chatHandler.SendMessage) // SendMessage 保持原样，使用SSE
		chat.GET("/conversations/:id/settings", utils.WrapHandler(r.chatHandler.GetConversationSettings))
		chat.PUT("/conversations/:id/settings", utils.WrapHandler(r.chatHandler.UpdateConversationSettings))
	}

	// 默认配置相关路由
	settings := api.Group("/settings")
	{
		settings.GET("/defaults", utils.WrapHandler(r.settingsHandler.GetDefaultSettings))
		settings.PUT("/defaults", utils.WrapHandler(r.settingsHandler.UpdateDefaultSettings))
		settings.POST("/defaults/reset", utils.WrapHandler(r.settingsHandler.ResetDefaultSettings))
	}

	// 知识库相关路由
	knowledge := api.Group("/knowledge")
	{
		knowledge.GET("/bases", utils.WrapHandler(r.knowledgeHandler.ListKnowledgeBases))
		knowledge.POST("/bases", utils.WrapHandler(r.knowledgeHandler.CreateKnowledgeBase))
		knowledge.PUT("/bases/:id", utils.WrapHandler(r.knowledgeHandler.UpdateKnowledgeBase))
		knowledge.DELETE("/bases/:id", utils.WrapHandler(r.knowledgeHandler.DeleteKnowledgeBase))
		knowledge.GET("/bases/:id/files", utils.WrapHandler(r.knowledgeHandler.GetKnowledgeBaseFiles))
		knowledge.POST("/bases/:id/files", r.knowledgeHandler.UploadFile) // UploadFile 保持原样，使用复杂逻辑

		// 文件操作路由（只需要文件ID）
		knowledge.DELETE("/files/:file_id", utils.WrapHandler(r.knowledgeHandler.DeleteFile))
		knowledge.PUT("/files/:file_id/toggle", utils.WrapHandler(r.knowledgeHandler.ToggleFileEnable))
	}

	// 模型相关路由
	models := api.Group("/models")
	{
		// 支持的模型列表
		models.GET("/supported/chat", utils.WrapHandler(r.modelHandler.ListSupportedChatModels))
		models.GET("/supported/vector", utils.WrapHandler(r.modelHandler.ListSupportedVectorModels))

		// 可用的模型列表
		models.GET("/available/all", utils.WrapHandler(r.modelHandler.ListAvailableAllModels))
		models.GET("/available/chat", utils.WrapHandler(r.modelHandler.ListAvailableChatModels))
		models.GET("/available/vector", utils.WrapHandler(r.modelHandler.ListAvailableVectorModels))

		// 模型管理
		models.POST("", utils.WrapHandler(r.modelHandler.SaveModel))
		models.DELETE("/:id", utils.WrapHandler(r.modelHandler.DeleteModel))
		models.PUT("/:id/status", utils.WrapHandler(r.modelHandler.SetModelStatus))
	}


	// 系统信息路由 - 为swag创建直接路由
	api.GET("/version", func(c *gin.Context) {
		result, err := r.versionHandler.GetVersion(c)
		if err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
		c.JSON(200, result)
	})


	api.GET("/swagger-test", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "swagger test"})
	})

	// 静态文件服务 - 使用嵌入的文件系统
	staticSubFS, err := fs.Sub(r.staticFS, "static")
	if err == nil {
		r.engine.StaticFS("/static", http.FS(staticSubFS))
	}

	// 日志查看
	logViewer := utils.NewLogViewerHandler("logs")
	r.engine.GET("/logs", logViewer.ListLogFiles)
	r.engine.GET("/logs/:filename", logViewer.ViewLogFile)
	r.engine.GET("/logs/:filename/download", logViewer.DownloadLogFile)
	r.engine.GET("/logs/:filename/stream", logViewer.StreamLogFile)

	// Swagger UI - 使用ginSwagger中间件，动态生成完整URL
	r.engine.GET("/swagger/*any", func(c *gin.Context) {
		// 构建完整的 swagger.json URL
		scheme := "http"
		if c.Request.TLS != nil {
			scheme = "https"
		}
		swaggerURL := scheme + "://" + c.Request.Host + "/swagger.json"

		// 使用动态URL创建handler
		handler := ginSwagger.WrapHandler(swaggerFiles.Handler,
			ginSwagger.URL(swaggerURL),
			ginSwagger.DefaultModelsExpandDepth(-1))

		handler(c)
	})

	// Swagger spec - 动态生成swagger.json，使用当前请求的host
	r.engine.GET("/swagger.json", func(c *gin.Context) {
		// 从嵌入的文件系统中读取 swagger.json
		data, err := r.docsFS.ReadFile("docs/swagger.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load swagger.json"})
			return
		}

		// 动态替换host字段
		var swaggerDoc map[string]interface{}
		if err := json.Unmarshal(data, &swaggerDoc); err != nil {
			c.JSON(500, gin.H{"error": "Failed to parse swagger.json"})
			return
		}

		// 获取当前请求的host
		host := c.Request.Host
		swaggerDoc["host"] = host

		c.JSON(200, swaggerDoc)
	})
}

// GetHandler 获取Gin引擎
func (r *Router) GetHandler() *gin.Engine {
	return r.engine
}
