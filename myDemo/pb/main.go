package main

import (
	"fmt"
	"github.com/aceld/zinx/examples/zinx_version_ex/protoDemo/pb"
	"github.com/golang/protobuf/proto"
)

func main() {
	person := &pb.Person{
		Name:   "silence",
		Age:    18,
		Emails: []string{""},
		Phones: []*pb.PhoneNumber{
			&pb.PhoneNumber{Number: "1111", Type: pb.PhoneType_HOME},
			&pb.PhoneNumber{Number: "222", Type: pb.PhoneType_WORK},
			&pb.PhoneNumber{Number: "333", Type: pb.PhoneType_MOBILE},
		},
	}
	// 编码
	// 将 Person 序列化，data 就是需要网络传输的数据
	data, err := proto.Marshal(person)
	if err != nil {
		fmt.Println("error", err)
	}

	//解码
	newData := &pb.Person{}
	err = proto.Unmarshal(data, newData)
	if err != nil {
		fmt.Println("unmarshal err", err)
	}
	fmt.Println(newData)
}
