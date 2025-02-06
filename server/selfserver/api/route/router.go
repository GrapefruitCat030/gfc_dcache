package route

import (
	"github.com/GrapefruitCat030/gfc_dcache/pkg/cluster"
	"github.com/GrapefruitCat030/gfc_dcache/pkg/protocol"
)

// HandlerFunc 定义处理器函数类型
type HandlerFunc func(*protocol.Request) *protocol.Response

// Router 路由器结构体
type Router struct {
	handlers map[protocol.OpType]HandlerFunc
}

// NewRouter 创建新的路由器
func NewRouter() *Router {
	return &Router{
		handlers: make(map[protocol.OpType]HandlerFunc),
	}
}

// Register 注册路由处理器
func (r *Router) Register(op protocol.OpType, handler HandlerFunc) {
	r.handlers[op] = handler
}

// Dispatch 分发请求到对应的处理器
func (r *Router) Dispatch(req *protocol.Request) *protocol.Response {
	// 判断当前节点是否负责处理该请求
	if addr, ok := cluster.GlobalNode().ShouldProcess(string(req.Key)); !ok {
		return &protocol.Response{IsError: true, Data: []byte("redirect to " + addr)}
	}
	if handler, ok := r.handlers[req.Op]; ok {
		return handler(req)
	}
	return &protocol.Response{IsError: true, Data: []byte("unknown operation")}
}
