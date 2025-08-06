package github

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/google/go-github/v74/github"
	"golang.org/x/oauth2"
)

// Client 封装了 GitHub 客户端
type Client struct {
	*github.Client
	tokens   []string
	tokenIdx int
}

// NewClient 创建一个新的 GitHub 客户端，支持令牌轮换
func NewClient(tokens []string) *Client {
	if len(tokens) == 0 {
		return &Client{
			Client: github.NewClient(nil),
			tokens: tokens,
		}
	}

	return &Client{
		Client:   github.NewClient(nil).WithAuthToken(tokens[0]),
		tokens:   tokens,
		tokenIdx: 0,
	}
}

// SearchCode 在 GitHub 上搜索代码，支持速率限制处理和文本匹配
func (c *Client) SearchCode(ctx context.Context, query string, opts *github.SearchOptions) (*github.CodeSearchResult, *github.Response, error) {
	// 为上下文添加超时
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	// 添加延迟以避免速率限制
	time.Sleep(100 * time.Millisecond)
	opt := &github.SearchOptions{
		ListOptions: github.ListOptions{
			Page:    opts.Page,
			PerPage: opts.PerPage,
		},
		TextMatch: true,
	}
	// 手动构造请求以添加文本匹配媒体类型
	result, resp, err := c.Client.Search.Code(ctx, query, opt)
	// 执行请求
	if err == nil {
		return result, resp, nil
	}

	// 检查是否为速率限制错误
	if rateLimitErr, ok := err.(*github.RateLimitError); ok {
		// 计算等待速率限制重置的持续时间
		resetTime := rateLimitErr.Rate.Reset.Time
		sleepDuration := time.Until(resetTime) + (time.Second * 1) // 添加一秒以确保安全

		// 记录速率限制错误并等待
		fmt.Printf("达到主要速率限制。等待 %v 直到重置。\n", sleepDuration)
		time.Sleep(sleepDuration)

		// 重试请求
		return c.SearchCode(ctx, query, opts) // 使用相同函数重试
	}

	// 检查是否为滥用速率限制错误（次要速率限制）
	if abuseRateLimitErr, ok := err.(*github.AbuseRateLimitError); ok {
		// 获取重试后的持续时间
		var retryAfter time.Duration
		if abuseRateLimitErr.RetryAfter != nil {
			retryAfter = *abuseRateLimitErr.RetryAfter
		} else {
			// 如果没有提供 RetryAfter，则默认为 1 分钟
			retryAfter = time.Minute
		}

		// 记录滥用速率限制错误并等待
		fmt.Printf("达到次要速率限制。等待 %v 后重试。\n", retryAfter)
		time.Sleep(retryAfter)

		// 重试请求
		return c.SearchCode(ctx, query, opts) // 使用相同函数重试
	}

	// 如果不是速率限制错误，则返回错误
	return nil, resp, err
}

// WithToken 为客户端设置特定的令牌
func (c *Client) WithToken(token string) *Client {
	ctx := context.Background()
	ts := oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: token},
	)
	tc := oauth2.NewClient(ctx, ts)
	// 为 HTTP 客户端设置超时
	tc.Timeout = 30 * time.Second

	return &Client{
		Client:   github.NewClient(tc),
		tokens:   c.tokens,
		tokenIdx: c.tokenIdx,
	}
}

// RotateToken 轮换到列表中的下一个令牌
func (c *Client) RotateToken() *Client {
	if len(c.tokens) <= 1 {
		return c
	}

	// 移动到下一个令牌
	c.tokenIdx = (c.tokenIdx + 1) % len(c.tokens)

	// 使用新令牌
	return &Client{
		Client:   github.NewClient(nil).WithAuthToken(c.tokens[c.tokenIdx]),
		tokens:   c.tokens,
		tokenIdx: c.tokenIdx,
	}
}

// SetHTTPClient 设置自定义 HTTP 客户端
func (c *Client) SetHTTPClient(httpClient *http.Client) *Client {
	// 如果 HTTP 客户端没有设置超时，则设置一个超时
	if httpClient.Timeout == 0 {
		httpClient.Timeout = 30 * time.Second
	}
	return &Client{
		Client:   github.NewClient(httpClient),
		tokens:   c.tokens,
		tokenIdx: c.tokenIdx,
	}
}
