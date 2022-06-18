package ziface

// IDataPack 封包拆包模块，直接面向 TCP 连接的数据流，用于处理 TCP 的粘包问题
type IDataPack interface {
	// GetHeadLen 获取包头的长度
	GetHeadLen() uint32
	// Pack 封包方法
	Pack(msg IMessage) ([]byte, error)
	// Unpack 拆包方法
	Unpack([]byte) (IMessage, error)
}
