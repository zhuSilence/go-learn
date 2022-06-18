package znet

import "github.com/zhuSilence/go-learn/zinx/ziface"

type Request struct {
	// 建立好连接的 conn
	conn ziface.IConnection
	// 客户端请求的数据
	msg ziface.IMessage
}

func (r *Request) GetConnection() ziface.IConnection {
	return r.conn
}

func (r *Request) GetDate() []byte {
	return r.msg.GetData()
}

func (r *Request) GetMsgID() uint32 {
	return r.msg.GetMsgId()
}
