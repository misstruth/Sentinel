package fetcher

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"SuperBizAgent/internal/model"
)

// NVDFetcher NVD 漏洞数据抓取器
type NVDFetcher struct {
	httpClient *http.Client
	apiKey     string
	baseURL    string
}

// NewNVDFetcher 创建 NVD 抓取器
func NewNVDFetcher(apiKey string) *NVDFetcher {
	return &NVDFetcher{
		httpClient: &http.Client{Timeout: 60 * time.Second},
		apiKey:     apiKey,
		baseURL:    "https://services.nvd.nist.gov/rest/json/cves/2.0",
	}
}

// Type 返回抓取器类型
func (f *NVDFetcher) Type() model.SourceType {
	return model.SourceTypeNVD
}

// NVD API 响应结构
type nvdResponse struct {
	Vulnerabilities []nvdVuln `json:"vulnerabilities"`
}

type nvdVuln struct {
	CVE nvdCVE `json:"cve"`
}

type nvdCVE struct {
	ID           string           `json:"id"`
	Published    string           `json:"published"`
	Description  []nvdDescription `json:"descriptions"`
	Metrics      nvdMetrics       `json:"metrics"`
}

type nvdDescription struct {
	Lang  string `json:"lang"`
	Value string `json:"value"`
}

type nvdMetrics struct {
	CvssMetricV31 []nvdCVSS `json:"cvssMetricV31"`
}

type nvdCVSS struct {
	CvssData nvdCVSSData `json:"cvssData"`
}

type nvdCVSSData struct {
	BaseScore float64 `json:"baseScore"`
}

// Fetch 抓取 NVD 漏洞数据
func (f *NVDFetcher) Fetch(ctx context.Context, sub *model.Subscription) ([]model.SecurityEvent, error) {
	url := f.buildURL(sub)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	if err != nil {
		return nil, err
	}

	if f.apiKey != "" {
		req.Header.Set("apiKey", f.apiKey)
	}

	resp, err := f.httpClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("NVD API error: %d", resp.StatusCode)
	}

	var nvdResp nvdResponse
	if err := json.NewDecoder(resp.Body).Decode(&nvdResp); err != nil {
		return nil, err
	}

	return f.parseVulns(sub, nvdResp.Vulnerabilities), nil
}

// buildURL 构建请求 URL
func (f *NVDFetcher) buildURL(sub *model.Subscription) string {
	lastDay := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000")
	return fmt.Sprintf("%s?pubStartDate=%s", f.baseURL, lastDay)
}

// parseVulns 解析漏洞数据
func (f *NVDFetcher) parseVulns(sub *model.Subscription, vulns []nvdVuln) []model.SecurityEvent {
	var events []model.SecurityEvent
	for _, v := range vulns {
		desc := f.getDescription(v.CVE.Description)
		score := f.getCVSSScore(v.CVE.Metrics)
		severity := f.scoreSeverity(score)
		pubTime, _ := time.Parse(time.RFC3339, v.CVE.Published)

		event := model.SecurityEvent{
			SubscriptionID: sub.ID,
			Title:          v.CVE.ID,
			Description:    desc,
			SourceURL:      "https://nvd.nist.gov/vuln/detail/" + v.CVE.ID,
			Severity:       severity,
			CVEID:          v.CVE.ID,
			CVSSScore:      score,
			EventTime:      pubTime,
			UniqueHash:     generateHash("nvd", v.CVE.ID),
			CreatedAt:      time.Now(),
		}
		events = append(events, event)
	}
	return events
}

// getDescription 获取英文描述
func (f *NVDFetcher) getDescription(descs []nvdDescription) string {
	for _, d := range descs {
		if d.Lang == "en" {
			return d.Value
		}
	}
	if len(descs) > 0 {
		return descs[0].Value
	}
	return ""
}

// getCVSSScore 获取 CVSS 分数
func (f *NVDFetcher) getCVSSScore(m nvdMetrics) float64 {
	if len(m.CvssMetricV31) > 0 {
		return m.CvssMetricV31[0].CvssData.BaseScore
	}
	return 0
}

// scoreSeverity 根据 CVSS 分数判断严重级别
func (f *NVDFetcher) scoreSeverity(score float64) model.SeverityLevel {
	switch {
	case score >= 9.0:
		return model.SeverityCritical
	case score >= 7.0:
		return model.SeverityHigh
	case score >= 4.0:
		return model.SeverityMedium
	case score > 0:
		return model.SeverityLow
	default:
		return model.SeverityInfo
	}
}
