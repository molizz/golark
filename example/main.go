package main

import (
	"github.com/molizz/golark"
	"github.com/molizz/golark/api"
	"github.com/molizz/golark/message/messager"
	"github.com/molizz/golark/message/receiver"
	"github.com/molizz/golark/utils"
)

func main() {

	app := golark.NewAPP("xxxx",
		"xxxxx",
		"xxxxx",
		"xxxxx")

	// 注册「事件订阅」响应处理
	golark.RegisterEventProcessor(&MessageEvent{})

	// 该方法通常在http的路由方法中调用
	golark.OnEventCallback(app, nil)

	// print api error
	golark.OnApiError = func(app *api.App, errDesc string) error {
		return nil
	}

	// print log
	utils.DefaultLog.Println("hello lark")

	// send message
	recv := receiver.NewUserReceiver("ou_a64c5416846d409eb4faa549248150d3")
	msg := messager.NewExampleMessage("hello", "desc", "https://mozz.in")
	_ = golark.NewMessage(app, nil, recv, msg).Send()
}

// 当用户给机器人发送时，响应的事件type
type MessageEvent struct {
}

func (m *MessageEvent) TypeLabel() string {
	return "message"
}

func (m *MessageEvent) Process(eventBody []byte) error {
	// eventBody 消息body
	return nil
}
