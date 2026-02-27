package v1

import "github.com/gogf/gf/v2/frame/g"

type SupervisorChatReq struct {
	g.Meta `path:"/chat/v1/supervisor" method:"post"`
	Query  string `json:"query" v:"required"`
}

type SupervisorChatRes struct {
	g.Meta `mime:"text/event-stream"`
}
