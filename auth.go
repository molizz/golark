package golark

import (
	"fmt"
	"strconv"
	"time"

	"github.com/pkg/errors"
)

// OAuth2
//
func OAuth2URL(appID, redirectURL string, state *string) (string, error) {
	if len(redirectURL) == 0 {
		return "", errors.New("missing redirect url")
	}
	if len(appID) == 0 {
		return "", errors.New("missing app id")
	}

	if state == nil {
		now := strconv.FormatInt(time.Now().Unix(), 10)
		state = &now
	}

	authURL := fmt.Sprintf("https://open.feishu.cn/open-apis/authen/v1/index?redirect_uri=%s&app_id=%s&state=%s",
		redirectURL, appID, *state)

	return authURL, nil
}
