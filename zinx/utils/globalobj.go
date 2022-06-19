package utils

import (
	"encoding/json"
	"github.com/zhuSilence/go-learn/zinx/ziface"
	"io/ioutil"
)

type GlobalObj struct {
	TcpServer ziface.IServer // 当前 Zinx 全局的 Server 对象
	Host      string         // 当前服务器的 IP
	TcpPort   int            // 当前服务器的端口
	Name      string         // 当前服务器的名称

	Version          string // 当前 Zinx 的版本好
	MaxConn          int    // 当前服务器允许的最大链接数
	MaxPackageSize   int    // 当前 Zinx 框架数据包的最大值
	WorkerPoolSize   uint32 // 消息队列线程池大小
	MaxWorkerTaskLen uint32 // Zinx 线程池最大限制
}

// GlobalObject 定义一个全局的对象 GlobalObj
var GlobalObject *GlobalObj

// Reload 加载配置文件
func (g *GlobalObj) Reload() {
	// 读取配置文件
	data, err := ioutil.ReadFile("conf/zinx.json")
	if err != nil {
		panic(err)
	}
	// 将文件内容转换到 struct 中
	err = json.Unmarshal(data, &GlobalObject)
	if err != nil {
		panic(err)
	}
}

// init 初始化全局对象
func init() {
	// 默认配置
	GlobalObject = &GlobalObj{
		Name:             "ZinxServerApp",
		Version:          "V0.4",
		TcpPort:          8999,
		Host:             "0.0.0.0",
		MaxConn:          1000,
		MaxPackageSize:   4096,
		WorkerPoolSize:   10,
		MaxWorkerTaskLen: 1024,
	}

	// 读取配置文件
	GlobalObject.Reload()
}
