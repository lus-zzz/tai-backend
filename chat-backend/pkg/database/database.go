package database

import (
	"fmt"
	"time"

	"chat-backend/models"
	"chat-backend/utils"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// Database 数据库封装结构
type Database struct {
	db *gorm.DB
}

// NewDatabase 创建新的数据库连接
func NewDatabase(dbPath string) (*Database, error) {
	if dbPath == "" {
		dbPath = "./chat_history.db" // 默认路径
	}

	// 连接数据库
	db, err := gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		NowFunc: func() time.Time {
			return time.Now()
		},
	})
	if err != nil {
		return nil, fmt.Errorf("连接数据库失败: %w", err)
	}

	database := &Database{
		db: db,
	}

	// 自动迁移表结构
	if err := database.AutoMigrate(); err != nil {
		return nil, fmt.Errorf("自动迁移失败: %w", err)
	}

	// 初始化默认数据
	if err := database.initDefaultData(); err != nil {
		utils.WarnWith("初始化默认数据失败", "error", err.Error())
	}

	utils.InfoWith("数据库初始化完成", "db_path", dbPath)
	return database, nil
}

// AutoMigrate 自动迁移数据库表结构
func (d *Database) AutoMigrate() error {
	return d.db.AutoMigrate(
		&models.ModelGORM{},
		&models.KnowledgeBaseGORM{},
		&models.KnowledgeBaseFileGORM{},
	)
}

// GetDB 获取GORM数据库实例
func (d *Database) GetDB() *gorm.DB {
	return d.db
}

// Close 关闭数据库连接
func (d *Database) Close() error {
	sqlDB, err := d.db.DB()
	if err != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", err)
	}
	return sqlDB.Close()
}

// initDefaultData 初始化默认数据
func (d *Database) initDefaultData() error {
	// 初始化默认模型
	return d.initDefaultModels()
}

// initDefaultModels 初始化默认模型数据
func (d *Database) initDefaultModels() error {
	// 检查是否已有模型数据
	var count int64
	if err := d.db.Model(&models.ModelGORM{}).Count(&count).Error; err != nil {
		return fmt.Errorf("查询模型数量失败: %w", err)
	}

	// 如果已有数据，跳过初始化
	if count > 0 {
		return nil
	}

	// 插入默认模型
	defaultModels := []models.ModelGORM{
		{Name: "gpt-3.5-turbo", Type: "chat", Description: "OpenAI GPT-3.5 Turbo", Enabled: true},
		{Name: "gpt-4", Type: "chat", Description: "OpenAI GPT-4", Enabled: true},
		{Name: "bge-m3:latest", Type: "embedding", Description: "BGE-M3 嵌入模型", Enabled: true},
	}

	for _, model := range defaultModels {
		if err := d.db.Create(&model).Error; err != nil {
			return fmt.Errorf("插入默认模型失败: %w", err)
		}
	}

	utils.InfoWith("默认模型初始化完成", "count", len(defaultModels))
	return nil
}

// === 模型相关操作 ===

// GetAllModels 获取所有模型
func (d *Database) GetAllModels() ([]models.ModelInfo, error) {
	var gormModels []models.ModelGORM
	if err := d.db.Find(&gormModels).Error; err != nil {
		return nil, fmt.Errorf("查询所有模型失败: %w", err)
	}

	var modelInfos []models.ModelInfo
	for _, m := range gormModels {
		modelInfos = append(modelInfos, *m.ToModelInfo())
	}

	return modelInfos, nil
}

// GetModelByID 根据ID获取模型
func (d *Database) GetModelByID(id int) (*models.ModelInfo, error) {
	var gormModel models.ModelGORM
	if err := d.db.First(&gormModel, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("模型ID %d 不存在", id)
		}
		return nil, fmt.Errorf("查询模型失败: %w", err)
	}

	return gormModel.ToModelInfo(), nil
}

// GetModelsByType 根据类型获取模型
func (d *Database) GetModelsByType(modelType string) ([]models.ModelInfo, error) {
	var gormModels []models.ModelGORM
	if err := d.db.Where("type = ? AND enabled = ?", modelType, true).Find(&gormModels).Error; err != nil {
		return nil, fmt.Errorf("查询模型失败: %w", err)
	}

	var modelInfos []models.ModelInfo
	for _, m := range gormModels {
		modelInfos = append(modelInfos, *m.ToModelInfo())
	}

	return modelInfos, nil
}

