package golark

import (
	"encoding/json"

	"github.com/molizz/golark/api"
	"github.com/molizz/golark/callback"
	"github.com/molizz/golark/utils"
	"github.com/pkg/errors"
)

func init() {
	callbackHub.register(callback.NewVerification())
	callbackHub.register(callback.NewCallback(eventHub))
}

type EncryptedEventRequest struct {
	Encrypt string `json:"encrypt"`
}

// 事件处理器
type CallbackProcessor interface {
	TypeLabel() string
	Process(params string) (interface{}, error)
}

var callbackHub = &CallbackHub{
	processorMap: make(map[string]CallbackProcessor),
}

type CallbackHub struct {
	processorMap map[string]CallbackProcessor
}

func (c *CallbackHub) register(p CallbackProcessor) {
	c.processorMap[p.TypeLabel()] = p
}

func (c *CallbackHub) find(typeLabel string) CallbackProcessor {
	return c.processorMap[typeLabel]
}

type EventType struct {
	Type  string `json:"type"`
	Token string `json:"token"`
}

func OnEventCallback(app *api.App, req *EncryptedEventRequest) (resp interface{}, err error) {
	if app == nil {
		return nil, errors.New("app is required")
	}
	key := app.EncryptKey
	token := app.VerificationToken

	// 解码数据
	body, err := utils.AesDecrypt(key, req.Encrypt)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	utils.DefaultLog.Printf("lark event callback: '%s'\n", body)

	// 序列化数据
	var eventType EventType
	err = json.Unmarshal([]byte(body), &eventType)
	if err != nil {
		return nil, errors.WithStack(err)
	}

	if eventType.Token != token {
		return nil, errors.New("event token not match")
	}

	proc := callbackHub.find(eventType.Type)
	if proc == nil {
		utils.DefaultLog.Printf("invalid event type '%s'.\n", eventType.Type)
		return nil, nil
	}
	return proc.Process(body)
}
