package knowledge

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	"flowy-sdk/pkg/client"
	"flowy-sdk/pkg/errors"
	"flowy-sdk/pkg/models"
)

// ==================== Request/Response Type Definitions ====================

// KnowledgeBaseListRequest 知识库列表请求
// API: POST /knowledge/list
type KnowledgeBaseListRequest struct{}

// KnowledgeBaseDetailRequest 知识库配置详情请求
// API: POST /knowledge/detail
type KnowledgeBaseDetailRequest struct {
	ID int `json:"id"` // 知识库ID
}

// KnowledgeBaseConfig 知识库配置（公用部分）
type KnowledgeBaseConfig struct {
	Type           int    `json:"type"`           // 知识库类型: 0=file文件型, 1=product产品型
	Name           string `json:"name"`           // 知识库名称
	Desc           string `json:"desc"`           // 知识库描述
	AgentModel     int    `json:"agentModel"`     // 对话模型ID
	ChunkStrategy  string `json:"chunkStrategy"`  // 切片策略
	ChunkSize      int    `json:"chunkSize"`      // 切片大小
	RecallStrategy string `json:"recallStrategy"` // 召回策略
	RecallLimit    int    `json:"recallLimit"`    // 召回大小限制
	EnableLangMix  bool   `json:"enableLangMix"`  // 是否启用多语言混合检索
}

// KnowledgeBaseCreateRequest 新建知识库请求
// API: POST /knowledge/create
type KnowledgeBaseCreateRequest struct {
	KnowledgeBaseConfig     // 嵌入公用配置
	VectorModel         int `json:"vectorModel"` // 向量模型ID
}

// NewDefaultKnowledgeBaseCreateRequest 创建带默认值的知识库创建请求
func NewDefaultKnowledgeBaseCreateRequest(name, desc string) *KnowledgeBaseCreateRequest {
	return &KnowledgeBaseCreateRequest{
		KnowledgeBaseConfig: KnowledgeBaseConfig{
			Type:           0,       // 文件型知识库
			Name:           name,    // 知识库名称
			Desc:           desc,    // 知识库描述
			AgentModel:     2,       // 默认对话模型ID
			ChunkStrategy:  "fixed", // 固定切片策略
			ChunkSize:      512,     // 切片大小512
			RecallStrategy: "slice", // 切片召回策略
			RecallLimit:    512,     // 召回限制512
			EnableLangMix:  false,   // 不启用多语言混合检索
		},
		VectorModel: 1, // 默认向量模型ID
	}
}

// KnowledgeBaseUpdateRequest 更新知识库请求
// API: POST /knowledge/update
type KnowledgeBaseUpdateRequest struct {
	ID                  int `json:"id"` // 知识库ID
	KnowledgeBaseConfig     // 嵌入公用配置
}

// NewDefaultKnowledgeBaseUpdateRequest 创建带默认值的知识库更新请求
func NewDefaultKnowledgeBaseUpdateRequest(id int, name, desc string) *KnowledgeBaseUpdateRequest {
	return &KnowledgeBaseUpdateRequest{
		ID: id, // 知识库ID
		KnowledgeBaseConfig: KnowledgeBaseConfig{
			Type:           0,       // 文件型知识库
			Name:           name,    // 知识库名称
			Desc:           desc,    // 知识库描述
			AgentModel:     2,       // 默认对话模型ID
			ChunkStrategy:  "fixed", // 固定切片策略（与CreateRequest保持一致）
			ChunkSize:      512,     // 切片大小512
			RecallStrategy: "slice", // 切片召回策略（与CreateRequest保持一致）
			RecallLimit:    512,     // 召回限制512（与CreateRequest保持一致）
			EnableLangMix:  false,   // 不启用多语言混合检索
		},
	}
}

// KnowledgeBaseDeleteRequest 删除知识库请求
// API: POST /knowledge/delete
type KnowledgeBaseDeleteRequest struct {
	ID int `json:"id"` // 知识库ID
}

// KnowledgeBaseInfo 知识库信息
type KnowledgeBaseInfo struct {
	ID          int `json:"id"`
	FileCount   int `json:"fileCount"`
	VectorModel int `json:"vectorModel"`
	KnowledgeBaseConfig
}

// FileListRequest filelistrequest
type FileListRequest struct {
	ID   int    `json:"id"`   // knowledge baseid
	Lang string `json:"lang"` // 语言
}

