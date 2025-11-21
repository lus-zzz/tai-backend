package interfaces

import (
	"context"
	"io"

	"chat-backend/models"
)

// KnowledgeServiceInterface 知识库服务接口
type KnowledgeServiceInterface interface {
	// ListKnowledgeBases 获取知识库列表
	ListKnowledgeBases(ctx context.Context) ([]models.KnowledgeBase, error)

	// CreateKnowledgeBase 创建知识库
	CreateKnowledgeBase(ctx context.Context, req *models.KnowledgeBaseCreateRequest) (*models.KnowledgeBase, error)

	// UpdateKnowledgeBase 更新知识库
	UpdateKnowledgeBase(ctx context.Context, id int, req *models.UpdateKnowledgeBaseRequest) (*models.KnowledgeBase, error)

	// DeleteKnowledgeBase 删除知识库
	DeleteKnowledgeBase(ctx context.Context, id int) error

	// GetKnowledgeBaseFiles 获取知识库文件列表
	GetKnowledgeBaseFiles(ctx context.Context, id int) ([]models.KnowledgeFile, error)

	// UploadFile 上传文件到知识库（文件流上传）
	UploadFile(ctx context.Context, id int, filename string, reader io.Reader) (*models.KnowledgeFile, error)

	// UploadFileFromPath 从文件路径上传文件到知识库
	UploadFileFromPath(ctx context.Context, id int, filePath string) (*models.KnowledgeFile, error)

	// BatchUploadFilesFromPath 批量从文件路径上传文件到知识库
	BatchUploadFilesFromPath(ctx context.Context, id int, filePaths []string) (*models.BatchUploadResponse, error)

	// DeleteFile 删除知识库文件
	DeleteFile(ctx context.Context, fileID int) error

	// ToggleFileEnable 切换文件启用状态
	ToggleFileEnable(ctx context.Context, fileID int, enable bool) error
}
