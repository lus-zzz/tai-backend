package services

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"chat-backend/models"
	"chat-backend/utils"
)

// ShortcutService 快捷方式服务
type ShortcutService struct {
	baseURL    string
	httpClient *http.Client
}

// NewShortcutService 创建快捷方式服务
func NewShortcutService(baseURL string) *ShortcutService {
	return &ShortcutService{
		baseURL: baseURL,
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
	}
}

// RecommendSettings 调用推荐设置API
func (s *ShortcutService) RecommendSettings(input *models.ChatInput) (*models.RecommendData, error) {
	utils.LogInfo("调用推荐设置API: %s", input.UserInput)

	// 构建请求URL
	url := fmt.Sprintf("%s/shortcut/recommend", s.baseURL)

	// 序列化请求体
	requestBody, err := json.Marshal(input)
	if err != nil {
		utils.LogError("序列化请求体失败: %v", err)
		return nil, fmt.Errorf("序列化请求体失败: %v", err)
	}

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(requestBody))
	if err != nil {
		utils.LogError("创建HTTP请求失败: %v", err)
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		utils.LogError("发送HTTP请求失败: %v", err)
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogError("读取响应体失败: %v", err)
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		utils.LogError("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBody))
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result models.RecommendData
	if err := json.Unmarshal(responseBody, &result); err != nil {
		utils.LogError("解析响应失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	utils.LogInfo("推荐设置API调用成功，返回%d个推荐", len(result.SettingName))
	return &result, nil
}

// GetSupportedSettings 调用获取支持设置API
func (s *ShortcutService) GetSupportedSettings() (*models.SettingName, error) {
	utils.LogInfo("调用获取支持设置API")

	// 构建请求URL
	url := fmt.Sprintf("%s/shortcut/supportedSetting", s.baseURL)

	// 创建HTTP请求
	req, err := http.NewRequest("POST", url, nil)
	if err != nil {
		utils.LogError("创建HTTP请求失败: %v", err)
		return nil, fmt.Errorf("创建HTTP请求失败: %v", err)
	}

	// 设置请求头
	req.Header.Set("Content-Type", "application/json")

	// 发送请求
	resp, err := s.httpClient.Do(req)
	if err != nil {
		utils.LogError("发送HTTP请求失败: %v", err)
		return nil, fmt.Errorf("发送HTTP请求失败: %v", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		utils.LogError("读取响应体失败: %v", err)
		return nil, fmt.Errorf("读取响应体失败: %v", err)
	}

	// 检查HTTP状态码
	if resp.StatusCode != http.StatusOK {
		utils.LogError("API请求失败，状态码: %d, 响应: %s", resp.StatusCode, string(responseBody))
		return nil, fmt.Errorf("API请求失败，状态码: %d", resp.StatusCode)
	}

	// 解析响应
	var result models.SettingName
	if err := json.Unmarshal(responseBody, &result); err != nil {
		utils.LogError("解析响应失败: %v", err)
		return nil, fmt.Errorf("解析响应失败: %v", err)
	}

	utils.LogInfo("获取支持设置API调用成功，返回%d个设置", len(result.SupportedSettingName))
	return &result, nil
}
