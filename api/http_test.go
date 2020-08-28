package api

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestClient_Do(t *testing.T) {
	app := NewApp("123", "321", "", "")

	outStruct := struct {
		Code int    `json:"code"`
		Msg  string `json:"msg"`
	}{}

	c := NewClient(app)
	_, err := c.do("POST", "auth/v3/tenant_access_token/internal/", app, &outStruct)

	assert.Equal(t, nil, err, fmt.Sprintf("Post() error = %v", err))
	assert.Equal(t, 10003, outStruct.Code, fmt.Sprintf("Post() got = %d expected = %d", outStruct.Code, 10003))

	_, err = c.Post("/auth/v3/tenant_access_token/internal/", app, &outStruct)

	assert.Equal(t, nil, err, fmt.Sprintf("Post() error = %v", err))
	assert.Equal(t, 10003, outStruct.Code, fmt.Sprintf("Post() got = %d expected = %d", outStruct.Code, 10003))
}
