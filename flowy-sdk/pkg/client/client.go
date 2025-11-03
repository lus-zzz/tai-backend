package client

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"strings"
	"time"

	"flowy-sdk/pkg/config"
	"flowy-sdk/pkg/errors"
	"flowy-sdk/pkg/models"
)

// HTTPClient HTTP客户端接口
type HTTPClient interface {
	Get(ctx context.Context, path string, params map[string]string) (*models.BaseResponse, error)
	Post(ctx context.Context, path string, body interface{}) (*models.BaseResponse, error)
	PostSSE(ctx context.Context, path string, body interface{}) (io.ReadCloser, error)
	Put(ctx context.Context, path string, body interface{}) (*models.BaseResponse, error)
	Delete(ctx context.Context, path string) (*models.BaseResponse, error)
	Upload(ctx context.Context, path string, fieldName string, filename string, fileContent io.Reader, params map[string]string) (*models.BaseResponse, error)
	Download(ctx context.Context, path string) (io.ReadCloser, error)
}

// Client FLOWY HTTP客户端实现
type Client struct {
	config        *config.Config
	httpClient    *http.Client
	httpClientSSE *http.Client // SSE 专用客户端，无超时限制
	baseURL       string
}

// New 创建新的HTTP客户端
func New(cfg *config.Config) *Client {
	if cfg == nil {
		cfg = config.DefaultConfig()
	}

	// 验证配置
	if err := cfg.Validate(); err != nil {
		panic(fmt.Sprintf("invalid config: %v", err))
	}

	// 创建普通HTTP客户端（有超时）
	httpClient := &http.Client{
		Timeout: cfg.GetTimeoutDuration(),
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipTLSVerify,
			},
		},
	}

	// 创建SSE专用HTTP客户端（无超时限制，适用于流式响应）
	httpClientSSE := &http.Client{
		Timeout: 0, // 无超时限制
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: cfg.SkipTLSVerify,
			},
		},
	}

	// 设置代理
	if cfg.ProxyURL != "" {
		if proxyURL, err := url.Parse(cfg.ProxyURL); err == nil {
			httpClient.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
			httpClientSSE.Transport.(*http.Transport).Proxy = http.ProxyURL(proxyURL)
		}
	}

	return &Client{
		config:        cfg,
		httpClient:    httpClient,
		httpClientSSE: httpClientSSE,
		baseURL:       strings.TrimRight(cfg.BaseURL, "/"),
	}
}

// Get 发送GET请求
func (c *Client) Get(ctx context.Context, path string, params map[string]string) (*models.BaseResponse, error) {
	url := c.buildURL(path, params)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	return c.doRequest(req)
}

// Post 发送POST请求
func (c *Client) Post(ctx context.Context, path string, body interface{}) (*models.BaseResponse, error) {
	url := c.buildURL(path, nil)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to marshal request body").WithDetails(err.Error())
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequest(req)
}

// PostSSE 发送POST请求并返回SSE流
func (c *Client) PostSSE(ctx context.Context, path string, body interface{}) (io.ReadCloser, error) {
	url := c.buildURL(path, nil)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to marshal request body").WithDetails(err.Error())
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, reqBody)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}
	req.Header.Set("Accept", "text/event-stream")

	c.setAuthHeaders(req)

	// 使用 SSE 专用客户端（无超时限制）
	resp, err := c.httpClientSSE.Do(req)
	if err != nil {
		return nil, errors.New(errors.ErrCodeNetworkError, "network request failed").WithDetails(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.FromHTTPStatus(resp.StatusCode, "SSE request failed")
	}

	return resp.Body, nil
}

// Put 发送PUT请求
func (c *Client) Put(ctx context.Context, path string, body interface{}) (*models.BaseResponse, error) {
	url := c.buildURL(path, nil)

	var reqBody io.Reader
	if body != nil {
		jsonData, err := json.Marshal(body)
		if err != nil {
			return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to marshal request body").WithDetails(err.Error())
		}
		reqBody = bytes.NewReader(jsonData)
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPut, url, reqBody)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	if body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	return c.doRequest(req)
}

