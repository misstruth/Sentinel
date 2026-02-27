package supervisor

import (
	"context"

	"github.com/cloudwego/eino/compose"
)

const (
	NodeRouter   = "Router"
	NodeExecutor = "Executor"
)

// BuildSupervisorGraph 构建 Supervisor Graph
// 流程: START -> Router -> [Branch] -> Executor -> END
func BuildSupervisorGraph(ctx context.Context) (compose.Runnable[*SupervisorInput, *SupervisorOutput], error) {
	g := compose.NewGraph[*SupervisorInput, *SupervisorOutput]()

	// 添加路由节点 - 使用 LLM 判断应该路由到哪个 Agent
	routerLambda, err := newRouterLambda(ctx)
	if err != nil {
		return nil, err
	}
	if err := g.AddLambdaNode(NodeRouter, routerLambda, compose.WithNodeName("IntentRouter")); err != nil {
		return nil, err
	}

	// 添加执行节点 - 执行对应的 SubAgent
	executorLambda := newExecutorLambda()
	if err := g.AddLambdaNode(NodeExecutor, executorLambda, compose.WithNodeName("AgentExecutor")); err != nil {
		return nil, err
	}

	// 连接边
	if err := g.AddEdge(compose.START, NodeRouter); err != nil {
		return nil, err
	}
	if err := g.AddEdge(NodeRouter, NodeExecutor); err != nil {
		return nil, err
	}
	if err := g.AddEdge(NodeExecutor, compose.END); err != nil {
		return nil, err
	}

	return g.Compile(ctx, compose.WithGraphName("SupervisorGraph"))
}
