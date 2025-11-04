package main

import (
	"chat-backend/handlers"
	"chat-backend/routes"
	"chat-backend/services"
	"chat-backend/utils"
	"embed"
	"flag"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"

	"flowy-sdk"
	"flowy-sdk/pkg/config"
)

// 嵌入静态文件
//
//go:embed static/*
var staticFiles embed.FS

// 嵌入文档文件
//
//go:embed docs/*
var docsFiles embed.FS

// 版本信息，在编译时通过 -ldflags 注入
var (
	Version   = "dev"
	BuildTime = "unknown"
	GitCommit = "unknown"
	GitBranch = "unknown"
	GitTag    = ""
)

// Server HTTP服务器结构
type Server struct {
	router           *routes.Router
	chatService      *services.ChatService
	knowledgeService *services.KnowledgeService
	modelService     *services.ModelService
}

// NewServer 创建新的服务器实例
func NewServer() *Server {
	// 配置Flowy SDK
	flowyConfig := &config.Config{
		BaseURL: utils.GetEnvOrDefault("FLOWY_BASE_URL", "http://10.18.13.10:8888/api/v1"),
		APIKey:  utils.GetEnvOrDefault("FLOWY_API_KEY", ""),
		Token:   utils.GetEnvOrDefault("FLOWY_TOKEN", "Basic c3dvcmQ6c3dvcmRfc2VjcmV0"),
		Timeout: 30,
	}

	// 创建Flowy SDK实例
	sdk := flowy.New(flowyConfig)

	// 创建默认配置服务(不使用嵌入文件)
	defaultSettingsService := services.NewDefaultSettingsService()

	// 创建服务层
	chatService := services.NewChatService(sdk, defaultSettingsService)
	knowledgeService := services.NewKnowledgeService(sdk)
	modelService := services.NewModelService(sdk)

	// 创建快捷方式服务 - 这里需要配置实际的API地址
	shortcutAPIURL := utils.GetEnvOrDefault("SHORTCUT_API_URL", "http://10.18.13.157:26034")
	shortcutService := services.NewShortcutService(shortcutAPIURL)

	// 创建处理器层
	chatHandler := handlers.NewChatHandler(chatService, defaultSettingsService)
	settingsHandler := handlers.NewSettingsHandler(defaultSettingsService)
	knowledgeHandler := handlers.NewKnowledgeHandler(knowledgeService)
	modelHandler := handlers.NewModelHandler(modelService)
	versionHandler := handlers.NewVersionHandler(Version, BuildTime, GitCommit, GitBranch, GitTag)
	shortcutHandler := handlers.NewShortcutHandler(shortcutService)

	// 创建路由，传入嵌入的文件系统
	router := routes.NewRouter(chatHandler, settingsHandler, knowledgeHandler, modelHandler, versionHandler, shortcutHandler, staticFiles, docsFiles)

	return &Server{
		router:           router,
		chatService:      chatService,
		knowledgeService: knowledgeService,
		modelService:     modelService,
	}
}

// Start 启动服务器
func (s *Server) Start() {
	port := utils.GetEnvOrDefault("PORT", "9090")
	addr := ":" + port

	utils.LogInfo("服务器监听端口: %s", port)
	utils.LogInfo("========================================")

	if err := http.ListenAndServe(addr, s.router.GetHandler()); err != nil {
		utils.LogError("启动服务器失败: %v", err)
		slog.Error("启动服务器失败", "error", err)
		os.Exit(1)
	}
}

func main() {
	// 定义命令行参数
	showVersion := flag.Bool("version", false, "显示版本信息")
	flag.Parse()

	// 如果指定了 --version，显示版本信息后退出
	if *showVersion {
		printVersion()
		os.Exit(0)
	}

	// 初始化日志系统
	logConfig := &utils.LogConfig{
		LogDir:       "logs",
		MaxFileSize:  10 * 1024 * 1024, // 10MB
		MaxFiles:     10,
		EnableStdout: true,
	}

	if err := utils.InitLogger(logConfig); err != nil {
		slog.Error("初始化日志系统失败", "error", err)
		os.Exit(1)
	}
	defer utils.CloseLogger()

	// 显示启动信息
	printBanner()

	// 创建并启动服务器
	server := NewServer()
	server.Start()
}

// printVersion 打印版本信息
func printVersion() {
	fmt.Printf("Chat Backend\n")
	fmt.Printf("Version:    %s\n", Version)
	fmt.Printf("Build Time: %s\n", BuildTime)
	fmt.Printf("Git Commit: %s\n", GitCommit)
	fmt.Printf("Git Branch: %s\n", GitBranch)
	if GitTag != "" {
		fmt.Printf("Git Tag:    %s\n", GitTag)
	}
}

// printBanner 打印启动横幅
func printBanner() {
	port := utils.GetEnvOrDefault("PORT", "9090")

	// 所有信息只写入日志
	utils.LogInfo("========================================")
	utils.LogInfo("Flowy 聊天后端服务启动")
	utils.LogInfo("版本: %s", Version)
	utils.LogInfo("构建时间: %s", BuildTime)
	utils.LogInfo("Git Commit: %s", GitCommit)
	utils.LogInfo("Git Branch: %s", GitBranch)
	if GitTag != "" {
		utils.LogInfo("Git Tag: %s", GitTag)
	}
	utils.LogInfo("工作目录: %s", getCurrentDir())
	utils.LogInfo("Flowy API: %s", utils.GetEnvOrDefault("FLOWY_BASE_URL", "http://10.18.13.10:8888/api/v1"))
	utils.LogInfo("服务端口: %s", port)
	utils.LogInfo("API 文档: http://localhost:%s/swagger/index.html", port)
	utils.LogInfo("版本信息: http://localhost:%s/api/v1/version", port)
	utils.LogInfo("日志查看器: http://localhost:%s/static/log-viewer.html", port)
}

func getCurrentDir() string {
	dir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		return "unknown"
	}
	return dir
}
