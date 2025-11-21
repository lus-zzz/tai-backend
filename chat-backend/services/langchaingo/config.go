package langchaingo

import (
	"fmt"
	"strconv"
	
	"chat-backend/utils"
)

// LangchaingoConfig Langchaingo 配置结构
type LangchaingoConfig struct {
	// LLM 配置
	LLM LLMConfig `json:"llm"`
	
	// Embedding 配置
	Embedding EmbeddingConfig `json:"embedding"`
	
	// Qdrant 配置
	Qdrant QdrantConfig `json:"qdrant"`
	
	// Docling 配置
	Docling DoclingConfig `json:"docling"`
	
	// SQLite 配置
	SQLite SQLiteConfig `json:"sqlite"`
}

// LLMConfig LLM 配置
type LLMConfig struct {
	BaseURL string `json:"base_url"`
	Token   string `json:"token"`
	Model   string `json:"model"`
}

// EmbeddingConfig Embedding 配置
type EmbeddingConfig struct {
	BaseURL string `json:"base_url"`
	Model   string `json:"model"`
}

// QdrantConfig Qdrant 配置
type QdrantConfig struct {
	URL       string `json:"url"`
	APIKey    string `json:"api_key"`
	VectorSize int    `json:"vector_size"`
}

// DoclingConfig Docling 配置
type DoclingConfig struct {
	BaseURL string `json:"base_url"`
	APIKey  string `json:"api_key"`
}

// SQLiteConfig SQLite 配置
type SQLiteConfig struct {
	DBPath   string `json:"db_path"`
	Password string `json:"password"`
}

// GetLangchaingoConfig 从统一环境配置获取 Langchaingo 配置
func GetLangchaingoConfig() *LangchaingoConfig {
	envConfig := utils.GetGlobalEnvConfig()
	
	// 解析向量维度
	vectorSize := 1024 // 默认值
	if sizeStr := envConfig.Get("LANGCHAINO_QDRANT_VECTOR_SIZE"); sizeStr != "" {
		if parsed, err := strconv.Atoi(sizeStr); err == nil {
			vectorSize = parsed
		}
	}
	
	return &LangchaingoConfig{
		LLM: LLMConfig{
			BaseURL: envConfig.Get("LANGCHAINO_LLM_BASE_URL"),
			Token:   envConfig.Get("LANGCHAINO_LLM_API_KEY"),
			Model:   envConfig.Get("LANGCHAINO_LLM_MODEL"),
		},
		Embedding: EmbeddingConfig{
			BaseURL: envConfig.Get("LANGCHAINO_EMBEDDING_URL"),
			Model:   envConfig.Get("LANGCHAINO_EMBEDDING_MODEL"),
		},
		Qdrant: QdrantConfig{
			URL:       envConfig.Get("LANGCHAINO_QDRANT_URL"),
			APIKey:    envConfig.Get("LANGCHAINO_QDRANT_API_KEY"),
			VectorSize: vectorSize,
		},
		Docling: DoclingConfig{
			BaseURL: envConfig.Get("LANGCHAINO_DOCLING_URL"),
			APIKey:  envConfig.Get("LANGCHAINO_DOCLING_API_KEY"),
		},
		SQLite: SQLiteConfig{
			DBPath:   envConfig.Get("LANGCHAINO_SQLITE_DB_PATH"),
			Password: envConfig.Get("LANGCHAINO_SQLITE_PASSWORD"),
		},
	}
}

// ValidateConfig 验证配置
func (c *LangchaingoConfig) ValidateConfig() error {
	// 验证必需的配置项
	if c.LLM.Token == "" {
		return fmt.Errorf("LANGCHAINO_LLM_API_KEY 不能为空")
	}
	
	if c.Embedding.BaseURL == "" {
		return fmt.Errorf("LANGCHAINO_EMBEDDING_URL 不能为空")
	}
	
	if c.Qdrant.URL == "" {
		return fmt.Errorf("LANGCHAINO_QDRANT_URL 不能为空")
	}
	
	if c.Docling.BaseURL == "" {
		return fmt.Errorf("LANGCHAINO_DOCLING_URL 不能为空")
	}
	
	// 设置默认值
	if c.LLM.BaseURL == "" {
		c.LLM.BaseURL = "https://api.openai.com/v1"
	}
	
	if c.LLM.Model == "" {
		c.LLM.Model = "gpt-3.5-turbo"
	}
	
	if c.Embedding.Model == "" {
		c.Embedding.Model = "bge-m3:latest"
	}
	
	if c.SQLite.DBPath == "" {
		c.SQLite.DBPath = "./chat_history.db"
	}
	
	return nil
}
