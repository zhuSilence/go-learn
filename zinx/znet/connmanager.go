package znet

import (
	"errors"
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"sync"
)

type ConnManager struct {
	connections map[uint32]ziface.IConnection // 管理的连接集合
	connLock    sync.RWMutex                  // 保护连接的读写锁
}

// NewConnManager 创建连接管理器
func NewConnManager() *ConnManager {
	return &ConnManager{
		connections: make(map[uint32]ziface.IConnection),
	}
}

func (connMgr *ConnManager) Add(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	connMgr.connections[conn.GetConnId()] = conn

	fmt.Println("connId = ", conn.GetConnId(), " add to ConnManager successful: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Remove(conn ziface.IConnection) {
	// 保护共享资源 map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()
	delete(connMgr.connections, conn.GetConnId())
	fmt.Println("connId = ", conn.GetConnId(), " remove from ConnManager successful: conn num = ", connMgr.Len())
}

func (connMgr *ConnManager) Get(connID uint32) (ziface.IConnection, error) {
	// 保护共享资源 map，加读锁
	connMgr.connLock.RLock()
	defer connMgr.connLock.RUnlock()

	if conn, ok := connMgr.connections[connID]; ok {
		// 找到了
		return conn, nil
	} else {
		return nil, errors.New("connection not find")
	}
}

func (connMgr *ConnManager) Len() int {
	return len(connMgr.connections)
}

func (connMgr *ConnManager) ClearConn() {
	// 保护共享资源 map，加写锁
	connMgr.connLock.Lock()
	defer connMgr.connLock.Unlock()

	for connID, conn := range connMgr.connections {
		// 停止
		conn.Stop()
		// 删除
		delete(connMgr.connections, connID)
	}
	fmt.Println("Clear All connection successful conn num = ", connMgr.Len())
}
