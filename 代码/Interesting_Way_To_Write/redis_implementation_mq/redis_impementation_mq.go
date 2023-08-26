package redis_implementation_mq

import (
	"fmt"
	"golang.org/x/net/context"
)

//消息队列
//1、消息所处的什么队列
//2、消息类型
//3、消息如何接收发送
//4、必须保证是先进后出

type INotify interface {
	Name() string
	Receive()
	Pop() <-chan []byte
	Push([]byte) error
	StopPop()
}

type Notify struct {
	*MQ
	name string
	data chan []byte
	stop chan bool
}

func NewNotify(name string) (notify *Notify, err error) {
	//获取队列
	mq := NewMQ(Config().A.GetAddr(), name)

	notify = &Notify{
		MQ:   mq,
		name: name,
		data: make(chan []byte),
		stop: make(chan bool),
	}

	return
}

func (n *Notify) GetMqName() string {
	return n.name
}

func (n *Notify) Pop() []byte {
	return <-n.data
}

func (n *Notify) Push(data []byte) error {
	intCmd := n.MQ.LPush(context.Background(), n.GetMqName(), string(data))
	if intCmd.Err() != nil {
		fmt.Println(fmt.Sprintf("%s push, data: %s, fail reason : %s ", n.GetMqName(), string(data), intCmd.Err().Error()))
	}
	return intCmd.Err()
}
