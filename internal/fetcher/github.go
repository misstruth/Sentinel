package fetcher

import (
	"context"
	"crypto/sha256"
	"fmt"
	"strings"
	"time"

	"SuperBizAgent/internal/client/github"
	"SuperBizAgent/internal/model"
)

// GitHubFetcher GitHub 数据抓取器
type GitHubFetcher struct {
	client *github.Client
}

// NewGitHubFetcher 创建 GitHub 抓取器
func NewGitHubFetcher() *GitHubFetcher {
	return &GitHubFetcher{
		client: github.NewClient(""),
	}
}

// Type 返回抓取器类型
func (f *GitHubFetcher) Type() model.SourceType {
	return model.SourceTypeGitHubRepo
}

// parseRepoURL 解析仓库 URL
func parseRepoURL(url string) (owner, repo string, err error) {
	url = strings.TrimPrefix(url, "https://github.com/")
	url = strings.TrimPrefix(url, "http://github.com/")
	parts := strings.Split(url, "/")
	if len(parts) < 2 {
		return "", "", fmt.Errorf("invalid repo URL: %s", url)
	}
	return parts[0], parts[1], nil
}

// generateHash 生成事件唯一哈希
func generateHash(source, id string) string {
	h := sha256.Sum256([]byte(source + ":" + id))
	return fmt.Sprintf("%x", h)
}

// Fetch 抓取 GitHub 仓库事件
func (f *GitHubFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	owner, repo, err := parseRepoURL(sub.SourceURL)
	if err != nil {
		return nil, err
	}

	var events []model.SecurityEvent

	// 抓取 Release
	releases, err := f.fetchReleases(ctx, sub, owner, repo)
	if err == nil {
		events = append(events, releases...)
	}

	// 抓取安全公告
	advisories, err := f.fetchAdvisories(ctx, sub, owner, repo)
	if err == nil {
		events = append(events, advisories...)
	}

	return events, nil
}

// fetchReleases 抓取 Release 事件
func (f *GitHubFetcher) fetchReleases(ctx context.Context, sub *model.Subscription, owner, repo string) ([]model.SecurityEvent, error) {
	releases, err := f.client.GetReleases(ctx, owner, repo, 10)
	if err != nil {
		return nil, err
	}

	var events []model.SecurityEvent
	for _, r := range releases {
		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          fmt.Sprintf("[Release] %s %s", repo, r.TagName),
			Description:    r.Body,
			SourceURL:      r.HTMLURL,
			Severity:       model.SeverityInfo,
			EventTime:      r.PublishedAt,
			UniqueHash:     generateHash("github-release", fmt.Sprintf("%d", r.ID)),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events, nil
}

// mapGitHubSeverity 映射 GitHub 严重级别
func mapGitHubSeverity(s string) model.SeverityLevel {
	switch strings.ToLower(s) {
	case "critical":
		return model.SeverityCritical
	case "high":
		return model.SeverityHigh
	case "medium":
		return model.SeverityMedium
	case "low":
		return model.SeverityLow
	default:
		return model.SeverityInfo
	}
}

// fetchAdvisories 抓取安全公告
func (f *GitHubFetcher) fetchAdvisories(ctx context.Context, sub *model.Subscription, owner, repo string) ([]model.SecurityEvent, error) {
	advisories, err := f.client.GetSecurityAdvisories(ctx, owner, repo)
	if err != nil {
		return nil, err
	}

	var events []model.SecurityEvent
	for _, a := range advisories {
		severity := mapGitHubSeverity(a.Severity)
		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          fmt.Sprintf("[Security] %s", a.Summary),
			Description:    a.Description,
			SourceURL:      a.HTMLURL,
			Severity:       severity,
			CVEID:          a.CVEID,
			EventTime:      a.PublishedAt,
			UniqueHash:     generateHash("github-advisory", a.GHSAID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events, nil
}
