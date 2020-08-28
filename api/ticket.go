package api

import "github.com/pkg/errors"

type Ticket struct {
	app *App
}

func NewTicket(app *App) *Ticket {
	return &Ticket{
		app: app,
	}
}

func (t *Ticket) Resend() error {
	resp := new(BaseResponse)
	_, err := NewClient(t.app).Post("auth/v3/app_ticket/resend/", t.app, resp)
	if err != nil {
		return errors.WithStack(err)
	}
	if !resp.OK() {
		return resp.Error()
	}
	return nil
}
