package event_pipeline

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildEventPipeline(ctx context.Context) (compose.Runnable[*PipelineInput, *PipelineResult], error) {
	g := compose.NewGraph[*PipelineInput, *PipelineResult]()

	pipeline := NewPipeline(ctx)
	lambda := compose.InvokableLambda(func(ctx context.Context, input *PipelineInput) (*PipelineResult, error) {
		return pipeline.Process(input)
	})

	_ = g.AddLambdaNode("Pipeline", lambda, compose.WithNodeName("EventPipeline"))
	_ = g.AddEdge(compose.START, "Pipeline")
	_ = g.AddEdge("Pipeline", compose.END)

	return g.Compile(ctx, compose.WithGraphName("EventProcessingPipeline"))
}