// FileBaseInfo 文件基础信息（公用部分）
type FileBaseInfo struct {
	ID                int         `json:"id"`                // 文件在数据库中的唯一标识符
	PID               int         `json:"pid"`               // 父级文件夹ID，0表示在根目录
	KnowledgeID       int         `json:"knowledgeId"`       // 所属知识库ID
	Name              string      `json:"name"`              // 文件名称
	OSS               string      `json:"oss"`               // OSS存储的文件名（UUID格式）
	FileConvertedOss  string      `json:"fileConvertedOss"`  // 转换后文件的OSS路径，空表示未转换
	FileSize          int         `json:"fileSize"`          // 文件大小（字节），0表示未获取到大小
	Type              string      `json:"type"`              // 文件类型，空字符串表示未设置,只支持*.doc,*.docx,*.pdf,*.md,*.txt,*.json
	Enable            bool        `json:"enable"`            // 文件是否启用，false表示未启用
	Wn                int         `json:"wn"`                // 词数统计，0表示未统计
	Sha1              string      `json:"sha1"`              // 文件SHA1哈希值，空表示未计算
	CreateTime        string      `json:"createTime"`        // 创建时间，null表示未设置
	RecallPrompt      string      `json:"recallPrompt"`      // 召回提示词，空表示未设置
	Status            int         `json:"status"`            // 处理状态: 0=待处理/构建中, 1=完成, 2=失败
	IndexPercent      int         `json:"indexPercent"`      // 索引进度百分比，0表示未开始索引
	UserData          interface{} `json:"userData"`          // 用户自定义数据，null表示无自定义数据
	PlainText         string      `json:"plainText"`         // 纯文本内容，空表示未提取
	MdText            string      `json:"mdText"`            // Markdown格式文本，空表示未转换
	MdTextByParagraph interface{} `json:"mdTextByParagraph"` // 按段落分割的Markdown文本，null表示未处理
	Labels            []string    `json:"labels"`            // 文件标签，null表示未设置标签
	Lang              string      `json:"lang"`              // 文件语言，zh_CN表示中文
	Children          interface{} `json:"children"`          // 子文件列表，null表示无子文件
	ChunkStrategy     string      `json:"chunkStrategy"`     // 分块策略，空表示使用默认策略
	ChunkSize         int         `json:"chunkSize"`         // 分块大小，0表示使用默认大小
	RecallStrategy    string      `json:"recallStrategy"`    // 召回策略，空表示使用默认策略
	RecallLimit       int         `json:"recallLimit"`       // 召回限制，0表示无限制
	ErrorMessage      string      `json:"errorMessage"`      // 错误信息，空字符串表示无错误
}

// FileInfo 文件信息
type FileInfo = FileBaseInfo

// FileUploadRequest 文件上传请求 (multipart/form-data)
// API: POST /knowledge/file/upload
type FileUploadRequest struct {
	KnowledgeID int    // 知识库ID
	PID         int    // 父级文件夹ID
	Lang        string // 语言
}

// FileUploadData 文件上传返回数据
type FileUploadData = FileBaseInfo

// FileDeleteRequest 文件删除请求
// API: POST /knowledge/file/delete
type FileDeleteRequest struct {
	ID int `json:"id"` // 文件ID
}

// FileToggleEnableRequest 切换文件启用/未启用请求
// API: POST /knowledge/file/toggleEnable
type FileToggleEnableRequest struct {
	ID     int  `json:"id"`     // 文件ID
	Enable bool `json:"enable"` // true=启用, false=不启用
}

// FileUpdateRecallPromptRequest 更新文件召回prompt请求
// API: POST /knowledge/file/updateRecallPrompt
type FileUpdateRecallPromptRequest struct {
	ID           int    `json:"id"`           // 文件ID
	RecallPrompt string `json:"recallPrompt"` // 召回提示词
}

// QAVFileSaveRequest 新建/修改问答文件请求 (虚拟文件)
// API: POST /knowledge/file/saveQAVFile
type QAVFileSaveRequest struct {
	ID          int      `json:"id"`          // 文件ID: 0=新建, >0=更新
	KnowledgeID int      `json:"knowledgeId"` // 知识库ID
	PID         int      `json:"pid"`         // 父级文件夹ID
	Name        string   `json:"name"`        // 文件名称
	QAList      []QAItem `json:"qaList"`      // 问答集合
	Lang        string   `json:"lang"`        // 语言
}

// QAItem 问答项
type QAItem struct {
	ID        int      `json:"id"`        // 问答块的ID: 0=新建, >0=更新
	FileID    int      `json:"fileId"`    // 问答文档的ID
	Question  string   `json:"question"`  // 问题
	Labels    []string `json:"labels"`    // 标签列表
	Answer    string   `json:"answer"`    // 答案
	ReferLink string   `json:"referLink"` // 参考链接
}

// QAVFileDetailRequest 问答文件详情请求
// API: POST /knowledge/file/qaVFileDetail
type QAVFileDetailRequest struct {
	ID int `json:"id"` // 文件ID
}

// QAVFileDetailData 问答文件详情数据
type QAVFileDetailData struct {
	ID          int            `json:"id"`          // 文件ID
	Name        string         `json:"name"`        // 文件名称
	QAList      []QADetailItem `json:"qaList"`      // 问答对列表
	KnowledgeID int            `json:"knowledgeId"` // 所属知识库ID
	PID         int            `json:"pid"`         // 父级ID
	Lang        string         `json:"lang"`        // 语言: zh_CN=中文, en_US=英文
}

// QADetailItem 问答详情项
type QADetailItem struct {
	ID        int      `json:"id"`        // 切片ID
	Title     string   `json:"title"`     // 问答对标题
	Answer    string   `json:"answer"`    // 答案
	Question  string   `json:"question"`  // 问题
	ReferLink string   `json:"referLink"` // 参考链接
	FileID    int      `json:"fileId"`    // 所属文件ID
	Labels    []string `json:"labels"`    // 标签列表
	VecID     int      `json:"vecId"`     // 向量ID: 0=未向量化, >0=已向量化
}

// FileSliceDetailRequest 查询单个文件切片请求
// API: POST /knowledge/file/sliceDetail
type FileSliceDetailRequest struct {
	ID int `json:"id"` // 文件类型文件的ID
}

// FileSliceDetailData 文件切片详情数据
type FileSliceDetailData struct {
	ID          int         `json:"id"`          // 文件ID (唯一标识)
	VecID       int         `json:"vecId"`       // 向量ID: 0=未向量化, >0=向量数据库中的ID
	Content     string      `json:"content"`     // 切片内容
	Title       string      `json:"title"`       // 切片标题
	IsLeaf      bool        `json:"isLeaf"`      // 是否为叶子节点
	PID         int         `json:"pid"`         // 父级切片ID
	FileID      int         `json:"fileId"`      // 所属文件ID
	FullContent string      `json:"fullContent"` // 完整内容
	Children    interface{} `json:"children"`    // 子切片列表 (null表示无)
}