// SaveModel 保存模型（新建或更新）
func (d *Database) SaveModel(req *models.ModelSaveRequest) (int, error) {
	var gormModel models.ModelGORM
	
	if req.ID > 0 {
		// 更新现有模型
		if err := d.db.First(&gormModel, req.ID).Error; err != nil {
			return 0, fmt.Errorf("模型ID %d 不存在: %w", req.ID, err)
		}
		
		// 更新字段
		gormModel.Name = req.Name
		gormModel.Enabled = req.Enable
		// 注意：GORM模型中只有基本的name, type, description, enabled字段
		// endpoint, symbol, credentials等需要扩展表结构或使用其他方式存储
		
		if err := d.db.Save(&gormModel).Error; err != nil {
			return 0, fmt.Errorf("更新模型失败: %w", err)
		}
	} else {
		// 创建新模型
		typeStr := "chat" // 默认为chat
		if req.Type == 1 {
			typeStr = "embedding"
		}
		
		gormModel = models.ModelGORM{
			Name:        req.Name,
			Type:        typeStr,
			Description: "", // 暂时留空
			Enabled:     req.Enable,
		}
		
		if err := d.db.Create(&gormModel).Error; err != nil {
			return 0, fmt.Errorf("创建模型失败: %w", err)
		}
	}
	
	return int(gormModel.ID), nil
}

// DeleteModel 删除模型
func (d *Database) DeleteModel(id int) error {
	if err := d.db.Delete(&models.ModelGORM{}, id).Error; err != nil {
		return fmt.Errorf("删除模型失败: %w", err)
	}
	return nil
}

// UpdateModelStatus 更新模型状态
func (d *Database) UpdateModelStatus(id int, enable bool) error {
	result := d.db.Model(&models.ModelGORM{}).Where("id = ?", id).Update("enabled", enable)
	if result.Error != nil {
		return fmt.Errorf("更新模型状态失败: %w", result.Error)
	}
	if result.RowsAffected == 0 {
		return fmt.Errorf("模型ID %d 不存在", id)
	}
	return nil
}

// GetSupportedChatModels 获取支持的聊天模型列表（硬编码）
func (d *Database) GetSupportedChatModels() ([]models.SupportedChatModel, error) {
	// 返回 langchaingo 支持的聊天模型
	return []models.SupportedChatModel{
		{
			Name: "gpt-3.5-turbo",
			Identify: "openai-chat",
			LLMProperty: models.LLMProperty{
				Model:           "gpt-3.5-turbo",
				MaxToken:        4096,
				ContextLength:   4096,
				TokenProportion: 1,
				Stream:          true,
				SupportTool:     true,
				FixQwqThink:    false,
			},
		},
		{
			Name: "gpt-4",
			Identify: "openai-chat",
			LLMProperty: models.LLMProperty{
				Model:           "gpt-4",
				MaxToken:        8192,
				ContextLength:   8192,
				TokenProportion: 1,
				Stream:          true,
				SupportTool:     true,
				FixQwqThink:    false,
			},
		},
	}, nil
}

// GetSupportedVectorModels 获取支持的向量模型列表（硬编码）
func (d *Database) GetSupportedVectorModels() ([]models.SupportedVectorModel, error) {
	// 返回 langchaingo 支持的向量模型
	return []models.SupportedVectorModel{
		{
			Name: "bge-m3:latest",
			Identify: "ollama-embedding",
			LLMProperty: models.EmbeddingProperty{
				Model:              "bge-m3:latest",
				EmbeddingMaxLength:  8192,
				EmbeddingDimension: 1024,
				BatchLimit:         100,
			},
		},
	}, nil
}

// GetAvailableModelsByType 根据类型获取可用模型（启用状态）
func (d *Database) GetAvailableModelsByType(modelType int) ([]models.ModelInfo, error) {
	var typeStr string
	if modelType == 0 {
		typeStr = "chat"
	} else if modelType == 1 {
		typeStr = "embedding"
	} else {
		return nil, fmt.Errorf("不支持的模型类型: %d", modelType)
	}

	var gormModels []models.ModelGORM
	if err := d.db.Where("type = ? AND enabled = ?", typeStr, true).Find(&gormModels).Error; err != nil {
		return nil, fmt.Errorf("查询模型失败: %w", err)
	}

	var modelInfos []models.ModelInfo
	for _, m := range gormModels {
		modelInfos = append(modelInfos, *m.ToModelInfo())
	}

	return modelInfos, nil
}

