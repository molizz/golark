package api

import (
	"github.com/pkg/errors"
)

var (
	TicketFinder TicketHandler
)

type TicketHandler interface {
	GetTicket() (ticket string, err error)
}

type App struct {
	isStoreApp bool

	AppID     string `json:"app_id"`
	AppSecret string `json:"app_secret"`

	VerificationToken string `json:"verification_token,omitempty"`
	EncryptKey        string `json:"encrypt_key,omitempty"`

	Ticket string `json:"app_ticket,omitempty"`
}

func NewApp(appID, appSecret, encryptKey, verificationToken string) *App {
	return &App{
		isStoreApp:        false,
		AppID:             appID,
		AppSecret:         appSecret,
		EncryptKey:        encryptKey,
		VerificationToken: verificationToken,
	}
}

func (a *App) InitTicket() (err error) {
	a.Ticket, err = TicketFinder.GetTicket()
	if err != nil {
		return errors.WithStack(err)
	}
	return nil
}

// 飞书的「应用商城」模式
func (a *App) UseStoreAppModel() {
	a.isStoreApp = true
}

func (a *App) IsStoreApp() bool {
	return a.isStoreApp
}

func (a *App) Api(tenantKey *string) *LarkApi {
	return NewLarkApi(a, tenantKey)
}
