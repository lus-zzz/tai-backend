package config

import (
	"os"
	"strconv"
	"time"
)

// Config Flowy SDK 配置
type Config struct {
	// 基础配置
	BaseURL string `json:"base_url" yaml:"base_url"`
	Timeout int    `json:"timeout" yaml:"timeout"` // 超时时间（秒）

	// 认证配置
	APIKey    string `json:"api_key" yaml:"api_key"`
	SecretKey string `json:"secret_key" yaml:"secret_key"`
	Token     string `json:"token" yaml:"token"`

	// HTTP配置
	MaxRetries    int  `json:"max_retries" yaml:"max_retries"`
	RetryInterval int  `json:"retry_interval" yaml:"retry_interval"` // 重试间隔（秒）
	SkipTLSVerify bool `json:"skip_tls_verify" yaml:"skip_tls_verify"`

	// 代理配置
	ProxyURL string `json:"proxy_url" yaml:"proxy_url"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		BaseURL:       "http://192.168.1.2:8888/api/v1",
		Timeout:       30,
		MaxRetries:    3,
		RetryInterval: 1,
		SkipTLSVerify: false,
	}
}

// LoadFromEnv 从环境变量加载配置
func (c *Config) LoadFromEnv() *Config {
	if baseURL := os.Getenv("FLOWY_BASE_URL"); baseURL != "" {
		c.BaseURL = baseURL
	}

	if timeout := os.Getenv("EAM_TIMEOUT"); timeout != "" {
		if t, err := strconv.Atoi(timeout); err == nil {
			c.Timeout = t
		}
	}

	if apiKey := os.Getenv("FLOWY_API_KEY"); apiKey != "" {
		c.APIKey = apiKey
	}

	if secretKey := os.Getenv("EAM_SECRET_KEY"); secretKey != "" {
		c.SecretKey = secretKey
	}

	if token := os.Getenv("FLOWY_TOKEN"); token != "" {
		c.Token = token
	}

	if maxRetries := os.Getenv("EAM_MAX_RETRIES"); maxRetries != "" {
		if r, err := strconv.Atoi(maxRetries); err == nil {
			c.MaxRetries = r
		}
	}

	if retryInterval := os.Getenv("EAM_RETRY_INTERVAL"); retryInterval != "" {
		if r, err := strconv.Atoi(retryInterval); err == nil {
			c.RetryInterval = r
		}
	}

	if skipTLS := os.Getenv("EAM_SKIP_TLS_VERIFY"); skipTLS != "" {
		if skip, err := strconv.ParseBool(skipTLS); err == nil {
			c.SkipTLSVerify = skip
		}
	}

	if proxyURL := os.Getenv("EAM_PROXY_URL"); proxyURL != "" {
		c.ProxyURL = proxyURL
	}

	return c
}

// GetTimeoutDuration 获取超时时间
func (c *Config) GetTimeoutDuration() time.Duration {
	return time.Duration(c.Timeout) * time.Second
}

// GetRetryIntervalDuration 获取重试间隔时间
func (c *Config) GetRetryIntervalDuration() time.Duration {
	return time.Duration(c.RetryInterval) * time.Second
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.BaseURL == "" {
		return &ConfigError{Field: "base_url", Message: "base URL is required"}
	}

	if c.Timeout <= 0 {
		return &ConfigError{Field: "timeout", Message: "timeout must be greater than 0"}
	}

	if c.MaxRetries < 0 {
		return &ConfigError{Field: "max_retries", Message: "max retries cannot be negative"}
	}

	if c.RetryInterval < 0 {
		return &ConfigError{Field: "retry_interval", Message: "retry interval cannot be negative"}
	}

	return nil
}

// ConfigError 配置错误
type ConfigError struct {
	Field   string `json:"field"`
	Message string `json:"message"`
}

// Error 实现 error 接口
func (e *ConfigError) Error() string {
	return e.Field + ": " + e.Message
}
