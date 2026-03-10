package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// AttackActivityFetcher 攻击活动抓取器
type AttackActivityFetcher struct {
	httpClient *http.Client
}

func NewAttackActivityFetcher() *AttackActivityFetcher {
	return &AttackActivityFetcher{
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
}

func (f *AttackActivityFetcher) Type() model.SourceType {
	return model.SourceTypeAttackActivity
}

type attackResponse struct {
	Objects []attackObject `json:"objects"`
}

type attackObject struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Created     string `json:"created"`
}

func (f *AttackActivityFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
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

	var attackResp attackResponse
	if err := json.NewDecoder(resp.Body).Decode(&attackResp); err != nil {
		return nil, fmt.Errorf("解析失败: %w", err)
	}

	return f.parseItems(sub, attackResp.Objects), nil
}

func (f *AttackActivityFetcher) parseItems(sub *model.Subscription, items []attackObject) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, item := range items {
		eventTime, _ := time.Parse(time.RFC3339, item.Created)
		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          item.Name,
			Description:    item.Description,
			SourceURL:      sub.SourceURL,
			Severity:       model.SeverityHigh,
			EventTime:      eventTime,
			UniqueHash:     generateHash("attack", item.ID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}
