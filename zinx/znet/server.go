package znet

import (
	"encoding/json"
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/utils"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"net"
)

type Server struct {
	// server name
	Name string
	// server ip version
	IPVersion string
	// 	server ip
	IP string
	//server port
	Port int
	// 当前 Server 的消息管理模块
	MsgHandler ziface.IMsgHandler
}

func (s *Server) Start() {
	marshal, err := json.Marshal(utils.GlobalObject)
	if err != nil {
		fmt.Println("json Marshal err", err)
	}
	fmt.Printf("[zinx config]: %s\n", marshal)
	fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
		//0 开启消息处理池
		s.MsgHandler.StartWorkerPool()

		//1. 获取一个 TCP 的 addr
		addr, err := net.ResolveTCPAddr(s.IPVersion, fmt.Sprintf("%s:%d", s.IP, s.Port))
		if err != nil {
			fmt.Println("resolve tcp addr error:", err)
			return
		}
		//2 监听服务器地址
		listener, err := net.ListenTCP(s.IPVersion, addr)
		if err != nil {
			fmt.Println("listen", s.IPVersion, " err ", err)
			return
		}

		fmt.Println("start Zinx server succ, ", s.Name, " listening")
		var cid uint32
		cid = 0
		//3 阻塞等待客户端链接
		for {
			conn, err := listener.AcceptTCP()
			if err != nil {
				fmt.Println("accept err ", err)
				continue
			}

			// 客户端已经与服务器建立链接，业务处理
			dealConn := NewConnection(conn, cid, s.MsgHandler)
			cid++
			go dealConn.Start()
			//go func() {
			//	for {
			//		buf := make([]byte, 512)
			//		cnt, err2 := conn.Read(buf)
			//		if err2 != nil {
			//			fmt.Println("recv buf err ", err2)
			//			continue
			//		}
			//		fmt.Printf("server recv buf, %s, cnt: %d\n", buf, cnt)
			//
			//		if _, err := conn.Write(buf[:cnt]); err != nil {
			//			fmt.Println("write back buf err ", err)
			//			continue
			//		}
			//	}
			//}()
		}
	}()
}

func (s *Server) Stop() {
	//todo 资源回收

}

func (s *Server) Server() {
	//启动 server
	s.Start()

	// 阻塞
	select {}
}

func (s *Server) AddRouter(msgId uint32, router ziface.IRouter) {
	s.MsgHandler.AddRouter(msgId, router)
	fmt.Println("Add router success")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:       utils.GlobalObject.Name,
		IPVersion:  "tcp4",
		IP:         utils.GlobalObject.Host,
		Port:       utils.GlobalObject.TcpPort,
		MsgHandler: NewMsgHandle(),
	}
	return s
}
