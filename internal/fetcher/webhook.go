package fetcher

import (
	"context"
	"fmt"

	"SuperBizAgent/internal/model"
)

// WebHookFetcher WebHook 抓取器（被动接收）
type WebHookFetcher struct{}

func NewWebHookFetcher() *WebHookFetcher {
	return &WebHookFetcher{}
}

func (f *WebHookFetcher) Type() model.SourceType {
	return model.SourceTypeWebHook
}

// Fetch WebHook 不主动抓取，返回空
func (f *WebHookFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	return nil, fmt.Errorf("webhook 是被动接收模式，不支持主动抓取")
}
