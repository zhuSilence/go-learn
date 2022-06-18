package znet

import (
	"bytes"
	"encoding/binary"
	"github.com/zhuSilence/go-learn/zinx/ziface"
)

type DataPack struct{}

//NewDataPack 封包拆包实例初始化方法
func NewDataPack() ziface.IDataPack {
	return &DataPack{}
}

func (dp *DataPack) GetHeadLen() uint32 {
	// dataLen uint32 (4个字节)
	// id uint32 (4个字节)
	return 8
}

func (dp *DataPack) Pack(msg ziface.IMessage) ([]byte, error) {
	// 创建一个存放 bytes 字节的缓冲
	dataBuf := bytes.NewBuffer([]byte{})
	// 将 dataLen 写进 dataBuf 中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetDataLen()); err != nil {
		return nil, err
	}
	// 将 msgId 写入 dataBuf 中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetMsgId()); err != nil {
		return nil, err
	}
	// 将 data 数据写入 dataBuf 中
	if err := binary.Write(dataBuf, binary.LittleEndian, msg.GetData()); err != nil {
		return nil, err
	}
	return dataBuf.Bytes(), nil
}

func (dp *DataPack) Unpack(binaryData []byte) (ziface.IMessage, error) {
	// 创建一个从输入二进制数据的 ioReader
	reader := bytes.NewReader(binaryData)
	// 创建 msg，用于存放数据和返回
	msg := &Message{}
	// 读取 dataLen
	if err := binary.Read(reader, binary.LittleEndian, &msg.DataLen); err != nil {
		return nil, err
	}
	// 读取 msgId
	if err := binary.Read(reader, binary.LittleEndian, &msg.Id); err != nil {
		return nil, err
	}
	// 判断 dataLen 是否超过范围
	//if utils.GlobalObject.MaxPackageSize > 0 && msg.DataLen > utils.GlobalObject.MaxPackageSize {
	//	return nil, errors.New("too large msg data")
	//}
	return msg, nil
}
