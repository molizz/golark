package message

import (
	"github.com/molizz/golark/api"
	"github.com/molizz/golark/message/card"
	"github.com/pkg/errors"
)

type Receiver interface {
	Type() string               // 接收者类型（可能是用户ids、部门ids、Open ids）
	ParamName() string          // 接收者http中的body参数名（例如 open_ids、user_ids、department_ids）
	List() (receivers []string) //
}

type Messager interface {
	Kind() int        // 消息分类
	Type() string     // 消息类型（通常是交互消息（interactive））
	Card() *card.Card // 消息卡片内容
}

type Message struct {
	tenantKey *string
	app       *api.App
	receiver  Receiver
	msg       Messager
}

func New(app *api.App, tenantKey *string, receiver Receiver, msg Messager) *Message {
	return &Message{
		tenantKey: tenantKey,
		app:       app,
		receiver:  receiver,
		msg:       msg,
	}
}

func (m *Message) Send() error {
	msgApi, err := m.app.Api(m.tenantKey).Message()
	if err != nil {
		return errors.WithStack(err)
	}
	return m.send(msgApi, m.receiver, m.msg)
}

// TODO 推送信息数量可能较多，理论上应该做队列处理
//
func (m *Message) send(msgApi *api.Message, receiver Receiver, msg Messager) error {
	_, err := msgApi.Send(msg.Type(), msg.Card(), receiver.List()...)
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}