// FileSliceListRequest 产品文件详情请求
// API: POST /knowledge/file/sliceList
type FileSliceListRequest struct {
	ID int `json:"id"` // 文件ID
}

// FileUpdateSliceRequest 更新文件切片请求
// API: POST /knowledge/file/updateSlice
type FileUpdateSliceRequest struct {
	ID    int      `json:"id"`    // 文件ID
	Texts []string `json:"texts"` // 切片文本列表
}

// FileToggleChunkStrategyRequest 修改文件切片策略请求
// API: POST /knowledge/file/toggleChunkStrategy
type FileToggleChunkStrategyRequest struct {
	ID       int    `json:"id"`       // 文件ID
	Strategy string `json:"strategy"` // 切片策略
}

// FileToggleRecallStrategyRequest 修改文件召回策略请求
// API: POST /knowledge/file/toggleRecallStrategy
type FileToggleRecallStrategyRequest struct {
	ID       int    `json:"id"`       // 文件ID
	Strategy string `json:"strategy"` // 召回策略
}

// FileUpdateRecallLimitRequest 修改文件召回限制请求
// API: POST /knowledge/file/updateRecallLimit
type FileUpdateRecallLimitRequest struct {
	ID    int `json:"id"`    // 文件ID
	Limit int `json:"limit"` // 召回限制
}

// FileModifyRequest 文件设置修改请求
// API: POST /knowledge/file/modify
type FileModifyRequest struct {
	ID             int      `json:"id"`             // 文件ID
	Name           string   `json:"name"`           // 文件名称
	ChunkStrategy  string   `json:"chunkStrategy"`  // 切片策略
	ChunkSize      int      `json:"chunkSize"`      // 切片大小
	RecallStrategy string   `json:"recallStrategy"` // 召回策略
	RecallLimit    int      `json:"recallLimit"`    // 召回限制
	RecallPrompt   string   `json:"recallPrompt"`   // 召回提示词
	Labels         []string `json:"labels"`         // 标签列表
}

// FileMdContentRequest 获取文件markdown格式内容请求
// API: POST /knowledge/file/mdContent
type FileMdContentRequest struct {
	ID int `json:"id"` // 文件ID
}

// FileLogsRequest 文件日志请求
// API: POST /knowledge/file/logs
type FileLogsRequest struct {
	ID int `json:"id"` // 文件ID
}

// FileLogItem 文件日志项
type FileLogItem struct {
	Message string `json:"message"` // 日志消息
	Time    string `json:"time"`    // 时间
}

// QASaveRequest 新建/修改问答请求
// API: POST /knowledge/qa/save
type QASaveRequest struct {
	FileID    int      `json:"fileId"`    // 虚拟文件ID
	ID        int      `json:"id"`        // 问答模块的ID: 0=新建, >0=修改
	Question  string   `json:"question"`  // 问题
	Labels    []string `json:"labels"`    // 标签
	Answer    string   `json:"answer"`    // 答案
	ReferLink string   `json:"referLink"` // 引用链接
}

// QADeleteRequest 删除问答块请求
// API: POST /knowledge/qa/delete
type QADeleteRequest struct {
	ID int `json:"id"` // 问答块ID
}

// QAListByPageRequest 问答列表分页查询请求
// API: POST /knowledge/qa/listByPage
type QAListByPageRequest struct {
	FileID int `json:"fileId"` // 虚拟文件ID
	Page   int `json:"page"`   // 页码
	Size   int `json:"size"`   // 页尺寸
}

// QAListByPageData 问答列表分页数据
type QAListByPageData struct {
	Total   int            `json:"total"`   // 总数
	Records []QARecordItem `json:"records"` // 记录列表
}

// QARecordItem 问答记录项
type QARecordItem struct {
	ID        int      `json:"id"`        // 问答ID
	Title     string   `json:"title"`     // 标题
	Answer    string   `json:"answer"`    // 答案
	Question  string   `json:"question"`  // 问题
	ReferLink string   `json:"referLink"` // 参考链接
	FileID    int      `json:"fileId"`    // 文件ID
	Labels    []string `json:"labels"`    // 标签列表
	VecID     int      `json:"vecId"`     // 向量ID
}

// ProductListRequest 产品列表请求
// API: POST /knowledge/product/list
type ProductListRequest struct {
	KnowledgeID int `json:"knowledgeId"` // 知识库ID
}

// ProductListLastRequest 产品最新列表请求
// API: POST /knowledge/product/listLast
type ProductListLastRequest struct {
	Limit int `json:"limit"` // 限制数量
}

// ProductUploadFilesRequest 上传产品文件请求 (multipart/form-data)
// API: POST /knowledge/product/uploadFiles
// 注意：这个是复杂的请求，包含文件和其他字段

// ProductFileItem 产品文件项
type ProductFileItem struct {
	Extension string `json:"extension"` // 文件扩展名
	File      string `json:"file"`      // 文件路径
	Origin    string `json:"origin"`    // 原始文件名
	Desc      string `json:"desc"`      // 描述
	Update    string `json:"update"`    // 更新时间
}

// ProductFiles 产品文件集合
type ProductFiles struct {
	Images    []ProductFileItem `json:"images"`    // 图片列表
	Documents []ProductFileItem `json:"documents"` // 文档列表
}

// ProductDetailRequest 产品详情请求
// API: POST /knowledge/product/detail
type ProductDetailRequest struct {
	ID          int `json:"id"`          // 产品ID
	KnowledgeID int `json:"knowledgeId"` // 知识库ID
}

