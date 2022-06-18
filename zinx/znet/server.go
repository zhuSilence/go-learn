package znet

import (
	"fmt"
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
	// 当前 Server 添加 Router
	Router ziface.IRouter
}

// CallBackToClient 当前客户端绑定的方法
//func CallBackToClient(conn *net.TCPConn, data []byte, cnt int) error {
//	// 回写功能
//	fmt.Println("[Conn Handle] CallBackToClient...")
//	if _, err := conn.Write(data[:cnt]); err != nil {
//		fmt.Println("write back buf err", err)
//		return errors.New("CallBackToClient Error")
//	}
//	return nil
//}

func (s *Server) Start() {

	fmt.Printf("[Start] Server Listenner at IP: %s, Port %d, is starting\n", s.IP, s.Port)

	go func() {
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
			dealConn := NewConnection(conn, cid, s.Router)
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

func (s *Server) AddRouter(router ziface.IRouter) {
	s.Router = router
	fmt.Println("Add router success")
}

func NewServer(name string) ziface.IServer {
	s := &Server{
		Name:      name,
		IPVersion: "tcp4",
		IP:        "0.0.0.0",
		Port:      8999,
		Router:    nil,
	}
	return s
}
