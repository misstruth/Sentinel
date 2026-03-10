package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// VendorAdvisoryFetcher 厂商安全公告抓取器
type VendorAdvisoryFetcher struct {
	httpClient *http.Client
}

func NewVendorAdvisoryFetcher() *VendorAdvisoryFetcher {
	return &VendorAdvisoryFetcher{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (f *VendorAdvisoryFetcher) Type() model.SourceType {
	return model.SourceTypeVendorAdvisory
}

type advisoryResponse struct {
	Advisories []advisoryItem `json:"value"`
}

type advisoryItem struct {
	ID          string `json:"ID"`
	Title       string `json:"Title"`
	Description string `json:"Description"`
	Published   string `json:"InitialReleaseDate"`
	Severity    string `json:"Severity"`
}

func (f *VendorAdvisoryFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
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
		return nil, fmt.Errorf("API 错误: %d", resp.StatusCode)
	}

	var advResp advisoryResponse
	if err := json.NewDecoder(resp.Body).Decode(&advResp); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return f.parseItems(sub, advResp.Advisories), nil
}

func (f *VendorAdvisoryFetcher) parseItems(sub *model.Subscription, items []advisoryItem) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, item := range items {
		eventTime, _ := time.Parse(time.RFC3339, item.Published)
		severity := f.parseSeverity(item.Severity)

		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          item.Title,
			Description:    item.Description,
			SourceURL:      sub.SourceURL,
			Severity:       severity,
			EventTime:      eventTime,
			UniqueHash:     generateHash("advisory", item.ID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}

func (f *VendorAdvisoryFetcher) parseSeverity(sev string) model.SeverityLevel {
	switch sev {
	case "Critical":
		return model.SeverityCritical
	case "Important", "High":
		return model.SeverityHigh
	case "Moderate", "Medium":
		return model.SeverityMedium
	case "Low":
		return model.SeverityLow
	default:
		return model.SeverityInfo
	}
}
