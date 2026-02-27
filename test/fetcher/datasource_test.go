package fetcher_test

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"testing"
	"time"
)

// 测试所有数据源的连通性和数据获取

// TestNVDDataSource 测试 NVD 漏洞数据库
func TestNVDDataSource(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)
	defer cancel()

	// 获取最近7天的漏洞
	lastWeek := time.Now().AddDate(0, 0, -7).Format("2006-01-02T15:04:05.000")
	url := fmt.Sprintf("https://services.nvd.nist.gov/rest/json/cves/2.0?pubStartDate=%s&resultsPerPage=5", lastWeek)

	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 60 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("NVD 请求失败: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("NVD 返回状态码: %d", resp.StatusCode)
	}

	var result struct {
		TotalResults    int `json:"totalResults"`
		Vulnerabilities []struct {
			CVE struct {
				ID          string `json:"id"`
				Published   string `json:"published"`
				Description []struct {
					Lang  string `json:"lang"`
					Value string `json:"value"`
				} `json:"descriptions"`
				Metrics struct {
					CvssMetricV31 []struct {
						CvssData struct {
							BaseScore float64 `json:"baseScore"`
						} `json:"cvssData"`
					} `json:"cvssMetricV31"`
				} `json:"metrics"`
			} `json:"cve"`
		} `json:"vulnerabilities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("NVD 解析失败: %v", err)
	}

	t.Logf("✅ NVD 数据源测试成功")
	t.Logf("   总漏洞数: %d", result.TotalResults)
	t.Logf("   返回条数: %d", len(result.Vulnerabilities))

	for i, v := range result.Vulnerabilities {
		if i >= 3 {
			break
		}
		desc := ""
		for _, d := range v.CVE.Description {
			if d.Lang == "en" {
				desc = d.Value
				break
			}
		}
		if len(desc) > 100 {
			desc = desc[:100] + "..."
		}
		score := 0.0
		if len(v.CVE.Metrics.CvssMetricV31) > 0 {
			score = v.CVE.Metrics.CvssMetricV31[0].CvssData.BaseScore
		}
		t.Logf("   [%d] %s (CVSS: %.1f)", i+1, v.CVE.ID, score)
		t.Logf("       %s", desc)
	}
}

// TestCISAKEVDataSource 测试 CISA KEV 已知被利用漏洞
func TestCISAKEVDataSource(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("CISA KEV 请求失败: %v", err)
	}
	defer resp.Body.Close()

	var result struct {
		Title           string `json:"title"`
		CatalogVersion  string `json:"catalogVersion"`
		Count           int    `json:"count"`
		Vulnerabilities []struct {
			CVEID            string `json:"cveID"`
			VendorProject    string `json:"vendorProject"`
			Product          string `json:"product"`
			VulnerabilityName string `json:"vulnerabilityName"`
			DateAdded        string `json:"dateAdded"`
			ShortDescription string `json:"shortDescription"`
			DueDate          string `json:"dueDate"`
		} `json:"vulnerabilities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("CISA KEV 解析失败: %v", err)
	}

	t.Logf("✅ CISA KEV 数据源测试成功")
	t.Logf("   目录版本: %s", result.CatalogVersion)
	t.Logf("   总漏洞数: %d", result.Count)

	for i, v := range result.Vulnerabilities {
		if i >= 3 {
			break
		}
		desc := v.ShortDescription
		if len(desc) > 80 {
			desc = desc[:80] + "..."
		}
		t.Logf("   [%d] %s - %s/%s", i+1, v.CVEID, v.VendorProject, v.Product)
		t.Logf("       添加日期: %s, 修复截止: %s", v.DateAdded, v.DueDate)
		t.Logf("       %s", desc)
	}
}

