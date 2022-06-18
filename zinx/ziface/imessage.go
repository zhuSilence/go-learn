package ziface

type IMessage interface {
	// GetMsgId 获取消息的 id
	GetMsgId() uint32
	// GetMsgLen 获取消息的长度
	GetDataLen() uint32
	// GetData 获取消息的内容
	GetData() []byte
	// SetMsgId 设置消息的 id
	SetMsgId(uint32)
	// SetMsgLen 设置消息的长度
	SetMsgLen(uint32)
	// SetData 设置消息的内容
	SetData([]byte)
}