// ProductDetailData 产品详情数据
type ProductDetailData struct {
	ID            int          `json:"id"`            // 产品ID
	Name          string       `json:"name"`          // 产品名称
	ReferLink     string       `json:"referLink"`     // 参考链接
	KnowledgeID   int          `json:"knowledgeId"`   // 知识库ID
	Labels        []string     `json:"labels"`        // 标签列表
	VecID         int          `json:"vecId"`         // 向量ID
	OriginContent string       `json:"originContent"` // 原始内容
	Properties    []string     `json:"properties"`    // 属性列表
	Files         ProductFiles `json:"files"`         // 文件集合
}

// ProductSaveRequest 新建/修改产品请求
// API: POST /knowledge/product/save
type ProductSaveRequest struct {
	KnowledgeID   int          `json:"knowledgeId"`   // 知识库ID
	ID            int          `json:"id"`            // 产品ID: 0=新建, >0=修改
	Name          string       `json:"name"`          // 产品名称
	OriginContent string       `json:"originContent"` // 原始内容
	Labels        []string     `json:"labels"`        // 标签列表
	ReferLink     string       `json:"referLink"`     // 参考链接
	Properties    []string     `json:"properties"`    // 属性列表
	Files         ProductFiles `json:"files"`         // 文件集合
}

// ProductDeleteRequest 删除产品请求
// API: POST /knowledge/product/delete
type ProductDeleteRequest struct {
	KnowledgeID int `json:"knowledgeId"` // 知识库ID
	ID          int `json:"id"`          // 产品ID
}

// ProductListByPageRequest 产品列表分页查询请求
// API: POST /knowledge/product/listByPage
type ProductListByPageRequest struct {
	KnowledgeID int `json:"knowledgeId"` // 知识库ID
	Page        int `json:"page"`        // 页码
	Size        int `json:"size"`        // 页尺寸
}

// ProductListByPageData 产品列表分页数据
type ProductListByPageData struct {
	Total   int                 `json:"total"`   // 总数
	Records []ProductRecordItem `json:"records"` // 记录列表
}

// ProductRecordItem 产品记录项
type ProductRecordItem struct {
	ID            int          `json:"id"`            // 产品ID
	Name          string       `json:"name"`          // 产品名称
	ReferLink     string       `json:"referLink"`     // 参考链接
	KnowledgeID   int          `json:"knowledgeId"`   // 知识库ID
	Labels        []string     `json:"labels"`        // 标签列表
	VecID         int          `json:"vecId"`         // 向量ID
	OriginContent string       `json:"originContent"` // 原始内容
	Properties    []string     `json:"properties"`    // 属性列表
	Files         ProductFiles `json:"files"`         // 文件集合
}

// ProductSchemaDetailRequest 产品文件详情请求
// API: POST /knowledge/product/schema/detail
type ProductSchemaDetailRequest struct {
	KnowledgeID int `json:"knowledgeId"` // 知识库ID
}

// ProductSchemaDetailData 产品文件详情数据
type ProductSchemaDetailData struct {
	Name        string            `json:"name"`        // 名称
	KnowledgeID int               `json:"knowledgeId"` // 知识库ID
	Properties  []ProductProperty `json:"properties"`  // 属性列表
	PID         int               `json:"pid"`         // 父级ID
	Lang        string            `json:"lang"`        // 语言
	Products    interface{}       `json:"products"`    // 产品列表
}

// ProductProperty 产品属性
type ProductProperty struct {
	ID          string            `json:"id"`          // 属性ID (前端生成的唯一ID)
	Name        string            `json:"name"`        // 属性名称
	DataType    string            `json:"dataType"`    // 数据类型
	Values      interface{}       `json:"values"`      // 值
	Enums       interface{}       `json:"enums"`       // 枚举值
	Definition  string            `json:"definition"`  // 定义
	NearSynonym interface{}       `json:"nearSynonym"` // 近义词
	Children    []ProductProperty `json:"children"`    // 子属性
	Multiple    bool              `json:"multiple"`    // 是否多选
	Group       bool              `json:"group"`       // 是否分组
}

// ProductSchemaUpdateRequest 更新产品文件详情请求
// API: POST /knowledge/product/schema/save
type ProductSchemaUpdateRequest struct {
	KnowledgeID int               `json:"knowledgeId"` // 知识库ID
	PID         int               `json:"pid"`         // 父级ID
	Name        string            `json:"name"`        // 名称
	Properties  []ProductProperty `json:"properties"`  // 属性列表
	Lang        string            `json:"lang"`        // 语言
}

// ProductSchemaAIGenerateRequest AI生成索引请求
// API: POST /knowledge/product/schema/AIGenerate
type ProductSchemaAIGenerateRequest struct {
	KnowledgeID int    `json:"knowledgeId"` // 知识库ID
	Content     string `json:"content"`     // 内容
}

// ProductSchemaSaveRequest 索引保存请求
// API: POST /knowledge/product/schema/save
type ProductSchemaSaveRequest struct {
	KnowledgeID int               `json:"knowledgeId"` // 知识库ID
	PID         int               `json:"pid"`         // 父级ID
	Name        string            `json:"name"`        // 名称
	Properties  []ProductProperty `json:"properties"`  // 属性列表 (ID是前端生成的唯一ID)
	Lang        string            `json:"lang"`        // 语言
}

