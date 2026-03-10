package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// ThreatIntelFetcher 威胁情报抓取器
type ThreatIntelFetcher struct {
	httpClient *http.Client
}

func NewThreatIntelFetcher() *ThreatIntelFetcher {
	return &ThreatIntelFetcher{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (f *ThreatIntelFetcher) Type() model.SourceType {
	return model.SourceTypeThreatIntel
}

type threatResponse struct {
	Results []threatItem `json:"results"`
}

type threatItem struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Created     string `json:"created"`
	ThreatType  string `json:"threat_type"`
}

func (f *ThreatIntelFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sub.SourceURL, nil)
	if err != nil {
		return nil, fmt.Errorf("创建请求失败: %w", err)
	}

	if sub.AuthToken != "" {
		req.Header.Set("Authorization", "Bearer "+sub.AuthToken)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("请求失败: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("API 错误: %d", resp.StatusCode)
	}

	var threatResp threatResponse
	if err := json.NewDecoder(resp.Body).Decode(&threatResp); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return f.parseItems(sub, threatResp.Results), nil
}

func (f *ThreatIntelFetcher) parseItems(sub *model.Subscription, items []threatItem) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, item := range items {
		eventTime, _ := time.Parse(time.RFC3339, item.Created)
		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          item.Name,
			Description:    item.Description,
			SourceURL:      sub.SourceURL + "/" + item.ID,
			Severity:       model.SeverityMedium,
			EventTime:      eventTime,
			UniqueHash:     generateHash("threat", item.ID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}
