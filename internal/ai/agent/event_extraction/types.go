package event_extraction

// RawEventInput 原始事件输入
type RawEventInput struct {
	Source      string `json:"source"`       // 来源: rss, github, webhook
	RawContent  string `json:"raw_content"`  // 原始内容
	SourceURL   string `json:"source_url"`   // 来源URL
	FetchedAt   string `json:"fetched_at"`   // 抓取时间
}

// ExtractedEvent 提取后的结构化事件
type ExtractedEvent struct {
	Title       string   `json:"title"`
	Description string   `json:"description"`
	Severity    string   `json:"severity"`    // critical, high, medium, low, info
	EventType   string   `json:"event_type"`  // vulnerability, attack, advisory, threat
	CVEIDs      []string `json:"cve_ids"`     // 关联的CVE
	Tags        []string `json:"tags"`        // 标签
	AffectedProducts []string `json:"affected_products"`
	Source      string   `json:"source"`
	SourceURL   string   `json:"source_url"`
}

// ExtractionResult 提取结果
type ExtractionResult struct {
	Events  []*ExtractedEvent `json:"events"`
	Success bool              `json:"success"`
	Error   string            `json:"error,omitempty"`
}
