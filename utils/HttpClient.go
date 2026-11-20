package utils

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

// ClientConfig HTTP客户端配置
type ClientConfig struct {
	Timeout             time.Duration
	MaxIdleConns        int
	MaxIdleConnsPerHost int
	IdleConnTimeout     time.Duration
	InsecureSkipVerify  bool // 跳过SSL验证（仅测试环境使用）
}

// DefaultConfig 默认配置
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		Timeout:             30 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		InsecureSkipVerify:  false,
	}
}

// HTTPClient 封装的HTTP客户端
type HTTPClient struct {
	client  *http.Client
	baseURL string
	headers map[string]string
}

// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(config *ClientConfig) *HTTPClient {
	if config == nil {
		config = DefaultConfig()
	}

	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     config.IdleConnTimeout,
	}

	if config.InsecureSkipVerify {
		transport.TLSClientConfig = &tls.Config{InsecureSkipVerify: true}
	}

	return &HTTPClient{
		client: &http.Client{
			Timeout:   config.Timeout,
			Transport: transport,
		},
		headers: make(map[string]string),
	}
}

// SetBaseURL 设置基础URL
func (c *HTTPClient) SetBaseURL(baseURL string) {
	c.baseURL = baseURL
}

// SetHeader 设置默认请求头
func (c *HTTPClient) SetHeader(key, value string) {
	c.headers[key] = value
}

// SetHeaders 批量设置请求头
func (c *HTTPClient) SetHeaders(headers map[string]string) {
	for k, v := range headers {
		c.headers[k] = v
	}
}

// Response 封装的响应结构
type Response struct {
	StatusCode int
	Body       []byte
	Header     http.Header
}

// 请求选项
type RequestOption func(*http.Request)

// WithContext 设置上下文
func WithContext(ctx context.Context) RequestOption {
	return func(req *http.Request) {
		req = req.WithContext(ctx)
	}
}

// WithHeader 添加请求头
func WithHeader(key, value string) RequestOption {
	return func(req *http.Request) {
		req.Header.Set(key, value)
	}
}

// WithHeaders 批量添加请求头
func WithHeaders(headers map[string]string) RequestOption {
	return func(req *http.Request) {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}
}

// WithQuery 添加查询参数
func WithQuery(params map[string]string) RequestOption {
	return func(req *http.Request) {
		q := req.URL.Query()
		for k, v := range params {
			q.Set(k, v)
		}
		req.URL.RawQuery = q.Encode()
	}
}

// doRequest 执行HTTP请求的核心方法
func (c *HTTPClient) doRequest(method, path string, body io.Reader, opts ...RequestOption) (*Response, error) {
	// 构建完整URL
	fullURL := path
	if c.baseURL != "" {
		// 简单拼接，生产环境建议使用 url.JoinPath
		fullURL = c.baseURL + path
	}

	// 创建请求
	req, err := http.NewRequest(method, fullURL, body)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	// 设置默认请求头
	for k, v := range c.headers {
		req.Header.Set(k, v)
	}

	// 如果没有设置Content-Type，默认设置为application/json
	if req.Header.Get("Content-Type") == "" && body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	// 应用选项
	for _, opt := range opts {
		opt(req)
	}

	// 执行请求
	resp, err := c.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求执行失败: %w", err)
	}
	defer resp.Body.Close()

	// 读取响应体
	bodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("读取响应体失败: %w", err)
	}

	return &Response{
		StatusCode: resp.StatusCode,
		Body:       bodyBytes,
		Header:     resp.Header,
	}, nil
}

// Get 发送GET请求
func (c *HTTPClient) Get(path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(http.MethodGet, path, nil, opts...)
}

// Post 发送POST请求（JSON格式）
func (c *HTTPClient) Post(path string, data interface{}, opts ...RequestOption) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}
	return c.doRequest(http.MethodPost, path, bytes.NewReader(body), opts...)
}

// PostRaw 发送原始POST请求
func (c *HTTPClient) PostRaw(path string, body io.Reader, opts ...RequestOption) (*Response, error) {
	return c.doRequest(http.MethodPost, path, body, opts...)
}

// Put 发送PUT请求（JSON格式）
func (c *HTTPClient) Put(path string, data interface{}, opts ...RequestOption) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}
	return c.doRequest(http.MethodPut, path, bytes.NewReader(body), opts...)
}

// PutRaw 发送原始PUT请求
func (c *HTTPClient) PutRaw(path string, body io.Reader, opts ...RequestOption) (*Response, error) {
	return c.doRequest(http.MethodPut, path, body, opts...)
}

