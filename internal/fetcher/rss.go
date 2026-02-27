package fetcher

import (
	"context"
	"encoding/xml"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// RSSFetcher RSS 数据抓取器
type RSSFetcher struct {
	httpClient *http.Client
}

// NewRSSFetcher 创建 RSS 抓取器
func NewRSSFetcher() *RSSFetcher {
	return &RSSFetcher{
		httpClient: &http.Client{Timeout: 30 * time.Second},
	}
}

// Type 返回抓取器类型
func (f *RSSFetcher) Type() model.SourceType {
	return model.SourceTypeRSS
}

// RSS Feed 结构
type rssFeed struct {
	XMLName xml.Name   `xml:"rss"`
	Channel rssChannel `xml:"channel"`
}

type rssChannel struct {
	Title string    `xml:"title"`
	Items []rssItem `xml:"item"`
}

type rssItem struct {
	Title       string `xml:"title"`
	Link        string `xml:"link"`
	Description string `xml:"description"`
	PubDate     string `xml:"pubDate"`
	GUID        string `xml:"guid"`
	Author      string `xml:"author"`
	Source      string `xml:"source"`
}

// Fetch 抓取 RSS 数据
func (f *RSSFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, sub.SourceURL, nil)
	if err != nil {
		return nil, err
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("RSS fetch error: %d", resp.StatusCode)
	}

	var feed rssFeed
	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		return nil, err
	}

	return f.parseItems(sub, feed.Channel.Items), nil
}

// parseItems 解析 RSS 条目
func (f *RSSFetcher) parseItems(sub *model.Subscription, items []rssItem) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, item := range items {
		eventTime := parseRSSDate(item.PubDate)

		// 优先使用link，否则用guid
		sourceURL := item.Link
		if sourceURL == "" {
			sourceURL = item.GUID
		}

		guid := item.GUID
		if guid == "" {
			guid = item.Link
		}

		// 描述为空时用标题
		desc := item.Description
		if desc == "" {
			desc = item.Title
		}

		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          item.Title,
			Description:    desc,
			SourceURL:      sourceURL,
			Severity:       model.SeverityInfo,
			EventTime:      eventTime,
			UniqueHash:     generateHash("rss", guid),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}

// parseRSSDate 解析 RSS 日期
func parseRSSDate(s string) time.Time {
	formats := []string{
		time.RFC1123Z,
		time.RFC1123,
		time.RFC822Z,
		time.RFC822,
		"Mon, 02 Jan 06 15:04:05 -0700", // 2位年份
		"Mon, 02 Jan 06 15:04:05 MST",
		"2006-01-02T15:04:05Z",
		"2006-01-02 15:04:05",
		"2006-01-02",
	}
	for _, f := range formats {
		if t, err := time.Parse(f, s); err == nil {
			return t
		}
	}
	return time.Now()
}
