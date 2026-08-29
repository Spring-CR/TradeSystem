package post_client

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"sync/atomic"
	"time"
)

// ClientConfig 客户端配置参数
type ClientConfig struct {
	MaxRetries          int           // 最大重试次数
	InitialRetryDelay   time.Duration // 初始重试延迟
	MaxRetryDelay       time.Duration // 最大重试延迟
	RequestTimeout      time.Duration // 单个请求超时时间
	ConnectionTimeout   time.Duration // 连接超时时间
	MaxIdleConns        int           // 最大空闲连接数
	MaxIdleConnsPerHost int           // 每个主机最大空闲连接数
	QueueSize           int           // 请求队列大小
	WorkerCount         int           // 工作协程数量
	EnableJitter        bool          // 是否启用重试抖动
}
// 默认配置
func DefaultConfig() *ClientConfig {
	return &ClientConfig{
		MaxRetries:          3,
		InitialRetryDelay:   100 * time.Millisecond,
		MaxRetryDelay:       5 * time.Second,
		RequestTimeout:      10 * time.Second,
		ConnectionTimeout:   3 * time.Second,
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 50,
		QueueSize:           1000,
		WorkerCount:         50,
		EnableJitter:        true,
	}
}
// HTTPClient 高性能HTTP客户端
type HTTPClient struct {
	client      *http.Client
	requestChan chan *Request
	wg          sync.WaitGroup
	closed      atomic.Bool
	config      *ClientConfig
}
// Request 封装HTTP请求
type Request struct {
	URL      string
	Payload  interface{}
	Headers  map[string]string
	Retries  int
	Callback func(*Result) // 结果回调函数
}
// Result 请求结果
type Result struct {
	StatusCode int
	Body       []byte
	Err        error
}
// Future 用于等待结果的Future对象
type Future struct {
	done   chan struct{}
	result *Result
}
// NewHTTPClient 创建新的HTTP客户端
func NewHTTPClient(config *ClientConfig) *HTTPClient {
	// 使用默认配置（如果未提供）
	if config == nil {
		config = DefaultConfig()
	}
	// 配置HTTP传输层
	transport := &http.Transport{
		MaxIdleConns:        config.MaxIdleConns,
		MaxIdleConnsPerHost: config.MaxIdleConnsPerHost,
		IdleConnTimeout:     90 * time.Second,
		DisableCompression:  false,
	}
	// 创建HTTP客户端
	client := &http.Client{
		Transport: transport,
		Timeout:   config.ConnectionTimeout,
	}
	hc := &HTTPClient{
		client:      client,
		requestChan: make(chan *Request, config.QueueSize),
		config:      config,
	}
	// 启动工作协程池
	for i := 0; i < config.WorkerCount; i++ {
		hc.wg.Add(1)
		go hc.worker()
	}
	return hc
}
// worker 处理请求的工作协程
func (hc *HTTPClient) worker() {
	defer hc.wg.Done()
	for req := range hc.requestChan {
		// 检查客户端是否已关闭
		if hc.closed.Load() {
			if req.Callback != nil {
				req.Callback(&Result{Err: errors.New("client closed")})
			}
			continue
		}
		// 执行请求（带重试）
		result := hc.doRequestWithRetry(req)
		
		// 调用回调函数
		if req.Callback != nil {
			req.Callback(result)
		}
	}
}
// doRequestWithRetry 执行请求并处理重试
func (hc *HTTPClient) doRequestWithRetry(req *Request) *Result {
	// 创建上下文（带超时）
	ctx, cancel := context.WithTimeout(context.Background(), hc.config.RequestTimeout)
	defer cancel()
	// 准备请求体
	payload, err := json.Marshal(req.Payload)
	if err != nil {
		return &Result{Err: fmt.Errorf("payload marshal error: %w", err)}
	}
	// 指数退避重试
	retryDelay := hc.config.InitialRetryDelay
	for attempt := 0; attempt <= req.Retries; attempt++ {
		// 创建新请求（每次重试都需要创建新请求）
		httpReq, err := http.NewRequestWithContext(ctx, "POST", req.URL, bytes.NewBuffer(payload))
		if err != nil {
			return &Result{Err: fmt.Errorf("create request error: %w", err)}
		}
		// 设置请求头
		httpReq.Header.Set("Content-Type", "application/json")
		for k, v := range req.Headers {
			httpReq.Header.Set(k, v)
		}
		// 执行请求
		resp, err := hc.client.Do(httpReq)
		if err != nil {
			// 如果是超时错误或可重试错误
			if isRetryableError(err) && attempt < req.Retries {
				// 计算退避时间（带可选抖动）
				delay := calculateBackoff(retryDelay, attempt, hc.config)
				
				select {
				case <-time.After(delay):
					retryDelay = time.Duration(math.Min(
						float64(hc.config.MaxRetryDelay), 
						float64(2*retryDelay),
					))
					continue
				case <-ctx.Done():
					return &Result{Err: ctx.Err()}
				}
			}
			return &Result{Err: fmt.Errorf("request failed: %w", err)}
		}
		// 处理响应
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			return &Result{StatusCode: resp.StatusCode, Err: fmt.Errorf("read body error: %w", err)}
		}
		// 检查状态码 - 5xx错误重试
		if resp.StatusCode >= 500 && attempt < req.Retries {
			delay := calculateBackoff(retryDelay, attempt, hc.config)
			
			select {
			case <-time.After(delay):
				retryDelay = time.Duration(math.Min(
					float64(hc.config.MaxRetryDelay), 
					float64(2*retryDelay),
				))
				continue
			case <-ctx.Done():
				return &Result{StatusCode: resp.StatusCode, Body: body, Err: ctx.Err()}
			}
		}
		return &Result{
			StatusCode: resp.StatusCode,
			Body:       body,
		}
	}
	return &Result{Err: fmt.Errorf("max retries exceeded (%d attempts)", req.Retries)}
}
// 计算退避时间（带抖动）
func calculateBackoff(baseDelay time.Duration, attempt int, config *ClientConfig) time.Duration {
	// 指数退避
	delay := time.Duration(float64(baseDelay) * math.Pow(2, float64(attempt)))
	
	// 限制最大延迟
	if delay > config.MaxRetryDelay {
		delay = config.MaxRetryDelay
	}
	
	// 添加随机抖动
	if config.EnableJitter {
		jitter := time.Duration(rand.Int63n(int64(delay / 2)))
		delay += jitter
	}
	
	return delay
}
// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	// 超时错误
	if errors.Is(err, context.DeadlineExceeded) {
		return true
	}
	
	// 网络错误
	var netErr interface{ Timeout() bool }
	if errors.As(err, &netErr) && netErr.Timeout() {
		return true
	}
	
	// 连接错误
	if errors.Is(err, context.Canceled) {
		return false // 用户取消，不重试
	}
	
	// 其他网络错误
	return true
}
// PostAsync 发送POST请求（异步回调方式）
func (hc *HTTPClient) PostAsync(url string, payload interface{}, headers map[string]string, callback func(*Result)) error {
	// 如果客户端已关闭，立即返回错误
	if hc.closed.Load() {
		return errors.New("client closed")
	}
	// 创建请求对象
	req := &Request{
		URL:      url,
		Payload:  payload,
		Headers:  headers,
		Retries:  hc.config.MaxRetries,
		Callback: callback,
	}
	// 将请求发送到队列
	select {
	case hc.requestChan <- req:
		return nil // 成功入队
	default:
		return errors.New("request queue full")
	}
}
// Post 发送POST请求（返回Future对象）
func (hc *HTTPClient) Post(url string, payload interface{}, headers map[string]string) (*Future, error) {
	// 如果客户端已关闭，立即返回错误
	if hc.closed.Load() {
		return nil, errors.New("client closed")
	}
	// 创建Future对象
	future := &Future{
		done: make(chan struct{}),
	}
	// 创建请求对象
	req := &Request{
		URL:     url,
		Payload: payload,
		Headers: headers,
		Retries: hc.config.MaxRetries,
		Callback: func(result *Result) {
			future.result = result
			close(future.done)
		},
	}
	// 将请求发送到队列
	select {
	case hc.requestChan <- req:
		return future, nil
	default:
		return nil, errors.New("request queue full")
	}
}
// Get 获取Future的结果（阻塞）
func (f *Future) Get() *Result {
	<-f.done
	return f.result
}
// GetWithTimeout 获取Future的结果（带超时）
func (f *Future) GetWithTimeout(timeout time.Duration) (*Result, error) {
	select {
	case <-f.done:
		return f.result, nil
	case <-time.After(timeout):
		return nil, errors.New("timeout waiting for response")
	}
}
// Close 关闭客户端
func (hc *HTTPClient) Close() {
	// 标记为已关闭
	hc.closed.Store(true)
	
	// 关闭请求通道
	close(hc.requestChan)
	
	// 等待所有工作协程完成
	hc.wg.Wait()
}
