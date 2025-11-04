package routes

import (
	"chat-backend/handlers"
	"chat-backend/utils"
	"embed"
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
	shortcutHandler  *handlers.ShortcutHandler
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
	shortcutHandler *handlers.ShortcutHandler,
	staticFS embed.FS,
	docsFS embed.FS,
) *Router {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()

	// 添加中间件
	engine.Use(gin.Logger())
	engine.Use(gin.Recovery())

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
		shortcutHandler:  shortcutHandler,
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
		chat.POST("/conversations", r.chatHandler.CreateConversation)
		chat.GET("/conversations", r.chatHandler.GetConversations)
		chat.DELETE("/conversations/:id", r.chatHandler.DeleteConversation)
		chat.GET("/conversations/:id/history", r.chatHandler.GetConversationHistory)
		chat.POST("/messages", r.chatHandler.SendMessage)
		chat.GET("/conversations/:id/settings", r.chatHandler.GetConversationSettings)
		chat.PUT("/conversations/:id/settings", r.chatHandler.UpdateConversationSettings)
	}

	// 默认配置相关路由
	settings := api.Group("/settings")
	{
		settings.GET("/defaults", r.settingsHandler.GetDefaultSettings)
		settings.PUT("/defaults", r.settingsHandler.UpdateDefaultSettings)
		settings.POST("/defaults/reset", r.settingsHandler.ResetDefaultSettings)
	}

	// 知识库相关路由
	knowledge := api.Group("/knowledge")
	{
		knowledge.GET("/bases", r.knowledgeHandler.ListKnowledgeBases)
		knowledge.POST("/bases", r.knowledgeHandler.CreateKnowledgeBase)
		knowledge.PUT("/bases/:id", r.knowledgeHandler.UpdateKnowledgeBase)
		knowledge.DELETE("/bases/:id", r.knowledgeHandler.DeleteKnowledgeBase)
		knowledge.GET("/bases/:id/files", r.knowledgeHandler.GetKnowledgeBaseFiles)
		knowledge.POST("/bases/:id/files", r.knowledgeHandler.UploadFile)

		// 文件操作路由（只需要文件ID）
		knowledge.DELETE("/files/:file_id", r.knowledgeHandler.DeleteFile)
		knowledge.PUT("/files/:file_id/toggle", r.knowledgeHandler.ToggleFileEnable)
	}

	// 模型相关路由
	models := api.Group("/models")
	{
		// 支持的模型列表
		models.GET("/supported/chat", r.modelHandler.ListSupportedChatModels)
		models.GET("/supported/vector", r.modelHandler.ListSupportedVectorModels)

		// 可用的模型列表
		models.GET("/available/all", r.modelHandler.ListAvailableAllModels)
		models.GET("/available/chat", r.modelHandler.ListAvailableChatModels)
		models.GET("/available/vector", r.modelHandler.ListAvailableVectorModels)

		// 模型管理
		models.POST("", r.modelHandler.SaveModel)
		models.DELETE("/:id", r.modelHandler.DeleteModel)
		models.PUT("/:id/status", r.modelHandler.SetModelStatus)
	}

	// 快捷方式相关路由
	shortcut := api.Group("/shortcut")
	{
		shortcut.POST("/recommend", r.shortcutHandler.RecommendSettings)
		shortcut.POST("/supportedSetting", r.shortcutHandler.GetSupportedSettings)
	}

	// 系统信息路由
	api.GET("/version", r.versionHandler.GetVersion)

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

	// Swagger UI - 使用ginSwagger中间件，指定spec URL
	r.engine.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler,
		ginSwagger.URL("http://localhost:9090/swagger.json"),
		ginSwagger.DefaultModelsExpandDepth(-1)))

	// Swagger spec - 使用嵌入的文件
	r.engine.GET("/swagger.json", func(c *gin.Context) {
		// 从嵌入的文件系统中读取 swagger.json
		data, err := r.docsFS.ReadFile("docs/swagger.json")
		if err != nil {
			c.JSON(500, gin.H{"error": "Failed to load swagger.json"})
			return
		}
		c.Data(200, "application/json", data)
	})
}

// GetHandler 获取Gin引擎
func (r *Router) GetHandler() *gin.Engine {
	return r.engine
}