// Service Knowledge service接口
type Service interface {
	// ==================== knowledge base基础管理 ====================
	// knowledge base列表
	ListKnowledgeBases(ctx context.Context) ([]KnowledgeBaseInfo, error)
	// knowledge base配置详情
	GetKnowledgeBaseDetail(ctx context.Context, id int) (*KnowledgeBaseInfo, error)
	// 新建知识库
	CreateKnowledgeBase(ctx context.Context, req *KnowledgeBaseCreateRequest) (int, error)
	// 更新知识库
	UpdateKnowledgeBase(ctx context.Context, req *KnowledgeBaseUpdateRequest) (int, error)
	// 删除知识库
	DeleteKnowledgeBase(ctx context.Context, id int) error

	// ==================== 文件操作 ====================
	// 文件列表(不分页)
	// API: POST /knowledge/file/list
	ListFiles(ctx context.Context, knowledgeID int, lang string) ([]FileInfo, error)

	// 文件上传
	// API: POST /knowledge/file/upload (multipart/form-data)
	UploadFile(ctx context.Context, file io.Reader, filename string, knowledgeID, pid int, lang string) (*FileUploadData, error)

	// 文件删除
	// API: POST /knowledge/file/delete
	DeleteFile(ctx context.Context, id int) error

	// 切换文件启用/未启用
	// API: POST /knowledge/file/toggleEnable
	ToggleFileEnable(ctx context.Context, id int, enable bool) error

	// 更新文件召回prompt
	// API: POST /knowledge/file/updateRecallPrompt
	UpdateFileRecallPrompt(ctx context.Context, id int, recallPrompt string) error

	// 下载原文件
	// API: GET /api/v1/knowledge/file/download
	DownloadOriginalFile(ctx context.Context, id int) (io.ReadCloser, error)

	// 获取照片或者文件
	// API: GET /api/v1/knowledge/filePreview/{ossFileName}
	GetFilePreview(ctx context.Context, ossFileName string) (io.ReadCloser, error)

	// 获取文件
	// API: GET /knowledge/filePreview/{ossFileName}
	GetFile(ctx context.Context, ossFileName string) (io.ReadCloser, error)

	// 新建问答文件(虚拟文件)
	// API: POST /knowledge/file/saveQAVFile
	SaveQAVFile(ctx context.Context, req *QAVFileSaveRequest) (int, error)

	// 问答文件详情
	// API: POST /knowledge/file/qaVFileDetail
	GetQAVFileDetail(ctx context.Context, id int) (*QAVFileDetailData, error)

	// 查询单个文件切片
	// API: POST /knowledge/file/sliceDetail
	GetFileSliceDetail(ctx context.Context, id int) (*FileSliceDetailData, error)

	// 产品文件详情（切片列表）
	// API: POST /knowledge/file/sliceList
	GetFileSliceList(ctx context.Context, id int) ([]FileSliceDetailData, error)

	// 更新文件切片
	// API: POST /knowledge/file/updateSlice
	UpdateFileSlice(ctx context.Context, id int, texts []string) error

	// 修改文件切片策略
	// API: POST /knowledge/file/toggleChunkStrategy
	ToggleFileChunkStrategy(ctx context.Context, id int, strategy string) error

	// 修改文件召回策略
	// API: POST /knowledge/file/toggleRecallStrategy
	ToggleFileRecallStrategy(ctx context.Context, id int, strategy string) error

	// 修改文件召回限制
	// API: POST /knowledge/file/updateRecallLimit
	UpdateFileRecallLimit(ctx context.Context, id int, limit int) error

	// 文件设置修改
	// API: POST /knowledge/file/modify
	ModifyFile(ctx context.Context, req *FileModifyRequest) error

	// 获取文件markdown格式内容
	// API: POST /knowledge/file/mdContent
	GetFileMdContent(ctx context.Context, id int) (string, error)

	// 文件日志
	// API: POST /knowledge/file/logs
	GetFileLogs(ctx context.Context, id int) ([]FileLogItem, error)

	// 导出知识库
	// API: GET /knowledge/exportFiles
	ExportKnowledgeBase(ctx context.Context, id int) (io.ReadCloser, error)

	// ==================== 问答管理 ====================
	// 新建/修改问答
	// API: POST /knowledge/qa/save
	SaveQA(ctx context.Context, req *QASaveRequest) (int, error)

	// 删除问答块
	// API: POST /knowledge/qa/delete
	DeleteQA(ctx context.Context, id int) error

	// 问答列表分页查询
	// API: POST /knowledge/qa/listByPage
	ListQAByPage(ctx context.Context, fileID, page, size int) (*QAListByPageData, error)

	// ==================== 产品管理 ====================
	// 产品列表(不分页)
	// API: POST /knowledge/product/list
	ListProducts(ctx context.Context, knowledgeID int) ([]ProductRecordItem, error)

	// 产品最新列表
	// API: POST /knowledge/product/listlast
	ListLastProducts(ctx context.Context, limit int) ([]ProductFileItem, error)

	// 上传产品文件
	// API: POST /knowledge/product/uploadFiles
	UploadProductFiles(ctx context.Context, files []io.Reader, filenames []string, req map[string]interface{}) ([]ProductFileItem, error)

	// 产品详情
	// API: POST /knowledge/product/detail
	GetProductDetail(ctx context.Context, id, knowledgeID int) (*ProductDetailData, error)

	// 新建/修改产品
	// API: POST /knowledge/product/save
	SaveProduct(ctx context.Context, req *ProductSaveRequest) error

	// 删除产品
	// API: POST /knowledge/product/delete
	DeleteProduct(ctx context.Context, knowledgeID, id int) error

	// 产品列表分页查询
	// API: POST /knowledge/product/listByPage
	ListProductsByPage(ctx context.Context, knowledgeID, page, size int) (*ProductListByPageData, error)

	// 产品文件详情(Schema详情)
	// API: POST /knowledge/product/schema/detail
	GetProductSchemaDetail(ctx context.Context, knowledgeID int) (*ProductSchemaDetailData, error)

	// 更新产品文件详情
	// API: POST /knowledge/product/schema/save
	UpdateProductSchema(ctx context.Context, req *ProductSchemaUpdateRequest) error

	// AI生成索引
	// API: POST /knowledge/product/schema/AIGenerate
	AIGenerateProductSchema(ctx context.Context, knowledgeID int, content string) ([]ProductProperty, error)

	// 索引保存
	// API: POST /knowledge/product/schema/save
	SaveProductSchema(ctx context.Context, req *ProductSchemaSaveRequest) error
}

