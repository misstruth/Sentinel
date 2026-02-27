package event_extraction

const ExtractionPrompt = `你是安全事件提取专家。从原始内容中提取结构化安全事件信息。

原始内容:
{{.RawContent}}

来源: {{.Source}}
URL: {{.SourceURL}}

请提取以下信息并以JSON格式返回:
{
  "events": [{
    "title": "事件标题",
    "description": "详细描述",
    "severity": "critical/high/medium/low/info",
    "event_type": "vulnerability/attack/advisory/threat",
    "cve_ids": ["CVE-xxxx-xxxx"],
    "tags": ["标签1", "标签2"],
    "affected_products": ["产品1"]
  }]
}

注意:
1. severity根据影响范围和危害程度判断
2. 提取所有CVE编号
3. 如果无法提取，返回空events数组`
