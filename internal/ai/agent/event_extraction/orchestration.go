package event_extraction

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildExtractionAgent(ctx context.Context) (compose.Runnable[*RawEventInput, *ExtractionResult], error) {
	g := compose.NewGraph[*RawEventInput, *ExtractionResult]()

	extractor := NewExtractor(ctx)
	extractLambda := compose.InvokableLambda(func(ctx context.Context, input *RawEventInput) (*ExtractionResult, error) {
		return extractor.Extract(input)
	})

	_ = g.AddLambdaNode("Extractor", extractLambda, compose.WithNodeName("EventExtractor"))
	_ = g.AddEdge(compose.START, "Extractor")
	_ = g.AddEdge("Extractor", compose.END)

	return g.Compile(ctx, compose.WithGraphName("EventExtractionAgent"))
}
