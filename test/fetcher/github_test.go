package fetcher_test

import (
	"testing"
)

func TestParseRepoURL(t *testing.T) {
	tests := []struct {
		url   string
		owner string
		repo  string
	}{
		{"https://github.com/owner/repo", "owner", "repo"},
		{"http://github.com/test/project", "test", "project"},
	}

	for _, tt := range tests {
		t.Run(tt.url, func(t *testing.T) {
			// 测试 URL 解析
		})
	}
}
