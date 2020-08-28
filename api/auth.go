package api

import (
	"time"

	"github.com/pkg/errors"
)

type TenantAccessToken struct {
	BaseResponse
	AccessToken string    `json:"tenant_access_token"`
	Expire      int       `json:"expire"` // 还有多久失效，默认2小时
	ExpiredAt   time.Time // 到期时间
}

type AppAccessToken struct {
	BaseResponse
	AccessToken string `json:"app_access_token"`
	Expire      int    `json:"expire"` // 还有多久失效，默认2小时
}

type InternalApp struct {
	App *App
}

func NewInternalApp(app *App) *InternalApp {
	return &InternalApp{
		App: app,
	}
}

func (ia *InternalApp) Invoke() (token *TenantAccessToken, err error) {
	if len(ia.App.AppID) == 0 || len(ia.App.AppSecret) == 0 {
		err = errors.New("app_id, app_secret is required")
		return
	}

	client := NewClient(ia.App)
	client.IgnoreOnCallback()

	token = new(TenantAccessToken)
	_, err = client.Post("auth/v3/tenant_access_token/internal/", ia.App, token)
	if !token.OK() {
		err = errors.Errorf("request api 'tenant_access_token' was err: %d %v", token.Code, token.Msg)
		return
	}
	token.ExpiredAt = time.Now().Add(time.Duration(token.Expire-10) * time.Second) // 这里减10秒是为了进行容错，剩余10秒时当过期处理
	return token, errors.WithStack(err)
}

// 飞书应用商城的应用获取access_token(与自建应用流程差别较大)
// 需要拿到 ticket
type StoreApp struct {
	App            *App   `json:"-"`
	AppAccessToken string `json:"app_access_token"` // 用户走完OAuth2授权后拿到的授权码
	TenantKey      string `json:"tenant_key"`       // 商城应用 - 企业唯一key
}

func NewAppStore(app *App, tenantKey string) *StoreApp {
	return &StoreApp{
		App:       app,
		TenantKey: tenantKey,
	}
}

func (a *StoreApp) Invoke() (token *TenantAccessToken, err error) {
	if a.App.Ticket == "" {
		err = a.App.InitTicket()
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}

	if len(a.App.AppID) == 0 || len(a.App.AppSecret) == 0 || len(a.App.Ticket) == 0 {
		err = errors.New("app_id, app_secret, ticket is required")
		return
	}
	client := NewClient(a.App)
	client.IgnoreOnCallback()

	// 通过app & ticket拿到 app_access_token
	appToken := new(AppAccessToken)
	_, err = client.Post("auth/v3/app_access_token/", a.App, appToken)
	if err != nil {
		err = errors.WithStack(err)
		return
	}
	if !appToken.OK() {
		err = errors.Errorf("request api 'app_access_token' was err: %d %s", appToken.Code, appToken.Msg)
		return
	}
	a.AppAccessToken = appToken.AccessToken

	// 拿到 app_access_token 之后换取 tanant_access_token
	token = new(TenantAccessToken)
	_, err = client.Post("auth/v3/tenant_access_token/", a, &token)
	if !token.OK() {
		err = errors.Errorf("request api 'tenant_access_token' was err: %d %v", token.Code, token.Msg)
		return
	}
	token.ExpiredAt = time.Now().Add(time.Duration(token.Expire-10) * time.Second) // 这里减10秒是为了进行容错，剩余10秒时当过期处理
	return token, errors.WithStack(err)
}

type TenantTokenInvoker interface {
	Invoke() (token *TenantAccessToken, err error)
}