// TestCVEFeedRSS 测试 CVEfeed RSS
func TestCVEFeedRSS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://cvefeed.io/rssfeed/severity/high.xml"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("CVEfeed RSS 请求失败: %v", err)
	}
	defer resp.Body.Close()

	var feed struct {
		Channel struct {
			Title string `xml:"title"`
			Items []struct {
				Title       string `xml:"title"`
				Link        string `xml:"link"`
				Description string `xml:"description"`
				PubDate     string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		t.Fatalf("CVEfeed RSS 解析失败: %v", err)
	}

	t.Logf("✅ CVEfeed RSS 数据源测试成功")
	t.Logf("   频道: %s", feed.Channel.Title)
	t.Logf("   条目数: %d", len(feed.Channel.Items))

	for i, item := range feed.Channel.Items {
		if i >= 3 {
			break
		}
		t.Logf("   [%d] %s", i+1, item.Title)
		t.Logf("       发布时间: %s", item.PubDate)
	}
}

// TestGitHubAdvisory 测试 GitHub Advisory Database
func TestGitHubAdvisory(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://api.github.com/advisories?per_page=5"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	req.Header.Set("Accept", "application/vnd.github+json")
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("GitHub Advisory 请求失败: %v", err)
	}
	defer resp.Body.Close()

	var advisories []struct {
		GHSAID      string `json:"ghsa_id"`
		CVEID       string `json:"cve_id"`
		Summary     string `json:"summary"`
		Severity    string `json:"severity"`
		PublishedAt string `json:"published_at"`
		Vulnerabilities []struct {
			Package struct {
				Ecosystem string `json:"ecosystem"`
				Name      string `json:"name"`
			} `json:"package"`
		} `json:"vulnerabilities"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&advisories); err != nil {
		t.Fatalf("GitHub Advisory 解析失败: %v", err)
	}

	t.Logf("✅ GitHub Advisory 数据源测试成功")
	t.Logf("   返回条数: %d", len(advisories))

	for i, a := range advisories {
		if i >= 3 {
			break
		}
		pkg := ""
		if len(a.Vulnerabilities) > 0 {
			pkg = fmt.Sprintf("%s/%s", a.Vulnerabilities[0].Package.Ecosystem, a.Vulnerabilities[0].Package.Name)
		}
		t.Logf("   [%d] %s (%s) - %s", i+1, a.GHSAID, a.Severity, a.CVEID)
		t.Logf("       包: %s", pkg)
		summary := a.Summary
		if len(summary) > 60 {
			summary = summary[:60] + "..."
		}
		t.Logf("       %s", summary)
	}
}

// TestAnquankeRSS 测试安全客 RSS
func TestAnquankeRSS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://api.anquanke.com/data/v1/rss"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("安全客 RSS 请求失败: %v", err)
	}
	defer resp.Body.Close()

	var feed struct {
		Channel struct {
			Title string `xml:"title"`
			Items []struct {
				Title   string `xml:"title"`
				Link    string `xml:"guid"`
				PubDate string `xml:"pubDate"`
				Source  string `xml:"source"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.NewDecoder(resp.Body).Decode(&feed); err != nil {
		t.Fatalf("安全客 RSS 解析失败: %v", err)
	}

	t.Logf("✅ 安全客 RSS 数据源测试成功")
	t.Logf("   频道: %s", feed.Channel.Title)
	t.Logf("   条目数: %d", len(feed.Channel.Items))

	for i, item := range feed.Channel.Items {
		if i >= 5 {
			break
		}
		t.Logf("   [%d] %s", i+1, item.Title)
		t.Logf("       来源: %s, 时间: %s", item.Source, item.PubDate)
	}
}

// TestFreeBufRSS 测试 FreeBuf RSS
func TestFreeBufRSS(t *testing.T) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()

	url := "https://www.freebuf.com/feed"
	req, _ := http.NewRequestWithContext(ctx, http.MethodGet, url, nil)
	client := &http.Client{Timeout: 30 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("FreeBuf RSS 请求失败: %v", err)
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(resp.Body)

	var feed struct {
		Channel struct {
			Title string `xml:"title"`
			Items []struct {
				Title    string `xml:"title"`
				Link     string `xml:"link"`
				Category string `xml:"category"`
				PubDate  string `xml:"pubDate"`
			} `xml:"item"`
		} `xml:"channel"`
	}

	if err := xml.Unmarshal(body, &feed); err != nil {
		t.Fatalf("FreeBuf RSS 解析失败: %v", err)
	}

	t.Logf("✅ FreeBuf RSS 数据源测试成功")
	t.Logf("   频道: %s", feed.Channel.Title)
	t.Logf("   条目数: %d", len(feed.Channel.Items))

	for i, item := range feed.Channel.Items {
		if i >= 5 {
			break
		}
		t.Logf("   [%d] [%s] %s", i+1, item.Category, item.Title)
		t.Logf("       时间: %s", item.PubDate)
	}
}