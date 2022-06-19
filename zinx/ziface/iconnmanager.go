package ziface

type IConnManager interface {
	// Add 添加连接
	Add(conn IConnection)
	// Remove 删除连接
	Remove(conn IConnection)
	// Get 根据 connId 获取一个连接
	Get(connID uint32) (IConnection, error)
	// Len 获取连接的个数
	Len() int
	// ClearConn 清除所有的连接
	ClearConn()
}
