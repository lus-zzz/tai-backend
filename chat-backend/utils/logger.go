package utils

import (
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"sync"
	"time"
)

var (
	logger     *slog.Logger
	logFile    *os.File
	logConfig  *LogConfig
	loggerOnce sync.Once
	logMutex   sync.RWMutex
)

// LogConfig 日志配置
type LogConfig struct {
	LogDir       string // 日志目录
	MaxFileSize  int64  // 最大文件大小(字节)
	MaxFiles     int    // 最大文件数量
	EnableStdout bool   // 是否输出到标准输出
}

// InitLogger 初始化日志系统
func InitLogger(config *LogConfig) error {
	var initErr error
	loggerOnce.Do(func() {
		if config == nil {
			config = &LogConfig{
				LogDir:       "logs",
				MaxFileSize:  10 * 1024 * 1024, // 10MB
				MaxFiles:     10,
				EnableStdout: true,
			}
		}

		logConfig = config

		// 创建日志目录
		if err := os.MkdirAll(config.LogDir, 0755); err != nil {
			initErr = fmt.Errorf("创建日志目录失败: %w", err)
			return
		}

		// 创建日志文件
		logPath := filepath.Join(config.LogDir, fmt.Sprintf("chat-backend-%s.log", time.Now().Format("2006-01-02")))
		file, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
		if err != nil {
			initErr = fmt.Errorf("创建日志文件失败: %w", err)
			return
		}

		logFile = file

		// 创建日志writer
		var writer io.Writer
		if config.EnableStdout {
			writer = io.MultiWriter(os.Stdout, file)
		} else {
			writer = file
		}

		// 创建JSON格式的logger
		logger = slog.New(slog.NewJSONHandler(writer, &slog.HandlerOptions{
			Level: slog.LevelInfo,
		}))

		slog.SetDefault(logger)
	})

	return initErr
}

// CloseLogger 关闭日志系统
func CloseLogger() {
	logMutex.Lock()
	defer logMutex.Unlock()

	if logFile != nil {
		logFile.Close()
		logFile = nil
	}
}

// GetLogPath 获取当前日志文件路径
func GetLogPath() string {
	if logConfig == nil {
		return ""
	}
	return filepath.Join(logConfig.LogDir, fmt.Sprintf("chat-backend-%s.log", time.Now().Format("2006-01-02")))
}

// LogInfo 记录信息日志
func LogInfo(format string, args ...interface{}) {
	if logger != nil {
		logger.Info(fmt.Sprintf(format, args...))
	}
}

// LogError 记录错误日志
func LogError(format string, args ...interface{}) {
	if logger != nil {
		logger.Error(fmt.Sprintf(format, args...))
	}
}

// LogWarn 记录警告日志
func LogWarn(format string, args ...interface{}) {
	if logger != nil {
		logger.Warn(fmt.Sprintf(format, args...))
	}
}

// LogDebug 记录调试日志
func LogDebug(format string, args ...interface{}) {
	if logger != nil {
		logger.Debug(fmt.Sprintf(format, args...))
	}
}

// InfoWith 记录带字段的信息日志
func InfoWith(msg string, args ...interface{}) {
	if logger != nil {
		logger.Info(msg, args...)
	}
}

// ErrorWith 记录带字段的错误日志
func ErrorWith(msg string, args ...interface{}) {
	if logger != nil {
		logger.Error(msg, args...)
	}
}

// WarnWith 记录带字段的警告日志
func WarnWith(msg string, args ...interface{}) {
	if logger != nil {
		logger.Warn(msg, args...)
	}
}

// DebugWith 记录带字段的调试日志
func DebugWith(msg string, args ...interface{}) {
	if logger != nil {
		logger.Debug(msg, args...)
	}
}
