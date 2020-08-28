package api

import (
	"fmt"

	"github.com/pkg/errors"
)

type BaseResponse struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (b *BaseResponse) OK() bool {
	return b.Code == 0
}

func (b *BaseResponse) Error() error {
	if !b.OK() {
		msg, ok := LarkErrorSet[b.Code]
		if !ok {
			msg = fmt.Sprintf("未知错误：%d:%s", b.Code, b.Msg)
		}
		return errors.New(msg)
	}
	return nil
}

func (b *BaseResponse) Self() *BaseResponse {
	return b
}

// for 断言
type SelfInvoker interface {
	Self() *BaseResponse
}
