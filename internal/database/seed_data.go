package database

import "SuperBizAgent/internal/model"

// getDefaultSources 返回初次部署的默认订阅源列表
func getDefaultSources() []defaultSource {
	return []defaultSource{
		// ========== 漏洞情报 ==========
		{
			Name:        "NVD 漏洞数据库",
			Type:        model.SourceTypeNVD,
			URL:         "https://services.nvd.nist.gov/rest/json/cves/2.0",
			Cron:        "0 */4 * * *",
			Description: "NIST国家漏洞数据库，提供CVE详细信息和CVSS评分",
			Enabled:     true,
		},
		{
			Name:        "CISA KEV 已知利用漏洞",
			Type:        model.SourceTypeVulnerability,
			URL:         "https://www.cisa.gov/sites/default/files/feeds/known_exploited_vulnerabilities.json",
			Cron:        "0 */6 * * *",
			Description: "CISA已知被积极利用的漏洞目录，高优先级修复",
			Enabled:     true,
		},
		{
			Name:        "CVE.org 官方漏洞库",
			Type:        model.SourceTypeCVE,
			URL:         "https://cveawg.mitre.org/api/cve",
			Cron:        "0 */4 * * *",
			Description: "CVE官方漏洞数据库",
			Enabled:     true,
		},
		{
			Name:        "CVEfeed 高危漏洞订阅",
			Type:        model.SourceTypeRSS,
			URL:         "https://cvefeed.io/rssfeed/severity/high.xml",
			Cron:        "0 */2 * * *",
			Description: "高危和严重漏洞RSS订阅",
			Enabled:     true,
		},

		// ========== 威胁情报 ==========
		{
			Name:        "AlienVault OTX 威胁情报",
			Type:        model.SourceTypeThreatIntel,
			URL:         "https://otx.alienvault.com/api/v1/pulses/subscribed",
			Cron:        "0 */3 * * *",
			Description: "开放威胁情报交换平台，提供IoC指标",
			Enabled:     true,
		},
		{
			Name:        "MISP 威胁情报共享",
			Type:        model.SourceTypeThreatIntel,
			URL:         "https://your-misp-instance/events/restSearch",
			Cron:        "0 */4 * * *",
			Description: "恶意软件信息共享平台，需自建或接入实例",
			Enabled:     false,
		},

		// ========== 厂商安全公告 ==========
		{
			Name:        "Microsoft 安全更新指南",
			Type:        model.SourceTypeVendorAdvisory,
			URL:         "https://api.msrc.microsoft.com/cvrf/v2.0/updates",
			Cron:        "0 8 * * *",
			Description: "微软安全响应中心安全更新",
			Enabled:     true,
		},
		{
			Name:        "Adobe 安全公告",
			Type:        model.SourceTypeRSS,
			URL:         "https://helpx.adobe.com/security/rss/security-bulletins.xml",
			Cron:        "0 9 * * *",
			Description: "Adobe产品安全公告",
			Enabled:     true,
		},
		{
			Name:        "Oracle 关键补丁更新",
			Type:        model.SourceTypeRSS,
			URL:         "https://www.oracle.com/security-alerts/cpujan2026.html",
			Cron:        "0 10 * * *",
			Description: "Oracle季度关键补丁更新",
			Enabled:     true,
		},

		// ========== GitHub 安全公告 ==========
		{
			Name:        "GitHub 安全公告数据库",
			Type:        model.SourceTypeGitHubRepo,
			URL:         "https://api.github.com/advisories",
			Cron:        "0 */4 * * *",
			Description: "GitHub全球安全公告数据库",
			Enabled:     true,
		},

		// ========== 攻击活动追踪 ==========
		{
			Name:        "MITRE ATT&CK 框架",
			Type:        model.SourceTypeAttackActivity,
			URL:         "https://raw.githubusercontent.com/mitre/cti/master/enterprise-attack/enterprise-attack.json",
			Cron:        "0 0 * * 1",
			Description: "MITRE ATT&CK企业攻击技术框架",
			Enabled:     true,
		},

		// ========== 安全新闻 ==========
		{
			Name:        "The Hacker News",
			Type:        model.SourceTypeRSS,
			URL:         "https://feeds.feedburner.com/TheHackersNews",
			Cron:        "0 */2 * * *",
			Description: "知名安全新闻网站",
			Enabled:     true,
		},
		{
			Name:        "BleepingComputer 安全新闻",
			Type:        model.SourceTypeRSS,
			URL:         "https://www.bleepingcomputer.com/feed/",
			Cron:        "0 */2 * * *",
			Description: "安全新闻和技术支持网站",
			Enabled:     true,
		},
		{
			Name:        "Krebs on Security",
			Type:        model.SourceTypeRSS,
			URL:         "https://krebsonsecurity.com/feed/",
			Cron:        "0 */3 * * *",
			Description: "知名安全研究员Brian Krebs的博客",
			Enabled:     true,
		},
		{
			Name:        "SecurityWeek",
			Type:        model.SourceTypeRSS,
			URL:         "https://www.securityweek.com/feed/",
			Cron:        "0 */3 * * *",
			Description: "企业安全新闻",
			Enabled:     true,
		},

		// ========== 国内安全数据源 ==========
		{
			Name:        "CNVD 漏洞库",
			Type:        model.SourceTypeVulnerability,
			URL:         "https://www.cnvd.org.cn",
			Cron:        "0 */6 * * *",
			Description: "国家信息安全漏洞共享平台",
			Enabled:     true,
		},
		{
			Name:        "CNNVD 漏洞库",
			Type:        model.SourceTypeVulnerability,
			URL:         "https://www.cnnvd.org.cn",
			Cron:        "0 */6 * * *",
			Description: "国家信息安全漏洞库",
			Enabled:     true,
		},
		{
			Name:        "安全客",
			Type:        model.SourceTypeRSS,
			URL:         "https://api.anquanke.com/data/v1/rss",
			Cron:        "0 */2 * * *",
			Description: "国内安全资讯平台",
			Enabled:     true,
		},
		{
			Name:        "FreeBuf",
			Type:        model.SourceTypeRSS,
			URL:         "https://www.freebuf.com/feed",
			Cron:        "0 */2 * * *",
			Description: "国内安全媒体",
			Enabled:     true,
		},
	}
}
