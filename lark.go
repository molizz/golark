package golark

import (
	"github.com/molizz/golark/api"
	"github.com/molizz/golark/message"
)

var (
	OnApiError   = api.OnApiError
	OnApiSuccess = api.OnApiSuccess
	OnResponse   = api.OnResponse
)

// example
//
// func InitLark() error {
// 	if err := platform.LoadDefaultLarkConfig(); err != nil {
// 		return errors.Trace(err)
// 	}
// 	api.OnApiError = func(app *api.App, errDesc string) error {
// 		return nil
// 	}
// 	api.OnApiSuccess = func(app *api.App) error {
// 		return nil
// 	}
// 	api.OnResponse = func(method, url string, in, out interface{}) {
// 	}

// 	api.IsStoreApp = true // NOTE
// 	return nil
// }

func NewAPP(appID, appSecret, encryptKey, verificationToken string) *api.App {
	return api.NewApp(appID, appSecret, encryptKey, verificationToken)
}

func NewMessage(app *api.App, tenantKey *string, recv message.Receiver, msg message.Messager) *message.Message {
	return message.New(app, tenantKey, recv, msg)
}