// === 知识库相关操作 ===

// CreateKnowledgeBase 创建知识库
func (d *Database) CreateKnowledgeBase(req *models.KnowledgeBaseCreateRequest) (*models.KnowledgeBase, error) {
	// 验证模型是否存在
	var vectorModel models.ModelGORM
	if err := d.db.First(&vectorModel, req.VectorModel).Error; err != nil {
		return nil, fmt.Errorf("向量模型ID %d 不存在", req.VectorModel)
	}
	if vectorModel.Type != "embedding" {
		return nil, fmt.Errorf("向量模型ID %d 不是嵌入类型", req.VectorModel)
	}

	var agentModel models.ModelGORM
	if err := d.db.First(&agentModel, req.AgentModel).Error; err != nil {
		return nil, fmt.Errorf("对话模型ID %d 不存在", req.AgentModel)
	}
	if agentModel.Type != "chat" {
		return nil, fmt.Errorf("对话模型ID %d 不是聊天类型", req.AgentModel)
	}

	// 创建知识库
	gormKB := models.NewKnowledgeBaseGORM(req)
	if err := d.db.Create(gormKB).Error; err != nil {
		return nil, fmt.Errorf("创建知识库失败: %w", err)
	}

	// 查询完整记录
	var result models.KnowledgeBaseGORM
	if err := d.db.Preload("Files").First(&result, gormKB.ID).Error; err != nil {
		return nil, fmt.Errorf("查询创建的知识库失败: %w", err)
	}

	// 计算文件数量
	var fileCount int64
	d.db.Model(&models.KnowledgeBaseFileGORM{}).Where("knowledge_base_id = ?", result.ID).Count(&fileCount)

	return result.ToKnowledgeBase(int(fileCount)), nil
}

// GetKnowledgeBaseByID 根据ID获取知识库
func (d *Database) GetKnowledgeBaseByID(id int) (*models.KnowledgeBase, error) {
	var gormKB models.KnowledgeBaseGORM
	if err := d.db.First(&gormKB, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("知识库ID %d 不存在", id)
		}
		return nil, fmt.Errorf("查询知识库失败: %w", err)
	}

	// 计算文件数量
	var fileCount int64
	d.db.Model(&models.KnowledgeBaseFileGORM{}).Where("knowledge_base_id = ?", gormKB.ID).Count(&fileCount)

	return gormKB.ToKnowledgeBase(int(fileCount)), nil
}

// ListKnowledgeBases 获取知识库列表
func (d *Database) ListKnowledgeBases() ([]models.KnowledgeBase, error) {
	var gormKBs []models.KnowledgeBaseGORM
	if err := d.db.Find(&gormKBs).Error; err != nil {
		return nil, fmt.Errorf("查询知识库列表失败: %w", err)
	}

	var knowledgeBases []models.KnowledgeBase
	for _, gormKB := range gormKBs {
		// 计算每个知识库的文件数量
		var fileCount int64
		d.db.Model(&models.KnowledgeBaseFileGORM{}).Where("knowledge_base_id = ?", gormKB.ID).Count(&fileCount)
		
		knowledgeBases = append(knowledgeBases, *gormKB.ToKnowledgeBase(int(fileCount)))
	}

	return knowledgeBases, nil
}

// UpdateKnowledgeBase 更新知识库
func (d *Database) UpdateKnowledgeBase(id int, req *models.UpdateKnowledgeBaseRequest) (*models.KnowledgeBase, error) {
	// 验证模型是否存在
	var vectorModel models.ModelGORM
	if err := d.db.First(&vectorModel, req.VectorModel).Error; err != nil {
		return nil, fmt.Errorf("向量模型ID %d 不存在", req.VectorModel)
	}
	if vectorModel.Type != "embedding" {
		return nil, fmt.Errorf("向量模型ID %d 不是嵌入类型", req.VectorModel)
	}

	var agentModel models.ModelGORM
	if err := d.db.First(&agentModel, req.AgentModel).Error; err != nil {
		return nil, fmt.Errorf("对话模型ID %d 不存在", req.AgentModel)
	}
	if agentModel.Type != "chat" {
		return nil, fmt.Errorf("对话模型ID %d 不是聊天类型", req.AgentModel)
	}

	// 更新知识库
	var gormKB models.KnowledgeBaseGORM
	if err := d.db.First(&gormKB, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("知识库ID %d 不存在", id)
		}
		return nil, fmt.Errorf("查询知识库失败: %w", err)
	}

	models.UpdateKnowledgeBaseGORM(&gormKB, req)
	if err := d.db.Save(&gormKB).Error; err != nil {
		return nil, fmt.Errorf("更新知识库失败: %w", err)
	}

	// 计算文件数量
	var fileCount int64
	d.db.Model(&models.KnowledgeBaseFileGORM{}).Where("knowledge_base_id = ?", gormKB.ID).Count(&fileCount)

	return gormKB.ToKnowledgeBase(int(fileCount)), nil
}