// Delete 发送DELETE请求
func (c *Client) Delete(ctx context.Context, path string) (*models.BaseResponse, error) {
	url := c.buildURL(path, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodDelete, url, nil)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	return c.doRequest(req)
}

// Upload 上传文件
func (c *Client) Upload(ctx context.Context, path string, fieldName string, filename string, fileContent io.Reader, params map[string]string) (*models.BaseResponse, error) {
	url := c.buildURL(path, nil)

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)

	// 添加文件
	part, err := writer.CreateFormFile(fieldName, filename)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create form file").WithDetails(err.Error())
	}

	if _, err := io.Copy(part, fileContent); err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to copy file content").WithDetails(err.Error())
	}

	// 添加其他参数
	for key, value := range params {
		if err := writer.WriteField(key, value); err != nil {
			return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to write field").WithDetails(err.Error())
		}
	}

	if err := writer.Close(); err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to close writer").WithDetails(err.Error())
	}

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, url, &buf)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	req.Header.Set("Content-Type", writer.FormDataContentType())

	return c.doRequest(req)
}

// Download 下载文件
func (c *Client) Download(ctx context.Context, path string) (io.ReadCloser, error) {
	url := c.buildURL(path, nil)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInvalidRequest, "failed to create request").WithDetails(err.Error())
	}

	c.setAuthHeaders(req)

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		resp.Body.Close()
		return nil, errors.FromHTTPStatus(resp.StatusCode, "download failed")
	}

	return resp.Body, nil
}

// doRequest 执行HTTP请求
func (c *Client) doRequest(req *http.Request) (*models.BaseResponse, error) {
	c.setAuthHeaders(req)

	resp, err := c.doRequestWithRetry(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return c.parseResponse(resp)
}

// doRequestWithRetry 带重试的请求执行
func (c *Client) doRequestWithRetry(req *http.Request) (*http.Response, error) {
	var resp *http.Response
	var err error

	for i := 0; i <= c.config.MaxRetries; i++ {
		if i > 0 {
			time.Sleep(c.config.GetRetryIntervalDuration())
		}

		resp, err = c.httpClient.Do(req)
		if err == nil && resp.StatusCode < 500 {
			break
		}

		if resp != nil {
			resp.Body.Close()
		}
	}

	if err != nil {
		return nil, errors.New(errors.ErrCodeNetworkError, "network request failed").WithDetails(err.Error())
	}

	return resp, nil
}

// parseResponse 解析响应
func (c *Client) parseResponse(resp *http.Response) (*models.BaseResponse, error) {
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New(errors.ErrCodeInternalError, "failed to read response body").WithDetails(err.Error())
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.FromHTTPStatus(resp.StatusCode, string(body))
	}

	var response models.BaseResponse
	if err := json.Unmarshal(body, &response); err != nil {
		return nil, errors.New(errors.ErrCodeInternalError, "failed to parse response").WithDetails(err.Error())
	}

	return &response, nil
}

// setAuthHeaders 设置认证头
func (c *Client) setAuthHeaders(req *http.Request) {
	if c.config.Token != "" {
		req.Header.Set("Authorization", "Bearer "+c.config.Token)
	} else if c.config.APIKey != "" {
		req.Header.Set("X-API-Key", c.config.APIKey)
		if c.config.SecretKey != "" {
			req.Header.Set("X-Secret-Key", c.config.SecretKey)
		}
	}

	req.Header.Set("User-Agent", "FLOWY-SDK-Go/1.0.0")
}

// buildURL 构建完整URL
func (c *Client) buildURL(path string, params map[string]string) string {
	url := c.baseURL + "/" + strings.TrimLeft(path, "/")

	if len(params) > 0 {
		values := make([]string, 0, len(params))
		for key, value := range params {
			values = append(values, key+"="+value)
		}
		url += "?" + strings.Join(values, "&")
	}

	return url
}
