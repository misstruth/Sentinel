package main

import (
	_ "SuperBizAgent/internal/ai/skills/builtin"
	"SuperBizAgent/internal/controller/chat"
	"SuperBizAgent/internal/controller/event"
	"SuperBizAgent/internal/controller/report"
	"SuperBizAgent/internal/controller/skill"
	"SuperBizAgent/internal/controller/subscription"
	"SuperBizAgent/internal/database"
	"SuperBizAgent/internal/service"
	"SuperBizAgent/utility/common"
	"SuperBizAgent/utility/middleware"

	"github.com/gogf/gf/v2/frame/g"
	"github.com/gogf/gf/v2/net/ghttp"
	"github.com/gogf/gf/v2/os/gctx"
)

func main() {
	ctx := gctx.New()
	fileDir, err := g.Cfg().Get(ctx, "file_dir")
	if err != nil {
		panic(err)
	}
	common.FileDir = fileDir.String()

	// 初始化数据库
	if err := database.Init(); err != nil {
		panic(err)
	}

	// 启动抓取服务
	service.InitFetcher()

	s := g.Server()
	// SSE路由（不经过中间件）
	s.BindHandler("POST:/api/event/pipeline/stream", event.New().PipelineStream)

	s.Group("/api", func(group *ghttp.RouterGroup) {
		group.Middleware(middleware.CORSMiddleware)
		group.Middleware(middleware.ResponseMiddleware)
		group.Bind(chat.NewV1())
		group.Bind(subscription.NewV1())
		group.Bind(report.New())
		group.Bind(report.NewTemplateController())
		group.Bind(event.New())
		group.Bind(skill.NewV1())
	})
	s.SetPort(6872)
	s.Run()
}
