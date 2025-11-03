package handlers

import (
	"chat-backend/utils"

	"github.com/gin-gonic/gin"
)

// VersionHandler 处理版本信息相关的HTTP请求。
type VersionHandler struct {
	version   string
	buildTime string
	gitCommit string
	gitBranch string
	gitTag    string
}

// NewVersionHandler 创建并返回一个新的版本处理器实例。
func NewVersionHandler(version, buildTime, gitCommit, gitBranch, gitTag string) *VersionHandler {
	return &VersionHandler{
		version:   version,
		buildTime: buildTime,
		gitCommit: gitCommit,
		gitBranch: gitBranch,
		gitTag:    gitTag,
	}
}

// VersionInfo 版本信息响应结构
// swagger:model
type VersionInfo struct {
	// 版本号
	// required: true
	Version string `json:"version" example:"1.0.0"`
	// 构建时间
	// required: true
	BuildTime string `json:"build_time" example:"2024-01-01 12:00:00"`
	// Git提交哈希
	// required: true
	GitCommit string `json:"git_commit" example:"abc1234"`
	// Git分支
	// required: true
	GitBranch string `json:"git_branch" example:"main"`
	// Git标签
	// required: true
	GitTag string `json:"git_tag" example:"v1.0.0"`
}

// VersionInfoSuccessResponse 版本信息成功响应
// swagger:response VersionInfoSuccessResponse
type VersionInfoSuccessResponse struct {
	// 版本信息响应
	// in: body
	Body struct {
		// 请求是否成功
		// required: true
		Success bool `json:"success"`
		// 响应消息
		// required: true
		Message string `json:"message"`
		// 版本信息数据
		// required: true
		Data VersionInfo `json:"data"`
		// 时间戳
		// required: true
		Timestamp string `json:"timestamp"`
	}
}

// GetVersion 返回应用程序的版本信息。
//
// swagger:route GET /version System getVersion
//
// 获取版本信息
//
// 返回应用程序的详细版本信息，包括版本号、构建时间和Git信息
//
// Produces:
// - application/json
//
// Responses:
//
//	200: VersionInfoSuccessResponse
func (h *VersionHandler) GetVersion(c *gin.Context) {
	versionInfo := VersionInfo{
		Version:   h.version,
		BuildTime: h.buildTime,
		GitCommit: h.gitCommit,
		GitBranch: h.gitBranch,
		GitTag:    h.gitTag,
	}

	utils.RespondWithSuccess(c, versionInfo, "获取版本信息成功")
}
