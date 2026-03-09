package retriever

import (
	"context"
	"time"
)

// FallbackRetriever 降级检索器
type FallbackRetriever struct {
	primary   Retriever
	secondary Retriever
	timeout   time.Duration
}

func NewFallbackRetriever(primary, secondary Retriever, timeout time.Duration) *FallbackRetriever {
	return &FallbackRetriever{
		primary:   primary,
		secondary: secondary,
		timeout:   timeout,
	}
}

func (f *FallbackRetriever) Retrieve(ctx context.Context, query string) ([]Document, error) {
	ctx, cancel := context.WithTimeout(ctx, f.timeout)
	defer cancel()

	resultCh := make(chan []Document, 1)
	errCh := make(chan error, 1)

	go func() {
		docs, err := f.primary.Retrieve(ctx, query)
		if err != nil {
			errCh <- err
			return
		}
		resultCh <- docs
	}()

	select {
	case docs := <-resultCh:
		return docs, nil
	case <-ctx.Done():
		// 超时或失败，降级到备用检索器
		return f.secondary.Retrieve(context.Background(), query)
	case <-errCh:
		return f.secondary.Retrieve(context.Background(), query)
	}
}
