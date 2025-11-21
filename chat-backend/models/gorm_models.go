package models

import (
	"time"

	"gorm.io/gorm"
)

// GORMModel 基础模型，包含通用字段
type GORMModel struct {
	ID        uint           `gorm:"primaryKey"`
	CreatedAt time.Time       `gorm:"autoCreateTime"`
	UpdatedAt time.Time       `gorm:"autoUpdateTime"`
	DeletedAt gorm.DeletedAt `gorm:"index"`
}

// ModelGORM 模型表对应的GORM结构
type ModelGORM struct {
	GORMModel
	Name        string `gorm:"not null;size:255"`
	Type        string `gorm:"not null;size:50"`
	Description string `gorm:"type:text"`
	Enabled     bool   `gorm:"default:true"`
}

// TableName 指定表名
func (ModelGORM) TableName() string {
	return "models"
}

// KnowledgeBaseGORM 知识库表对应的GORM结构
type KnowledgeBaseGORM struct {
	GORMModel
	Name          string `gorm:"not null;size:255"`
	Desc          string `gorm:"type:text"`
	VectorModel    int    `gorm:"not null;column:vector_model"`
	AgentModel     int    `gorm:"not null;column:agent_model"`
	ChunkStrategy string `gorm:"not null;size:50;column:chunk_strategy"`
	ChunkSize    int    `gorm:"not null;column:chunk_size"`
	FileCount     int    `gorm:"default:0"`
}

// TableName 指定表名
func (KnowledgeBaseGORM) TableName() string {
	return "knowledge_bases"
}

// KnowledgeBaseFileGORM 知识库文件表对应的GORM结构
type KnowledgeBaseFileGORM struct {
	GORMModel
	KnowledgeBaseID int    `gorm:"not null;column:knowledge_base_id;index"`
	Name           string `gorm:"not null;size:255"`
	Size           int64  `gorm:"not null"`
	Enable         bool   `gorm:"default:true"`
	Status         int    `gorm:"default:0"` // 0=构建中, 1=完成, 2=失败
	UploadedAt     time.Time `gorm:"autoCreateTime"`
	IndexPercent   int    `gorm:"default:0;column:index_percent"`
	ErrorMessage   string `gorm:"type:text"`
}

// TableName 指定表名
func (KnowledgeBaseFileGORM) TableName() string {
	return "knowledge_base_files"
}

// 转换函数：GORM模型 -> API模型

// ToModelInfo 将GORM模型转换为ModelInfo
func (m *ModelGORM) ToModelInfo() *ModelInfo {
	// 将字符串类型的 Type 转换为 int 类型
	var typeInt int
	role := "chat" // 默认角色
	if m.Type == "chat" {
		typeInt = 0
		role = "chat"
	} else if m.Type == "embedding" {
		typeInt = 1
		role = "embedding"
	}

	return &ModelInfo{
		ID:       int(m.ID),
		Name:     m.Name,
		Symbol:   m.Type, // 使用 type 作为 symbol
		Endpoint: "",     // GORM 模型中没有 endpoint 字段，留空
		Enable:   m.Enabled,
		Type:     typeInt,
		Role:     role,
	}
}

// ToKnowledgeBase 将GORM模型转换为KnowledgeBase
func (kb *KnowledgeBaseGORM) ToKnowledgeBase(fileCount int) *KnowledgeBase {
	return &KnowledgeBase{
		ID: int(kb.ID),
		KnowledgeBaseConfig: KnowledgeBaseConfig{
			Name:          kb.Name,
			Desc:          kb.Desc,
			VectorModel:    kb.VectorModel,
			AgentModel:     kb.AgentModel,
			ChunkStrategy:  kb.ChunkStrategy,
			ChunkSize:      kb.ChunkSize,
		},
		FileCount: fileCount,
	}
}

// ToKnowledgeFile 将GORM模型转换为KnowledgeFile
func (kf *KnowledgeBaseFileGORM) ToKnowledgeFile() *KnowledgeFile {
	return &KnowledgeFile{
		ID:           int(kf.ID),
		Name:         kf.Name,
		Size:         int(kf.Size),
		Enable:       kf.Enable,
		Status:       kf.Status,
		UploadedAt:   kf.UploadedAt,
		IndexPercent:  kf.IndexPercent,
		ErrorMessage: kf.ErrorMessage,
	}
}

// 转换函数：API模型 -> GORM模型

// NewModelGORM 从ModelInfo创建GORM模型
func NewModelGORM(mi *ModelInfo) *ModelGORM {
	// 将 int 类型的 Type 转换为 string 类型
	var typeStr string
	if mi.Type == 0 {
		typeStr = "chat"
	} else if mi.Type == 1 {
		typeStr = "embedding"
	} else {
		typeStr = "unknown"
	}

	return &ModelGORM{
		Name:        mi.Name,
		Type:        typeStr,
		Description: "", // ModelInfo 中没有 Description 字段，留空
		Enabled:     mi.Enable,
	}
}

// NewKnowledgeBaseGORM 从KnowledgeBaseCreateRequest创建GORM模型
func NewKnowledgeBaseGORM(req *KnowledgeBaseCreateRequest) *KnowledgeBaseGORM {
	return &KnowledgeBaseGORM{
		Name:          req.Name,
		Desc:          req.Desc,
		VectorModel:    req.VectorModel,
		AgentModel:     req.AgentModel,
		ChunkStrategy:  req.ChunkStrategy,
		ChunkSize:      req.ChunkSize,
	}
}

// UpdateKnowledgeBaseGORM 更新GORM模型
func UpdateKnowledgeBaseGORM(gormModel *KnowledgeBaseGORM, req *UpdateKnowledgeBaseRequest) {
	gormModel.Name = req.Name
	gormModel.Desc = req.Desc
	gormModel.VectorModel = req.VectorModel
	gormModel.AgentModel = req.AgentModel
	gormModel.ChunkStrategy = req.ChunkStrategy
	gormModel.ChunkSize = req.ChunkSize
}

// NewKnowledgeBaseFileGORM 创建新的知识库文件GORM模型
func NewKnowledgeBaseFileGORM(knowledgeBaseID int, filename string, size int64) *KnowledgeBaseFileGORM {
	return &KnowledgeBaseFileGORM{
		KnowledgeBaseID: knowledgeBaseID,
		Name:           filename,
		Size:           size,
		Enable:         true,
		Status:         1, // 完成
		IndexPercent:    100,
		ErrorMessage:   "",
	}
}
