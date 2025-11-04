package utils

import (
	"bufio"
	"fmt"
	"os"
	"strings"
	"sync"
)

// EnvConfig 环境配置管理器
type EnvConfig struct {
	configs map[string]string
	mu      sync.RWMutex
}

// 默认配置项（移除日志相关配置）
var defaultConfigs = map[string]string{
	"PORT":             "9090",
	"FLOWY_BASE_URL":   "http://10.18.13.10:8888/api/v1",
	"FLOWY_API_KEY":    "",
	"FLOWY_TOKEN":      "Basic c3dvcmQ6c3dvcmRfc2VjcmV0",
	"SHORTCUT_API_URL": "http://10.18.13.157:26034",
}

var globalEnvConfig *EnvConfig
var envOnce sync.Once

// GetGlobalEnvConfig 获取全局环境配置实例
func GetGlobalEnvConfig() *EnvConfig {
	envOnce.Do(func() {
		globalEnvConfig = NewEnvConfig()
	})
	return globalEnvConfig
}

// NewEnvConfig 创建新的环境配置管理器
func NewEnvConfig() *EnvConfig {
	config := &EnvConfig{
		configs: make(map[string]string),
	}

	// 加载配置：优先级 默认值 < .env文件 < OS环境变量
	config.load()

	// 确保 .env 文件存在
	config.ensureEnvFile()

	return config
}

// load 从 .env 文件和环境变量加载配置
func (c *EnvConfig) load() {
	c.mu.Lock()
	defer c.mu.Unlock()

	// 首先设置默认值
	for key, defaultValue := range defaultConfigs {
		c.configs[key] = defaultValue
	}

	// 然后从 .env 文件覆盖
	c.loadFromFile(".env")

	// 最后从 OS 环境变量覆盖（优先级最高）
	for key := range defaultConfigs {
		if value := os.Getenv(key); value != "" {
			c.configs[key] = value
		}
	}
}

// loadFromFile 从指定文件加载配置
func (c *EnvConfig) loadFromFile(filename string) {
	file, err := os.Open(filename)
	if err != nil {
		// 文件不存在，跳过
		return
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// 跳过空行和注释
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// 解析 KEY=VALUE 格式
		parts := strings.SplitN(line, "=", 2)
		if len(parts) == 2 {
			key := strings.TrimSpace(parts[0])
			value := strings.TrimSpace(parts[1])

			// 移除引号
			value = strings.Trim(value, `"'`)

			// 只加载我们关心的配置项
			if _, exists := defaultConfigs[key]; exists {
				c.configs[key] = value
			}
		}
	}
}

// ensureEnvFile 确保 .env 文件存在，如果不存在则创建
func (c *EnvConfig) ensureEnvFile() {
	// 检查 .env 文件是否存在
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		// 创建 .env 文件
		file, err := os.Create(".env")
		if err != nil {
			// 创建失败，跳过
			return
		}
		defer file.Close()

		// 写入默认配置
		fmt.Fprintf(file, "# 环境配置文件\n")
		fmt.Fprintf(file, "# 可以修改以下配置项，OS 环境变量会覆盖这些设置\n\n")

		fmt.Fprintf(file, "# 服务器配置\n")
		fmt.Fprintf(file, "PORT=%s\n\n", defaultConfigs["PORT"])

		fmt.Fprintf(file, "# Flowy SDK 配置\n")
		fmt.Fprintf(file, "FLOWY_BASE_URL=%s\n", defaultConfigs["FLOWY_BASE_URL"])
		fmt.Fprintf(file, "FLOWY_API_KEY=%s\n", defaultConfigs["FLOWY_API_KEY"])
		fmt.Fprintf(file, "FLOWY_TOKEN=%s\n\n", defaultConfigs["FLOWY_TOKEN"])

		fmt.Fprintf(file, "# 快捷方式服务配置\n")
		fmt.Fprintf(file, "SHORTCUT_API_URL=%s\n", defaultConfigs["SHORTCUT_API_URL"])
	}
}

// Get 获取配置值
func (c *EnvConfig) Get(key string) string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	if value, exists := c.configs[key]; exists {
		return value
	}

	return ""
}

// GetAll 获取所有配置
func (c *EnvConfig) GetAll() map[string]string {
	c.mu.RLock()
	defer c.mu.RUnlock()

	result := make(map[string]string)
	for key, value := range c.configs {
		result[key] = value
	}
	return result
}

// GetEnvOrDefault 兼容原有接口的函数
func GetEnvOrDefault(key, defaultValue string) string {
	config := GetGlobalEnvConfig()
	if value := config.Get(key); value != "" {
		return value
	}
	return defaultValue
}
