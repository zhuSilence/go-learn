package znet

import (
	"github.com/zhuSilence/go-learn/zinx/ziface"
)

// BaseRouter 实现 router 时，先嵌入这个基类
// 实现空方法，方便子类可以实现单个方法
type BaseRouter struct {
}

func (b *BaseRouter) PreHandle(request ziface.IRequest) {

}

func (b *BaseRouter) Handle(request ziface.IRequest) {

}

func (b *BaseRouter) PostHandle(request ziface.IRequest) {

}
