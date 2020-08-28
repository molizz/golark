package api

import (
	"github.com/pkg/errors"
)

type LarkApi struct {
	authInvoker TenantTokenInvoker

	app       *App
	tenantKey string // SaaS 下飞书团队的唯一标识
}

// NewLarkApi
// 	isInternal 是否为自建应用
//	tenantKey 应用商城的企业key
//
func NewLarkApi(app *App, tenantKey *string) *LarkApi {
	return newLarkApi(app, tenantKey)
}

func NewLarkApiByApp(app *App) *LarkApi {
	return newLarkApi(app, nil)
}

func newLarkApi(app *App, tenantKey *string) *LarkApi {
	mustTenantKey := ""
	if tenantKey != nil {
		mustTenantKey = *tenantKey
	}

	lark := &LarkApi{}
	if app.IsStoreApp() {
		lark.authInvoker = NewAppStore(app, mustTenantKey)
	} else {
		lark.authInvoker = NewInternalApp(app)
	}
	lark.app = app
	lark.tenantKey = mustTenantKey
	return lark
}

func (lark *LarkApi) User() (*User, error) {
	token, err := lark.token()
	if err != nil {
		return nil, errors.WithStack(err)
	}
	return NewUser(lark.app, token.AccessToken), nil
}

func (lark *LarkApi) Ticket() *Ticket {
	return NewTicket(lark.app)
}

func (lark *LarkApi) Message() (*Message, error) {
	token, err := lark.token()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	return NewMessage(lark.app, token.AccessToken), nil
}

func (lark *LarkApi) token() (token *TenantAccessToken, err error) {
	if lark.app.IsStoreApp() && len(lark.tenantKey) == 0 {
		return nil, errors.New("missing tenant_key")
	}

	tokenKey := lark.app.AppID + lark.tenantKey

	token = tokenManager.Get(tokenKey)
	if token == nil {
		token, err = lark.authInvoker.Invoke()
		if err != nil {
			err = errors.WithStack(err)
			return
		}
		tokenManager.Set(tokenKey, token)
	}
	return
}
