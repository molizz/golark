package api

type SendMessageResponse struct {
	BaseResponse
	Data struct {
		MessageID string `json:"message_id"`
	} `json:"data"`
}

type Message struct {
	app         *App
	accessToken string
}

func NewMessage(app *App, accessToken string) *Message {
	return &Message{
		app:         app,
		accessToken: accessToken,
	}
}

func (m *Message) Send(msgType string, card interface{}, openIDs ...string) (resp *SendMessageResponse, err error) {
	if len(openIDs) == 0 {
		return nil, nil
	}
	body := map[string]interface{}{
		"open_ids": openIDs, // TODO 这里应该是open_ids还是部门ids，应该由上级提供
		"msg_type": msgType,
		"card":     card,
	}

	resp = new(SendMessageResponse)
	_, err = NewClientWithBearer(m.app, m.accessToken).Post("message/v4/batch_send/", body, resp)
	if !resp.OK() {
		return nil, resp.Error()
	}
	return
}