// Delete 发送DELETE请求
func (c *HTTPClient) Delete(path string, opts ...RequestOption) (*Response, error) {
	return c.doRequest(http.MethodDelete, path, nil, opts...)
}

// Patch 发送PATCH请求（JSON格式）
func (c *HTTPClient) Patch(path string, data interface{}, opts ...RequestOption) (*Response, error) {
	body, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("JSON序列化失败: %w", err)
	}
	return c.doRequest(http.MethodPatch, path, bytes.NewReader(body), opts...)
}

// PatchRaw 发送原始PATCH请求
func (c *HTTPClient) PatchRaw(path string, body io.Reader, opts ...RequestOption) (*Response, error) {
	return c.doRequest(http.MethodPatch, path, body, opts...)
}

// 辅助方法

// ToJSON 将响应体解析为JSON
func (r *Response) ToJSON(v interface{}) error {
	if err := json.Unmarshal(r.Body, v); err != nil {
		return fmt.Errorf("JSON解析失败: %w", err)
	}
	return nil
}

// ToString 将响应体转为字符串
func (r *Response) ToString() string {
	return string(r.Body)
}

// IsSuccess 判断是否成功（2xx状态码）
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsError 判断是否错误（4xx或5xx状态码）
func (r *Response) IsError() bool {
	return r.StatusCode >= 400
}

//使用示例
// package main

// import (
//     "context"
//     "fmt"
//     "log"
//     "time"

//     "your-project/httputil"
// )

// func main() {
//     // 1. 创建客户端
//     config := &httputil.ClientConfig{
//         Timeout: 10 * time.Second,
//         MaxIdleConns: 50,
//     }
//     client := httputil.NewHTTPClient(config)

//     // 2. 设置基础URL和默认请求头
//     client.SetBaseURL("https://api.example.com")
//     client.SetHeader("User-Agent", "MyApp/1.0")
//     client.SetHeader("Accept", "application/json")

//     // 3. GET请求示例
//     getExample(client)

//     // 4. POST请求示例
//     postExample(client)

//     // 5. 带查询参数和上下文的请求
//     queryExample(client)
// }

// func getExample(client *httputil.HTTPClient) {
//     resp, err := client.Get("/users/1")
//     if err != nil {
//         log.Printf("请求失败: %v", err)
//         return
//     }

//     if !resp.IsSuccess() {
//         log.Printf("HTTP错误: %d, 响应: %s", resp.StatusCode, resp.ToString())
//         return
//     }

//     var user map[string]interface{}
//     if err := resp.ToJSON(&user); err != nil {
//         log.Printf("解析JSON失败: %v", err)
//         return
//     }

//     fmt.Printf("用户信息: %+v\n", user)
// }

// func postExample(client *httputil.HTTPClient) {
//     // 准备请求数据
//     data := map[string]interface{}{
//         "name":  "张三",
//         "email": "zhangsan@example.com",
//         "age":   25,
//     }

//     // 添加自定义请求头
//     resp, err := client.Post("/users", data,
//         httputil.WithHeader("X-Request-ID", "12345"),
//         httputil.WithHeader("Authorization", "Bearer token123"),
//     )

//     if err != nil {
//         log.Printf("POST请求失败: %v", err)
//         return
//     }

//     if resp.IsSuccess() {
//         fmt.Printf("创建成功: %s\n", resp.ToString())
//     } else {
//         fmt.Printf("创建失败: %d, %s\n", resp.StatusCode, resp.ToString())
//     }
// }

// func queryExample(client *httputil.HTTPClient) {
//     // 带查询参数和上下文
//     ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
//     defer cancel()

//     resp, err := client.Get("/search",
//         httputil.WithContext(ctx),
//         httputil.WithQuery(map[string]string{
//             "q":     "golang",
//             "page":  "1",
//             "limit": "10",
//         }),
//     )

//     if err != nil {
//         log.Printf("查询失败: %v", err)
//         return
//     }

//     fmt.Printf("查询结果: %s\n", resp.ToString())
// }

// // 流式下载示例
// func downloadExample(client *httputil.HTTPClient) {
//     // 使用Raw方法处理大文件
//     resp, err := client.GetRaw("/files/large.zip")
//     if err != nil {
//         log.Printf("下载失败: %v", err)
//         return
//     }

//     // 这里可以流式处理 resp.Body
//     // 注意：resp.Body已经在doRequest中被关闭，如果需要在外部处理，
//     // 需要修改doRequest方法，返回原始响应体
// }

// // 文件上传示例
// func uploadExample(client *httputil.HTTPClient) {
//     // 构造multipart/form-data请求
//     // 这里需要更复杂的处理，可以使用PostRaw方法
//     // 或者扩展工具函数支持multipart
// }
