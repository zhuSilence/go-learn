package main

import (
	"fmt"
	"github.com/zhuSilence/go-learn/zinx/znet"
	"io"
	"net"
	"time"
)

func main() {
	fmt.Println("client start...")
	time.Sleep(1 * time.Second)
	//1. 直接链接远程服务器，的带一个 conn 链接
	conn, err := net.Dial("tcp", "127.0.0.1:8999")
	if err != nil {
		fmt.Println("client start err:", err)
		return
	}
	//2. 链接调用 Write 写数据
	for {
		dp := znet.NewDataPack()
		binaryMsg, err := dp.Pack(znet.NewMsgPackage(0, []byte("ZinxV0.6 client0 test message")))
		if err != nil {
			fmt.Println("pack msg err", err)
			return
		}

		if _, err := conn.Write(binaryMsg); err != nil {
			fmt.Println("write msg err", err)
			return
		}

		// 读取服务器的返回
		binaryHead := make([]byte, dp.GetHeadLen())
		if _, err := io.ReadFull(conn, binaryHead); err != nil {
			fmt.Println("read msg head err", err)
			return
		}

		msgHead, err := dp.Unpack(binaryHead)
		if err != nil {
			fmt.Println("unpack msg head err", err)
			return
		}
		if msgHead.GetDataLen() > 0 {
			// message 有数据 进行二次读取
			msg := msgHead.(*znet.Message)
			msg.Data = make([]byte, msg.GetDataLen())

			if _, err := io.ReadFull(conn, msg.Data); err != nil {
				fmt.Println("read msg data err", err)
				return
			}

			fmt.Println(">>> recv MsgId:", msg.GetMsgId(), ", dataLen", msg.GetDataLen(), "data:", string(msg.GetData()))

		}

		time.Sleep(1 * time.Second)
	}

}
