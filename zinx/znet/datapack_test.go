package znet

import (
	"fmt"
	"io"
	"net"
	"testing"
)

func TestDataPack(t *testing.T) {
	// 模拟创建服务器
	// 1. 创建 socket tcp
	listener, err := net.Listen("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("server listen err", err)
		return
	}
	go func() {
		for {
			conn, err2 := listener.Accept()
			if err2 != nil {
				fmt.Println("server accept err", err2)
			}
			go func(conn net.Conn) {
				// 处理客户端请求，拆包，先取 dataLen，再读取数据
				dp := NewDataPack()
				for {
					headData := make([]byte, dp.GetHeadLen())
					_, err3 := io.ReadFull(conn, headData)
					if err3 != nil {
						fmt.Println("read head err", err3)
						break
					}

					msgHead, err := dp.Unpack(headData)
					if err != nil {
						fmt.Println("unpack err", err)
						return
					}
					// msg 有数据，进行第二次读取
					if msgHead.GetDataLen() > 0 {
						// 进行强转
						msg := msgHead.(*Message)
						msg.Data = make([]byte, msg.GetDataLen())
						// 根据 dataLen 从 io 中读取
						_, err := io.ReadFull(conn, msg.Data)
						if err != nil {
							fmt.Println("ReadFull err", err)
							return
						}
						// 完整的消息已经读取完毕
						fmt.Println(">>> recv MsgId:", msg.GetMsgId(), ", dataLen", msg.GetDataLen(), "data:", string(msg.GetData()))
					}

				}

			}(conn)
		}
	}()

	// 2. 从客户端读取数据，拆包处理
	dp := NewDataPack()

	conn, err := net.Dial("tcp", "127.0.0.1:7777")
	if err != nil {
		fmt.Println("client dial err:", err)
		return
	}

	msg1 := &Message{
		Id:      1,
		DataLen: 4,
		Data:    []byte{'z', 'i', 'n', 'x'},
	}
	sendData1, err := dp.Pack(msg1)
	if err != nil {
		fmt.Println("client pack msg1 err:", err)
		return
	}

	msg2 := &Message{
		Id:      2,
		DataLen: 7,
		Data:    []byte{'n', 'i', 'h', 'a', 'o', '!', '!'},
	}
	sendData2, err := dp.Pack(msg2)
	if err != nil {
		fmt.Println("client pack msg2 err:", err)
		return
	}
	sendData1 = append(sendData1, sendData2...)
	conn.Write(sendData1)
	// 阻塞
	select {}
}
