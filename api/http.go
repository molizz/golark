package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/molizz/golark/utils"
	"github.com/pkg/errors"
)

const (
	baseURL = "https://open.feishu.cn/open-apis/"
)

var (
	// 当接口请求成功
	OnApiSuccess func(app *App) error
	// 当接口请求失败
	OnApiError func(app *App, errDesc string) error
	// 当接口请求完成后
	OnResponse func(method, url string, in, out interface{})
)

type Client struct {
	// 忽略调用OnApiSuccess、OnApiError等消息
	ignoreCallback bool

	client *http.Client

	app *App
}

func NewClient(app *App) *Client {
	return newClient(app, "")
}

func NewClientWithBearer(app *App, bearerToken string) *Client {
	return newClient(app, bearerToken)
}

func newClient(app *App, bearerToken string) *Client {
	var trans = http.DefaultTransport

	if len(bearerToken) > 0 {
		trans = &BearerToken{
			Base:  http.DefaultTransport,
			Token: bearerToken,
		}
	}

	client := &http.Client{
		Transport: trans,
		Timeout:   10 * time.Second,
	}

	return &Client{
		ignoreCallback: false,
		client:         client,
		app:            app,
	}
}

func (c *Client) IgnoreOnCallback() {
	c.ignoreCallback = true
}

func (c *Client) Post(path string, body, out interface{}) (*http.Response, error) {
	return c.do("POST", path, body, out)
}

func (c *Client) Get(path string, out interface{}) (*http.Response, error) {
	return c.do("GET", path, nil, out)
}

func (c *Client) do(method, path string, in, out interface{}) (resp *http.Response, err error) {
	if strings.HasPrefix(path, "/") {
		path = path[1:]
	}
	fullURL := baseURL + path

	req, _ := http.NewRequest(method, fullURL, nil)
	if in != nil {
		buf := new(bytes.Buffer)
		_ = json.NewEncoder(buf).Encode(in)
		req.Header.Set("Content-Type", "application/json")
		req.Body = ioutil.NopCloser(buf)
	}

	resp, err = c.client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	defer func() {
		_ = resp.Body.Close()
		c.onDone(method, fullURL, in, out)
	}()

	if out != nil {
		if w, ok := out.(io.Writer); ok {
			_, _ = io.Copy(w, resp.Body)
			return resp, nil
		}
		return resp, json.NewDecoder(resp.Body).Decode(out)
	}

	return resp, nil
}

func (c *Client) onDone(method, fullURL string, in, out interface{}) {
	c.onResponse(method, fullURL, in, out)

	if c.ignoreCallback {
		utils.DefaultLog.Printf("ignore, lark request %s\n", fullURL)
		return
	}
	c.onApiError(out)
	c.onApiSuccess(out)
}

func (c *Client) onResponse(method, url string, in, out interface{}) {
	if OnResponse == nil {
		return
	}
	OnResponse(method, url, in, out)
}

func (c *Client) onApiError(out interface{}) {
	if OnApiError == nil {
		return
	}

	if e, ok := out.(SelfInvoker); ok && !e.Self().OK() {
		var errCode = e.Self().Code
		var errString, ok = LarkErrorSet[errCode]
		if !ok {
			errString = e.Self().Msg
		}
		_ = OnApiError(c.app, fmt.Sprintf("%d:%s", errCode, errString))
	}
}

func (c *Client) onApiSuccess(out interface{}) {
	if OnApiSuccess == nil {
		return
	}

	if e, ok := out.(SelfInvoker); ok && e.Self().OK() {
		_ = OnApiSuccess(c.app)
	}
}
