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
	router *routes.Router
}

// NewServer 创建新的服务器实例
func NewServer() *Server {
	// 获取服务类型配置
	envConfig := utils.GetGlobalEnvConfig()
	serviceTypeStr := envConfig.Get("SERVICE_TYPE")
	var serviceType services.ServiceType
	
	switch serviceTypeStr {
	case "langchaingo":
		serviceType = services.ServiceTypeLangchaingo
	default:
		serviceType = services.ServiceTypeFlowy
	}

	// 初始化全局服务
	if err := services.InitGlobalServices(serviceType); err != nil {
		utils.LogError("初始化全局服务失败: %v", err)
		slog.Error("初始化全局服务失败", "error", err)
		os.Exit(1)
	}

	// 创建处理器层
	chatHandler := handlers.NewChatHandlerFromGlobal()
	settingsHandler := handlers.NewSettingsHandlerFromGlobal()
	knowledgeHandler := handlers.NewKnowledgeHandlerFromGlobal()
	modelHandler := handlers.NewModelHandlerFromGlobal()
	versionHandler := handlers.NewVersionHandler(Version, BuildTime, GitCommit, GitBranch, GitTag)

	// 创建路由，传入嵌入的文件系统
	router := routes.NewRouter(chatHandler, settingsHandler, knowledgeHandler, modelHandler, versionHandler, staticFiles, docsFiles)

	return &Server{
		router: router,
	}
}

// Start 启动服务器
func (s *Server) Start() {
	envConfig := utils.GetGlobalEnvConfig()
	port := envConfig.Get("PORT")
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
	envConfig := utils.GetGlobalEnvConfig()
	port := envConfig.Get("PORT")
	serviceType := envConfig.Get("SERVICE_TYPE")

	// 所有信息只写入日志
	utils.LogInfo("========================================")
	utils.LogInfo("聊天后端服务启动")
	utils.LogInfo("版本: %s", Version)
	utils.LogInfo("构建时间: %s", BuildTime)
	utils.LogInfo("Git Commit: %s", GitCommit)
	utils.LogInfo("Git Branch: %s", GitBranch)
	if GitTag != "" {
		utils.LogInfo("Git Tag: %s", GitTag)
	}
	utils.LogInfo("工作目录: %s", getCurrentDir())
	
	// 显示服务类型和相关配置
	utils.LogInfo("服务类型: %s", serviceType)
	
	if serviceType == "langchaingo" {
		utils.LogInfo("LLM URL: %s", envConfig.Get("LANGCHAINO_LLM_BASE_URL"))
		utils.LogInfo("LLM Model: %s", envConfig.Get("LANGCHAINO_LLM_MODEL"))
		utils.LogInfo("Embedding URL: %s", envConfig.Get("LANGCHAINO_EMBEDDING_URL"))
		utils.LogInfo("Embedding Model: %s", envConfig.Get("LANGCHAINO_EMBEDDING_MODEL"))
		utils.LogInfo("Qdrant URL: %s", envConfig.Get("LANGCHAINO_QDRANT_URL"))
		utils.LogInfo("Docling URL: %s", envConfig.Get("LANGCHAINO_DOCLING_URL"))
		utils.LogInfo("SQLite DB: %s", envConfig.Get("LANGCHAINO_SQLITE_DB_PATH"))
		if password := envConfig.Get("LANGCHAINO_SQLITE_PASSWORD"); password != "" {
			utils.LogInfo("SQLite Password: [已设置]")
		} else {
			utils.LogInfo("SQLite Password: [未设置]")
		}
	} else {
		utils.LogInfo("Flowy API: %s", envConfig.Get("FLOWY_BASE_URL"))
	}
	
	utils.LogInfo("快捷方式 API: %s", envConfig.Get("SHORTCUT_API_URL"))
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
