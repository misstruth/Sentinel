package github

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Client GitHub API 客户端
type Client struct {
	httpClient *http.Client
	token      string
	baseURL    string
}

// NewClient 创建 GitHub 客户端
func NewClient(token string) *Client {
	return &Client{
		httpClient: &http.Client{Timeout: 30 * time.Second},
		token:      token,
		baseURL:    "https://api.github.com",
	}
}

// doRequest 执行 HTTP 请求
func (c *Client) doRequest(ctx context.Context, method, path string) (*http.Response, error) {
	url := c.baseURL + path
	req, err := http.NewRequestWithContext(ctx, method, url, nil)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Accept", "application/vnd.github.v3+json")
	if c.token != "" {
		req.Header.Set("Authorization", "Bearer "+c.token)
	}

	return c.httpClient.Do(req)
}

// decodeResponse 解析响应
func decodeResponse(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub API error: %d", resp.StatusCode)
	}
	return json.NewDecoder(resp.Body).Decode(v)
}

// GetRepository 获取仓库信息
func (c *Client) GetRepository(ctx context.Context, owner, repo string) (*Repository, error) {
	path := fmt.Sprintf("/repos/%s/%s", owner, repo)
	resp, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var result Repository
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

// GetReleases 获取仓库 Release 列表
func (c *Client) GetReleases(ctx context.Context, owner, repo string, limit int) ([]Release, error) {
	path := fmt.Sprintf("/repos/%s/%s/releases?per_page=%d", owner, repo, limit)
	resp, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var result []Release
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetCommits 获取仓库 Commit 列表
func (c *Client) GetCommits(ctx context.Context, owner, repo string, limit int) ([]Commit, error) {
	path := fmt.Sprintf("/repos/%s/%s/commits?per_page=%d", owner, repo, limit)
	resp, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var result []Commit
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}

// GetSecurityAdvisories 获取安全公告
func (c *Client) GetSecurityAdvisories(ctx context.Context, owner, repo string) ([]SecurityAdvisory, error) {
	path := fmt.Sprintf("/repos/%s/%s/security-advisories", owner, repo)
	resp, err := c.doRequest(ctx, http.MethodGet, path)
	if err != nil {
		return nil, err
	}
	var result []SecurityAdvisory
	if err := decodeResponse(resp, &result); err != nil {
		return nil, err
	}
	return result, nil
}
