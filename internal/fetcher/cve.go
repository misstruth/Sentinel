package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// CVEFetcher CVE 官方数据抓取器
type CVEFetcher struct {
	httpClient *http.Client
}

// NewCVEFetcher 创建 CVE 抓取器
func NewCVEFetcher() *CVEFetcher {
	return &CVEFetcher{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

// Type 返回抓取器类型
func (f *CVEFetcher) Type() model.SourceType {
	return model.SourceTypeCVE
}

type cveResponse struct {
	CVEItems []cveItem `json:"cveItems"`
}

type cveItem struct {
	CVEID       string `json:"cveId"`
	Description string `json:"description"`
	Published   string `json:"published"`
}

// Fetch 抓取 CVE 数据
func (f *CVEFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sub.SourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("CVE API 错误: %d", resp.StatusCode)
	}

	var cveResp cveResponse
	if err := json.NewDecoder(resp.Body).Decode(&cveResp); err != nil {
		return nil, fmt.Errorf("解析响应失败: %w", err)
	}

	return f.parseItems(sub, cveResp.CVEItems), nil
}

func (f *CVEFetcher) parseItems(sub *model.Subscription, items []cveItem) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, item := range items {
		pubTime, _ := time.Parse(time.RFC3339, item.Published)
		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          item.CVEID,
			Description:    item.Description,
			SourceURL:      "https://cve.mitre.org/cgi-bin/cvename.cgi?name=" + item.CVEID,
			Severity:       model.SeverityMedium,
			CVEID:          item.CVEID,
			EventTime:      pubTime,
			UniqueHash:     generateHash("cve", item.CVEID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}