// DeleteKnowledgeBase 删除知识库
func (d *Database) DeleteKnowledgeBase(id int) error {
	// 使用事务删除
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 删除相关文件
		if err := tx.Where("knowledge_base_id = ?", id).Delete(&models.KnowledgeBaseFileGORM{}).Error; err != nil {
			return fmt.Errorf("删除知识库文件失败: %w", err)
		}

		// 删除知识库
		if err := tx.Delete(&models.KnowledgeBaseGORM{}, id).Error; err != nil {
			return fmt.Errorf("删除知识库失败: %w", err)
		}

		return nil
	})
}

// === 知识库文件相关操作 ===

// CreateKnowledgeBaseFile 创建知识库文件
func (d *Database) CreateKnowledgeBaseFile(knowledgeBaseID int, filename string, size int64) (*models.KnowledgeFile, error) {
	gormFile := models.NewKnowledgeBaseFileGORM(knowledgeBaseID, filename, size)
	if err := d.db.Create(gormFile).Error; err != nil {
		return nil, fmt.Errorf("创建知识库文件失败: %w", err)
	}

	return gormFile.ToKnowledgeFile(), nil
}

// GetKnowledgeBaseFiles 获取知识库文件列表
func (d *Database) GetKnowledgeBaseFiles(knowledgeBaseID int) ([]models.KnowledgeFile, error) {
	var gormFiles []models.KnowledgeBaseFileGORM
	if err := d.db.Where("knowledge_base_id = ?", knowledgeBaseID).Order("uploaded_at DESC").Find(&gormFiles).Error; err != nil {
		return nil, fmt.Errorf("查询知识库文件列表失败: %w", err)
	}

	var knowledgeFiles []models.KnowledgeFile
	for _, gormFile := range gormFiles {
		knowledgeFiles = append(knowledgeFiles, *gormFile.ToKnowledgeFile())
	}

	return knowledgeFiles, nil
}

// DeleteKnowledgeBaseFile 删除知识库文件
func (d *Database) DeleteKnowledgeBaseFile(fileID int) error {
	// 使用事务删除
	return d.db.Transaction(func(tx *gorm.DB) error {
		// 查询文件所属的知识库ID
		var gormFile models.KnowledgeBaseFileGORM
		if err := tx.First(&gormFile, fileID).Error; err != nil {
			if err == gorm.ErrRecordNotFound {
				return fmt.Errorf("文件ID %d 不存在", fileID)
			}
			return fmt.Errorf("查询文件失败: %w", err)
		}

		// 删除文件
		if err := tx.Delete(&gormFile).Error; err != nil {
			return fmt.Errorf("删除文件记录失败: %w", err)
		}

		// 更新知识库文件计数
		var fileCount int64
		tx.Model(&models.KnowledgeBaseFileGORM{}).Where("knowledge_base_id = ?", gormFile.KnowledgeBaseID).Count(&fileCount)
		
		// 这里可以添加更新文件计数的逻辑，但为了简化，暂时跳过

		return nil
	})
}

// ToggleKnowledgeBaseFileEnable 切换文件启用状态
func (d *Database) ToggleKnowledgeBaseFileEnable(fileID int, enable bool) error {
	result := d.db.Model(&models.KnowledgeBaseFileGORM{}).Where("id = ?", fileID).Update("enable", enable)
	if result.Error != nil {
		return fmt.Errorf("更新文件启用状态失败: %w", result.Error)
	}

	if result.RowsAffected == 0 {
		return fmt.Errorf("文件ID %d 不存在", fileID)
	}

	return nil
}
