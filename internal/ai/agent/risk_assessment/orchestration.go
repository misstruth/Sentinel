package risk_assessment

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

func BuildAssessmentAgent(ctx context.Context) (compose.Runnable[*AssessmentInput, *AssessmentResult], error) {
	g := compose.NewGraph[*AssessmentInput, *AssessmentResult]()

	assessor := NewAssessor(ctx)
	lambda := compose.InvokableLambda(func(ctx context.Context, input *AssessmentInput) (*AssessmentResult, error) {
		return assessor.Assess(input)
	})

	_ = g.AddLambdaNode("Assessor", lambda, compose.WithNodeName("RiskAssessor"))
	_ = g.AddEdge(compose.START, "Assessor")
	_ = g.AddEdge("Assessor", compose.END)

	return g.Compile(ctx, compose.WithGraphName("RiskAssessmentAgent"))
}
