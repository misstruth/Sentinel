package event_dedup

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildDedupAgent(ctx context.Context) (compose.Runnable[*DedupInput, *DedupResult], error) {
	g := compose.NewGraph[*DedupInput, *DedupResult]()

	dedup := NewDeduplicator()
	lambda := compose.InvokableLambda(func(ctx context.Context, input *DedupInput) (*DedupResult, error) {
		return dedup.Dedup(ctx, input)
	})

	_ = g.AddLambdaNode("Dedup", lambda, compose.WithNodeName("EventDedup"))
	_ = g.AddEdge(compose.START, "Dedup")
	_ = g.AddEdge("Dedup", compose.END)

	return g.Compile(ctx, compose.WithGraphName("EventDedupAgent"))
}
