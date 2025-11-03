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
	Version   string `json:"version" example:"1.0.0"`                  // 版本号
	BuildTime string `json:"build_time" example:"2024-01-01 12:00:00"` // 构建时间
	GitCommit string `json:"git_commit" example:"abc1234"`             // Git提交哈希
	GitBranch string `json:"git_branch" example:"main"`                // Git分支
	GitTag    string `json:"git_tag,omitempty" example:"v1.0.0"`       // Git标签
}

// versionInfoResponseWrapper represents a version info response.
//
// swagger:response VersionInfoResponse
type versionInfoResponseWrapper struct {
	// in: body
	Body struct {
		Success bool        `json:"success"`
		Message string      `json:"message,omitempty"`
		Data    VersionInfo `json:"data,omitempty"`
	}
}

// GetVersion 返回应用程序的版本信息。
//
// swagger:route GET /api/v1/version System getVersion
//
// Returns the version information of the application.
//
// Produces:
// - application/json
//
// Responses:
//
//	200: VersionInfoResponse
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
