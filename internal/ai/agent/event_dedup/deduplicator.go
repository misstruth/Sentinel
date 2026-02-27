package event_dedup

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"strings"

	"SuperBizAgent/internal/ai/agent/event_extraction"
)

type Deduplicator struct{}

func NewDeduplicator() *Deduplicator {
	return &Deduplicator{}
}

func (d *Deduplicator) Dedup(ctx context.Context, input *DedupInput) (*DedupResult, error) {
	seen := make(map[string]bool)
	var unique []*event_extraction.ExtractedEvent

	for _, ev := range input.Events {
		hash := d.computeHash(ev)
		if !seen[hash] {
			seen[hash] = true
			unique = append(unique, ev)
		}
	}

	return &DedupResult{
		UniqueEvents: unique,
		DupCount:     len(input.Events) - len(unique),
	}, nil
}

func (d *Deduplicator) computeHash(ev *event_extraction.ExtractedEvent) string {
	data := strings.ToLower(ev.Title + "|" + strings.Join(ev.CVEIDs, ","))
	hash := md5.Sum([]byte(data))
	return hex.EncodeToString(hash[:])
}