// ServiceImpl 知识库服务实现
type ServiceImpl struct {
	client client.HTTPClient
}

// NewService 创建知识库服务
func NewService(client client.HTTPClient) Service {
	return &ServiceImpl{
		client: client,
	}
}

// ==================== knowledge base基础管理实现 ====================

// ListKnowledgeBases knowledge base列表
func (s *ServiceImpl) ListKnowledgeBases(ctx context.Context) ([]KnowledgeBaseInfo, error) {
	resp, err := s.client.Post(ctx, "/knowledge/list", map[string]interface{}{})
	if err != nil {
		return nil, err
	}

	var result []KnowledgeBaseInfo
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// GetKnowledgeBaseDetail knowledge base配置详情
func (s *ServiceImpl) GetKnowledgeBaseDetail(ctx context.Context, id int) (*KnowledgeBaseInfo, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := KnowledgeBaseDetailRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/detail", req)
	if err != nil {
		return nil, err
	}

	var result KnowledgeBaseInfo
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// CreateKnowledgeBase 新建知识库
func (s *ServiceImpl) CreateKnowledgeBase(ctx context.Context, req *KnowledgeBaseCreateRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "knowledge base name is required")
	}

	resp, err := s.client.Post(ctx, "/knowledge/create", req)
	if err != nil {
		return 0, err
	}

	var result int
	if err := s.parseResponseData(resp, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// UpdateKnowledgeBase 更新知识库
func (s *ServiceImpl) UpdateKnowledgeBase(ctx context.Context, req *KnowledgeBaseUpdateRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.ID <= 0 {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	resp, err := s.client.Post(ctx, "/knowledge/update", req)
	if err != nil {
		return 0, err
	}

	var result int
	if err := s.parseResponseData(resp, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// DeleteKnowledgeBase 删除知识库
func (s *ServiceImpl) DeleteKnowledgeBase(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := KnowledgeBaseDeleteRequest{ID: id}
	_, err := s.client.Post(ctx, "/knowledge/delete", req)
	return err
}

// ==================== 文件操作实现 ====================

// ListFiles 文件列表(不分页)
func (s *ServiceImpl) ListFiles(ctx context.Context, knowledgeID int, lang string) ([]FileInfo, error) {
	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := FileListRequest{
		ID:   knowledgeID,
		Lang: lang,
	}

	resp, err := s.client.Post(ctx, "/knowledge/file/list", req)
	if err != nil {
		return nil, err
	}

	var result []FileInfo
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UploadFile 文件上传
func (s *ServiceImpl) UploadFile(ctx context.Context, file io.Reader, filename string, knowledgeID, pid int, lang string) (*FileUploadData, error) {
	if file == nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file is required")
	}

	if filename == "" {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "filename is required")
	}

	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	params := map[string]string{
		"knowledgeId": fmt.Sprintf("%d", knowledgeID),
		"pid":         fmt.Sprintf("%d", pid),
		"lang":        lang,
	}

	resp, err := s.client.Upload(ctx, "/knowledge/file/upload", "file", filename, file, params)
	if err != nil {
		return nil, err
	}

	var result FileUploadData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// DeleteFile 文件删除
func (s *ServiceImpl) DeleteFile(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileDeleteRequest{ID: id}
	_, err := s.client.Post(ctx, "/knowledge/file/delete", req)
	return err
}

// ToggleFileEnable 切换文件启用/未启用
func (s *ServiceImpl) ToggleFileEnable(ctx context.Context, id int, enable bool) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileToggleEnableRequest{
		ID:     id,
		Enable: enable,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/toggleEnable", req)
	return err
}

// UpdateFileRecallPrompt 更新文件召回prompt
func (s *ServiceImpl) UpdateFileRecallPrompt(ctx context.Context, id int, recallPrompt string) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileUpdateRecallPromptRequest{
		ID:           id,
		RecallPrompt: recallPrompt,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/updateRecallPrompt", req)
	return err
}

// DownloadOriginalFile 下载原文件
func (s *ServiceImpl) DownloadOriginalFile(ctx context.Context, id int) (io.ReadCloser, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	// 构建带参数的URL
	path := fmt.Sprintf("/api/v1/knowledge/file/download?id=%d", id)
	return s.client.Download(ctx, path)
}

// GetFilePreview 获取照片或者文件
func (s *ServiceImpl) GetFilePreview(ctx context.Context, ossFileName string) (io.ReadCloser, error) {
	if ossFileName == "" {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "OSS file name is required")
	}

	return s.client.Download(ctx, fmt.Sprintf("/api/v1/knowledge/filePreview/%s", ossFileName))
}

// GetFile 获取文件
func (s *ServiceImpl) GetFile(ctx context.Context, ossFileName string) (io.ReadCloser, error) {
	if ossFileName == "" {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "OSS file name is required")
	}

	return s.client.Download(ctx, fmt.Sprintf("/knowledge/filePreview/%s", ossFileName))
}

// SaveQAVFile 新建问答文件(虚拟文件)
func (s *ServiceImpl) SaveQAVFile(ctx context.Context, req *QAVFileSaveRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Name == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "file name is required")
	}

	resp, err := s.client.Post(ctx, "/knowledge/file/saveQAVFile", req)
	if err != nil {
		return 0, err
	}

	var result int
	if err := s.parseResponseData(resp, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// GetQAVFileDetail 问答文件详情
func (s *ServiceImpl) GetQAVFileDetail(ctx context.Context, id int) (*QAVFileDetailData, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := QAVFileDetailRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/file/qaVFileDetail", req)
	if err != nil {
		return nil, err
	}

	var result QAVFileDetailData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFileSliceDetail 查询单个文件切片
func (s *ServiceImpl) GetFileSliceDetail(ctx context.Context, id int) (*FileSliceDetailData, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileSliceDetailRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/file/sliceDetail", req)
	if err != nil {
		return nil, err
	}

	var result FileSliceDetailData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetFileSliceList 产品文件详情（切片列表）
func (s *ServiceImpl) GetFileSliceList(ctx context.Context, id int) ([]FileSliceDetailData, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileSliceListRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/file/sliceList", req)
	if err != nil {
		return nil, err
	}

	var result []FileSliceDetailData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UpdateFileSlice 更新文件切片
func (s *ServiceImpl) UpdateFileSlice(ctx context.Context, id int, texts []string) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileUpdateSliceRequest{
		ID:    id,
		Texts: texts,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/updateSlice", req)
	return err
}

// ToggleFileChunkStrategy 修改文件切片策略
func (s *ServiceImpl) ToggleFileChunkStrategy(ctx context.Context, id int, strategy string) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileToggleChunkStrategyRequest{
		ID:       id,
		Strategy: strategy,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/toggleChunkStrategy", req)
	return err
}

// ToggleFileRecallStrategy 修改文件召回策略
func (s *ServiceImpl) ToggleFileRecallStrategy(ctx context.Context, id int, strategy string) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileToggleRecallStrategyRequest{
		ID:       id,
		Strategy: strategy,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/toggleRecallStrategy", req)
	return err
}

// UpdateFileRecallLimit 修改文件召回限制
func (s *ServiceImpl) UpdateFileRecallLimit(ctx context.Context, id int, limit int) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileUpdateRecallLimitRequest{
		ID:    id,
		Limit: limit,
	}

	_, err := s.client.Post(ctx, "/knowledge/file/updateRecallLimit", req)
	return err
}

// ModifyFile 文件设置修改
func (s *ServiceImpl) ModifyFile(ctx context.Context, req *FileModifyRequest) error {
	if req == nil {
		return errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.ID <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	_, err := s.client.Post(ctx, "/knowledge/file/modify", req)
	return err
}

// GetFileMdContent 获取文件markdown格式内容
func (s *ServiceImpl) GetFileMdContent(ctx context.Context, id int) (string, error) {
	if id <= 0 {
		return "", errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileMdContentRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/file/mdContent", req)
	if err != nil {
		return "", err
	}

	var result string
	if err := s.parseResponseData(resp, &result); err != nil {
		return "", err
	}

	return result, nil
}

// GetFileLogs 文件日志
func (s *ServiceImpl) GetFileLogs(ctx context.Context, id int) ([]FileLogItem, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := FileLogsRequest{ID: id}
	resp, err := s.client.Post(ctx, "/knowledge/file/logs", req)
	if err != nil {
		return nil, err
	}

	var result []FileLogItem
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ExportKnowledgeBase 导出知识库
func (s *ServiceImpl) ExportKnowledgeBase(ctx context.Context, id int) (io.ReadCloser, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	// 构建带参数的URL
	path := fmt.Sprintf("/knowledge/exportFiles?id=%d", id)
	return s.client.Download(ctx, path)
}

// ==================== 问答管理实现 ====================

// SaveQA 新建/修改问答
func (s *ServiceImpl) SaveQA(ctx context.Context, req *QASaveRequest) (int, error) {
	if req == nil {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Question == "" {
		return 0, errors.New(errors.ErrCodeInvalidRequest, "question is required")
	}

	resp, err := s.client.Post(ctx, "/knowledge/qa/save", req)
	if err != nil {
		return 0, err
	}

	var result int
	if err := s.parseResponseData(resp, &result); err != nil {
		return 0, err
	}

	return result, nil
}

// DeleteQA 删除问答块
func (s *ServiceImpl) DeleteQA(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "QA ID is required")
	}

	req := QADeleteRequest{ID: id}
	_, err := s.client.Post(ctx, "/knowledge/qa/delete", req)
	return err
}

// ListQAByPage 问答列表分页查询
func (s *ServiceImpl) ListQAByPage(ctx context.Context, fileID, page, size int) (*QAListByPageData, error) {
	if fileID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "file ID is required")
	}

	req := QAListByPageRequest{
		FileID: fileID,
		Page:   page,
		Size:   size,
	}

	resp, err := s.client.Post(ctx, "/knowledge/qa/listByPage", req)
	if err != nil {
		return nil, err
	}

	var result QAListByPageData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// ==================== 产品管理实现 ====================

// ListProducts 产品列表(不分页)
func (s *ServiceImpl) ListProducts(ctx context.Context, knowledgeID int) ([]ProductRecordItem, error) {
	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := ProductListRequest{KnowledgeID: knowledgeID}
	resp, err := s.client.Post(ctx, "/knowledge/product/list", req)
	if err != nil {
		return nil, err
	}

	var result []ProductRecordItem
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// ListLastProducts 产品最新列表
func (s *ServiceImpl) ListLastProducts(ctx context.Context, limit int) ([]ProductFileItem, error) {
	if limit <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "limit is required")
	}

	req := ProductListLastRequest{Limit: limit}
	resp, err := s.client.Post(ctx, "/knowledge/product/listlast", req)
	if err != nil {
		return nil, err
	}

	var result []ProductFileItem
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// UploadProductFiles 上传产品文件
func (s *ServiceImpl) UploadProductFiles(ctx context.Context, files []io.Reader, filenames []string, req map[string]interface{}) ([]ProductFileItem, error) {
	if len(files) == 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "files are required")
	}

	if len(files) != len(filenames) {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "files and filenames count mismatch")
	}

	// 注意：这个API需要multipart/form-data上传多个文件
	// 由于当前client不支持多文件上传，这里返回未实现错误
	// 实际使用时需要扩展client支持UploadMultiple方法
	return nil, errors.New(errors.ErrCodeInternalError, "multiple file upload not implemented yet, please extend HTTPClient interface with UploadMultiple method")
}

// GetProductDetail 产品详情
func (s *ServiceImpl) GetProductDetail(ctx context.Context, id, knowledgeID int) (*ProductDetailData, error) {
	if id <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "product ID is required")
	}

	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := ProductDetailRequest{
		ID:          id,
		KnowledgeID: knowledgeID,
	}

	resp, err := s.client.Post(ctx, "/knowledge/product/detail", req)
	if err != nil {
		return nil, err
	}

	var result ProductDetailData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// SaveProduct 新建/修改产品
func (s *ServiceImpl) SaveProduct(ctx context.Context, req *ProductSaveRequest) error {
	if req == nil {
		return errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.Name == "" {
		return errors.New(errors.ErrCodeInvalidRequest, "product name is required")
	}

	_, err := s.client.Post(ctx, "/knowledge/product/save", req)
	return err
}

// DeleteProduct 删除产品
func (s *ServiceImpl) DeleteProduct(ctx context.Context, knowledgeID, id int) error {
	if knowledgeID <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	if id <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "product ID is required")
	}

	req := ProductDeleteRequest{
		KnowledgeID: knowledgeID,
		ID:          id,
	}

	_, err := s.client.Post(ctx, "/knowledge/product/delete", req)
	return err
}

// ListProductsByPage 产品列表分页查询
func (s *ServiceImpl) ListProductsByPage(ctx context.Context, knowledgeID, page, size int) (*ProductListByPageData, error) {
	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := ProductListByPageRequest{
		KnowledgeID: knowledgeID,
		Page:        page,
		Size:        size,
	}

	resp, err := s.client.Post(ctx, "/knowledge/product/listByPage", req)
	if err != nil {
		return nil, err
	}

	var result ProductListByPageData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// GetProductSchemaDetail 产品文件详情(Schema详情)
func (s *ServiceImpl) GetProductSchemaDetail(ctx context.Context, knowledgeID int) (*ProductSchemaDetailData, error) {
	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	req := ProductSchemaDetailRequest{KnowledgeID: knowledgeID}
	resp, err := s.client.Post(ctx, "/knowledge/product/schema/detail", req)
	if err != nil {
		return nil, err
	}

	var result ProductSchemaDetailData
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return &result, nil
}

// UpdateProductSchema 更新产品文件详情
func (s *ServiceImpl) UpdateProductSchema(ctx context.Context, req *ProductSchemaUpdateRequest) error {
	if req == nil {
		return errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.KnowledgeID <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	_, err := s.client.Post(ctx, "/knowledge/product/schema/save", req)
	return err
}

// AIGenerateProductSchema AI生成索引
func (s *ServiceImpl) AIGenerateProductSchema(ctx context.Context, knowledgeID int, content string) ([]ProductProperty, error) {
	if knowledgeID <= 0 {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	if content == "" {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "content is required")
	}

	req := ProductSchemaAIGenerateRequest{
		KnowledgeID: knowledgeID,
		Content:     content,
	}

	resp, err := s.client.Post(ctx, "/knowledge/product/schema/AIGenerate", req)
	if err != nil {
		return nil, err
	}

	var result []ProductProperty
	if err := s.parseResponseData(resp, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// SaveProductSchema 索引保存
func (s *ServiceImpl) SaveProductSchema(ctx context.Context, req *ProductSchemaSaveRequest) error {
	if req == nil {
		return errors.New(errors.ErrCodeInvalidRequest, "request cannot be nil")
	}

	if req.KnowledgeID <= 0 {
		return errors.New(errors.ErrCodeInvalidRequest, "knowledge base ID is required")
	}

	_, err := s.client.Post(ctx, "/knowledge/product/schema/save", req)
	return err
}

// parseResponseData 解析响应数据
func (s *ServiceImpl) parseResponseData(resp *models.BaseResponse, target interface{}) error {
	if resp == nil {
		return errors.New(errors.ErrCodeInternalError, "response is nil")
	}

	if !resp.Success {
		return errors.New(errors.ErrCodeInternalError, resp.Message)
	}

	if resp.Data == nil {
		return nil // 某些API返回null data
	}

	// 将数据转换为JSON然后反序列化到目标结构
	data, err := json.Marshal(resp.Data)
	if err != nil {
		return errors.New(errors.ErrCodeInternalError, "failed to marshal response data").WithDetails(err.Error())
	}

	if err := json.Unmarshal(data, target); err != nil {
		return errors.New(errors.ErrCodeInternalError, "failed to unmarshal response data").WithDetails(err.Error())
	}

	return nil
}
