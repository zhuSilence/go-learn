package ziface

type IRequest interface {
	GetConnection() IConnection

	// GetDate 获取消息的内容
	GetDate() []byte

	// GetMsgID 获取消息的 id
	GetMsgID() uint32
}
