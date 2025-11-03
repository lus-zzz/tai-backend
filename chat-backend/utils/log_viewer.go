package utils

import (
	"bufio"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

// LogViewerHandler 日志查看器HTTP处理器
type LogViewerHandler struct {
	logDir string
}

// NewLogViewerHandler 创建日志查看器处理器
func NewLogViewerHandler(logDir string) *LogViewerHandler {
	return &LogViewerHandler{
		logDir: logDir,
	}
}

// ListLogFiles 列出所有日志文件
func (h *LogViewerHandler) ListLogFiles(c *gin.Context) {
	files, err := os.ReadDir(h.logDir)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "无法读取日志目录",
		})
		return
	}

	type logFileInfo struct {
		Name     string `json:"name"`
		Size     int64  `json:"size"`
		Modified string `json:"modified"`
	}

	var logFiles []logFileInfo
	for _, file := range files {
		if file.IsDir() || !strings.HasSuffix(file.Name(), ".log") {
			continue
		}

		info, err := file.Info()
		if err != nil {
			continue
		}

		logFiles = append(logFiles, logFileInfo{
			Name:     file.Name(),
			Size:     info.Size(),
			Modified: info.ModTime().Format("2006-01-02 15:04:05"),
		})
	}

	// 按修改时间倒序排序
	sort.Slice(logFiles, func(i, j int) bool {
		return logFiles[i].Modified > logFiles[j].Modified
	})

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    logFiles,
	})
}

// ViewLogFile 查看日志文件内容
func (h *LogViewerHandler) ViewLogFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文件名不能为空",
		})
		return
	}

	// 安全检查：防止路径遍历攻击
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "非法的文件名",
		})
		return
	}

	logPath := filepath.Join(h.logDir, filename)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "日志文件不存在",
		})
		return
	}

	// 获取查询参数
	maxLines := 1000
	if linesParam := c.Query("lines"); linesParam != "" {
		if lines, err := strconv.Atoi(linesParam); err == nil && lines > 0 {
			maxLines = lines
			if maxLines > 10000 { // 限制最大行数
				maxLines = 10000
			}
		}
	}

	level := c.Query("level")     // 日志级别过滤: INFO, WARN, ERROR
	keyword := c.Query("keyword") // 关键词搜索

	// 读取日志文件
	lines, err := h.readLogLinesWithFilter(logPath, maxLines, level, keyword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "读取日志文件失败",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data": map[string]interface{}{
			"filename": filename,
			"lines":    lines,
			"total":    len(lines),
		},
	})
}

// DownloadLogFile 下载日志文件
func (h *LogViewerHandler) DownloadLogFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文件名不能为空",
		})
		return
	}

	// 安全检查
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "非法的文件名",
		})
		return
	}

	logPath := filepath.Join(h.logDir, filename)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "日志文件不存在",
		})
		return
	}

	c.FileAttachment(logPath, filename)
}

// readLogLines 读取日志文件的最后N行
func (h *LogViewerHandler) readLogLines(logPath string, maxLines int) ([]string, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer file.Close()

	var lines []string
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		lines = append(lines, scanner.Text())
		if len(lines) > maxLines {
			lines = lines[1:] // 保持最后maxLines行
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取日志文件失败: %w", err)
	}

	return lines, nil
}

// readLogLinesWithFilter 读取日志文件的最后N行，支持过滤
func (h *LogViewerHandler) readLogLinesWithFilter(logPath string, maxLines int, level, keyword string) ([]string, error) {
	file, err := os.Open(logPath)
	if err != nil {
		return nil, fmt.Errorf("打开日志文件失败: %w", err)
	}
	defer file.Close()

	var allLines []string
	scanner := bufio.NewScanner(file)

	// 先读取所有行
	for scanner.Scan() {
		allLines = append(allLines, scanner.Text())
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("读取日志文件失败: %w", err)
	}

	// 应用过滤条件
	var filteredLines []string
	levelUpper := strings.ToUpper(level)
	keywordLower := strings.ToLower(keyword)

	for _, line := range allLines {
		// 日志级别过滤
		if levelUpper != "" {
			if !strings.Contains(strings.ToUpper(line), "level="+levelUpper) &&
				!strings.Contains(strings.ToUpper(line), levelUpper) {
				continue
			}
		}

		// 关键词过滤
		if keywordLower != "" {
			if !strings.Contains(strings.ToLower(line), keywordLower) {
				continue
			}
		}

		filteredLines = append(filteredLines, line)
	}

	// 保持最后maxLines行
	if len(filteredLines) > maxLines {
		filteredLines = filteredLines[len(filteredLines)-maxLines:]
	}

	return filteredLines, nil
}

// StreamLogFile 实时流式传输日志文件
func (h *LogViewerHandler) StreamLogFile(c *gin.Context) {
	filename := c.Param("filename")
	if filename == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "文件名不能为空",
		})
		return
	}

	// 安全检查
	if strings.Contains(filename, "..") || strings.Contains(filename, "/") || strings.Contains(filename, "\\") {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "非法的文件名",
		})
		return
	}

	logPath := filepath.Join(h.logDir, filename)
	if _, err := os.Stat(logPath); os.IsNotExist(err) {
		c.JSON(http.StatusNotFound, gin.H{
			"success": false,
			"error":   "日志文件不存在",
		})
		return
	}

	// 获取过滤参数
	level := c.Query("level")
	keyword := c.Query("keyword")

	// 设置SSE响应头
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Header("Access-Control-Allow-Origin", "*")

	// 打开文件
	file, err := os.Open(logPath)
	if err != nil {
		c.SSEvent("error", "无法打开日志文件")
		return
	}
	defer file.Close()

	// 跳转到文件末尾
	file.Seek(0, io.SeekEnd)

	// 创建通知通道
	clientGone := c.Request.Context().Done()
	ticker := time.NewTicker(500 * time.Millisecond)
	defer ticker.Stop()

	c.Stream(func(w io.Writer) bool {
		select {
		case <-clientGone:
			return false
		case <-ticker.C:
			// 读取新行
			scanner := bufio.NewScanner(file)
			for scanner.Scan() {
				line := scanner.Text()

				// 应用过滤
				if level != "" {
					levelUpper := strings.ToUpper(level)
					if !strings.Contains(strings.ToUpper(line), "level="+levelUpper) &&
						!strings.Contains(strings.ToUpper(line), levelUpper) {
						continue
					}
				}

				if keyword != "" {
					if !strings.Contains(strings.ToLower(line), strings.ToLower(keyword)) {
						continue
					}
				}

				c.SSEvent("message", line)
			}
			return true
		}
	})
}
